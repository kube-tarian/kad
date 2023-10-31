package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) GetStoreAppValues(ctx context.Context, request *serverpb.GetStoreAppValuesRequest) (
	*serverpb.GetStoreAppValuesResponse, error) {
	orgId, err := validateOrgWithArgs(ctx, request.ClusterID)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.GetStoreAppValuesResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Get store app [%s:%s] values request for cluster %s recieved, [org: %s]",
		request.AppName, request.Version, request.ClusterID, orgId)

	config, err := s.serverStore.GetAppFromStore(request.AppName, request.Version)
	if err != nil {
		s.log.Errorf("failed to get store app values, %v", err)
		return &serverpb.GetStoreAppValuesResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get store app values",
		}, nil
	}

	clusterGlobalValues, err := s.getClusterGlobalValues(orgId, request.ClusterID)
	if err != nil {
		s.log.Errorf("failed to get cluster global values, %v", err)
		return &serverpb.GetStoreAppValuesResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to fetch cluster global values values",
		}, nil
	}

	overrideValues, err := s.deriveTemplateOverrideValues(config.OverrideValues, clusterGlobalValues)
	if err != nil {
		s.log.Errorf("failed to update overrided store app values, %v", err)
		return &serverpb.GetStoreAppValuesResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to update overrided store app values",
		}, nil
	}

	s.log.Infof("Fetched store app [%s:%s] values request for cluster %s successful, [org: %s]",
		request.AppName, request.Version, request.ClusterID, orgId)
	return &serverpb.GetStoreAppValuesResponse{
		Status:         serverpb.StatusCode_OK,
		StatusMessage:  "store app values sucessfuly fetched",
		OverrideValues: overrideValues,
	}, nil
}
