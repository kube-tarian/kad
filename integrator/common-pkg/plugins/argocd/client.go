package argocd

import (
	"context"
	"encoding/json"
	"fmt"
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
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	argoappv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/argoproj/argo-cd/v2/util/clusterauth"
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
		return a.HandleRepo(payload)
	default:
		return nil, fmt.Errorf("unsupported action for argocd plugin: %v", payload.Action)
	}
}

func (a *ArgoCDCLient) HandleCluster(req interface{}) (json.RawMessage, error) {
	payload, _ := req.(model.ConfigPayload)
	switch payload.Action {
	case "add":
		return a.ClusterAdd(payload)
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
		// return a.RepoAdd(payload)
	case "delete":
		// return a.RepoDelete(payload)
	case "list":
		// return a.RepoList(payload)
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

func (a *ArgoCDCLient) ClusterAdd(req model.ConfigPayload) (json.RawMessage, error) {
	conn, clusterClient, err := a.client.NewClusterClient()
	if err != nil {
		a.logger.Errorf("Application client intilialization failed: %v", err)
		return nil, err
	}

	defer io.Close(conn)
	var clusterOpts cmdutil.ClusterOptions
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
	ctx := context.Background()
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

func (a *ArgoCDCLient) ClusterDelete(req model.ConfigPayload) (json.RawMessage, error) {
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
