package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/opentelemetry"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) GetClusters(ctx context.Context, request *serverpb.GetClustersRequest) (
	*serverpb.GetClustersResponse, error) {

	_, span := opentelemetry.GetTracer("Get Clusters").
		Start(opentelemetry.BuildContext(ctx), "CaptenServer")
	defer span.End()

	orgId, err := validateOrgWithArgs(ctx)

	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.GetClustersResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	s.log.Infof("GetClusters request recieved, [org: %s]", orgId)
	clusterDetails, err := s.serverStore.GetClusters(orgId)
	if err != nil {
		s.log.Errorf("failed to get clusters, %v", err)
		return &serverpb.GetClustersResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed get cluster details",
		}, err
	}

	var data []*serverpb.ClusterInfo
	for _, cluster := range clusterDetails {
		launchConfigList, err := s.getClusterAppLaunchesAgent(ctx, orgId, cluster.ClusterID)
		if err != nil {
			s.log.Errorf("failed to get cluster appLaunches for cluster: %s, %v", cluster.ClusterID, err)
			continue
		}

		attributes := []*serverpb.ClusterAttribute{}
		data = append(data, &serverpb.ClusterInfo{
			ClusterID:        cluster.ClusterID,
			ClusterName:      cluster.ClusterName,
			AgentEndpoint:    cluster.Endpoint,
			Attributes:       attributes,
			AppLaunchConfigs: mapAgentAppLaunchConfigsToServer(launchConfigList),
		})
	}

	s.log.Infof("Fetched %d clusters, [org: %s]", len(data), orgId)
	return &serverpb.GetClustersResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "get cluster details success",
		Data:          data,
	}, nil
}
