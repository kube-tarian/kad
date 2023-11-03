package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/captenpluginspb"
)

func (s *Server) GetManagedClusters(ctx context.Context, request *captenpluginspb.GetManagedClustersRequest) (
	*captenpluginspb.GetManagedClustersResponse, error) {
	s.log.Infof("Get Managed Clusters request recieved")
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.GetManagedClustersResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Get Managed Clusters request for cluster %s recieved, [org: %s]",
		clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.GetManagedClustersResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster GitProject",
		}, nil
	}

	response, err := agent.GetCaptenPluginsClient().GetManagedClusters(context.Background(), request)
	if err != nil {
		s.log.Errorf("failed to get the managed clusters, %v", err)
		return &captenpluginspb.GetManagedClustersResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Managed Clusters",
		}, nil
	}

	if response.Status != captenpluginspb.StatusCode_OK {
		s.log.Errorf("failed to get the managed clusters")
		return &captenpluginspb.GetManagedClustersResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the managed clusters",
		}, nil
	}

	s.log.Infof("Fetched %d managed clusters request for cluster %s processed, [org: %s]",
		len(response.Clusters), clusterId, orgId)
	return &captenpluginspb.GetManagedClustersResponse{
		Clusters:      response.Clusters,
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}
