package argocd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/repository"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/argoproj/argo-cd/v2/util/io"
	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"github.com/kube-tarian/kad/capten/common-pkg/plugins/fetcher"
	"github.com/kube-tarian/kad/capten/model"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ArgoCDCLient struct {
	conf   *Configuration
	logger logging.Logger
	client apiclient.Client
}

func NewClient(logger logging.Logger) (*ArgoCDCLient, error) {
	// cfg, err := fetchConfiguration(logger)
	// if err != nil {
	// 	return nil, err
	// }

	cfg := &Configuration{}
	err := envconfig.Process("", cfg)
	if err == nil {
		return nil, err
	}

	k8sClient, err := k8s.NewK8SClient(logger)
	if err != nil {
		return nil, err
	}
	res, err := k8sClient.FetchSecretDetails(&k8s.SecretDetailsRequest{Namespace: "argo-cd", SecretName: "argocd-initial-admin-secret"})
	if err != nil {
		return nil, err
	}
	fmt.Printf("secret res => %+v \n", *res)
	password := res.Data["password"]
	if len(password) == 0 {
		return nil, fmt.Errorf("credentials not found in the secret")
	}

	cfg.Username = password
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
		// return a.ClusterAdd(payload)
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

func (a *ArgoCDCLient) ListRepositories(ctx context.Context) (*v1alpha1.RepositoryList, error) {
	conn, appClient, err := a.client.NewRepoClient()
	if err != nil {
		return nil, err
	}
	defer io.Close(conn)

	list, err := appClient.ListRepositories(ctx, &repository.RepoQuery{})
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (a *ArgoCDCLient) CreateRepository(ctx context.Context, repo *Repository) (*v1alpha1.Repository, error) {
	conn, appClient, err := a.client.NewRepoClient()
	if err != nil {
		return nil, err
	}
	defer io.Close(conn)

	resp, err := appClient.CreateRepository(ctx, &repository.RepoCreateRequest{
		Repo: &v1alpha1.Repository{
			Project:       repo.Project,
			Repo:          repo.Repo,
			SSHPrivateKey: repo.SSHPrivateKey,
			Type:          repo.Type,
			ConnectionState: v1alpha1.ConnectionState{
				Status:  repo.ConnectionState.Status,
				Message: repo.ConnectionState.Message,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *ArgoCDCLient) DeleteRepository(ctx context.Context, repo string) (*repository.RepoResponse, error) {
	conn, appClient, err := a.client.NewRepoClient()
	if err != nil {
		return nil, err
	}
	defer io.Close(conn)

	encodedRepo := url.QueryEscape(repo)

	resp, err := appClient.DeleteRepository(ctx, &repository.RepoQuery{Repo: encodedRepo})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
