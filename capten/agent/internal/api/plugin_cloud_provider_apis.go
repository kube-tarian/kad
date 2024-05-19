package api

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/intelops/go-common/credentials"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/captenpluginspb"
)

const cloudProviderEntityName = "cloud-provider"

func (a *Agent) AddCloudProvider(ctx context.Context, request *captenpluginspb.AddCloudProviderRequest) (
	*captenpluginspb.AddCloudProviderResponse, error) {
	if err := validateArgs(request.GetCloudType(), request.GetCloudAttributes()); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.AddCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}

	a.log.Infof("Add Cloud Provider %s request received", request.CloudType)

	id := uuid.New()
	if err := a.storeCloudProviderCredential(ctx, id.String(), request.GetCloudAttributes()); err != nil {
		return &captenpluginspb.AddCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add cloud provider credential in vault",
		}, err
	}

	CloudProvider := captenpluginspb.CloudProvider{
		Id:        id.String(),
		CloudType: request.CloudType,
		Labels:    request.Labels,
	}
	if err := a.as.UpsertCloudProvider(&CloudProvider); err != nil {
		a.log.Errorf("failed to store cloud provider to DB, %v", err)
		return &captenpluginspb.AddCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add CloudProvider in db",
		}, err
	}

	a.log.Infof("Cloud Provider %s added with id %s", request.GetCloudType(), id.String())
	return &captenpluginspb.AddCloudProviderResponse{
		Id:            id.String(),
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) UpdateCloudProvider(ctx context.Context, request *captenpluginspb.UpdateCloudProviderRequest) (
	*captenpluginspb.UpdateCloudProviderResponse, error) {
	if err := validateArgs(request.GetCloudType(), request.GetId(), request.GetCloudAttributes()); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.UpdateCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}
	a.log.Infof("Update Cloud Provider %s, %s request received", request.CloudType, request.Id)

	id, err := uuid.Parse(request.Id)
	if err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.UpdateCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: fmt.Sprintf("invalid uuid: %s", request.Id),
		}, err
	}

	if err := a.storeCloudProviderCredential(ctx, request.Id, request.GetCloudAttributes()); err != nil {
		return &captenpluginspb.UpdateCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add CloudProvider credential in vault",
		}, err
	}

	CloudProvider := captenpluginspb.CloudProvider{
		Id:        id.String(),
		CloudType: request.CloudType,
		Labels:    request.Labels,
	}
	if err := a.as.UpsertCloudProvider(&CloudProvider); err != nil {
		a.log.Errorf("failed to update CloudProvider in db, %v", err)
		return &captenpluginspb.UpdateCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to update CloudProvider in db",
		}, err
	}

	a.log.Infof("Cloud Provider %s, %s updated", request.CloudType, request.Id)
	return &captenpluginspb.UpdateCloudProviderResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) DeleteCloudProvider(ctx context.Context, request *captenpluginspb.DeleteCloudProviderRequest) (
	*captenpluginspb.DeleteCloudProviderResponse, error) {
	if err := validateArgs(request.GetId()); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.DeleteCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}
	a.log.Infof("Delete Cloud Provider %s request recieved", request.Id)

	if err := a.deleteCloudProviderCredential(ctx, request.Id); err != nil {
		return &captenpluginspb.DeleteCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete cloud provider credential in vault",
		}, err
	}

	if err := a.as.DeleteCloudProviderById(request.Id); err != nil {
		a.log.Errorf("failed to delete CloudProvider from db, %v", err)
		return &captenpluginspb.DeleteCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete CloudProvider from db",
		}, err
	}

	a.log.Infof("Cloud Provider %s deleted", request.Id)
	return &captenpluginspb.DeleteCloudProviderResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) GetCloudProviders(ctx context.Context, request *captenpluginspb.GetCloudProvidersRequest) (
	*captenpluginspb.GetCloudProvidersResponse, error) {
	a.log.Infof("Get Cloud Provider request recieved")
	res, err := a.as.GetCloudProviders()
	if err != nil {
		a.log.Errorf("failed to get CloudProviders from db, %v", err)
		return &captenpluginspb.GetCloudProvidersResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch cloud providers",
		}, err
	}

	for _, r := range res {
		cloudAttributes, secretPath, secretKeys, err := a.getCloudProviderCredential(ctx, r.Id)
		if err != nil {
			a.log.Errorf("failed to get credential, %v", err)
			return &captenpluginspb.GetCloudProvidersResponse{
				Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to fetch cloud providers",
			}, err
		}
		r.CloudAttributes = cloudAttributes
		r.SecretePath = secretPath
		r.SecreteKeys = secretKeys
	}

	a.log.Infof("Found %d cloud providers", len(res))
	return &captenpluginspb.GetCloudProvidersResponse{
		Status:         captenpluginspb.StatusCode_OK,
		StatusMessage:  "successful",
		CloudProviders: res,
	}, nil

}

func (a *Agent) GetCloudProvidersWithFilter(ctx context.Context, request *captenpluginspb.GetCloudProvidersWithFilterRequest) (
	*captenpluginspb.GetCloudProvidersWithFilterResponse, error) {
	if len(request.GetLabels()) == 0 {
		a.log.Infof("request validation failed")
		return &captenpluginspb.GetCloudProvidersWithFilterResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, fmt.Errorf("labels cannot be empty")
	}
	a.log.Infof("Get Cloud providers with labels %v request recieved", request.Labels)

	res, err := a.as.GetCloudProvidersByLabelsAndCloudType(request.Labels, request.CloudType)
	if err != nil {
		a.log.Errorf("failed to get CloudProviders for labels from db, %v", err)
		return &captenpluginspb.GetCloudProvidersWithFilterResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch cloud providers",
		}, err
	}

	for _, r := range res {
		cloudAttributes, secretPath, secretKeys, err := a.getCloudProviderCredential(ctx, r.Id)
		if err != nil {
			a.log.Errorf("failed to get credential, %v", err)
			return &captenpluginspb.GetCloudProvidersWithFilterResponse{
				Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to fetch cloud providers",
			}, err
		}
		r.CloudAttributes = cloudAttributes
		r.SecretePath = secretPath
		r.SecreteKeys = secretKeys
	}

	a.log.Infof("Found %d cloud providers for lables %v and cloud type %v", len(res), request.Labels, request.CloudType)
	return &captenpluginspb.GetCloudProvidersWithFilterResponse{
		Status:         captenpluginspb.StatusCode_OK,
		StatusMessage:  "successful",
		CloudProviders: res,
	}, nil
}

func (a *Agent) getCloudProviderCredential(ctx context.Context, id string) (map[string]string, string, []string, error) {
	credPath := fmt.Sprintf("%s/%s/%s", credentials.GenericCredentialType, cloudProviderEntityName, id)
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to intialize credentials client for %s", credPath)
		a.log.Errorf("failed to get crendential for %s, %v", credPath, err)
		return nil, "", nil, err
	}

	cred, err := credAdmin.GetCredential(ctx, credentials.GenericCredentialType, cloudProviderEntityName, id)
	if err != nil {
		a.log.Errorf("failed to get credential for %s, %v", credPath, err)
		return nil, "", nil, err
	}

	secretKeys := []string{}
	for key := range cred {
		secretKeys = append(secretKeys, key)
	}
	return cred, credPath, secretKeys, nil
}

func (a *Agent) storeCloudProviderCredential(ctx context.Context, id string, credentialMap map[string]string) error {
	credPath := fmt.Sprintf("%s/%s/%s", credentials.GenericCredentialType, cloudProviderEntityName, id)
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to intialize credentials client for %s", credPath)
		a.log.Errorf("failed to store credential for %s, %v", credPath, err)
		return err
	}

	err = credAdmin.PutCredential(ctx, credentials.GenericCredentialType, cloudProviderEntityName,
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

func (a *Agent) deleteCloudProviderCredential(ctx context.Context, id string) error {
	credPath := fmt.Sprintf("%s/%s/%s", credentials.GenericCredentialType, cloudProviderEntityName, id)
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to intialize credentials client for %s", credPath)
		a.log.Errorf("failed to delete credential for %s, %v", credPath, err)
		return err
	}

	err = credAdmin.DeleteCredential(ctx, credentials.GenericCredentialType, cloudProviderEntityName, id)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to store crendential for %s", credPath)
		a.log.Errorf("failed to delete credential for %s, %v", credPath, err)
		return err
	}
	a.log.Audit("security", "storecred", "success", "system", "credential stored for %s", credPath)
	a.log.Infof("deleted credential for entity %s", credPath)
	return nil
}
