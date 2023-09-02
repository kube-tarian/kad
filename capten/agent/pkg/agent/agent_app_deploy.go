package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/workers"
)

func (a *Agent) InstallApp(ctx context.Context, request *agentpb.InstallAppRequest) (*agentpb.InstallAppResponse, error) {
	a.log.Infof("Recieved App Install request %+v", request)
	templateValues, err := deriveTemplateValues(request.AppValues.OverrideValues, request.AppValues.TemplateValues)
	if err != nil {
		a.log.Errorf("failed to derive template values, %v", err)
		return &agentpb.InstallAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to derive template values",
		}, err
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
			InstallStatus:       "Installing",
			DefualtApp:          request.AppConfig.DefualtApp,
		},
		Values: &agentpb.AppValues{
			OverrideValues: request.AppValues.OverrideValues,
			LaunchUIValues: request.AppValues.LaunchUIValues,
			TemplateValues: request.AppValues.TemplateValues,
		},
	}

	if err := a.as.UpsertAppConfig(syncConfig); err != nil {
		a.log.Errorf("failed to update sync app config, %v", err)
		return &agentpb.InstallAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to sync app config",
		}, err
	}

	go func() {
		wd := workers.NewDeployment(a.tc, a.log)
		_, err = wd.SendEvent(context.TODO(), "install",
			toAppDeployRequestFromSyncApp(syncConfig, templateValues))
		if err != nil {
			syncConfig.Config.InstallStatus = "Install Failed"
			if err := a.as.UpsertAppConfig(syncConfig); err != nil {
				a.log.Errorf("failed to update sync app config, %v", err)
				return
			}

			a.log.Errorf("failed to send store app install event, %v", err)
			return
		}

		syncConfig.Config.InstallStatus = "Installed"
		if err := a.as.UpsertAppConfig(syncConfig); err != nil {
			a.log.Errorf("failed to update sync app config, %v", err)
			return
		}
	}()

	a.log.Infof("Triggerred app [%s] deployment", request.AppConfig.ReleaseName)
	return &agentpb.InstallAppResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "Triggerred app deployment",
	}, nil
}
