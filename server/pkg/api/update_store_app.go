package api

import (
	"context"
	"encoding/hex"

	"github.com/kube-tarian/kad/server/pkg/opentelemetry"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
	"github.com/kube-tarian/kad/server/pkg/types"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Server) UpdateStoreApp(ctx context.Context, request *serverpb.UpdateStoreAppRequest) (
	*serverpb.UpdateStoreAppRsponse, error) {

	_, span := opentelemetry.GetTracer(request.AppConfig.AppName).
		Start(opentelemetry.BuildContext(ctx), "CaptenServer")
	defer span.End()

	span.SetAttributes(attribute.String("App Name", request.AppConfig.AppName))

	_, err := validateOrgWithArgs(ctx, request.AppConfig.AppName, request.AppConfig.Version)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.UpdateStoreAppRsponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Update store app [%s:%s] request recieved", request.AppConfig.AppName, request.AppConfig.Version)

	//TODO check store app exist in DB
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
		OverrideValues:      request.AppValues.OverrideValues,
		LaunchUIValues:      request.AppValues.LaunchUIValues,
		TemplateValues:      request.AppValues.TemplateValues,
	}

	if err := s.serverStore.AddOrUpdateStoreApp(config); err != nil {
		s.log.Errorf("failed to update app config in store, %v", err)
		return &serverpb.UpdateStoreAppRsponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to update app config in store",
		}, nil
	}

	s.log.Infof("Update store app [%s:%s] request successful", request.AppConfig.AppName, request.AppConfig.Version)
	return &serverpb.UpdateStoreAppRsponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "app config is sucessfuly updated",
	}, nil
}
