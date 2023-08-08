package api

import (
	"context"
	"encoding/json"

	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) GetClusterApps(ctx context.Context, request *serverpb.GetClusterAppsRequest) (
	*serverpb.GetClusterAppsResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if len(orgId) == 0 {
		s.log.Error("organizationID is missing in the request")
		return &serverpb.GetClusterAppsResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "organizationID is missing",
		}, nil
	}

	a, err := s.agentHandeler.GetAgent(orgId, request.ClusterID)
	if err != nil {
		return &serverpb.GetClusterAppsResponse{Status: serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to connect to agent"}, err
	}

	resp, err := a.GetClient().GetClusterApps(ctx, &agentpb.GetClusterAppsRequest{})
	if err != nil {
		return &serverpb.GetClusterAppsResponse{Status: serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "faield to get cluster application from agent"}, err
	}

	appConfigData, err := json.Marshal(resp.AppData)
	if err != nil {
		return &serverpb.GetClusterAppsResponse{Status: serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "faield to marshall appConfig"}, err
	}

	var clusterAppConfig []*serverpb.ClusterAppConfig
	err = json.Unmarshal(appConfigData, &clusterAppConfig)
	if err != nil {
		return &serverpb.GetClusterAppsResponse{Status: serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Unmarshall of appConfig failed"}, err
	}

	return &serverpb.GetClusterAppsResponse{Status: serverpb.StatusCode_OK, StatusMessage: "successfully fetched the data from agent",
		AppConfigs: clusterAppConfig}, nil
}

func (s *Server) GetClusterAppLaunchConfigs(ctx context.Context, request *serverpb.GetClusterAppLaunchConfigsRequest) (
	*serverpb.GetClusterAppLaunchConfigsResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if len(orgId) == 0 {
		s.log.Error("organizationID is missing in the request")
		return &serverpb.GetClusterAppLaunchConfigsResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "organizationID is missing",
		}, nil
	}

	a, err := s.agentHandeler.GetAgent(orgId, request.ClusterID)
	if err != nil {
		return &serverpb.GetClusterAppLaunchConfigsResponse{Status: serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to connect to agent"}, err
	}

	resp, err := a.GetClient().GetClusterAppLaunches(ctx, &agentpb.GetClusterAppLaunchesRequest{})
	if err != nil {
		return &serverpb.GetClusterAppLaunchConfigsResponse{Status: serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "faield to get cluster application launches from agent"}, err
	}

	appConfigData, err := json.Marshal(resp.LaunchConfigList)
	if err != nil {
		return &serverpb.GetClusterAppLaunchConfigsResponse{Status: serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "faield to marshall app launches"}, err
	}

	var clusterAppLaunchConfig []*serverpb.AppLaunchConfig
	err = json.Unmarshal(appConfigData, &clusterAppLaunchConfig)
	if err != nil {
		return &serverpb.GetClusterAppLaunchConfigsResponse{Status: serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Unmarshall of app launches failed"}, err
	}

	return &serverpb.GetClusterAppLaunchConfigsResponse{Status: serverpb.StatusCode_OK, StatusMessage: "successfully fetched the data from agent",
		AppLaunchConfig: clusterAppLaunchConfig}, nil
}

func (s *Server) GetClusterApp(ctx context.Context, request *serverpb.GetClusterAppRequest) (
	*serverpb.GetClusterAppResponse, error) {
	return &serverpb.GetClusterAppResponse{}, nil
}
