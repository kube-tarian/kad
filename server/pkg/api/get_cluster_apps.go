package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) GetClusterApps(ctx context.Context, request *serverpb.GetClusterAppsRequest) (
	*serverpb.GetClusterAppsResponse, error) {
	orgId, err := validateRequest(ctx, request.ClusterID)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.GetClusterAppsResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("GetClusterApps request recieved for cluster %s, [org: %s]", request.ClusterID, orgId)

	a, err := s.agentHandeler.GetAgent(orgId, request.ClusterID)
	if err != nil {
		s.log.Error("failed to connect to agent", err)
		return &serverpb.GetClusterAppsResponse{Status: serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to connect to agent"}, nil
	}

	resp, err := a.GetClient().GetClusterApps(ctx, &agentpb.GetClusterAppsRequest{})
	if err != nil {
		s.log.Error("failed to get cluster application from agent", err)
		return &serverpb.GetClusterAppsResponse{Status: serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get cluster application from agent"}, nil
	}

	s.log.Infof("Fetched %d installed apps from the cluster %s, [org: %s]", len(resp.AppData), request.ClusterID, orgId)
	return &serverpb.GetClusterAppsResponse{Status: serverpb.StatusCode_OK,
		StatusMessage: "successfully fetched the data from agent",
		AppConfigs:    mapAgentAppsToServerResp(resp.AppData)}, nil
}
