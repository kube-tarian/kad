package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) GetClusterApps(ctx context.Context, request *serverpb.GetClusterAppsRequest) (
	*serverpb.GetClusterAppsResponse, error) {
	return &serverpb.GetClusterAppsResponse{}, nil
}

func (s *Server) GetClusterAppLaunchConfigs(ctx context.Context, request *serverpb.GetClusterAppLaunchConfigsRequest) (
	*serverpb.GetClusterAppLaunchConfigsResponse, error) {
	return &serverpb.GetClusterAppLaunchConfigsResponse{}, nil
}

func (s *Server) GetClusterApp(ctx context.Context, request *serverpb.GetClusterAppRequest) (
	*serverpb.GetClusterAppResponse, error) {
	return &serverpb.GetClusterAppResponse{}, nil
}
