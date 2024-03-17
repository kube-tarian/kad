package activities

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/capten-sdk/db"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	pluginconfigstore "github.com/kube-tarian/kad/capten/common-pkg/pluginconfig-store"
	"github.com/kube-tarian/kad/capten/model"
)

type Configuration struct {
	AgentAddress string `envconfig:"AGENT_ADDRESSES" required:"true"`
}

type PluginActivities struct {
	config    *Configuration
	pas       *pluginconfigstore.Store
	k8sClient *k8s.K8SClient
}

func NewPluginActivities() (*PluginActivities, error) {
	if err != nil {
		logger.Errorf("failed to initialize plugin app store, %v", err)
		return nil, err
	}

	k8sclient, err := k8s.NewK8SClient(logger)
	if err != nil {
		logger.Errorf("failed to get k8s client, %v", err)
		return nil, err
	}

	pas, err := pluginconfigstore.NewStore(logger)
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
		config:    conf,
		pas:       pas,
		k8sClient: k8sclient,
	}, nil
}

func (p *PluginActivities) PluginDeployPreActionPostgresStoreActivity(ctx context.Context, req *model.ApplicationDeployRequest) (*model.ResponsePayload, error) {
	err := p.updateStatus(req.ReleaseName, "postgres-"+"initializing")
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"update status: %s\"}", err.Error())),
		}, err
	}

	// Call capten-sdk DB setup
	sdkDBClient := db.NewDBClientWithConfig(&db.DBConfig{
		DbOemName:         db.POSTGRES,
		PluginName:        req.PluginName,
		DbName:            req.Namespace + "-" + req.PluginName,
		DbServiceUserName: req.PluginName,
	})

	vaultPath, err := sdkDBClient.SetupDatabase()
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"setup database: %s\"}", err.Error())),
		}, err
	}

	err = p.k8sClient.CreateConfigmap(req.Namespace, req.PluginName+"-init-config", map[string]string{
		"vault-path": vaultPath,
	}, nil)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"update configmap: %s\"}", err.Error())),
		}, err
	}

	err = p.updateStatus(req.ReleaseName, "postgres-"+"initialized")
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"update status: %s\"}", err.Error())),
		}, err
	}
	return &model.ResponsePayload{
		Status: "SUCCESS",
	}, nil
}

func (p *PluginActivities) PluginDeployPreActionVaultStoreActivity(
	ctx context.Context,
	req *model.ApplicationDeployRequest,
) (*model.ResponsePayload, error) {
	err := p.updateStatus(req.ReleaseName, "vaultstore-"+"initializing")
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}
	// TODO: Call vault policy creation and path authorizations
	// Write the credentials in the vault
	logger.Infof("vault store activity Not implemented yet")

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
	logger.Infof("MTLS activity Not implemented yet")

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

// PluginDeployPostActionActivity... Updates the plugin deployment as "installed"
func (p *PluginActivities) PluginDeployPostActionActivity(ctx context.Context, req *model.ApplicationDeployRequest) (model.ResponsePayload, error) {
	err := p.updateStatus(req.ReleaseName, "deployed")
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

func (p *PluginActivities) PluginDeployUpdateStatusActivity(ctx context.Context, pluginName, status string) (model.ResponsePayload, error) {
	err := p.updateStatus(pluginName, status)
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
