package api

import (
	"context"
	"encoding/hex"

	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) GetStoreApps(ctx context.Context, request *serverpb.GetStoreAppsRequest) (
	*serverpb.GetStoreAppsResponse, error) {
	_, err := validateOrgWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.GetStoreAppsResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	configs, err := s.serverStore.GetAppsFromStore()
	if err != nil {
		s.log.Errorf("failed to get app config's from store, %v", err)
		return &serverpb.GetStoreAppsResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get app config's from store",
		}, nil
	}

	appsData := []*serverpb.StoreAppsData{}
	for _, config := range *configs {
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
			OverrideValues: config.OverrideValues,
		})
	}

	return &serverpb.GetStoreAppsResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "app config's are sucessfuly fetched from store",
		Data:          appsData,
	}, nil
}
