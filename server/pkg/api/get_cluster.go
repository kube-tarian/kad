package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

const (
	clusterNotFound = "no cluster found"
)

func (s *Server) GetClusterDetails(ctx context.Context, request *serverpb.GetClusterDetailsRequest) (
	*serverpb.GetClusterDetailsResponse, error) {
	orgId, err := validateOrgWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.GetClusterDetailsResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	s.log.Infof("GetClusterDetails request recieved, [org: %s]", orgId)
	cluster, err := s.serverStore.GetClusterForOrg(orgId)
	if err != nil {
		if strings.EqualFold(err.Error(), clusterNotFound) {
			s.log.Infof("cluster not found for org %s", orgId)
			return &serverpb.GetClusterDetailsResponse{
				Status:        serverpb.StatusCode_NOT_FOUND,
				StatusMessage: "cluster not found",
			}, nil
		}

		return &serverpb.GetClusterDetailsResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get cluster details",
		}, err
	}

	launchConfigList, err := s.getClusterAppLaunchesAgent(ctx, orgId, cluster.ClusterID)
	if err != nil {
		s.log.Error("failed to get cluster application launches from cache/agent: %v", err)
	}

	attributes := []*serverpb.ClusterAttribute{}
	data := &serverpb.ClusterInfo{
		ClusterID:        cluster.ClusterID,
		ClusterName:      cluster.ClusterName,
		AgentEndpoint:    cluster.Endpoint,
		Attributes:       attributes,
		AppLaunchConfigs: mapAgentAppLaunchConfigsToServer(launchConfigList),
	}

	s.log.Infof("GetClusterDetails request processed for cluster %s, [org: %s]", cluster.ClusterID, orgId)
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
