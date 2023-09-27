package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/agent/pkg/workers"
	"github.com/pkg/errors"
)

func (a *Agent) InstallApp(ctx context.Context, request *agentpb.InstallAppRequest) (*agentpb.InstallAppResponse, error) {
	a.log.Infof("Recieved App Install request %+v", request)
	templateValues, err := deriveTemplateValues(request.AppValues.OverrideValues, request.AppValues.TemplateValues)
	if err != nil {
		a.log.Errorf("failed to derive template values for app %s, %v", request.AppConfig.ReleaseName, err)
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
			LaunchURL:           request.AppConfig.LaunchURL,
			LaunchUIDescription: request.AppConfig.LaunchUIDescription,
			InstallStatus:       string(appIntallingStatus),
			DefualtApp:          request.AppConfig.DefualtApp,
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

	if request.Namespace == "" {
		a.log.Errorf("namespace is empty")
		return &agentpb.UnInstallAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "namespace is missing in request",
		}, nil
	}
	if request.ReleaseName == "" {
		a.log.Errorf("release name is empty")
		return &agentpb.UnInstallAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "release name is missing in request",
		}, nil
	}

	req := &model.ApplicationDeleteRequest{
		PluginName:  "helm",
		Namespace:   request.Namespace,
		ReleaseName: request.ReleaseName,
		ClusterName: "capten",
		Timeout:     10,
	}

	if err := a.unInstallAppWithWorkflow(req); err != nil {
		a.log.Errorf("failed to uninstall the app %s, %v", req.ReleaseName, err)
		return &agentpb.UnInstallAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to undeploy the app",
		}, nil
	}

	if err := a.as.DeleteAppConfigByReleaseName(request.ReleaseName); err != nil {
		a.log.Errorf("failed to delete installed app config record %s, %v", req.ReleaseName, err)
		return &agentpb.UnInstallAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to undeploy the app",
		}, nil
	}

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

	launchUiValues, err := GetSSOvalues(req.ReleaseName)
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

	appConfig.Config.InstallStatus = string(appUpgradingStatus)
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
	a.log.Infof("Received UpgradeApp request %+v", req.ReleaseName)

	appConfig, err := a.as.GetAppConfig(req.ReleaseName)
	if err != nil {
		a.log.Errorf("failed to read app %s config data, %v", req.ReleaseName, err)
		return &agentpb.UpgradeAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "failed to read app config data").Error(),
		}, nil
	}

	launchUiValues, err := GetSSOvalues(req.ReleaseName)
	if err != nil {
		a.log.Errorf("failed to SSO config for app %s, %v", req.ReleaseName, err)
		return &agentpb.UpgradeAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err in getLanchUiValues").Error(),
		}, nil
	}

	newAppConfig, marshaledOverrideValues, err := populateTemplateValues(appConfig, nil, launchUiValues, a.log)
	if err != nil {
		return &agentpb.UpgradeAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err populating template values").Error(),
		}, nil
	}

	installReq := prepareAppDeployRequestFromSyncApp(newAppConfig, marshaledOverrideValues)
	installReq.Version = req.GetVersion()
	go a.upgradeAppWithWorkflow(installReq, newAppConfig)

	a.log.Infof("Triggerred app [%s] upgrade", newAppConfig.Config.ReleaseName)
	return &agentpb.UpgradeAppResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "Triggerred app upgrade",
	}, nil
}

func (a *Agent) installAppWithWorkflow(req *model.ApplicationInstallRequest,
	appConfig *agentpb.SyncAppData) {
	wd := workers.NewDeployment(a.tc, a.log)
	_, err := wd.SendEvent(context.TODO(), string(appInstallAction), req)
	if err != nil {
		appConfig.Config.InstallStatus = string(appIntallFailedStatus)
		if err := a.as.UpsertAppConfig(appConfig); err != nil {
			a.log.Errorf("failed to update app config for app %s, %v", appConfig.Config.ReleaseName, err)
			return
		}
		a.log.Errorf("failed to send event to workflow for app %s, %v", appConfig.Config.ReleaseName, err)
		return
	}

	appConfig.Config.InstallStatus = string(appIntalledStatus)
	if err := a.as.UpsertAppConfig(appConfig); err != nil {
		a.log.Errorf("failed to update app config for app %s, %v", appConfig.Config.ReleaseName, err)
		return
	}
}

func (a *Agent) unInstallAppWithWorkflow(req *model.ApplicationDeleteRequest) error {
	wd := workers.NewDeployment(a.tc, a.log)
	_, err := wd.SendDeleteEvent(context.TODO(), string(appUnInstallAction), req)
	if err != nil {
		a.log.Errorf("failed to send delete event to workflow for app %s, %v", req.ReleaseName, err)
		return err
	}
	return nil
}

func (a *Agent) upgradeAppWithWorkflow(req *model.ApplicationInstallRequest,
	appConfig *agentpb.SyncAppData) {
	wd := workers.NewDeployment(a.tc, a.log)
	_, err := wd.SendEvent(context.TODO(), string(appUpgradeAction), req)
	if err != nil {
		appConfig.Config.InstallStatus = string(appUpgradeFaileddStatus)
		if err := a.as.UpsertAppConfig(appConfig); err != nil {
			a.log.Errorf("failed to update app config for app %s, %v", appConfig.Config.ReleaseName, err)
			return
		}
		a.log.Errorf("failed to send event to workflow for app %s, %v", appConfig.Config.ReleaseName, err)
		return
	}

	appConfig.Config.InstallStatus = string(appUpgradeAction)
	if err := a.as.UpsertAppConfig(appConfig); err != nil {
		a.log.Errorf("failed to update app config for app %s, %v", appConfig.Config.ReleaseName, err)
		return
	}
}
