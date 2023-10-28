package agent

import (
	"context"

	"github.com/google/uuid"
	"github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
)

func (a *Agent) AddCrossplanProvider(ctx context.Context, request *captenpluginspb.AddCrossplanProviderRequest) (
	*captenpluginspb.AddCrossplanProviderResponse, error) {
	if err := validateArgs(request.CloudType, request.ProviderName, request.CloudProviderId); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.AddCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	a.log.Infof("Add Crossplane Provider %s request recieved", request.CloudProviderId)
	id := uuid.New()
	provider := model.CrossplaneProvider{
		Id:              id.String(),
		CloudType:       request.CloudType,
		ProviderName:    request.ProviderName,
		CloudProviderId: request.CloudProviderId,
		Status:          "added",
	}

	if err := a.as.InsertCrossplaneProvider(&provider); err != nil {
		a.log.Errorf("failed to store crossplane provider to DB, %v", err)
		return &captenpluginspb.AddCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add crossplane provider in db",
		}, nil
	}

	a.log.Infof("Crossplane Provider %s added with id %s", request.ProviderName, id.String())
	return &captenpluginspb.AddCrossplanProviderResponse{
		Id:            id.String(),
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) DeleteCrossplanProvider(ctx context.Context, request *captenpluginspb.DeleteCrossplanProviderRequest) (
	*captenpluginspb.DeleteCrossplanProviderResponse, error) {

	if err := a.as.DeleteCrossplaneProviderById(request.Id); err != nil {
		a.log.Errorf("failed to delete crossplane provider from DB, %v", err)
		return &captenpluginspb.DeleteCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete crossplane provider from db",
		}, nil
	}

	a.log.Infof("Crossplane Provider with id %s deleted", request.Id)
	return &captenpluginspb.DeleteCrossplanProviderResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) GetCrossplanProviders(ctx context.Context, _ *captenpluginspb.GetCrossplanProvidersRequest) (
	*captenpluginspb.GetCrossplanProvidersResponse, error) {

	providers, err := a.as.GetCrossplaneProviders()
	if err != nil {
		a.log.Errorf("failed to fetch crossplane providers from DB, %v", err)
		return &captenpluginspb.GetCrossplanProvidersResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch crossplane providers from db",
		}, nil
	}

	if len(providers) == 0 {
		return &captenpluginspb.GetCrossplanProvidersResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "No crossplane providers found",
		}, nil
	}

	a.log.Infof("Fetched %d Crossplane Providers", len(providers))
	return &captenpluginspb.GetCrossplanProvidersResponse{
		Providers:     providers,
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}
