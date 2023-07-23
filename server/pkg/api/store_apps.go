package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) AddStoreApp(ctx context.Context, request *serverpb.AddStoreAppRequest) (
	*serverpb.AddStoreAppResponse, error) {
	return &serverpb.AddStoreAppResponse{}, nil
}

func (s *Server) UpdateStoreApp(ctx context.Context, request *serverpb.UpdateStoreAppRequest) (
	*serverpb.UpdateStoreAppRsponse, error) {
	return &serverpb.UpdateStoreAppRsponse{}, nil
}

func (s *Server) DeleteStoreApp(ctx context.Context, request *serverpb.DeleteStoreAppRequest) (
	*serverpb.DeleteStoreAppResponse, error) {
	return &serverpb.DeleteStoreAppResponse{}, nil
}

func (s *Server) GetStoreApp(ctx context.Context, request *serverpb.GetStoreAppRequest) (
	*serverpb.GetStoreAppResponse, error) {
	return &serverpb.GetStoreAppResponse{}, nil
}

func (s *Server) GetStoreApps(ctx context.Context, request *serverpb.GetStoreAppsRequest) (
	*serverpb.GetStoreAppsResponse, error) {
	return &serverpb.GetStoreAppsResponse{}, nil
}
