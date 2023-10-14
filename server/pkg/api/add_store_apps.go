package api

import (
	"context"
	"encoding/hex"

	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
	"github.com/kube-tarian/kad/server/pkg/types"
)

func (s *Server) AddStoreApp(ctx context.Context, request *serverpb.AddStoreAppRequest) (
	*serverpb.AddStoreAppResponse, error) {
	_, err := validateRequest(ctx, request.AppConfig.AppName, request.AppConfig.Version)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.AddStoreAppResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Add store app [%s:%s] request recieved", request.AppConfig.AppName, request.AppConfig.Version)

	config := &types.StoreAppConfig{
		ReleaseName:         request.AppConfig.ReleaseName,
		AppName:             request.AppConfig.AppName,
		Version:             request.AppConfig.Version,
		Category:            request.AppConfig.Category,
		Description:         request.AppConfig.Description,
		ChartName:           request.AppConfig.ChartName,
		RepoName:            request.AppConfig.RepoName,
		RepoURL:             request.AppConfig.RepoURL,
		Namespace:           request.AppConfig.Namespace,
		CreateNamespace:     request.AppConfig.CreateNamespace,
		PrivilegedNamespace: request.AppConfig.PrivilegedNamespace,
		Icon:                hex.EncodeToString(request.AppConfig.Icon),
		LaunchURL:           request.AppConfig.LaunchURL,
		LaunchUIDescription: request.AppConfig.LaunchUIDescription,
		OverrideValues:      encodeBase64BytesToString(request.AppValues.OverrideValues),
		LaunchUIValues:      encodeBase64BytesToString(request.AppValues.LaunchUIValues),
		TemplateValues:      encodeBase64BytesToString(request.AppValues.TemplateValues),
	}

	if err := s.serverStore.AddOrUpdateStoreApp(config); err != nil {
		s.log.Errorf("failed to add app config to store, %v", err)
		return &serverpb.AddStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed add app config to store",
		}, nil
	}

	s.log.Infof("Add store app [%s:%s] request successful", request.AppConfig.AppName, request.AppConfig.Version)
	return &serverpb.AddStoreAppResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "app config is sucessfuly added to store",
	}, nil
}
