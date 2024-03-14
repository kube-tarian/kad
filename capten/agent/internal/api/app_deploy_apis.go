package api

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/internal/pb/agentpb"
	"github.com/kube-tarian/kad/capten/agent/internal/workers"
	"github.com/kube-tarian/kad/capten/model"
	"github.com/pkg/errors"
)

func (a *Agent) InstallApp(ctx context.Context, request *agentpb.InstallAppRequest) (*agentpb.InstallAppResponse, error) {
	a.log.Infof("Recieved App Install request for appName %s, version %+v", request.AppConfig.AppName, request.AppConfig.Version)
	templateValues, err := deriveTemplateValues(request.AppValues.OverrideValues, request.AppValues.TemplateValues)
	if err != nil {
		a.log.Errorf("failed to derive template values for app %s, %v", request.AppConfig.ReleaseName, err)
		return &agentpb.InstallAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to prepare app values",
		}, nil
	}

	launchURL, err := executeStringTemplateValues(request.AppConfig.LaunchURL, request.AppValues.OverrideValues)
	if err != nil {
		a.log.Errorf("failed to derive template launch URL for app %s, %v", request.AppConfig.ReleaseName, err)
		return &agentpb.InstallAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to prepare app values",
		}, nil
	}

	apiEndpoint, err := executeStringTemplateValues(request.AppConfig.ApiEndpoint, request.AppValues.OverrideValues)
	if err != nil {
		a.log.Errorf("failed to derive template launch URL for app %s, %v", request.AppConfig.ReleaseName, err)
		return &agentpb.InstallAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to prepare app values",
		}, nil
	}

	syncConfig := &agentpb.SyncAppData{
		Config: &agentpb.AppConfig{
			ReleaseName:         request.AppConfig.ReleaseName,
			AppName:             request.AppConfig.AppName,
			Version:             request.AppConfig.Version,
			Category:            request.AppConfig.Category,
			Description:         request.AppConfig.Description,
			ChartName:           request.AppConfig.ChartName,
			RepoName:            request.AppConfig.RepoName,
			RepoURL:             request.AppConfig.RepoURL,
			Namespace:           request.AppConfig.Namespace,
			CreateNamespace:     request.AppConfig.CreateNamespace,
			PrivilegedNamespace: request.AppConfig.PrivilegedNamespace,
			Icon:                request.AppConfig.Icon,
			LaunchURL:           launchURL,
			LaunchUIDescription: request.AppConfig.LaunchUIDescription,
			InstallStatus:       string(model.AppIntallingStatus),
			DefualtApp:          request.AppConfig.DefualtApp,
			PluginName:          request.AppConfig.PluginName,
			PluginDescription:   request.AppConfig.PluginDescription,
			ApiEndpoint:         apiEndpoint,
		},
		Values: &agentpb.AppValues{
			OverrideValues: request.AppValues.OverrideValues,
			LaunchUIValues: request.AppValues.LaunchUIValues,
			TemplateValues: request.AppValues.TemplateValues,
		},
	}

	if err := a.as.UpsertAppConfig(syncConfig); err != nil {
		a.log.Errorf("failed to update app config data for app %s, %v", request.AppConfig.ReleaseName, err)
		return &agentpb.InstallAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to update app config data",
		}, nil
	}

	deployReq := prepareAppDeployRequestFromSyncApp(syncConfig, templateValues)
	go a.installAppWithWorkflow(deployReq, syncConfig)

	a.log.Infof("Triggerred app [%s] install", request.AppConfig.ReleaseName)
	return &agentpb.InstallAppResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "Triggerred app install",
	}, nil
}

func (a *Agent) UnInstallApp(ctx context.Context, request *agentpb.UnInstallAppRequest) (*agentpb.UnInstallAppResponse, error) {
	a.log.Infof("Recieved App UnInstall request %+v", request)

	if request.ReleaseName == "" {
		a.log.Errorf("release name is empty")
		return &agentpb.UnInstallAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "release name is missing in request",
		}, nil
	}

	appConfigdata, err := a.as.GetAppConfig(request.ReleaseName)
	if err != nil {
		a.log.Errorf("failed to fetch app config record %s, %v", request.ReleaseName, err)
		return &agentpb.UnInstallAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to fetch app config",
		}, nil
	}

	req := &model.ApplicationDeleteRequest{
		PluginName:  "helm",
		Namespace:   appConfigdata.Config.Namespace,
		ReleaseName: request.ReleaseName,
		ClusterName: "capten",
		Timeout:     10,
	}

	appConfigdata.Config.InstallStatus = string(model.AppUnInstallingStatus)
	if err := a.as.UpsertAppConfig(appConfigdata); err != nil {
		a.log.Errorf("failed to update app config status with UnInstalling for app %s, %v", req.ReleaseName, err)
		return &agentpb.UnInstallAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to undeploy the app",
		}, nil
	}

	go a.unInstallAppWithWorkflow(req, appConfigdata)

	a.log.Infof("Triggerred app [%s] un install", request.ReleaseName)
	return &agentpb.UnInstallAppResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "app is successfully undeployed",
	}, nil
}

func (a *Agent) UpdateAppValues(ctx context.Context, req *agentpb.UpdateAppValuesRequest) (*agentpb.UpdateAppValuesResponse, error) {
	a.log.Infof("Received UpdateAppValues request %+v", req)

	// Get the config templates for release name
	appConfig, err := a.as.GetAppConfig(req.ReleaseName)
	if err != nil {
		a.log.Errorf("failed to read app %s config data, %v", req.ReleaseName, err)
		return &agentpb.UpdateAppValuesResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "failed to read app config data").Error(),
		}, nil
	}

	launchUiValues, err := getAppLaunchSSOvalues(req.ReleaseName)
	if err != nil {
		a.log.Errorf("failed to SSO config for app %s, %v", req.ReleaseName, err)
		return &agentpb.UpdateAppValuesResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "failed to read app SSO config").Error(),
		}, nil
	}

	// populate template values, overriding with launchUiValues if needed
	updateAppConfig, marshaledOverrideValues, err := populateTemplateValues(appConfig, req.OverrideValues, launchUiValues, a.log)
	if err != nil {
		a.log.Errorf("failed to populate template values for app %s, %v", req.ReleaseName, err)
		return &agentpb.UpdateAppValuesResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "failed to prepare app values").Error(),
		}, nil
	}

	launchURL, err := executeStringTemplateValues(appConfig.Config.LaunchURL, req.OverrideValues)
	if err != nil {
		a.log.Errorf("failed to derive template launch URL for app %s, %v", req.ReleaseName, err)
		return &agentpb.UpdateAppValuesResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to prepare app values",
		}, nil
	}

	appConfig.Config.LaunchURL = launchURL
	appConfig.Config.InstallStatus = string(model.AppUpgradingStatus)
	if err := a.as.UpsertAppConfig(appConfig); err != nil {
		a.log.Errorf("failed to update app config data for app %s, %v", req.ReleaseName, err)
		return &agentpb.UpdateAppValuesResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "failed to update app config data").Error(),
		}, nil
	}

	deployReq := prepareAppDeployRequestFromSyncApp(updateAppConfig, marshaledOverrideValues)
	go a.upgradeAppWithWorkflow(deployReq, updateAppConfig)

	a.log.Infof("Triggerred app [%s] upgrade with values", updateAppConfig.Config.ReleaseName)
	return &agentpb.UpdateAppValuesResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "Triggerred app upgrade",
	}, nil

}

func (a *Agent) UpgradeApp(ctx context.Context, req *agentpb.UpgradeAppRequest) (*agentpb.UpgradeAppResponse, error) {
	a.log.Infof("Received UpgradeApp request %+v", req.AppConfig.ReleaseName)

	_, err := a.as.GetAppConfig(req.AppConfig.ReleaseName)
	if err != nil {
		a.log.Errorf("failed to read app %s config data, %v", req.AppConfig.ReleaseName, err)
		return &agentpb.UpgradeAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "failed to read app config data").Error(),
		}, nil
	}

	launchUiValues, err := getAppLaunchSSOvalues(req.AppConfig.ReleaseName)
	if err != nil {
		a.log.Errorf("failed to SSO config for app %s, %v", req.AppConfig.ReleaseName, err)
		return &agentpb.UpgradeAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err in getLanchUiValues").Error(),
		}, nil
	}

	templateValues, err := deriveTemplateValues(req.AppValues.OverrideValues, req.AppValues.TemplateValues)
	if err != nil {
		a.log.Errorf("failed to derive template values for app %s, %v", req.AppConfig.ReleaseName, err)
		return &agentpb.UpgradeAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to prepare app values",
		}, nil
	}

	launchURL, err := executeStringTemplateValues(req.AppConfig.LaunchURL, req.AppValues.OverrideValues)
	if err != nil {
		a.log.Errorf("failed to derive template launch URL for app %s, %v", req.AppConfig.ReleaseName, err)
		return &agentpb.UpgradeAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to prepare app values",
		}, nil
	}

	syncConfig := &agentpb.SyncAppData{
		Config: &agentpb.AppConfig{
			ReleaseName:         req.AppConfig.ReleaseName,
			AppName:             req.AppConfig.AppName,
			Version:             req.AppConfig.Version,
			Category:            req.AppConfig.Category,
			Description:         req.AppConfig.Description,
			ChartName:           req.AppConfig.ChartName,
			RepoName:            req.AppConfig.RepoName,
			RepoURL:             req.AppConfig.RepoURL,
			Namespace:           req.AppConfig.Namespace,
			CreateNamespace:     req.AppConfig.CreateNamespace,
			PrivilegedNamespace: req.AppConfig.PrivilegedNamespace,
			Icon:                req.AppConfig.Icon,
			LaunchURL:           launchURL,
			LaunchUIDescription: req.AppConfig.LaunchUIDescription,
			InstallStatus:       string(model.AppIntallingStatus),
			DefualtApp:          req.AppConfig.DefualtApp,
			PluginName:          req.AppConfig.PluginName,
			PluginDescription:   req.AppConfig.PluginDescription,
		},
		Values: &agentpb.AppValues{
			OverrideValues: req.AppValues.OverrideValues,
			LaunchUIValues: launchUiValues,
			TemplateValues: req.AppValues.TemplateValues,
		},
	}

	if err := a.as.UpsertAppConfig(syncConfig); err != nil {
		a.log.Errorf("failed to update app config data for app %s, %v", req.AppConfig.ReleaseName, err)
		return &agentpb.UpgradeAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to update app config data",
		}, nil
	}

	deployReq := prepareAppDeployRequestFromSyncApp(syncConfig, templateValues)

	go a.upgradeAppWithWorkflow(deployReq, syncConfig)

	a.log.Infof("Triggerred app [%s] upgrade", syncConfig.Config.ReleaseName)
	return &agentpb.UpgradeAppResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "Triggerred app upgrade",
	}, nil
}

func (a *Agent) installAppWithWorkflow(req *model.ApplicationInstallRequest,
	appConfig *agentpb.SyncAppData) {
	wd := workers.NewDeployment(a.tc, a.log)
	_, err := wd.SendEvent(context.TODO(), wd.GetWorkflowName(), string(model.AppInstallAction), req)
	if err != nil {
		appConfig.Config.InstallStatus = string(model.AppIntallFailedStatus)
		if err := a.as.UpsertAppConfig(appConfig); err != nil {
			a.log.Errorf("failed to update app config for app %s, %v", appConfig.Config.ReleaseName, err)
			return
		}
		a.log.Errorf("failed to send event to workflow for app %s, %v", appConfig.Config.ReleaseName, err)
		return
	}

	appConfig.Config.InstallStatus = string(model.AppIntalledStatus)
	if err := a.as.UpsertAppConfig(appConfig); err != nil {
		a.log.Errorf("failed to update app config for app %s, %v", appConfig.Config.ReleaseName, err)
		return
	}
}

func (a *Agent) unInstallAppWithWorkflow(req *model.ApplicationDeleteRequest, appConfig *agentpb.SyncAppData) {
	wd := workers.NewDeployment(a.tc, a.log)
	_, err := wd.SendDeleteEvent(context.TODO(), wd.GetWorkflowName(), string(model.AppUnInstallAction), req)
	if err != nil {
		a.log.Errorf("failed to send delete event to workflow for app %s, %v", req.ReleaseName, err)

		appConfig.Config.InstallStatus = string(model.AppIntalledStatus)
		if err := a.as.UpsertAppConfig(appConfig); err != nil {
			a.log.Errorf("failed to update app config status with Installed for app %s, %v", req.ReleaseName, err)
		}
		return
	}

	if err := a.as.DeleteAppConfigByReleaseName(req.ReleaseName); err != nil {
		a.log.Errorf("failed to delete installed app config record %s, %v", req.ReleaseName, err)
		return
	}
}

func (a *Agent) upgradeAppWithWorkflow(req *model.ApplicationInstallRequest,
	appConfig *agentpb.SyncAppData) {
	wd := workers.NewDeployment(a.tc, a.log)
	_, err := wd.SendEvent(context.TODO(), wd.GetWorkflowName(), string(model.AppUpgradeAction), req)
	if err != nil {
		appConfig.Config.InstallStatus = string(model.AppUpgradeFaileddStatus)
		if err := a.as.UpsertAppConfig(appConfig); err != nil {
			a.log.Errorf("failed to update app config for app %s, %v", appConfig.Config.ReleaseName, err)
			return
		}
		a.log.Errorf("failed to send event to workflow for app %s, %v", appConfig.Config.ReleaseName, err)
		return
	}

	appConfig.Config.InstallStatus = string(model.AppUpgradedStatus)
	if err := a.as.UpsertAppConfig(appConfig); err != nil {
		a.log.Errorf("failed to update app config for app %s, %v", appConfig.Config.ReleaseName, err)
		return
	}
}
