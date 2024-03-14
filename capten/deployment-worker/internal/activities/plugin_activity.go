package activities

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	pluginappstore "github.com/kube-tarian/kad/capten/common-pkg/pluginapp-store"
	dbstore "github.com/kube-tarian/kad/capten/deployment-worker/internal/db-store"
	"github.com/kube-tarian/kad/capten/model"
)

type PluginActivities struct {
	store     *dbstore.Store
	pas       *pluginappstore.Store
	k8sClient *k8s.K8SClient
}

func NewPluginActivities() (*PluginActivities, error) {
	store, err := dbstore.NewStore(logger)
	if err != nil {
		return nil, err
	}

	pas, err := pluginappstore.NewStore(logger)
	if err != nil {
		logger.Errorf("failed to initialize plugin app store, %v", err)
		return nil, err
	}

	k8sclient, err := k8s.NewK8SClient(logger)
	if err != nil {
		logger.Errorf("failed to get k8s client, %v", err)
		return nil, err
	}

	return &PluginActivities{
		store:     store,
		pas:       pas,
		k8sClient: k8sclient,
	}, nil
}

func (p *PluginActivities) PluginDeployPreActionPostgresStoreActivity(ctx context.Context, req *model.ApplicationDeployRequest) (*model.ResponsePayload, error) {
	err := p.updateStatus(req.ReleaseName, "postgres-"+"initializing")
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}

	// TODO: Call capten-sdk DB setup
	// Write the credentials in the vault

	err = p.updateStatus(req.ReleaseName, "postgres-"+"initialized")
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}
	return &model.ResponsePayload{
		Status: "SUCCESS",
	}, nil
}

func (p *PluginActivities) PluginDeployPreActionVaultStoreActivity(ctx context.Context, req *model.ApplicationDeployRequest) (*model.ResponsePayload, error) {
	err := p.updateStatus(req.ReleaseName, "vaultstore-"+"initializing")
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}
	// TODO: Call vault policy creation and path authorizations
	// Write the credentials in the vault

	err = p.updateStatus(req.ReleaseName, "vaultstore-"+"initialized")
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}
	return &model.ResponsePayload{
		Status: "SUCCESS",
	}, nil
}

func (p *PluginActivities) PluginDeployPreActionMTLSActivity(ctx context.Context, req *model.ApplicationDeployRequest) (*model.ResponsePayload, error) {
	err := p.updateStatus(req.ReleaseName, "mtls-"+"initializing")
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}
	// TODO: Call MTLS creation
	// Write the mtls in the vault/conigmap

	err = p.updateStatus(req.ReleaseName, "mtls-"+"initialized")
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}
	return &model.ResponsePayload{
		Status: "SUCCESS",
	}, nil
}

func (p *PluginActivities) PluginDeployActivity(ctx context.Context, req *model.ApplicationDeployRequest) (model.ResponsePayload, error) {
	err := p.updateStatus(req.ReleaseName, "plugin-install-"+"initializing")
	if err != nil {
		return model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}

	// Call application install
	resp, err := installApplication(req)
	if err != nil {
		status := "plugin-install-" + "failed"
		_ = p.updateStatus(req.ReleaseName, status)
		return resp, err
	}

	err = p.updateStatus(req.ReleaseName, "plugin-install-"+"initialized")
	if err != nil {
		return model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}
	return model.ResponsePayload{
		Status: "SUCCESS",
	}, nil
}

// PluginDeployPostActionActivity... Updates the plugin deployment as "installed"
func (p *PluginActivities) PluginDeployPostActionActivity(ctx context.Context, req *model.ApplicationDeployRequest) (model.ResponsePayload, error) {
	err := p.updateStatus(req.ReleaseName, "installed")
	if err != nil {
		return model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}
	return model.ResponsePayload{
		Status: "SUCCESS",
	}, nil
}

func (p *PluginActivities) PluginUndeployActivity(ctx context.Context, req *model.DeployerDeleteRequest) (model.ResponsePayload, error) {
	err := p.updateStatus(req.ReleaseName, "delete-"+"initialized")
	if err != nil {
		return model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}

	resp, err := uninstallApplication(req)
	if err != nil {
		status := "delete-" + "failed"
		_ = p.updateStatus(req.ReleaseName, status)
		return resp, err
	}

	err = p.updateStatus(req.ReleaseName, "delete-"+"success")
	if err != nil {
		return model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}

	return model.ResponsePayload{
		Status: "SUCCESS",
	}, nil
}

func (p *PluginActivities) PluginUndeployPostActionsActivity(ctx context.Context, req *model.DeployerDeleteRequest) (model.ResponsePayload, error) {
	err := p.pas.DeletePluginConfigByReleaseName(req.ReleaseName)
	if err != nil {
		return model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}

	return model.ResponsePayload{
		Status: "SUCCESS",
	}, nil
}

func (p *PluginActivities) updateStatus(releaseName, status string) error {
	plugin, err := p.pas.GetPluginConfig(releaseName)
	if err != nil {
		return fmt.Errorf("plugin application %s not found in database", releaseName)
	}
	plugin.InstallStatus = status
	p.pas.UpsertPluginConfig(plugin)
	return nil
}

func (p *PluginActivities) updateConfigmap(namespace, cmName string, data map[string]string) error {
	cm, err := p.k8sClient.GetConfigmap(namespace, cmName)
	if err != nil {
		return fmt.Errorf("plugin configmap %s not found", cmName)
	}
	for k, v := range data {
		cm[k] = v
	}
	err = p.k8sClient.UpdateConfigmap(namespace, cmName, cm)
	if err != nil {
		return fmt.Errorf("plugin configmap %s not found", cmName)
	}
	return nil
}
