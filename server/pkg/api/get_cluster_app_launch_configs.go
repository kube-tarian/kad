package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) GetClusterAppLaunchConfigs(ctx context.Context, request *serverpb.GetClusterAppLaunchConfigsRequest) (
	*serverpb.GetClusterAppLaunchConfigsResponse, error) {
	orgId, err := validateRequest(ctx, request.ClusterID)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.GetClusterAppLaunchConfigsResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("GetClusterAppLaunchConfigs request recieved for cluster %s, [org: %s]", request.ClusterID, orgId)

	resp, err := s.getClusterAppLaunchesFromCacheOrAgent(ctx, orgId, request.ClusterID)
	if err != nil {
		s.log.Error("failed to get cluster application launches from agent", err)
		return &serverpb.GetClusterAppLaunchConfigsResponse{Status: serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get cluster application launches from agent"}, err
	}

	s.log.Infof("Fetched %d app launch UIs from the cluster %s, [org: %s]", len(resp.LaunchConfigList), request.ClusterID, orgId)
	return &serverpb.GetClusterAppLaunchConfigsResponse{Status: serverpb.StatusCode_OK,
		StatusMessage:   "successfully fetched the data from agent",
		AppLaunchConfig: mapAgentAppLauncesToServerResp(resp.LaunchConfigList)}, nil
}
