package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/model"
)

const (
	objectNotFoundErrorMessage = "object not found"
)

func (a *Agent) AddCrossplanProvider(ctx context.Context, request *captenpluginspb.AddCrossplanProviderRequest) (
	*captenpluginspb.AddCrossplanProviderResponse, error) {
	if err := validateArgs(request.CloudType, request.CloudProviderId); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.AddCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	a.log.Infof("Add Crossplane Provider type %s with cloud provider %s request recieved", request.CloudType, request.CloudProviderId)

	project, err := a.as.GetCrossplanProviderByCloudType(request.CloudType)
	if err != nil {
		a.log.Infof("failed to get crossplane provider", err)
		return &captenpluginspb.AddCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get crossplane provider for " + request.CloudType,
		}, nil
	}
	if project != nil {
		return &captenpluginspb.AddCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Crossplane provider is already available",
		}, nil
	}

	id := uuid.New()
	provider := model.CrossplaneProvider{
		Id:              id.String(),
		CloudType:       request.CloudType,
		ProviderName:    model.PrepareCrossplaneProviderName(request.CloudType),
		CloudProviderId: request.CloudProviderId,
		Status:          string(model.CrossPlaneProviderOutofSynch),
	}

	if err := a.as.UpsertCrossplaneProvider(&provider); err != nil {
		a.log.Errorf("failed to store crossplane provider to DB, %v", err)
		return &captenpluginspb.AddCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add crossplane provider in db",
		}, nil
	}

	a.log.Infof("Crossplane Provider type %s added with id %s", request.CloudType, id.String())
	return &captenpluginspb.AddCrossplanProviderResponse{
		Id:            id.String(),
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) DeleteCrossplanProvider(ctx context.Context, request *captenpluginspb.DeleteCrossplanProviderRequest) (
	*captenpluginspb.DeleteCrossplanProviderResponse, error) {

	if err := validateArgs(request.Id); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.DeleteCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	a.log.Infof("Delete Crossplane Provider %s request recieved", request.Id)

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
	a.log.Infof("Get Crossplane Providers request received")

	providers, err := a.as.GetCrossplaneProviders()
	if err != nil {
		if err.Error() == objectNotFoundErrorMessage {
			a.log.Info("No crossplane providers in DB")
			return &captenpluginspb.GetCrossplanProvidersResponse{
				Status:        captenpluginspb.StatusCode_NOT_FOUND,
				StatusMessage: "No crossplane providers found",
			}, nil
		}
		a.log.Errorf("failed to fetch crossplane providers from DB, %v", err)
		return &captenpluginspb.GetCrossplanProvidersResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch crossplane providers from db",
		}, nil
	}

	a.log.Infof("Fetched %d Crossplane Providers", len(providers))
	return &captenpluginspb.GetCrossplanProvidersResponse{
		Providers:     providers,
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) UpdateCrossplanProvider(ctx context.Context, request *captenpluginspb.UpdateCrossplanProviderRequest) (
	*captenpluginspb.UpdateCrossplanProviderResponse, error) {

	if err := validateArgs(request.Id, request.CloudType, request.CloudProviderId); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.UpdateCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	a.log.Infof("Update Crossplane Provider %s, %s, %s request recieved", request.CloudType, request.Id, request.CloudProviderId)

	project, err := a.as.GetCrossplanProviderById(request.Id)
	if err != nil {
		a.log.Infof("failed to get crossplane provider for "+request.Id, err)
		return &captenpluginspb.UpdateCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get crossplane provider for " + request.Id,
		}, nil
	} else if project == nil {
		return &captenpluginspb.UpdateCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Crossplane provider is not available for" + request.Id,
		}, nil
	} else if project.CloudType != request.CloudType {
		return &captenpluginspb.UpdateCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "Crossplane provider cloud type change is not supported for" + request.Id,
		}, nil
	}

	provider := model.CrossplaneProvider{
		Id:              request.Id,
		CloudType:       request.CloudType,
		ProviderName:    model.PrepareCrossplaneProviderName(request.CloudType),
		CloudProviderId: request.CloudProviderId,
		Status:          string(model.CrossPlaneProviderOutofSynch),
	}

	if err := a.as.UpdateCrossplaneProvider(&provider); err != nil {
		a.log.Errorf("failed to update crossplane provider in DB, %v", err)
		return &captenpluginspb.UpdateCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to update crossplane provider in db",
		}, nil
	}

	a.log.Infof("Crossplane Provider type %s with id %s updated", request.CloudType, request.Id)
	return &captenpluginspb.UpdateCrossplanProviderResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}
