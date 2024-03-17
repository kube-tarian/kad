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
	a.log.Infof("Recieved App Install request for appName %s, version %+v", request.Plugin.PluginName, request.Plugin.Version)

	pluginCofnig := &pluginconfigstore.PluginConfig{
		Plugin: request.Plugin,
	}

	apiEndpoint, err := executeStringTemplateValues(pluginCofnig.ApiEndpoint, pluginCofnig.Values)
	if err != nil {
		a.log.Errorf("failed to derive template launch URL for app %s, %v", pluginCofnig.PluginName, err)
		return &clusterpluginspb.DeployClusterPluginResponse{
			Status:        clusterpluginspb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to prepare app values",
		}, nil
	}

	pluginCofnig.InstallStatus = string(model.AppIntallingStatus)
	pluginCofnig.ApiEndpoint = apiEndpoint

	if err := a.pas.UpsertPluginConfig(pluginCofnig); err != nil {
		a.log.Errorf("failed to update app config data for app %s, %v", pluginCofnig.PluginName, err)
		return &clusterpluginspb.DeployClusterPluginResponse{
			Status:        clusterpluginspb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to update app config data",
		}, nil
	}

	// deployReq := prepareAppDeployRequestFromPlugin(plugin)
	go a.deployPluginWithWorkflow(request.Plugin, pluginCofnig)

	a.log.Infof("Triggerred app [%s] install", request.Plugin.PluginName)
	return &clusterpluginspb.DeployClusterPluginResponse{
		Status:        clusterpluginspb.StatusCode_OK,
		StatusMessage: "Triggerred app install",
	}, nil
}

func (a *Agent) UnDeployClusterPlugin(ctx context.Context, request *clusterpluginspb.UnDeployClusterPluginRequest) (
	*clusterpluginspb.UnDeployClusterPluginResponse, error) {
	a.log.Infof("Recieved App UnInstall request %+v", request)

	if request.PluginName == "" {
		a.log.Errorf("release name is empty")
		return &clusterpluginspb.UnDeployClusterPluginResponse{
			Status:        clusterpluginspb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "release name is missing in request",
		}, nil
	}

	appConfigdata, err := a.pas.GetPluginConfig(request.PluginName)
	if err != nil {
		a.log.Errorf("failed to fetch app config record %s, %v", request.PluginName, err)
		return &clusterpluginspb.UnDeployClusterPluginResponse{
			Status:        clusterpluginspb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to fetch app config",
		}, nil
	}

	req := &model.ApplicationDeleteRequest{
		PluginName:  "helm",
		Namespace:   appConfigdata.DefaultNamespace,
		ReleaseName: request.PluginName,
		ClusterName: "capten",
		Timeout:     10,
	}

	appConfigdata.InstallStatus = string(model.AppUnInstallingStatus)
	if err := a.pas.UpsertPluginConfig(appConfigdata); err != nil {
		a.log.Errorf("failed to update app config status with UnInstalling for app %s, %v", req.ReleaseName, err)
		return &clusterpluginspb.UnDeployClusterPluginResponse{
			Status:        clusterpluginspb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to undeploy the app",
		}, nil
	}

	go a.unInstallPluginWithWorkflow(req, appConfigdata)

	a.log.Infof("Triggerred app [%s] un install", request.PluginName)
	return &clusterpluginspb.UnDeployClusterPluginResponse{
		Status:        clusterpluginspb.StatusCode_OK,
		StatusMessage: "app is successfully undeployed",
	}, nil
}

func (a *Agent) deployPluginWithWorkflow(plugin *clusterpluginspb.Plugin, pluginConfig *pluginconfigstore.PluginConfig) {
	wd := workers.NewDeployment(a.tc, a.log)
	_, err := wd.SendEventV2(context.TODO(), wd.GetPluginWorkflowName(), string(model.AppInstallAction), plugin)
	if err != nil {
		pluginConfig.InstallStatus = string(model.AppIntallFailedStatus)
		if err := a.pas.UpsertPluginConfig(pluginConfig); err != nil {
			a.log.Errorf("failed to update app config for app %s, %v", pluginConfig.PluginName, err)
			return
		}
		a.log.Errorf("failed to send event to workflow for app %s, %v", pluginConfig.PluginName, err)
		return
	}

	pluginConfig.InstallStatus = string(model.AppIntalledStatus)
	if err := a.pas.UpsertPluginConfig(pluginConfig); err != nil {
		a.log.Errorf("failed to update app config for app %s, %v", pluginConfig.PluginName, err)
		return
	}
}

func (a *Agent) unInstallPluginWithWorkflow(req *model.ApplicationDeleteRequest, appConfig *pluginconfigstore.PluginConfig) {
	wd := workers.NewDeployment(a.tc, a.log)
	_, err := wd.SendDeleteEvent(context.TODO(), wd.GetPluginWorkflowName(), string(model.AppUnInstallAction), req, appConfig.Capabilities)
	if err != nil {
		a.log.Errorf("failed to send delete event to workflow for app %s, %v", req.ReleaseName, err)

		appConfig.InstallStatus = string(model.AppIntalledStatus)
		if err := a.pas.UpsertPluginConfig(appConfig); err != nil {
			a.log.Errorf("failed to update app config status with Installed for app %s, %v", req.ReleaseName, err)
		}
		return
	}

	if err := a.as.DeleteAppConfigByReleaseName(req.ReleaseName); err != nil {
		a.log.Errorf("failed to delete installed app config record %s, %v", req.ReleaseName, err)
		return
	}
}

func prepareAppDeployRequestFromPlugin(data *clusterpluginspb.Plugin) *model.ApplicationInstallRequest {
	return &model.ApplicationInstallRequest{
		PluginName:     "helm",
		RepoName:       data.ChartName,
		RepoURL:        data.ChartRepo,
		ChartName:      data.ChartName,
		Namespace:      data.DefaultNamespace,
		ReleaseName:    data.ChartName,
		Version:        data.Version,
		ClusterName:    "capten",
		OverrideValues: string(data.Values),
		Timeout:        10,
	}
}
