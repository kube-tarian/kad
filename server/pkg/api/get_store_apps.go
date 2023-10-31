package api

import (
	"context"
	"encoding/hex"

	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) GetStoreApps(ctx context.Context, request *serverpb.GetStoreAppsRequest) (
	*serverpb.GetStoreAppsResponse, error) {
	orgId, err := validateOrgWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.GetStoreAppsResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Get store apps request recieved, [org: %s]", orgId)

	cluster, err := s.serverStore.GetClusterForOrg(orgId)
	if err != nil {
		s.log.Errorf("failed to get clusterID for org %s, %v", orgId, err)
		return &serverpb.GetStoreAppsResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed get cluster details",
		}, err
	}

	configs, err := s.serverStore.GetAppsFromStore()
	if err != nil {
		s.log.Errorf("failed to get app config's from store, %v", err)
		return &serverpb.GetStoreAppsResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get app config's from store",
		}, nil
	}

	clusterGlobalValues, err := s.getClusterGlobalValues(orgId, cluster.ClusterID)
	if err != nil {
		s.log.Errorf("failed to get cluster global values, %v", err)
		return &serverpb.GetStoreAppsResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to fetch cluster global values values",
		}, nil
	}

	appsData := []*serverpb.StoreAppsData{}
	for _, config := range *configs {
		overrideValues, err := s.deriveTemplateOverrideValues(config.OverrideValues, clusterGlobalValues)
		if err != nil {
			s.log.Errorf("failed to update overrided store app values for app %s, %v", config.ReleaseName, err)
		}

		decodedIconBytes, _ := hex.DecodeString(config.Icon)
		appsData = append(appsData, &serverpb.StoreAppsData{
			AppConfigs: &serverpb.StoreAppConfig{
				AppName:             config.Name,
				Version:             config.Version,
				Category:            config.Category,
				Description:         config.Description,
				ChartName:           config.ChartName,
				RepoName:            config.RepoName,
				RepoURL:             config.RepoURL,
				Namespace:           config.Namespace,
				CreateNamespace:     config.CreateNamespace,
				PrivilegedNamespace: config.PrivilegedNamespace,
				Icon:                decodedIconBytes,
				LaunchURL:           config.LaunchURL,
				LaunchUIDescription: config.LaunchUIDescription,
				ReleaseName:         config.ReleaseName,
			},
			OverrideValues: overrideValues,
		})
	}

	s.log.Infof("Fetched %d store apps, [org: %s]", len(appsData), orgId)
	return &serverpb.GetStoreAppsResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "app config's are sucessfuly fetched from store",
		Data:          appsData,
	}, nil
}
