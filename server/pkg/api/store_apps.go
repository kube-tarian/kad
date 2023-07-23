package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (a *Api) AddStoreApp(ctx context.Context, request *serverpb.AddStoreAppRequest) (
	*serverpb.AddStoreAppResponse, error) {
	return &serverpb.AddStoreAppResponse{}, nil
}

func (a *Api) UpdateStoreApp(ctx context.Context, request *serverpb.UpdateStoreAppRequest) (
	*serverpb.UpdateStoreAppRsponse, error) {
	return &serverpb.UpdateStoreAppRsponse{}, nil
}

func (a *Api) DeleteStoreApp(ctx context.Context, request *serverpb.DeleteStoreAppRequest) (
	*serverpb.DeleteStoreAppResponse, error) {
	return &serverpb.DeleteStoreAppResponse{}, nil
}

func (a *Api) GetStoreApp(ctx context.Context, request *serverpb.GetStoreAppRequest) (
	*serverpb.GetStoreAppResponse, error) {
	return &serverpb.GetStoreAppResponse{}, nil
}

func (a *Api) GetStoreApps(ctx context.Context, request *serverpb.GetStoreAppsRequest) (
	*serverpb.GetStoreAppsResponse, error) {
	return &serverpb.GetStoreAppsResponse{}, nil
}
