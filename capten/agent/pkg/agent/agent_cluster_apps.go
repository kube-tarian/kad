package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/pb/agentpb"
)

func (a *Agent) SyncApp(ctx context.Context, request *agentpb.SyncAppRequest) (
	*agentpb.SyncAppResponse, error) {
	if request.Data == nil {
		return &agentpb.SyncAppResponse{
			Status:        agentpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "invalid data passed",
		}, nil
	}

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
		appConfig := r.GetConfig()
		if len(appConfig.LaunchURL) == 0 {
			continue
		}

		cfgs = append(cfgs, &agentpb.AppLaunchConfig{
			ReleaseName: r.GetConfig().GetReleaseName(),
			Category:    r.GetConfig().GetCategory(),
			Description: r.GetConfig().GetLaunchUIDescription(),
			Icon:        r.GetConfig().GetIcon(),
			LaunchURL:   r.GetConfig().GetLaunchURL(),
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
