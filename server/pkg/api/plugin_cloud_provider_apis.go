package api

import (
	"context"
	"fmt"

	"github.com/kube-tarian/kad/server/pkg/pb/captenpluginspb"
)

func (s *Server) AddCloudProvider(ctx context.Context, request *captenpluginspb.AddCloudProviderRequest) (
	*captenpluginspb.AddCloudProviderResponse, error) {
	return &captenpluginspb.AddCloudProviderResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}

func (s *Server) UpdateCloudProviders(ctx context.Context, request *captenpluginspb.UpdateCloudProviderRequest) (
	*captenpluginspb.UpdateCloudProviderResponse, error) {
	return &captenpluginspb.UpdateCloudProviderResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}

func (s *Server) DeleteCloudProvider(ctx context.Context, request *captenpluginspb.DeleteCloudProviderRequest) (
	*captenpluginspb.DeleteCloudProviderResponse, error) {
	return &captenpluginspb.DeleteCloudProviderResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}

func (s *Server) GetCloudProviders(ctx context.Context, request *captenpluginspb.GetCloudProvidersRequest) (
	*captenpluginspb.GetCloudProvidersResponse, error) {
	return &captenpluginspb.GetCloudProvidersResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}

func (s *Server) GetCloudProvidersForLabels(ctx context.Context, request *captenpluginspb.GetCloudProvidersForLabelsRequest) (
	*captenpluginspb.GetCloudProvidersForLabelsResponse, error) {
	return &captenpluginspb.GetCloudProvidersForLabelsResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}
