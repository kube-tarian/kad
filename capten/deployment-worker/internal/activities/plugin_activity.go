package activities

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/capten-sdk/db"
	"github.com/kube-tarian/kad/capten/common-pkg/cluster-plugins/clusterpluginspb"
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
	conf := &Configuration{}
	if err := envconfig.Process("", conf); err != nil {
		return nil, fmt.Errorf("cassandra config read faile, %v", err)
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

func (p *PluginActivities) PluginDeployPreActionPostgresStoreActivity(ctx context.Context, req *clusterpluginspb.Plugin) (*model.ResponsePayload, error) {
	err := p.updateStatus(req.PluginName, "postgres-"+"initializing")
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
		DbName:            req.DefaultNamespace + "-" + req.PluginName,
		DbServiceUserName: req.PluginName,
	})

	vaultPath, err := sdkDBClient.SetupDatabase()
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"setup database: %s\"}", err.Error())),
		}, err
	}

	err = p.createUpdateConfigmap(req.DefaultNamespace, req.PluginName+"-init-config", map[string]string{
		"vault-path": vaultPath,
	})
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"update configmap: %s\"}", err.Error())),
		}, err
	}

	err = p.updateStatus(req.PluginName, "postgres-"+"initialized")
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

func (p *PluginActivities) PluginUndeployPreActionPostgresStoreActivity(ctx context.Context, req *pluginconfigstore.PluginConfig) (*model.ResponsePayload, error) {
	err := p.updateStatus(req.PluginName, "postgres-"+"uninitializing")
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"update status: %s\"}", err.Error())),
		}, err
	}

	// Call capten-sdk DB setup
	db.NewDBClientWithConfig(&db.DBConfig{
		DbOemName:         db.POSTGRES,
		PluginName:        req.PluginName,
		DbName:            req.DefaultNamespace + "-" + req.PluginName,
		DbServiceUserName: req.PluginName,
	})
	// TODO: Invoke  captensdk DBDestroy

	err = p.pas.DeletePluginConfigByPluginName(req.DefaultNamespace)
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
	req *clusterpluginspb.Plugin,
) (*model.ResponsePayload, error) {
	err := p.updateStatus(req.PluginName, "vaultstore-"+"initializing")
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"update status: %s\"}", err.Error())),
		}, err
	}
	// TODO: Call vault policy creation and path authorizations
	// Write the credentials in the vault
	logger.Infof("vault store activity Not implemented yet")

	err = p.createUpdateConfigmap(req.DefaultNamespace, req.PluginName+"-init-config", map[string]string{})
	if err != nil {
		logger.Errorf("createupdate configmap failed: %v", err)
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"update configmap: %s\"}", req.PluginName+"-init-config")),
		}, err
	}

	err = p.updateStatus(req.PluginName, "vaultstore-"+"initialized")
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

func (p *PluginActivities) PluginUndeployPreActionVaultStoreActivity(
	ctx context.Context,
	req *pluginconfigstore.PluginConfig,
) (*model.ResponsePayload, error) {
	err := p.updateStatus(req.PluginName, "vaultstore-"+"uninitializing")
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}
	// TODO: Call vault policy creation and path authorizations
	// Write the credentials in the vault
	logger.Infof("vault store activity Not implemented yet")

	err = p.updateStatus(req.PluginName, "vaultstore-"+"uninitialized")
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

func (p *PluginActivities) PluginDeployPreActionMTLSActivity(ctx context.Context, req *clusterpluginspb.Plugin) (*model.ResponsePayload, error) {
	err := p.updateStatus(req.PluginName, "mtls-"+"initializing")
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}
	// TODO: Call MTLS creation
	// Write the mtls in the vault/conigmap
	logger.Infof("MTLS activity Not implemented yet")

	err = p.createUpdateConfigmap(req.DefaultNamespace, req.PluginName+"-init-config", map[string]string{})
	if err != nil {
		logger.Errorf("createupdate configmap failed: %v", err)
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"update configmap failed, %s\"}", req.PluginName+"-init-config")),
		}, err
	}

	err = p.updateStatus(req.PluginName, "mtls-"+"initialized")
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

func (p *PluginActivities) PluginUndeployPreActionMTLSActivity(ctx context.Context, req *pluginconfigstore.PluginConfig) (*model.ResponsePayload, error) {
	err := p.updateStatus(req.PluginName, "mtls-"+"uninitializing")
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"%s\"}", err.Error())),
		}, err
	}
	// TODO: Call MTLS creation
	// Write the mtls in the vault/conigmap
	logger.Infof("MTLS activity Not implemented yet")

	err = p.updateStatus(req.PluginName, "mtls-"+"uninitialized")
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
func (p *PluginActivities) PluginDeployPostActionActivity(ctx context.Context, req *clusterpluginspb.Plugin) (model.ResponsePayload, error) {
	err := p.updateStatus(req.PluginName, "deployed")
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
func (p *PluginActivities) PluginUndeployPostActionActivity(ctx context.Context, req *pluginconfigstore.PluginConfig) (model.ResponsePayload, error) {
	err := p.k8sClient.DeleteConfigmap(req.DefaultNamespace, req.PluginName+"-init-config")
	if err != nil {
		return model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"delete configmap %s faled\"}", req.PluginName+"-init-config")),
		}, err
	}

	err = p.pas.DeletePluginConfigByPluginName(req.PluginName)
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
	err := p.updateStatus(req.ReleaseName, "delete-"+"uninitialized")
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

func (p *PluginActivities) updateStatus(releaseName, status string) error {
	plugin, err := p.pas.GetPluginConfig(releaseName)
	if err != nil {
		return fmt.Errorf("plugin application %s not found in database", releaseName)
	}
	plugin.InstallStatus = status
	p.pas.UpsertPluginConfig(plugin)
	return nil
}

func (p *PluginActivities) createUpdateConfigmap(namespace, cmName string, data map[string]string) error {
	err := p.k8sClient.CreateNamespace(namespace)
	if err != nil {
		logger.Errorf("Creation of namespace failed: %v", err)
		return fmt.Errorf("creation of namespace faield")
	}
	cm, err := p.k8sClient.GetConfigmap(namespace, cmName)
	if err != nil {
		logger.Infof("plugin configmap %s not found", cmName)
		err = p.k8sClient.CreateConfigmap(namespace, cmName, data, nil)
		if err != nil {
			return fmt.Errorf("failed to create configmap %v", cmName)
		}
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
