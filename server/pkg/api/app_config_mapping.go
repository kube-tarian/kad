package api

import (
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func mapAgentAppsToServerResp(appDataList []*agentpb.AppData) []*serverpb.ClusterAppConfig {
	clusterAppConfigs := make([]*serverpb.ClusterAppConfig, len(appDataList))
	for index, appConfig := range appDataList {
		var clusterAppConfig serverpb.ClusterAppConfig
		clusterAppConfig.ReleaseName = appConfig.Config.ReleaseName
		clusterAppConfig.AppName = appConfig.Config.AppName
		clusterAppConfig.Version = appConfig.Config.Version
		clusterAppConfig.Category = appConfig.Config.Category
		clusterAppConfig.Description = appConfig.Config.Description
		clusterAppConfig.ChartName = appConfig.Config.ChartName
		clusterAppConfig.RepoName = appConfig.Config.RepoName
		clusterAppConfig.RepoURL = appConfig.Config.RepoURL
		clusterAppConfig.Namespace = appConfig.Config.Namespace
		clusterAppConfig.CreateNamespace = appConfig.Config.CreateNamespace
		clusterAppConfig.PrivilegedNamespace = appConfig.Config.PrivilegedNamespace
		clusterAppConfig.Icon = appConfig.Config.Icon
		clusterAppConfig.UiEndpoint = appConfig.Config.UiEndpoint
		clusterAppConfig.UiModuleEndpoint = appConfig.Config.UiModuleEndpoint
		clusterAppConfig.InstallStatus = appConfig.Config.InstallStatus
		clusterAppConfig.RuntimeStatus = ""
		clusterAppConfig.DefualtApp = appConfig.Config.DefualtApp
		clusterAppConfig.PluginName = appConfig.Config.PluginName
		clusterAppConfig.PluginStoreType = serverpb.PluginStoreType(appConfig.Config.PluginStoreType)
		clusterAppConfig.ApiEndpoint = appConfig.Config.ApiEndpoint
		clusterAppConfigs[index] = &clusterAppConfig
	}
	return clusterAppConfigs
}

func mapAgentAppLaunchConfigsToServer(appLaunchCfgs []*agentpb.AppLaunchConfig) []*serverpb.AppLaunchConfig {
	svrAppLaunchCfg := make([]*serverpb.AppLaunchConfig, len(appLaunchCfgs))
	for index, cfg := range appLaunchCfgs {
		var launchCfg serverpb.AppLaunchConfig
		launchCfg.ReleaseName = cfg.ReleaseName
		launchCfg.Category = cfg.Category
		launchCfg.Description = cfg.Description
		launchCfg.Icon = cfg.Icon
		launchCfg.UiEndpoint = cfg.UiEndpoint
		svrAppLaunchCfg[index] = &launchCfg
	}
	return svrAppLaunchCfg
}
