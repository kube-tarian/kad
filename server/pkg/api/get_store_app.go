package api

import (
	"context"
	"encoding/hex"

	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) GetStoreApp(ctx context.Context, request *serverpb.GetStoreAppRequest) (
	*serverpb.GetStoreAppResponse, error) {
	_, err := validateRequest(ctx, request.AppName, request.Version)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.GetStoreAppResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	config, err := s.serverStore.GetAppFromStore(request.AppName, request.Version)
	if err != nil {
		s.log.Errorf("failed to get app config from store, %v", err)
		return &serverpb.GetStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get app config from store",
		}, nil
	}

	decodedIconBytes, _ := hex.DecodeString(config.Icon)
	appConfig := &serverpb.StoreAppConfig{
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
	}

	return &serverpb.GetStoreAppResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "app config is sucessfuly fetched from store",
		AppConfig:     appConfig,
	}, nil

}
