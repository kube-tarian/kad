package argocd

import (
	"context"
	"encoding/json"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/argoproj/argo-cd/v2/util/io"
	"github.com/kube-tarian/kad/capten/model"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (a *ArgoCDClient) Create(req *model.CreteRequestPayload) (json.RawMessage, error) {
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

func (a *ArgoCDClient) Delete(req *model.DeleteRequestPayload) (json.RawMessage, error) {
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

func (a *ArgoCDClient) List(req *model.ListRequestPayload) (json.RawMessage, error) {
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

func (a *ArgoCDClient) TriggerAppSync(ctx context.Context, namespace, name string) (*v1alpha1.Application, error) {
	conn, app, err := a.client.NewApplicationClient()
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	pruneApp := true
	resp, err := app.Sync(ctx, &application.ApplicationSyncRequest{
		Name:         &name,
		AppNamespace: &namespace,
		Prune:        &pruneApp,
		RetryStrategy: &v1alpha1.RetryStrategy{
			Limit: 3,
		}})
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (a *ArgoCDClient) GetAppSyncStatus(ctx context.Context, namespace, name string) (*v1alpha1.Application, error) {
	conn, app, err := a.client.NewApplicationClient()
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	resp, err := app.Get(ctx, &application.ApplicationQuery{Name: &name, AppNamespace: &namespace})
	if err != nil {
		return nil, err
	}

	return resp, err
}
