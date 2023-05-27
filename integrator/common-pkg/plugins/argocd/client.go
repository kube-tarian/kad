package argocd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	cmdutil "github.com/argoproj/argo-cd/v2/cmd/util"
	"github.com/argoproj/argo-cd/v2/common"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	clusterpkg "github.com/argoproj/argo-cd/v2/pkg/apiclient/cluster"
	repocredspkg "github.com/argoproj/argo-cd/v2/pkg/apiclient/repocreds"
	repositorypkg "github.com/argoproj/argo-cd/v2/pkg/apiclient/repository"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	appsv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	argoappv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/argoproj/argo-cd/v2/util/cli"
	"github.com/argoproj/argo-cd/v2/util/clusterauth"
	"github.com/argoproj/argo-cd/v2/util/errors"
	"github.com/argoproj/argo-cd/v2/util/git"
	"github.com/argoproj/argo-cd/v2/util/io"
	"github.com/argoproj/argo-cd/v2/util/text/label"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/kube-tarian/kad/integrator/common-pkg/plugins/fetcher"
	"github.com/kube-tarian/kad/integrator/model"
)

type ArgoCDCLient struct {
	conf   *Configuration
	logger logging.Logger
	client apiclient.Client
}

func NewClient(logger logging.Logger) (*ArgoCDCLient, error) {
	cfg, err := fetchConfiguration(logger)
	if err != nil {
		return nil, err
	}

	if cfg.IsSSLEnabled {
		// TODO: Configure SSL certificates
		logger.Errorf("SSL not yet supported, continuing with insecure verify true")
	}

	client, err := getNewAPIClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ArgoCDCLient{
		conf:   cfg,
		logger: logger,
		client: client,
	}, nil
}

func (a *ArgoCDCLient) ConfigurationActivities(req interface{}) (json.RawMessage, error) {
	payload, _ := req.(model.ConfigPayload)
	switch payload.Resource {
	case "cluster":
		return a.HandleCluster(req)
	case "repo":
		return a.HandleRepo(req)
	default:
		return nil, fmt.Errorf("unsupported action for argocd plugin: %v", payload.Action)
	}
}

func (a *ArgoCDCLient) HandleCluster(req interface{}) (json.RawMessage, error) {
	payload, _ := req.(model.ConfigPayload)
	switch payload.Action {
	case "add":
		return a.ClusterAdd(&payload)
	case "delete":
		// return a.ClusterDelete(payload)
	case "list":
		// return a.ClusterList(payload)
	default:
		return nil, fmt.Errorf("unsupported action for argocd plugin: %v", payload.Action)
	}
	return nil, nil
}

func (a *ArgoCDCLient) HandleRepo(req interface{}) (json.RawMessage, error) {
	payload, _ := req.(model.ConfigPayload)
	switch payload.Action {
	case "add":
		return a.RepoAdd(&payload)
	case "delete":
		return a.RepoDelete(&payload)
	case "list":
		// return a.RepoList(payload)
	default:
		return nil, fmt.Errorf("unsupported action for argocd plugin: %v", payload.Action)
	}
	return nil, nil
}

func (a *ArgoCDCLient) HandleRepoCreds(req interface{}) (json.RawMessage, error) {
	payload, _ := req.(model.ConfigPayload)
	switch payload.Action {
	case "add":
		return a.RepoCredsAdd(&payload)
	case "delete":
		return a.RepoCredsDelete(&payload)
	case "list":
		// return a.RepoCRedsList(payload)
	default:
		return nil, fmt.Errorf("unsupported action for argocd plugin: %v", payload.Action)
	}
	return nil, nil
}

func fetchConfiguration(log logging.Logger) (*Configuration, error) {
	// If ARGOCD_PASSWORD env variable is configured then it will use local default configuration
	// Else it uses fetched to get the plugin details and prepares the configuration
	cfg := &Configuration{}
	err := envconfig.Process("", cfg)
	if err == nil {
		return cfg, err
	}

	fetcherClient, err := fetcher.NewCredentialFetcher(log)
	if err != nil {
		log.Errorf("fetcher client initialization failed: %v", err)
		return nil, err
	}

	response, err := fetcherClient.FetchPluginDetails(&fetcher.PluginRequest{
		PluginName: "argocd",
	})
	if err != nil {
		log.Errorf("Failed to get the plugin details: %v", err)
		return nil, err
	}
	cfg = &Configuration{
		ServiceURL:   response.ServiceURL,
		IsSSLEnabled: response.IsSSLEnabled,
		Username:     response.Username,
		Password:     response.Password,
	}
	return cfg, err
}

func (a *ArgoCDCLient) Create(req *model.CreteRequestPayload) (json.RawMessage, error) {
	conn, appClient, err := a.client.NewApplicationClient()
	if err != nil {
		a.logger.Errorf("Application client intilialization failed: %v", err)
		return nil, err
	}
	defer io.Close(conn)

	resp, err := appClient.Create(
		context.Background(),
		&application.ApplicationCreateRequest{
			Application: &v1alpha1.Application{
				ObjectMeta: v1.ObjectMeta{
					Name: req.ReleaseName,
				},
				Spec: v1alpha1.ApplicationSpec{
					Destination: v1alpha1.ApplicationDestination{
						Namespace: req.Namespace,
						Server:    "https://kubernetes.default.svc",
					},
					Source: v1alpha1.ApplicationSource{
						RepoURL:        req.RepoURL,
						Path:           req.ChartName,
						TargetRevision: "HEAD",
					},
					SyncPolicy: &v1alpha1.SyncPolicy{
						Automated: &v1alpha1.SyncPolicyAutomated{
							Prune:    false,
							SelfHeal: false,
						},
					},
					Project: "default",
				},
			},
		})
	if err != nil {
		a.logger.Errorf("Application %s install failed: %v", req.ReleaseName, err)
		return nil, err
	}

	respMsg, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	// a.logger.Infof("argo-cd msg: %s", string(respMsg))
	return respMsg, nil
}

func (a *ArgoCDCLient) Delete(req *model.DeleteRequestPayload) (json.RawMessage, error) {
	conn, appClient, err := a.client.NewApplicationClient()
	if err != nil {
		return nil, err
	}
	defer io.Close(conn)

	resp, err := appClient.Delete(
		context.Background(),
		&application.ApplicationDeleteRequest{
			Name:         &req.ReleaseName,
			AppNamespace: &req.Namespace,
		},
	)
	if err != nil {
		return nil, err
	}

	respMsg, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return respMsg, nil
}

func (a *ArgoCDCLient) List(req *model.ListRequestPayload) (json.RawMessage, error) {
	conn, appClient, err := a.client.NewApplicationClient()
	if err != nil {
		return nil, err
	}
	defer io.Close(conn)

	list, err := appClient.List(context.Background(), &application.ApplicationQuery{})
	if err != nil {
		return nil, err
	}

	listMsg, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}
	return listMsg, nil
}

func (a *ArgoCDCLient) ClusterAdd(req *model.ConfigPayload) (json.RawMessage, error) {
	opts, ok := req.Data["clusterOpts"]
	if !ok {
		return nil, fmt.Errorf("repo options not present in data")
	}
	clusterOpts := opts.(*cmdutil.ClusterOptions)
	ctx := context.Background()

	conn, clusterClient, err := a.client.NewClusterClient()
	if err != nil {
		a.logger.Errorf("Application client intilialization failed: %v", err)
		return nil, err
	}
	defer io.Close(conn)
	var labels []string
	var annotations []string
	cn, ok := req.Data["contextName"]
	if !ok {
		return nil, fmt.Errorf("failed to get context name from config payload")
	}
	contextName, ok := cn.(string)
	if !ok {
		return nil, fmt.Errorf("context name type is not string")
	}
	po, ok := req.Data["pathOpts"]
	if !ok {
		return nil, fmt.Errorf("failed to get path options from config payload")
	}
	pathOpts, ok := po.(*clientcmd.PathOptions)
	if !ok {
		return nil, fmt.Errorf("path options type is not string")
	}
	conf, err := getRestConfig(pathOpts, contextName)
	if err != nil {
		a.logger.Errorf("failed to get rest config %w", err)
		return nil, err
	}
	managerBearerToken := ""
	// Install RBAC resources for managing the cluster
	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		a.logger.Errorf("failed to get rest config %w", err)
		return nil, err
	}

	accessLevel := "cluster"
	if len(clusterOpts.Namespaces) > 0 {
		accessLevel = "namespace"
	}
	fmt.Printf("WARNING: This will create a service account `argocd-manager` on the cluster referenced by context `%s` with full %s level privileges. Do you want to continue [y/N]? ", contextName, accessLevel)
	managerBearerToken, err = clusterauth.InstallClusterManagerRBAC(clientset, clusterOpts.SystemNamespace, clusterOpts.Namespaces, common.BearerTokenTimeout)

	labelsMap, err := label.Parse(labels)
	if err != nil {
		a.logger.Errorf("failed to parse labels %w", err)
		return nil, err
	}
	annotationsMap, err := label.Parse(annotations)
	if err != nil {
		a.logger.Errorf("failed to parse annotations %w", err)
		return nil, err
	}

	//conn, clusterIf := headless.NewClientOrDie(clientOpts, c).NewClusterClientOrDie()
	//defer io.Close(conn)
	if clusterOpts.Name != "" {
		contextName = clusterOpts.Name
	}
	clst := cmdutil.NewCluster(contextName, clusterOpts.Namespaces, clusterOpts.ClusterResources, conf, managerBearerToken, nil, nil, labelsMap, annotationsMap)
	if clusterOpts.InCluster {
		clst.Server = argoappv1.KubernetesInternalAPIServerAddr
	}
	if clusterOpts.Shard >= 0 {
		clst.Shard = &clusterOpts.Shard
	}
	if clusterOpts.Project != "" {
		clst.Project = clusterOpts.Project
	}
	clstCreateReq := clusterpkg.ClusterCreateRequest{
		Cluster: clst,
		Upsert:  clusterOpts.Upsert,
	}
	c, err := clusterClient.Create(ctx, &clstCreateReq)
	if err != nil {
		a.logger.Errorf("failed to create cluster client %w", err)
		return nil, err
	}
	fmt.Printf("Cluster '%s' added\n", clst.Server)
	cluster, err := json.Marshal(c)
	if err != nil {
		a.logger.Errorf("failed to marshal newly created cluster  %w", err)
		return nil, err
	}

	return cluster, nil
}

func (a *ArgoCDCLient) ClusterDelete(req *model.ConfigPayload) (json.RawMessage, error) {
	ctx := context.Background()
	conn, clusterClient, err := a.client.NewClusterClient()
	if err != nil {
		a.logger.Errorf("Application client intilialization failed: %v", err)
		return nil, err
	}
	defer io.Close(conn)

	cs, ok := req.Data["clusterSelector"]
	if !ok {
		return nil, fmt.Errorf("failed to get clusterSelector from config payload")
	}
	clusterSelector, ok := cs.(string)
	if !ok {
		return nil, fmt.Errorf("clusterSelector type is not string")
	}

	clusterQuery := getQueryBySelector(clusterSelector)

	// send server or cluster name in data field in configPayload
	clst, err := clusterClient.Get(ctx, clusterQuery)
	if err != nil {
		a.logger.Errorf("failed to get cluster client %w", err)
		return nil, err
	}

	// remove cluster
	clusterResp, err := clusterClient.Delete(ctx, clusterQuery)
	if err != nil {
		a.logger.Errorf("failed to delete cluster %w", err)
		return nil, err
	}
	fmt.Printf("Cluster '%s' removed\n", clusterSelector)

	po, ok := req.Data["pathOpts"]
	if !ok {
		return nil, fmt.Errorf("failed to get clusterSelector from config payload")
	}
	pathOpts, ok := po.(*clientcmd.PathOptions)
	if !ok {
		return nil, fmt.Errorf("clusterSelector type is not string")
	}

	// remove RBAC from cluster
	conf, err := getRestConfig(pathOpts, clst.Name)
	if err != nil {
		a.logger.Errorf("failed to get rest config %w", err)
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		a.logger.Errorf("failed to get client config %w", err)
		return nil, err
	}

	err = clusterauth.UninstallClusterManagerRBAC(clientset)
	if err != nil {
		a.logger.Errorf("failed to uninstall cluster manager RBAC %w", err)
		return nil, err
	}
	cr, err := json.Marshal(clusterResp)
	if err != nil {
		a.logger.Errorf("failed to marshal newly created cluster  %w", err)
		return nil, err
	}
	return cr, nil
}

func (a *ArgoCDCLient) RepoAdd(req *model.ConfigPayload) (json.RawMessage, error) {
	opts, ok := req.Data["repoOpts"]
	if !ok {
		return nil, fmt.Errorf("repo options not present in data")
	}
	repoOpts := opts.(cmdutil.RepoOptions)
	ctx := context.Background()

	// specify Repository URL as below in Data field of configPayload
	//repoOpts.Repo.Repo = repoURL

	// Specifying ssh-private-key-path is only valid for SSH repositories
	if repoOpts.SshPrivateKeyPath != "" {
		if ok, _ := git.IsSSHURL(repoOpts.Repo.Repo); ok {
			keyData, err := os.ReadFile(repoOpts.SshPrivateKeyPath)
			if err != nil {
				return nil, err
			}
			repoOpts.Repo.SSHPrivateKey = string(keyData)
		} else {
			err := fmt.Errorf("--ssh-private-key-path is only supported for SSH repositories.")
			return nil, err
		}
	}

	// tls-client-cert-path and tls-client-cert-key-key-path must always be
	// specified together
	if (repoOpts.TlsClientCertPath != "" && repoOpts.TlsClientCertKeyPath == "") || (repoOpts.TlsClientCertPath == "" && repoOpts.TlsClientCertKeyPath != "") {
		err := fmt.Errorf("--tls-client-cert-path and --tls-client-cert-key-path must be specified together")
		return nil, err
	}

	// Specifying tls-client-cert-path is only valid for HTTPS repositories
	if repoOpts.TlsClientCertPath != "" {
		if git.IsHTTPSURL(repoOpts.Repo.Repo) {
			tlsCertData, err := os.ReadFile(repoOpts.TlsClientCertPath)
			errors.CheckError(err)
			tlsCertKey, err := os.ReadFile(repoOpts.TlsClientCertKeyPath)
			errors.CheckError(err)
			repoOpts.Repo.TLSClientCertData = string(tlsCertData)
			repoOpts.Repo.TLSClientCertKey = string(tlsCertKey)
		} else {
			err := fmt.Errorf("--tls-client-cert-path is only supported for HTTPS repositories")
			errors.CheckError(err)
		}
	}

	// Specifying github-app-private-key-path is only valid for HTTPS repositories
	if repoOpts.GithubAppPrivateKeyPath != "" {
		if git.IsHTTPSURL(repoOpts.Repo.Repo) {
			githubAppPrivateKey, err := os.ReadFile(repoOpts.GithubAppPrivateKeyPath)
			if err != nil {
				return nil, err
			}
			repoOpts.Repo.GithubAppPrivateKey = string(githubAppPrivateKey)
		} else {
			err := fmt.Errorf("--github-app-private-key-path is only supported for HTTPS repositories")
			return nil, err
		}
	}

	// Set repository connection properties only when creating repository, not
	// when creating repository credentials.
	repoOpts.Repo.Insecure = repoOpts.InsecureSkipServerVerification
	repoOpts.Repo.EnableLFS = repoOpts.EnableLfs
	repoOpts.Repo.EnableOCI = repoOpts.EnableOci
	repoOpts.Repo.GithubAppId = repoOpts.GithubAppId
	repoOpts.Repo.GithubAppInstallationId = repoOpts.GithubAppInstallationId
	repoOpts.Repo.GitHubAppEnterpriseBaseURL = repoOpts.GitHubAppEnterpriseBaseURL
	repoOpts.Repo.Proxy = repoOpts.Proxy

	if repoOpts.Repo.Type == "helm" && repoOpts.Repo.Name == "" {
		err := fmt.Errorf("Must specify --name for repos of type 'helm'")
		return nil, err
	}

	conn, repoClient, err := a.client.NewRepoClient()
	if err != nil {
		a.logger.Errorf("Application client intilialization failed: %v", err)
		return nil, err
	}

	defer io.Close(conn)

	// If the user set a username, but didn't supply password via --password,
	// then we prompt for it
	if repoOpts.Repo.Username != "" && repoOpts.Repo.Password == "" {
		repoOpts.Repo.Password = cli.PromptPassword(repoOpts.Repo.Password)
	}

	// We let the server check access to the repository before adding it. If
	// it is a private repo, but we cannot access with the credentials
	// that were supplied, we bail out.
	//
	// Skip validation if we are just adding credentials template, chances
	// are high that we do not have the given URL pointing to a valid Git
	// repo anyway.
	repoAccessReq := repositorypkg.RepoAccessQuery{
		Repo:                       repoOpts.Repo.Repo,
		Type:                       repoOpts.Repo.Type,
		Name:                       repoOpts.Repo.Name,
		Username:                   repoOpts.Repo.Username,
		Password:                   repoOpts.Repo.Password,
		SshPrivateKey:              repoOpts.Repo.SSHPrivateKey,
		TlsClientCertData:          repoOpts.Repo.TLSClientCertData,
		TlsClientCertKey:           repoOpts.Repo.TLSClientCertKey,
		Insecure:                   repoOpts.Repo.IsInsecure(),
		EnableOci:                  repoOpts.Repo.EnableOCI,
		GithubAppPrivateKey:        repoOpts.Repo.GithubAppPrivateKey,
		GithubAppID:                repoOpts.Repo.GithubAppId,
		GithubAppInstallationID:    repoOpts.Repo.GithubAppInstallationId,
		GithubAppEnterpriseBaseUrl: repoOpts.Repo.GitHubAppEnterpriseBaseURL,
		Proxy:                      repoOpts.Proxy,
		Project:                    repoOpts.Repo.Project,
	}
	_, err = repoClient.ValidateAccess(ctx, &repoAccessReq)
	// credentials can be added later on, so we just log the error in case repo cannot be accessed
	a.logger.Errorf("failed to access repo %w", err)

	repoCreateReq := repositorypkg.RepoCreateRequest{
		Repo:   &repoOpts.Repo,
		Upsert: repoOpts.Upsert,
	}

	createdRepo, err := repoClient.CreateRepository(ctx, &repoCreateReq)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Repository '%s' added\n", createdRepo.Repo)
	cr, err := json.Marshal(createdRepo)
	if err != nil {
		a.logger.Errorf("failed to marshal newly added repo details  %w", err)
		return nil, err
	}
	return cr, nil
}

func (a *ArgoCDCLient) RepoDelete(req *model.ConfigPayload) (json.RawMessage, error) {
	opts, ok := req.Data["repoOpts"]
	if !ok {
		return nil, fmt.Errorf("repo options not present in data")
	}
	repoOpts := opts.(cmdutil.RepoOptions)
	ctx := context.Background()

	repoURL := repoOpts.Repo.Repo
	if repoURL == "" {
		return nil, fmt.Errorf("repoURL is empty")
	}
	conn, repoClient, err := a.client.NewRepoClient()
	if err != nil {
		a.logger.Errorf("Application client intilialization failed: %v", err)
		return nil, err
	}

	defer io.Close(conn)

	deletedRepo, err := repoClient.DeleteRepository(ctx, &repositorypkg.RepoQuery{Repo: repoURL})
	if err != nil {
		return nil, err
	}
	fmt.Printf("Repository '%s' removed\n", repoURL)
	dr, err := json.Marshal(deletedRepo)
	if err != nil {
		a.logger.Errorf("failed to marshal deleted repo details  %w", err)
		return nil, err
	}
	return dr, nil
}

func (a *ArgoCDCLient) RepoCredsAdd(req *model.ConfigPayload) (json.RawMessage, error) {
	var (
		repo   appsv1.RepoCreds
		upsert bool
	)
	opts, ok := req.Data["repoOpts"]
	if !ok {
		return nil, fmt.Errorf("repo options not present in data")
	}
	repoOpts := opts.(cmdutil.RepoOptions)
	ctx := context.Background()
	//ru, ok := req.Data["repoURL"]
	//if !ok {
	//	return nil, fmt.Errorf("failed to get repoURL from config payload")
	//}
	//repoURL, ok := ru.(string)
	//if !ok {
	//	return nil, fmt.Errorf("repoURL type is not string")
	//}
	// Repository URL
	repo.URL = repoOpts.Repo.Repo
	// Specifying ssh-private-key-path is only valid for SSH repositories
	if repoOpts.SshPrivateKeyPath != "" {
		if ok, _ := git.IsSSHURL(repo.URL); ok {
			keyData, err := os.ReadFile(repoOpts.SshPrivateKeyPath)
			if err != nil {
				return nil, err
			}
			repo.SSHPrivateKey = string(keyData)
		} else {
			err := fmt.Errorf("--ssh-private-key-path is only supported for SSH repositories.")
			return nil, err
		}
	}

	// tls-client-cert-path and tls-client-cert-key-key-path must always be
	// specified together
	if (repoOpts.TlsClientCertPath != "" && repoOpts.TlsClientCertKeyPath == "") || (repoOpts.TlsClientCertPath == "" && repoOpts.TlsClientCertKeyPath != "") {
		err := fmt.Errorf("--tls-client-cert-path and --tls-client-cert-key-path must be specified together")
		return nil, err
	}

	// Specifying tls-client-cert-path is only valid for HTTPS repositories
	if repoOpts.TlsClientCertPath != "" {
		if git.IsHTTPSURL(repo.URL) {
			tlsCertData, err := os.ReadFile(repoOpts.TlsClientCertPath)
			if err != nil {
				return nil, err
			}
			tlsCertKey, err := os.ReadFile(repoOpts.TlsClientCertKeyPath)
			if err != nil {
				return nil, err
			}
			repo.TLSClientCertData = string(tlsCertData)
			repo.TLSClientCertKey = string(tlsCertKey)
		} else {
			err := fmt.Errorf("--tls-client-cert-path is only supported for HTTPS repositories")
			return nil, err
		}
	}

	// Specifying github-app-private-key-path is only valid for HTTPS repositories
	if repoOpts.GithubAppPrivateKeyPath != "" {
		if git.IsHTTPSURL(repo.URL) {
			githubAppPrivateKey, err := os.ReadFile(repoOpts.GithubAppPrivateKeyPath)
			if err != nil {
				return nil, err
			}
			repo.GithubAppPrivateKey = string(githubAppPrivateKey)
		} else {
			err := fmt.Errorf("--github-app-private-key-path is only supported for HTTPS repositories")
			return nil, err
		}
	}

	conn, repoClient, err := a.client.NewRepoCredsClient()
	if err != nil {
		a.logger.Errorf("Application client intilialization failed: %v", err)
		return nil, err
	}

	defer io.Close(conn)
	// If the user set a username, but didn't supply password via --password,
	// then we prompt for it
	if repoOpts.Repo.Username != "" && repoOpts.Repo.Password == "" {
		repo.Password = cli.PromptPassword(repo.Password)
	}

	repoCreateReq := repocredspkg.RepoCredsCreateRequest{
		Creds:  &repo,
		Upsert: upsert,
	}

	createdRepo, err := repoClient.CreateRepositoryCredentials(ctx, &repoCreateReq)
	errors.CheckError(err)
	fmt.Printf("Repository credentials for '%s' added\n", createdRepo.URL)
	cr, err := json.Marshal(createdRepo)
	if err != nil {
		a.logger.Errorf("failed to marshal newly added repo creds details  %w", err)
		return nil, err
	}
	return cr, nil
}

func (a *ArgoCDCLient) RepoCredsDelete(req *model.ConfigPayload) (json.RawMessage, error) {
	ctx := context.Background()
	ru, ok := req.Data["repoURL"]
	if !ok {
		return nil, fmt.Errorf("failed to get repoURL from config payload")
	}
	repoURL, ok := ru.(string)
	if !ok {
		return nil, fmt.Errorf("repoURL type is not string")
	}
	conn, repoClient, err := a.client.NewRepoCredsClient()
	if err != nil {
		a.logger.Errorf("Application client intilialization failed: %v", err)
		return nil, err
	}
	defer io.Close(conn)

	deletedRepoCreds, err := repoClient.DeleteRepositoryCredentials(ctx, &repocredspkg.RepoCredsDeleteRequest{Url: repoURL})
	if err != nil {
		return nil, err
	}
	fmt.Printf("Repository '%s' removed\n", repoURL)
	dr, err := json.Marshal(deletedRepoCreds)
	if err != nil {
		a.logger.Errorf("failed to marshal deleted repo details  %w", err)
		return nil, err
	}
	return dr, nil
}

// Returns cluster query for getting cluster depending on the cluster selector
func getQueryBySelector(clusterSelector string) *clusterpkg.ClusterQuery {
	var query clusterpkg.ClusterQuery
	isServer, err := regexp.MatchString(`^https?://`, clusterSelector)
	if isServer || err != nil {
		query.Server = clusterSelector
	} else {
		query.Name = clusterSelector
	}
	return &query
}

func getRestConfig(pathOpts *clientcmd.PathOptions, ctxName string) (*rest.Config, error) {
	config, err := pathOpts.GetStartingConfig()
	if err != nil {
		return nil, err
	}

	clstContext := config.Contexts[ctxName]
	if clstContext == nil {
		return nil, fmt.Errorf("Context %s does not exist in kubeconfig", ctxName)
	}

	overrides := clientcmd.ConfigOverrides{
		Context: *clstContext,
	}

	clientConfig := clientcmd.NewDefaultClientConfig(*config, &overrides)
	conf, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	return conf, nil
}
