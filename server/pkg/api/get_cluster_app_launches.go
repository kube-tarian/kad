package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/opentelemetry"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Server) GetClusterAppLaunchConfigs(ctx context.Context, request *serverpb.GetClusterAppLaunchConfigsRequest) (
	*serverpb.GetClusterAppLaunchConfigsResponse, error) {

	_, span := opentelemetry.GetTracer(request.ClusterID).
		Start(opentelemetry.BuildContext(ctx), "CaptenServer")
	defer span.End()

	span.SetAttributes(attribute.String("Cluster Id", request.ClusterID))

	orgId, err := validateOrgWithArgs(ctx, request.ClusterID)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.GetClusterAppLaunchConfigsResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("GetClusterAppLaunchConfigs request recieved for cluster %s, [org: %s]", request.ClusterID, orgId)

	launchConfigList, err := s.getClusterAppLaunchesAgent(ctx, orgId, request.ClusterID)
	if err != nil {
		s.log.Error("failed to get cluster application launches from agent", err)
		return &serverpb.GetClusterAppLaunchConfigsResponse{Status: serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get cluster application launches from agent"}, err
	}

	s.log.Infof("Fetched %d app launch UIs from the cluster %s, [org: %s]", len(launchConfigList), request.ClusterID, orgId)
	return &serverpb.GetClusterAppLaunchConfigsResponse{Status: serverpb.StatusCode_OK,
		StatusMessage:   "successfully fetched the data from agent",
		AppLaunchConfig: mapAgentAppLaunchConfigsToServer(launchConfigList)}, nil
}
