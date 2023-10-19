package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/credential"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) DeleteClusterRegistration(ctx context.Context, request *serverpb.DeleteClusterRegistrationRequest) (
	*serverpb.DeleteClusterRegistrationResponse, error) {
	orgId, err := validateOrgWithArgs(ctx, request.ClusterID)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.DeleteClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	s.log.Infof("Delete cluster registration request for cluster %s recieved, [org: %s]", request.ClusterID, orgId)
	s.agentHandeler.RemoveAgent(request.ClusterID)
	err = credential.DeleteClusterCerts(ctx, request.ClusterID)
	if err != nil {
		s.log.Errorf("failed to delete cert in vault for cluster %s, %v", request.ClusterID, err)
		return &serverpb.DeleteClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed delete register cluster",
		}, nil
	}

	err = s.serverStore.DeleteFullClusterAppLaunches(orgId, request.ClusterID)
	if err != nil {
		s.log.Errorf("failed to delete clusterappLaunches %s from db, %v", request.ClusterID, err)
		return &serverpb.DeleteClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed delete register cluster",
		}, nil
	}

	err = s.serverStore.DeleteCluster(orgId, request.ClusterID)
	if err != nil {
		s.log.Errorf("failed to delete cluster %s from db, %v", request.ClusterID, err)
		return &serverpb.DeleteClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed delete register cluster",
		}, nil
	}

	s.log.Infof("Delete cluster registration request for cluster %s successful, [org: %s]", request.ClusterID, orgId)
	return &serverpb.DeleteClusterRegistrationResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "cluster deletion success",
	}, nil
}
