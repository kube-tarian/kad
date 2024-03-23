package api

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/internal/workers"
	"github.com/kube-tarian/kad/capten/common-pkg/cluster-plugins/clusterpluginspb"
	pluginconfigstore "github.com/kube-tarian/kad/capten/common-pkg/pluginconfig-store"
	"github.com/kube-tarian/kad/capten/model"
)

func (a *Agent) GetClusterPlugins(ctx context.Context, request *clusterpluginspb.GetClusterPluginsRequest) (
	*clusterpluginspb.GetClusterPluginsResponse, error) {
	pluginConfigList, err := a.pas.GetAllPlugins()
	if err != nil {
		return &clusterpluginspb.GetClusterPluginsResponse{
			Status:        clusterpluginspb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to fetch plugins informations",
		}, nil
	}

	clusterPlugins := []*clusterpluginspb.ClusterPlugin{}
	for idx, pluginConfig := range pluginConfigList {
		clusterPlugins[idx] = &clusterpluginspb.ClusterPlugin{
			StoreType:     pluginConfig.StoreType,
			PluginName:    pluginConfig.PluginName,
			Description:   pluginConfig.Description,
			Category:      pluginConfig.Category,
			Icon:          pluginConfig.Icon,
			Version:       pluginConfig.Version,
			InstallStatus: pluginConfig.InstallStatus,
		}
	}
	return &clusterpluginspb.GetClusterPluginsResponse{
		Status:  clusterpluginspb.StatusCode_OK,
		Plugins: clusterPlugins,
	}, nil
}

func (a *Agent) DeployClusterPlugin(ctx context.Context, request *clusterpluginspb.DeployClusterPluginRequest) (
	*clusterpluginspb.DeployClusterPluginResponse, error) {
	a.log.Infof("Recieved Plugin Deploy request for plugin %s, version %+v", request.Plugin.PluginName, request.Plugin.Version)

	pluginConfig := &pluginconfigstore.PluginConfig{
		Plugin: request.Plugin,
	}

	apiEndpoint, err := executeStringTemplateValues(pluginConfig.ApiEndpoint, pluginConfig.Values)
	if err != nil {
		a.log.Errorf("failed to derive template launch URL for plugin %s, %v", pluginConfig.PluginName, err)
		return &clusterpluginspb.DeployClusterPluginResponse{
			Status:        clusterpluginspb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to prepare plugin values",
		}, nil
	}

	pluginConfig.InstallStatus = string(model.AppIntallingStatus)
	pluginConfig.ApiEndpoint = apiEndpoint

	if err := a.pas.UpsertPluginConfig(pluginConfig); err != nil {
		a.log.Errorf("failed to update plugin config data for plugin %s, %v", pluginConfig.PluginName, err)
		return &clusterpluginspb.DeployClusterPluginResponse{
			Status:        clusterpluginspb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to update plugin config data",
		}, nil
	}

	// deployReq := prepareAppDeployRequestFromPlugin(plugin)
	go a.deployPluginWithWorkflow(request.Plugin, pluginConfig)

	a.log.Infof("Triggerred plugin [%s] install", request.Plugin.PluginName)
	return &clusterpluginspb.DeployClusterPluginResponse{
		Status:        clusterpluginspb.StatusCode_OK,
		StatusMessage: "Triggerred plugin install",
	}, nil
}

func (a *Agent) UnDeployClusterPlugin(ctx context.Context, request *clusterpluginspb.UnDeployClusterPluginRequest) (
	*clusterpluginspb.UnDeployClusterPluginResponse, error) {
	a.log.Infof("Recieved Plugin UnInstall request %+v", request)

	if request.PluginName == "" {
		a.log.Errorf("release name is empty")
		return &clusterpluginspb.UnDeployClusterPluginResponse{
			Status:        clusterpluginspb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "release name is missing in request",
		}, nil
	}

	pluginConfigdata, err := a.pas.GetPluginConfig(request.PluginName)
	if err != nil {
		a.log.Errorf("failed to fetch plugin config record %s, %v", request.PluginName, err)
		return &clusterpluginspb.UnDeployClusterPluginResponse{
			Status:        clusterpluginspb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to fetch plugin config",
		}, nil
	}

	pluginConfigdata.InstallStatus = string(model.AppUnInstallingStatus)
	if err := a.pas.UpsertPluginConfig(pluginConfigdata); err != nil {
		a.log.Errorf("failed to update plugin config status with UnInstalling for plugin %s, %v", pluginConfigdata.PluginName, err)
		return &clusterpluginspb.UnDeployClusterPluginResponse{
			Status:        clusterpluginspb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to undeploy the plugin",
		}, nil
	}

	go a.unInstallPluginWithWorkflow(request, pluginConfigdata)

	a.log.Infof("Triggerred plugin [%s] un install", request.PluginName)
	return &clusterpluginspb.UnDeployClusterPluginResponse{
		Status:        clusterpluginspb.StatusCode_OK,
		StatusMessage: "plugin is successfully undeployed",
	}, nil
}

func (a *Agent) deployPluginWithWorkflow(plugin *clusterpluginspb.Plugin, pluginConfig *pluginconfigstore.PluginConfig) {
	wd := workers.NewDeployment(a.tc, a.log)
	_, err := wd.SendEventV2(context.TODO(), wd.GetPluginWorkflowName(), string(model.AppInstallAction), plugin)
	if err != nil {
		// pluginConfig.InstallStatus = string(model.AppIntallFailedStatus)
		// if err := a.pas.UpsertPluginConfig(pluginConfig); err != nil {
		// 	a.log.Errorf("failed to update plugin config for plugin %s, %v", pluginConfig.PluginName, err)
		// 	return
		// }
		a.log.Errorf("sendEventV2 failed, plugin: %s, reason: %v", pluginConfig.PluginName, err)
		return
	}
	// TODO: workflow will update the final status
	// Write a periodic scheduler which will go through all apps not in installed status and check the status till either success or failed.
	// Make SendEventV2 asynchrounous so that periodic scheduler will take care of monitoring.
}

func (a *Agent) unInstallPluginWithWorkflow(request *clusterpluginspb.UnDeployClusterPluginRequest, pluginConfig *pluginconfigstore.PluginConfig) {
	wd := workers.NewDeployment(a.tc, a.log)
	_, err := wd.SendDeleteEvent(context.TODO(), wd.GetPluginWorkflowName(), string(model.AppUnInstallAction), request)
	if err != nil {
		a.log.Errorf("failed to send delete event to workflow for plugin %s, %v", pluginConfig.PluginName, err)

		pluginConfig.InstallStatus = string(model.AppIntalledStatus)
		if err := a.pas.UpsertPluginConfig(pluginConfig); err != nil {
			a.log.Errorf("failed to update plugin config status with Installed for plugin %s, %v", pluginConfig.PluginName, err)
		}
		return
	}

	if err := a.as.DeleteAppConfigByReleaseName(pluginConfig.PluginName); err != nil {
		a.log.Errorf("failed to delete installed plugin config record %s, %v", pluginConfig.PluginName, err)
		return
	}
}
