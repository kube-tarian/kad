package argocd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/argoproj/argo-cd/v2/util/io"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/kube-tarian/kad/integrator/model"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ArgoCDCLient struct {
	conf   *Configuration
	logger logging.Logger
	client apiclient.Client
}

func NewClient(logger logging.Logger) (*ArgoCDCLient, error) {
	cfg, err := fetchConfiguration()
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

func (a *ArgoCDCLient) DeployActivities(req interface{}) (json.RawMessage, error) {
	payload, _ := req.(model.RequestPayload)
	switch payload.Action {
	case "install":
		return a.Create(payload)
	case "delete":
		return a.Delete(payload)
	case "list":
		return a.List(payload)
	default:
		return nil, fmt.Errorf("unsupported action for argocd plugin: %v", payload.Action)
	}
}

func fetchConfiguration() (*Configuration, error) {
	cfg := &Configuration{}
	err := envconfig.Process("", cfg)
	return cfg, err
}

func (a *ArgoCDCLient) Create(payload model.RequestPayload) (json.RawMessage, error) {
	req := &model.Request{}
	err := json.Unmarshal(payload.Data, req)
	if err != nil {
		a.logger.Errorf("payload unmarshal failed, %v", err)
		return nil, err
	}
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

func (a *ArgoCDCLient) Delete(payload model.RequestPayload) (json.RawMessage, error) {
	req := &model.Request{}
	err := json.Unmarshal(payload.Data, req)
	if err != nil {
		a.logger.Errorf("payload unmarshal failed, %v", err)
		return nil, err
	}

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

func (a *ArgoCDCLient) List(req model.RequestPayload) (json.RawMessage, error) {
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
