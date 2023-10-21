package api

import (
	"context"
	"fmt"

	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) GetClusterDetails(ctx context.Context, request *serverpb.GetClusterDetailsRequest) (
	*serverpb.GetClusterDetailsResponse, error) {
	orgId, err := validateOrgWithArgs(ctx, request.ClusterID)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.GetClusterDetailsResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	s.log.Infof("GetClusterDetails request recieved for cluster %s, [org: %s]", request.ClusterID, orgId)
	clusterDetails, err := s.serverStore.GetClusterDetails(orgId, request.ClusterID)
	if err != nil {
		s.log.Errorf("failed to get cluster %s, %v", request.ClusterID, err)
		return &serverpb.GetClusterDetailsResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed get cluster details",
		}, err
	}

	launchConfigList, err := s.getClusterAppLaunchesAgent(ctx, orgId, request.ClusterID)
	if err != nil {
		s.log.Error("failed to get cluster application launches from cache/agent: %v", err)
		return &serverpb.GetClusterDetailsResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed get cluster app launches",
		}, err
	}

	attributes := []*serverpb.ClusterAttribute{}
	data := &serverpb.ClusterInfo{
		ClusterID:        request.ClusterID,
		ClusterName:      clusterDetails.ClusterName,
		AgentEndpoint:    clusterDetails.Endpoint,
		Attributes:       attributes,
		AppLaunchConfigs: mapAgentAppLaunchConfigsToServer(launchConfigList),
	}

	s.log.Infof("GetClusterDetails request processed for cluster %s, [org: %s]", request.ClusterID, orgId)
	return &serverpb.GetClusterDetailsResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "get cluster details success",
		Data:          data,
	}, nil
}

func (s *Server) GetCluster(ctx context.Context, request *serverpb.GetClusterRequest) (
	*serverpb.GetClusterResponse, error) {
	orgId, err := validateOrgWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.GetClusterResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	s.log.Infof("GetCluster request recieved, [org: %s]", orgId)
	cluster, err := s.serverStore.GetClusterForOrg(orgId)
	if err != nil {
		s.log.Errorf("failed to get clusterID for org %s, %v", orgId, err)
		return &serverpb.GetClusterResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed get cluster details",
		}, err
	}

	s.log.Infof("GetCluster request processed for cluster %s, [org: %s]", cluster.ClusterID, orgId)
	return &serverpb.GetClusterResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "get cluster ID details success",
		ClusterID:     cluster.ClusterID,
		ClusterName:   cluster.ClusterName,
		AgentEndpoint: cluster.Endpoint,
	}, nil
}

func (s *Server) getClusterAppLaunchesAgent(ctx context.Context, orgId, clusterID string) (
	[]*agentpb.AppLaunchConfig, error) {
	agentClient, err := s.agentHandeler.GetAgent(orgId, clusterID)
	if err != nil {
		return nil, err
	}

	resp, err := agentClient.GetClient().GetClusterAppLaunches(ctx, &agentpb.GetClusterAppLaunchesRequest{})
	if err != nil {
		return nil, err
	}
	if resp.Status != agentpb.StatusCode_OK {
		return nil, fmt.Errorf(resp.StatusMessage)
	}

	return resp.LaunchConfigList, nil
}
