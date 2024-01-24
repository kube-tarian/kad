package api

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/intelops/go-common/credentials"
	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"
)

const containerRegEntityName = "container-registry"

func (a *Agent) AddContainerRegistry(ctx context.Context, request *captenpluginspb.AddContainerRegistryRequest) (
	*captenpluginspb.AddContainerRegistryResponse, error) {
	if err := validateArgs(request.RegistryUrl, request.RegistryType); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.AddContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	a.log.Infof("Add Container registry %s request received", request.RegistryUrl)

	id := uuid.New()

	if err := a.storeContainerRegCredential(ctx, id.String(), request.RegistryAttributes); err != nil {
		return &captenpluginspb.AddContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add Container registry credential in vault",
		}, nil
	}

	ContainerRegistry := captenpluginspb.ContainerRegistry{
		Id:           id.String(),
		RegistryUrl:  request.RegistryUrl,
		Labels:       request.Labels,
		RegistryType: request.RegistryType,
	}
	if err := a.as.UpsertContainerRegistry(&ContainerRegistry); err != nil {
		a.log.Errorf("failed to store Container registry to DB, %v", err)
		return &captenpluginspb.AddContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add Container registry in db",
		}, nil
	}

	a.log.Infof("Container registry %s added with id %s", request.RegistryUrl, id.String())
	return &captenpluginspb.AddContainerRegistryResponse{
		Id:            id.String(),
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) UpdateContainerRegistry(ctx context.Context, request *captenpluginspb.UpdateContainerRegistryRequest) (
	*captenpluginspb.UpdateContainerRegistryResponse, error) {
	if err := validateArgs(request.RegistryUrl); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.UpdateContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	a.log.Infof("Update container registry project %s, %s request recieved", request.RegistryUrl, request.Id)

	id, err := uuid.Parse(request.Id)
	if err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.UpdateContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: fmt.Sprintf("invalid uuid: %s", request.Id),
		}, nil
	}

	if err := a.storeContainerRegCredential(ctx, request.Id, request.RegistryAttributes); err != nil {
		return &captenpluginspb.UpdateContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add ContainerRegistry credential in vault",
		}, nil
	}

	ContainerRegistry := captenpluginspb.ContainerRegistry{
		Id:           id.String(),
		RegistryUrl:  request.RegistryUrl,
		Labels:       request.Labels,
		RegistryType: request.RegistryType,
	}

	if err := a.as.UpsertContainerRegistry(&ContainerRegistry); err != nil {
		a.log.Errorf("failed to update ContainerRegistry in db, %v", err)
		return &captenpluginspb.UpdateContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to update ContainerRegistry in db",
		}, nil
	}

	a.log.Infof("ContainerRegistry %s, %s updated", request.RegistryUrl, request.Id)
	return &captenpluginspb.UpdateContainerRegistryResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) DeleteContainerRegistry(ctx context.Context, request *captenpluginspb.DeleteContainerRegistryRequest) (
	*captenpluginspb.DeleteContainerRegistryResponse, error) {
	if err := validateArgs(request.Id); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.DeleteContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	a.log.Infof("Delete ContainerRegistry %s request recieved", request.Id)

	if err := a.deleteContainerRegCredential(ctx, request.Id); err != nil {
		return &captenpluginspb.DeleteContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete ContainerRegistry credential in vault",
		}, nil
	}

	if err := a.as.DeleteContainerRegistryById(request.Id); err != nil {
		a.log.Errorf("failed to delete ContainerRegistry from db, %v", err)
		return &captenpluginspb.DeleteContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete ContainerRegistry from db",
		}, nil
	}

	a.log.Infof("ContainerRegistry %s deleted", request.Id)
	return &captenpluginspb.DeleteContainerRegistryResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) GetContainerRegistry(ctx context.Context, request *captenpluginspb.GetContainerRegistryRequest) (
	*captenpluginspb.GetContainerRegistryResponse, error) {
	a.log.Infof("Get Git projects request recieved")
	res, err := a.as.GetContainerRegistries()
	if err != nil {
		a.log.Errorf("failed to get ContainerRegistry from db, %v", err)
		return &captenpluginspb.GetContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch git projects",
		}, nil
	}

	for _, r := range res {
		cred, err := a.getContainerRegCredential(ctx, r.Id)
		if err != nil {
			a.log.Errorf("failed to get credential, %v", err)
			return &captenpluginspb.GetContainerRegistryResponse{
				Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to fetch container registry",
			}, nil
		}

		r.RegistryAttributes = cred
	}

	a.log.Infof("Found %d container registry", len(res))
	return &captenpluginspb.GetContainerRegistryResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successful",
		Registries:    res,
	}, nil

}

func (a *Agent) getContainerRegCredential(ctx context.Context, id string) (map[string]string, error) {
	credPath := fmt.Sprintf("%s/%s/%s", credentials.GenericCredentialType, containerRegEntityName, id)
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to intialize credentials client for %s", credPath)
		a.log.Errorf("failed to get crendential for %s, %v", credPath, err)
		return nil, err
	}

	cred, err := credAdmin.GetCredential(ctx, credentials.GenericCredentialType, containerRegEntityName, id)
	if err != nil {
		a.log.Errorf("failed to get credential for %s, %v", credPath, err)
		return nil, err
	}
	return cred, nil
}

func (a *Agent) storeContainerRegCredential(ctx context.Context, id string, credentialMap map[string]string) error {
	credPath := fmt.Sprintf("%s/%s/%s", credentials.GenericCredentialType, containerRegEntityName, id)
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to intialize credentials client for %s", credPath)
		a.log.Errorf("failed to store credential for %s, %v", credPath, err)
		return err
	}

	err = credAdmin.PutCredential(ctx, credentials.GenericCredentialType, containerRegEntityName,
		id, credentialMap)

	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to store crendential for %s", credPath)
		a.log.Errorf("failed to store credential for %s, %v", credPath, err)
		return err
	}
	a.log.Audit("security", "storecred", "success", "system", "credential stored for %s", credPath)
	a.log.Infof("stored credential for entity %s", credPath)
	return nil
}

func (a *Agent) deleteContainerRegCredential(ctx context.Context, id string) error {
	credPath := fmt.Sprintf("%s/%s/%s", credentials.GenericCredentialType, containerRegEntityName, id)
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to intialize credentials client for %s", credPath)
		a.log.Errorf("failed to delete credential for %s, %v", credPath, err)
		return err
	}

	err = credAdmin.DeleteCredential(ctx, credentials.GenericCredentialType, containerRegEntityName, id)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to store crendential for %s", credPath)
		a.log.Errorf("failed to delete credential for %s, %v", credPath, err)
		return err
	}
	a.log.Audit("security", "storecred", "success", "system", "credential stored for %s", credPath)
	a.log.Infof("deleted credential for entity %s", credPath)
	return nil
}
