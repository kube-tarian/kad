package api

import (
	"context"

	"github.com/kube-tarian/kad/capten/common-pkg/pb/clusterpluginspb"
)

func (a *Agent) GetClusterPlugins(ctx context.Context, request *clusterpluginspb.GetClusterPluginsRequest) (
	*clusterpluginspb.GetClusterPluginsResponse, error) {
	pluginConfigList, err := a.as.GetAllClusterPluginConfigs()
	if err != nil {
		return &clusterpluginspb.GetClusterPluginsResponse{
			Status:        clusterpluginspb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to fetch plugins informations",
		}, nil
	}

	clusterPlugins := []*clusterpluginspb.ClusterPlugin{}
	for _, pluginConfig := range pluginConfigList {
		clusterPlugins = append(clusterPlugins, &clusterpluginspb.ClusterPlugin{
			StoreType:     pluginConfig.StoreType,
			PluginName:    pluginConfig.PluginName,
			Description:   pluginConfig.Description,
			Category:      pluginConfig.Category,
			Icon:          pluginConfig.Icon,
			Version:       pluginConfig.Version,
			InstallStatus: pluginConfig.InstallStatus,
		})
	}
	return &clusterpluginspb.GetClusterPluginsResponse{
		Status:  clusterpluginspb.StatusCode_OK,
		Plugins: clusterPlugins,
	}, nil
}

func (a *Agent) DeployClusterPlugin(ctx context.Context, request *clusterpluginspb.DeployClusterPluginRequest) (
	*clusterpluginspb.DeployClusterPluginResponse, error) {
	a.log.Infof("Recieved Plugin Deploy request for plugin %s, version %+v", request.Plugin.PluginName, request.Plugin.Version)

	err := a.plugin.DeployClusterPlugin(ctx, request.Plugin)
	if err != nil {
		a.log.Errorf("failed to deploy plugin [%s], %v", request.Plugin.PluginName, err)
		return &clusterpluginspb.DeployClusterPluginResponse{
			Status:        clusterpluginspb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to deploy plugin",
		}, err
	}

	a.log.Infof("Triggerred plugin [%s] install", request.Plugin.PluginName)
	return &clusterpluginspb.DeployClusterPluginResponse{
		Status:        clusterpluginspb.StatusCode_OK,
		StatusMessage: "Triggerred plugin install",
	}, nil
}

func (a *Agent) UnDeployClusterPlugin(ctx context.Context, request *clusterpluginspb.UnDeployClusterPluginRequest) (
	*clusterpluginspb.UnDeployClusterPluginResponse, error) {
	a.log.Infof("Recieved Plugin UnInstall request %+v", request)

	err := a.plugin.UnDeployClusterPlugin(ctx, request)
	if err != nil {
		a.log.Errorf("failed to undeploy plugin [%s], %v", request.PluginName, err)
		return &clusterpluginspb.UnDeployClusterPluginResponse{
			Status:        clusterpluginspb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to undeploy plugin",
		}, err
	}

	a.log.Infof("Triggerred plugin [%s] un install", request.PluginName)
	return &clusterpluginspb.UnDeployClusterPluginResponse{
		Status:        clusterpluginspb.StatusCode_OK,
		StatusMessage: "plugin is successfully undeployed",
	}, nil
}
