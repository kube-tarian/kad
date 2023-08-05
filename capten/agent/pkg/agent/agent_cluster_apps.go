package agent

import (
	"context"
	"fmt"

	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
)

func (a *Agent) SyncApp(ctx context.Context, request *agentpb.SyncAppRequest) (*agentpb.SyncAppResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("nil agentpb.SyncAppRequest")
	}
	if err := a.as.UpsertAppConfig(request.GetData()); err != nil {
		return nil, err
	}

	return &agentpb.SyncAppResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: agentpb.StatusCode_name[int32(agentpb.StatusCode_OK)],
	}, nil
}

func (a *Agent) GetClusterApps(ctx context.Context, request *agentpb.GetClusterAppsRequest) (*agentpb.GetClusterAppsResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("nil agentpb.GetClusterAppsRequest")
	}
	res, err := a.as.GetAllApps()
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return &agentpb.GetClusterAppsResponse{
			Status:        agentpb.StatusCode_NOT_FOUND,
			StatusMessage: agentpb.StatusCode_name[int32(agentpb.StatusCode_NOT_FOUND)],
		}, nil
	}

	appData := make([]*agentpb.AppData, 0)
	for _, r := range res {
		appData = append(appData, &agentpb.AppData{
			Config: r.GetConfig(),
		})
	}

	return &agentpb.GetClusterAppsResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: agentpb.StatusCode_name[int32(agentpb.StatusCode_OK)],
		AppData:       appData,
	}, nil

}

func (a *Agent) GetClusterAppLaunches(ctx context.Context, request *agentpb.GetClusterAppLaunchesRequest) (*agentpb.GetClusterAppLaunchesResponse, error) {

	res, err := a.GetClusterApps(context.TODO(), &agentpb.GetClusterAppsRequest{})
	if err != nil {
		return nil, err
	}

	if len(res.GetAppData()) == 0 {
		return &agentpb.GetClusterAppLaunchesResponse{
			Status:        agentpb.StatusCode_NOT_FOUND,
			StatusMessage: agentpb.StatusCode_name[int32(agentpb.StatusCode_NOT_FOUND)],
		}, nil
	}

	cfgs := make([]*agentpb.AppLaunchConfig, 0)
	for _, r := range res.GetAppData() {
		cfg := &agentpb.AppLaunchConfig{
			ReleaseName:       r.GetConfig().GetReleaseName(),
			Category:          r.GetConfig().GetCategory(),
			Description:       r.GetConfig().GetDescription(),
			Icon:              r.GetConfig().GetIcon(),
			LaunchURL:         r.GetConfig().GetLaunchURL(),
			LaunchRedirectURL: r.GetConfig().GetLaunchRedirectURL(),
		}
		cfgs = append(cfgs, cfg)
	}

	return &agentpb.GetClusterAppLaunchesResponse{
		LaunchConfigList: cfgs,
		Status:           agentpb.StatusCode_OK,
		StatusMessage:    agentpb.StatusCode_name[int32(agentpb.StatusCode_OK)],
	}, nil
}

func (a *Agent) GetClusterAppConfig(ctx context.Context, request *agentpb.GetClusterAppConfigRequest) (*agentpb.GetClusterAppConfigResponse, error) {

	res, err := a.as.GetAppConfig("release_name", request.GetReleaseName())

	if err != nil && err.Error() == "not found" {
		return &agentpb.GetClusterAppConfigResponse{
			Status:        agentpb.StatusCode_NOT_FOUND,
			StatusMessage: agentpb.StatusCode_name[int32(agentpb.StatusCode_NOT_FOUND)],
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return &agentpb.GetClusterAppConfigResponse{
		AppConfig:     res.GetConfig(),
		Status:        agentpb.StatusCode_OK,
		StatusMessage: agentpb.StatusCode_name[int32(agentpb.StatusCode_OK)],
	}, nil

}

func (a *Agent) GetClusterAppValues(ctx context.Context, request *agentpb.GetClusterAppValuesRequest) (*agentpb.GetClusterAppValuesResponse, error) {

	res, err := a.as.GetAppConfig("release_name", request.GetReleaseName())

	if err != nil && err.Error() == "not found" {
		return &agentpb.GetClusterAppValuesResponse{
			Status:        agentpb.StatusCode_NOT_FOUND,
			StatusMessage: agentpb.StatusCode_name[int32(agentpb.StatusCode_NOT_FOUND)],
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return &agentpb.GetClusterAppValuesResponse{
		Values:        res.GetValues(),
		Status:        agentpb.StatusCode_OK,
		StatusMessage: agentpb.StatusCode_name[int32(agentpb.StatusCode_OK)],
	}, nil
}
