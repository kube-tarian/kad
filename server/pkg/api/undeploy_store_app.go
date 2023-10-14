package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) UnDeployStoreApp(ctx context.Context, request *serverpb.UnDeployStoreAppRequest) (
	*serverpb.UnDeployStoreAppResponse, error) {
	orgId, err := validateRequest(ctx, request.ClusterID, request.ReleaseName)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.UnDeployStoreAppResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	s.log.Infof("UnDeploy store app %s request for cluster %s recieved, [org: %s]",
		request.ReleaseName, request.ClusterID, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, request.ClusterID)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &serverpb.UnDeployStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to undeploy the app",
		}, nil
	}

	req := &agentpb.UnInstallAppRequest{
		ReleaseName: request.ReleaseName,
	}
	resp, err := agent.GetClient().UnInstallApp(ctx, req)
	if err != nil {
		s.log.Errorf("failed to undeploy app, %v", err)
		return &serverpb.UnDeployStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to undeploy the app",
		}, nil
	}

	s.log.Infof("UnDeploy store app %s request request triggered for cluster %s, [org: %s]",
		request.ReleaseName, request.ClusterID, orgId)

	return &serverpb.UnDeployStoreAppResponse{
		Status:        serverpb.StatusCode(resp.Status),
		StatusMessage: resp.StatusMessage,
	}, nil
}
