package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/opentelemetry"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Server) DeleteStoreApp(ctx context.Context, request *serverpb.DeleteStoreAppRequest) (
	*serverpb.DeleteStoreAppResponse, error) {
	_, span := opentelemetry.GetTracer(request.AppName).
		Start(opentelemetry.BuildContext(ctx), "CaptenServer")
	defer span.End()

	span.SetAttributes(attribute.String("Cluster Name", request.AppName))
	span.SetAttributes(attribute.String("Agent EndPoint", request.Version))
	err := validateArgs(request.AppName, request.Version)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.DeleteStoreAppResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Delete store app [%s:%s] request recieved", request.AppName, request.Version)

	if err := s.serverStore.DeleteAppInStore(request.AppName, request.Version); err != nil {
		s.log.Errorf("failed to delete app config from store, %v", err)
		return &serverpb.DeleteStoreAppResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to delete app config from store",
		}, nil
	}

	s.log.Infof("Delete store app [%s:%s] request successful", request.AppName, request.Version)
	return &serverpb.DeleteStoreAppResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "app config is sucessfuly deleted",
	}, nil

}
