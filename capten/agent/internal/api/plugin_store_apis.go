package api

import (
	"context"

	"github.com/kube-tarian/kad/capten/common-pkg/pb/pluginstorepb"
)

func (a *Agent) ConfigurePluginStore(ctx context.Context, request *pluginstorepb.ConfigurePluginStoreRequest) (
	*pluginstorepb.ConfigurePluginStoreResponse, error) {
	if err := validateArgs(request.Config.GitProjectId, request.Config.GitProjectURL); err != nil {
		a.log.Infof("request validation failed", err)
		return &pluginstorepb.ConfigurePluginStoreResponse{
			Status:        pluginstorepb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}
	a.log.Infof("Configure plugin store request recieved for store type %d", request.Config.StoreType)
	err := a.plugin.ConfigureStore(request.Config)
	if err != nil {
		a.log.Errorf("Configure plugin store request failed, %w", err)
		return &pluginstorepb.ConfigurePluginStoreResponse{
			Status:        pluginstorepb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to configuere plugin store",
		}, err
	}

	a.log.Infof("Plugin store request processed")
	return &pluginstorepb.ConfigurePluginStoreResponse{
		Status: pluginstorepb.StatusCode_OK,
	}, nil
}

func (a *Agent) GetPluginStoreConfig(ctx context.Context, request *pluginstorepb.GetPluginStoreConfigRequest) (
	*pluginstorepb.GetPluginStoreConfigResponse, error) {
	a.log.Infof("Get plugin store config request recieved for store type %d", request.StoreType)

	config, err := a.plugin.GetStoreConfig(request.StoreType)
	if err != nil {
		a.log.Errorf("Get plugin store config request failed, %w", err)
		return &pluginstorepb.GetPluginStoreConfigResponse{
			Status:        pluginstorepb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get plugin store config",
		}, err
	}

	a.log.Infof("Get plugin store config request processed")
	return &pluginstorepb.GetPluginStoreConfigResponse{
		Status: pluginstorepb.StatusCode_OK,
		Config: config,
	}, nil
}

func (a *Agent) SyncPluginStore(ctx context.Context, request *pluginstorepb.SyncPluginStoreRequest) (
	*pluginstorepb.SyncPluginStoreResponse, error) {
	a.log.Infof("Sync plugin store request recieved for store type %d", request.StoreType)

	err := a.plugin.SyncPlugins(request.StoreType)
	if err != nil {
		a.log.Errorf("Sync plugin store request failed, %w", err)
		return &pluginstorepb.SyncPluginStoreResponse{
			Status:        pluginstorepb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: err.Error(),
		}, err
	}

	a.log.Infof("Sync plugin store request processed")
	return &pluginstorepb.SyncPluginStoreResponse{
		Status: pluginstorepb.StatusCode_OK,
	}, nil
}

func (a *Agent) GetPlugins(ctx context.Context, request *pluginstorepb.GetPluginsRequest) (
	*pluginstorepb.GetPluginsResponse, error) {
	a.log.Infof("Get plugins request recieved")

	plugins, err := a.plugin.GetPlugins(request.StoreType)
	if err != nil {
		a.log.Errorf("Get plugins request failed, %w", err)
		return &pluginstorepb.GetPluginsResponse{
			Status:        pluginstorepb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get plugins",
		}, err
	}

	a.log.Infof("Get plugins request processed")
	return &pluginstorepb.GetPluginsResponse{
		Status:  pluginstorepb.StatusCode_OK,
		Plugins: plugins,
	}, nil
}

func (a *Agent) GetPluginValues(ctx context.Context, request *pluginstorepb.GetPluginValuesRequest) (
	*pluginstorepb.GetPluginValuesResponse, error) {
	err := validateArgs(request.PluginName, request.Version)
	if err != nil {
		a.log.Infof("request validation failed", err)
		return &pluginstorepb.GetPluginValuesResponse{
			Status:        pluginstorepb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}
	a.log.Infof("Get plugin values request recieved for plugin %s-%s", request.PluginName, request.Version)

	values, err := a.plugin.GetPluginValues(request.StoreType, request.PluginName, request.Version)
	if err != nil {
		a.log.Errorf("Get plugin values request failed for plugin %s-%s, %w",
			request.PluginName, request.Version, err)
		return &pluginstorepb.GetPluginValuesResponse{
			Status:        pluginstorepb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get plugins",
		}, err
	}

	a.log.Infof("Get plugin values request processed for plugin %s-%s",
		request.PluginName, request.Version)
	return &pluginstorepb.GetPluginValuesResponse{
		Status: pluginstorepb.StatusCode_OK,
		Values: values,
	}, nil
}

func (a *Agent) GetPluginData(ctx context.Context, request *pluginstorepb.GetPluginDataRequest) (
	*pluginstorepb.GetPluginDataResponse, error) {
	err := validateArgs(request.PluginName)
	if err != nil {
		a.log.Infof("request validation failed", err)
		return &pluginstorepb.GetPluginDataResponse{
			Status:        pluginstorepb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}
	a.log.Infof("Get plugin data request recieved for plugin %s", request.PluginName)

	pluginData, err := a.plugin.GetPluginData(request.StoreType, request.PluginName)
	if err != nil {
		a.log.Errorf("Get plugin data request failed for plugin %s, %w", request.PluginName, err)
		return &pluginstorepb.GetPluginDataResponse{
			Status:        pluginstorepb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get plugins",
		}, err
	}

	a.log.Infof("Get plugin data request processed for plugin %s", request.PluginName)
	return &pluginstorepb.GetPluginDataResponse{
		Status:     pluginstorepb.StatusCode_OK,
		PluginData: pluginData,
	}, nil
}

func (a *Agent) DeployPlugin(ctx context.Context, request *pluginstorepb.DeployPluginRequest) (
	*pluginstorepb.DeployPluginResponse, error) {
	err := validateArgs(request.PluginName)
	if err != nil {
		a.log.Infof("request validation failed", err)
		return &pluginstorepb.DeployPluginResponse{
			Status:        pluginstorepb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}
	a.log.Infof("Deploy plugin request recieved for plugin %s-%s", request.PluginName, request.Version)

	err = a.plugin.DeployPlugin(request.StoreType, request.PluginName, request.Version, request.Values)
	if err != nil {
		a.log.Errorf("Deploy plugin request failed for plugin %s-%s, %v", request.PluginName, request.Version, err)
		return &pluginstorepb.DeployPluginResponse{
			Status:        pluginstorepb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to deploy plugin",
		}, err
	}

	a.log.Infof("Deploy plugin request processed for plugin %s-%s", request.PluginName, request.Version)
	return &pluginstorepb.DeployPluginResponse{
		Status: pluginstorepb.StatusCode_OK,
	}, nil
}

func (a *Agent) UnDeployPlugin(ctx context.Context, request *pluginstorepb.UnDeployPluginRequest) (
	*pluginstorepb.UnDeployPluginResponse, error) {
	err := validateArgs(request.PluginName)
	if err != nil {
		a.log.Infof("request validation failed", err)
		return &pluginstorepb.UnDeployPluginResponse{
			Status:        pluginstorepb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}
	a.log.Infof("UnDeploy plugin request recieved for plugin %s", request.PluginName)

	err = a.plugin.UnDeployPlugin(request.StoreType, request.PluginName)
	if err != nil {
		a.log.Errorf("UnDeploy plugin request failed for plugin %s, %v", request.PluginName, err)
		return &pluginstorepb.UnDeployPluginResponse{
			Status:        pluginstorepb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to undeploy plugin",
		}, err
	}

	a.log.Infof("UnDeploy plugin request processed for plugin %s", request.PluginName)
	return &pluginstorepb.UnDeployPluginResponse{
		Status: pluginstorepb.StatusCode_OK,
	}, nil
}
