package api

import (
	"context"
	"time"

	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) GetCluster(ctx context.Context, request *serverpb.GetClusterRequest) (
	*serverpb.GetClusterResponse, error) {
	orgId, err := validateRequest(ctx, request.ClusterID)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.GetClusterResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	s.log.Infof("GetCluster request recieved for cluster %s, [org: %s]", orgId, request.ClusterID)
	clusterDetails, err := s.serverStore.GetClusterDetails(orgId, request.ClusterID)
	if err != nil {
		s.log.Errorf("failed to get cluster %s, %v", request.ClusterID, err)
		return &serverpb.GetClusterResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed get cluster details",
		}, err
	}

	resp, err := s.getClusterAppLaunchesFromCacheOrAgent(ctx, orgId, request.ClusterID)
	if err != nil || resp == nil || resp.Status != agentpb.StatusCode_OK {
		s.log.Error("failed to get cluster application launches from cache/agent: %v", resp)
		return &serverpb.GetClusterResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed get cluster app lauches",
		}, err
	}

	attributes := []*serverpb.ClusterAttribute{}
	data := &serverpb.ClusterInfo{
		ClusterID:        request.ClusterID,
		ClusterName:      clusterDetails.ClusterName,
		AgentEndpoint:    clusterDetails.Endpoint,
		Attributes:       attributes,
		AppLaunchConfigs: mapAgentAppLauncesToServerResp(resp.LaunchConfigList),
	}

	s.log.Infof("GetCluster request processed for cluster %s, [org: %s]", request.ClusterID, orgId)
	return &serverpb.GetClusterResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "get cluster details success",
		Data:          data,
	}, nil
}

func (s *Server) getClusterAppLaunchesFromCacheOrAgent(ctx context.Context, orgId, clusterID string) (
	*agentpb.GetClusterAppLaunchesResponse, error) {
	currentTime := time.Now()
	lastFetchedTime := time.Unix(s.orgClusterIDCache[orgId+"-"+clusterID], 0)

	if currentTime.After(lastFetchedTime) {
		// cache expired re-trigger the cache
		agentClient, aErr := s.agentHandeler.GetAgent(orgId, clusterID)
		if aErr == nil {
			resp, err := agentClient.GetClient().GetClusterAppLaunches(ctx, &agentpb.GetClusterAppLaunchesRequest{})
			if err == nil {
				updateErr := s.serverStore.UpdateClusterAppLaunches(orgId, clusterID, resp.LaunchConfigList)
				if updateErr == nil {
					s.mutex.Lock()
					s.orgClusterIDCache[orgId+"-"+clusterID] = currentTime.Add(delayTimeinMin * time.Minute).Unix()
					s.mutex.Unlock()

					return resp, err
				}

			}
		}
	}

	// If any failure happens return from cache.
	return s.serverStore.GetClusterAppLaunches(orgId, clusterID)
}
