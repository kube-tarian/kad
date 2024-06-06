package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/kube-tarian/kad/capten/common-pkg/pb/agentpb"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/pluginstorepb"
)

const (
	deployedStatus = "deployed"
)

func (a *Agent) SyncApp(ctx context.Context, request *agentpb.SyncAppRequest) (
	*agentpb.SyncAppResponse, error) {
	if request.Data == nil {
		return &agentpb.SyncAppResponse{
			Status:        agentpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "invalid data passed",
		}, nil
	}

	request.Data.Config.PluginStoreType = agentpb.PluginStoreType_DEFAULT_CAPTEN_STORE
	if err := a.as.UpsertAppConfig(request.Data); err != nil {
		a.log.Errorf("failed to update sync app config, %v", err)
		return &agentpb.SyncAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to sync app config",
		}, nil
	}

	a.log.Infof("Sync app [%s] successful", request.Data.Config.ReleaseName)
	return &agentpb.SyncAppResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "successful",
	}, nil
}

func (a *Agent) GetClusterApps(ctx context.Context, request *agentpb.GetClusterAppsRequest) (
	*agentpb.GetClusterAppsResponse, error) {
	res, err := a.as.GetAllApps()
	if err != nil {
		return &agentpb.GetClusterAppsResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to fetch cluster app configs",
		}, nil
	}

	appData := make([]*agentpb.AppData, 0)
	for _, r := range res {
		appData = append(appData, &agentpb.AppData{
			Config: r.GetConfig(),
		})
	}

	a.log.Infof("Found %d apps", len(appData))
	return &agentpb.GetClusterAppsResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "successful",
		AppData:       appData,
	}, nil

}

func (a *Agent) GetClusterAppLaunches(ctx context.Context, request *agentpb.GetClusterAppLaunchesRequest) (
	*agentpb.GetClusterAppLaunchesResponse, error) {
	res, err := a.as.GetAllApps()
	if err != nil {
		return &agentpb.GetClusterAppLaunchesResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to fetch cluster app configs",
		}, nil
	}

	cfgs := make([]*agentpb.AppLaunchConfig, 0)
	for _, r := range res {
		ssoSupported := false
		if len(r.Values.LaunchUIValues) != 0 {
			ssoSupported = true
		}

		appConfig := r.GetConfig()
		if len(appConfig.UiEndpoint) == 0 {
			continue
		}

		cfgs = append(cfgs, &agentpb.AppLaunchConfig{
			ReleaseName:  r.GetConfig().GetReleaseName(),
			Category:     r.GetConfig().GetCategory(),
			Description:  r.GetConfig().GetDescription(),
			Icon:         r.GetConfig().GetIcon(),
			UiEndpoint:   r.GetConfig().GetUiEndpoint(),
			SsoSupported: ssoSupported,
		})
	}

	a.log.Infof("Found %d apps with launch configs", len(cfgs))
	return &agentpb.GetClusterAppLaunchesResponse{
		LaunchConfigList: cfgs,
		Status:           agentpb.StatusCode_OK,
		StatusMessage:    "successful",
	}, nil
}

func (a *Agent) GetClusterAppConfig(ctx context.Context, request *agentpb.GetClusterAppConfigRequest) (
	*agentpb.GetClusterAppConfigResponse, error) {
	res, err := a.as.GetAppConfig(request.ReleaseName)
	if err != nil && err.Error() == "not found" {
		return &agentpb.GetClusterAppConfigResponse{
			Status:        agentpb.StatusCode_NOT_FOUND,
			StatusMessage: "app not found",
		}, nil
	}

	if err != nil {
		return &agentpb.GetClusterAppConfigResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to fetch app config",
		}, nil
	}

	a.log.Infof("Fetched app config for app [%s]", request.ReleaseName)
	return &agentpb.GetClusterAppConfigResponse{
		AppConfig:     res.GetConfig(),
		Status:        agentpb.StatusCode_OK,
		StatusMessage: agentpb.StatusCode_name[int32(agentpb.StatusCode_OK)],
	}, nil

}

func (a *Agent) GetClusterAppValues(ctx context.Context, request *agentpb.GetClusterAppValuesRequest) (
	*agentpb.GetClusterAppValuesResponse, error) {
	res, err := a.as.GetAppConfig(request.ReleaseName)
	if err != nil && err.Error() == "not found" {
		return &agentpb.GetClusterAppValuesResponse{
			Status:        agentpb.StatusCode_NOT_FOUND,
			StatusMessage: "app not found",
		}, nil
	}

	if err != nil {
		return &agentpb.GetClusterAppValuesResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to fetch app config",
		}, nil
	}

	a.log.Infof("Fetched app values for app [%s]", request.ReleaseName)
	return &agentpb.GetClusterAppValuesResponse{
		Values:        res.GetValues(),
		Status:        agentpb.StatusCode_OK,
		StatusMessage: agentpb.StatusCode_name[int32(agentpb.StatusCode_OK)],
	}, nil
}

func (a *Agent) DeployDefaultApps(ctx context.Context, request *agentpb.DeployDefaultAppsRequest) (
	*agentpb.DeployDefaultAppsResponse, error) {
	if err := a.plugin.SyncPlugins(pluginstorepb.StoreType_DEFAULT_STORE); err != nil {
		a.log.Errorf("failed to synch providers, %v", err)
	}

	plugins, err := a.plugin.GetPlugins(pluginstorepb.StoreType_DEFAULT_STORE)
	if err != nil {
		a.log.Errorf("failed to get plugins, %v", err)
	}

	failedPlugins := []string{}
	for _, plugin := range plugins {
		if err := a.plugin.DeployPlugin(pluginstorepb.StoreType_DEFAULT_STORE, plugin.PluginName, plugin.Versions[0], []byte{}); err != nil {
			a.log.Errorf("failed to deploy plugin, %v", err)
			failedPlugins = append(failedPlugins, plugin.PluginName)
		}
	}

	if len(failedPlugins) == len(plugins) {
		return &agentpb.DeployDefaultAppsResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: fmt.Sprintf("failed to deploy all default apps: %s", strings.Join(failedPlugins, ",")),
		}, nil
	}

	statusMessage := agentpb.StatusCode_name[int32(agentpb.StatusCode_OK)]
	if len(failedPlugins) != 0 {
		statusMessage = fmt.Sprintf("failed to deploy default apps: %s", strings.Join(failedPlugins, ","))
	}

	return &agentpb.DeployDefaultAppsResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: statusMessage,
	}, nil
}

func (a *Agent) GetDefaultAppsStatus(ctx context.Context, request *agentpb.GetDefaultAppsStatusRequest) (
	*agentpb.GetDefaultAppsStatusResponse, error) {
	plugins, err := a.plugin.GetPlugins(pluginstorepb.StoreType_DEFAULT_STORE)
	if err != nil {
		a.log.Errorf("failed to get plugins, %v", err)
	}

	overallStatus := agentpb.DeploymentStatus_SUCCESS
	anyPluginFailed := false
	resp := []*agentpb.ApplicationStatus{}
	for _, plugin := range plugins {
		pluginData, err := a.plugin.GetClusterPluginData(plugin.PluginName)
		if err != nil {
			a.log.Errorf("failed to fetch plugin status, %v", err)
			resp = append(resp, &agentpb.ApplicationStatus{
				AppName:       plugin.PluginName,
				Version:       plugin.Versions[0],
				Category:      plugin.Category,
				InstallStatus: "failed to fetch status",
				RuntimeStatus: "Unknown",
			})
			continue
		}

		if !strings.Contains(pluginData.InstallStatus, "failed") {
			anyPluginFailed = true
		} else if pluginData.InstallStatus != deployedStatus {
			overallStatus = agentpb.DeploymentStatus_ONGOING
		}

		resp = append(resp, &agentpb.ApplicationStatus{
			AppName:       pluginData.PluginName,
			Version:       pluginData.Version,
			Category:      pluginData.Category,
			InstallStatus: pluginData.InstallStatus,
			RuntimeStatus: pluginData.InstallStatus,
		})
	}

	// if any plugin failed and other plugins status is success, set overall status will be failed
	if anyPluginFailed && overallStatus == agentpb.DeploymentStatus_SUCCESS {
		overallStatus = agentpb.DeploymentStatus_FAILED
	}

	return &agentpb.GetDefaultAppsStatusResponse{
		Status:            agentpb.StatusCode_OK,
		StatusMessage:     agentpb.StatusCode_name[int32(agentpb.StatusCode_OK)],
		DeploymentStatus:  overallStatus,
		DefaultAppsStatus: resp,
	}, nil
}
