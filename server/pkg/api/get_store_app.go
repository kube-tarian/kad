package api

import (
	"context"
	"encoding/hex"

	"github.com/kube-tarian/kad/server/pkg/opentelemetry"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Server) GetStoreApp(ctx context.Context, request *serverpb.GetStoreAppRequest) (
	*serverpb.GetStoreAppResponse, error) {
	_, span := opentelemetry.GetTracer(request.AppName).
		Start(opentelemetry.BuildContext(ctx), "CaptenServer")
	defer span.End()

	span.SetAttributes(attribute.String("App Name", request.AppName))

	orgId, err := validateOrgWithArgs(ctx, request.AppName, request.Version)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.GetStoreAppResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Get store app %s:%s request recieved, [org: %s]", request.AppName, request.Version, orgId)

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

	s.log.Infof("Fetched store app %s:%s, [org: %s]", request.AppName, request.Version, orgId)
	return &serverpb.GetStoreAppResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "app config is sucessfuly fetched from store",
		AppConfig:     appConfig,
	}, nil

}
