package agent

import (
	"context"
	"fmt"

	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
)

func (a *Agent) AddCloudProvider(ctx context.Context, request *captenpluginspb.AddCloudProviderRequest) (
	*captenpluginspb.AddCloudProviderResponse, error) {
	return &captenpluginspb.AddCloudProviderResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}

func (a *Agent) UpdateCloudProviders(ctx context.Context, request *captenpluginspb.UpdateCloudProviderRequest) (
	*captenpluginspb.UpdateCloudProviderResponse, error) {
	return &captenpluginspb.UpdateCloudProviderResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}

func (a *Agent) DeleteCloudProvider(ctx context.Context, request *captenpluginspb.DeleteCloudProviderRequest) (
	*captenpluginspb.DeleteCloudProviderResponse, error) {
	return &captenpluginspb.DeleteCloudProviderResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}

func (a *Agent) GetCloudProviders(ctx context.Context, request *captenpluginspb.GetCloudProvidersRequest) (
	*captenpluginspb.GetCloudProvidersResponse, error) {
	return &captenpluginspb.GetCloudProvidersResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}

func (a *Agent) GetCloudProvidersForLabels(ctx context.Context, request *captenpluginspb.GetCloudProvidersForLabelRequest) (
	*captenpluginspb.GetCloudProvidersForLabelResponse, error) {
	return &captenpluginspb.GetCloudProvidersForLabelResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}
