package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/intelops/go-common/credentials"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/captenpluginspb"
)

const containerRegEntityName = "container-registry"

type DockerConfigEntry struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty" datapolicy:"password"`
	Email    string `json:"email,omitempty"`
	Auth     string `json:"auth,omitempty" datapolicy:"token"`
}

type DockerConfig map[string]DockerConfigEntry

type DockerConfigJSON struct {
	Auths DockerConfig `json:"auths" datapolicy:"token"`
	// +optional
	HttpHeaders map[string]string `json:"HttpHeaders,omitempty" datapolicy:"token"`
}

func (a *Agent) AddContainerRegistry(ctx context.Context, request *captenpluginspb.AddContainerRegistryRequest) (
	*captenpluginspb.AddContainerRegistryResponse, error) {
	if err := validateArgs(request.RegistryUrl, request.RegistryType); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.AddContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}

	a.log.Infof("Add Container registry %s request received", request.RegistryUrl)

	id := uuid.New()
	configData, err := parseAndPrepareDockerConfigJSONContent(request.RegistryAttributes, request.RegistryUrl)
	if err != nil {
		return &captenpluginspb.AddContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add ContainerRegistry credential in vault",
		}, err
	}
	request.RegistryAttributes["config.json"] = string(configData)

	if err := a.storeContainerRegCredential(ctx, id.String(), request.RegistryAttributes); err != nil {
		return &captenpluginspb.AddContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add Container registry credential in vault",
		}, err
	}

	ContainerRegistry := captenpluginspb.ContainerRegistry{
		Id:           id.String(),
		RegistryUrl:  request.RegistryUrl,
		Labels:       request.Labels,
		RegistryType: request.RegistryType,
	}
	if err := a.as.AddContainerRegistry(&ContainerRegistry); err != nil {
		a.log.Errorf("failed to store Container registry to DB, %v", err)
		return &captenpluginspb.AddContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add Container registry in db",
		}, err
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
		}, err
	}
	a.log.Infof("Update container registry project %s, %s request recieved", request.RegistryUrl, request.Id)

	id, err := uuid.Parse(request.Id)
	if err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.UpdateContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: fmt.Sprintf("invalid uuid: %s", request.Id),
		}, err
	}

	configData, err := parseAndPrepareDockerConfigJSONContent(request.RegistryAttributes, request.RegistryUrl)
	if err != nil {
		return &captenpluginspb.UpdateContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add ContainerRegistry credential in vault",
		}, err
	}
	request.RegistryAttributes["config.json"] = string(configData)

	if err := a.storeContainerRegCredential(ctx, request.Id, request.RegistryAttributes); err != nil {
		return &captenpluginspb.UpdateContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add ContainerRegistry credential in vault",
		}, err
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
		}, err
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
		}, err
	}
	a.log.Infof("Delete ContainerRegistry %s request recieved", request.Id)

	if err := a.deleteContainerRegCredential(ctx, request.Id); err != nil {
		return &captenpluginspb.DeleteContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete ContainerRegistry credential in vault",
		}, err
	}

	if err := a.as.DeleteContainerRegistryById(request.Id); err != nil {
		a.log.Errorf("failed to delete ContainerRegistry from db, %v", err)
		return &captenpluginspb.DeleteContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete ContainerRegistry from db",
		}, err
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
		}, err
	}

	for _, r := range res {
		cred, secretPath, secretKeys, err := a.getContainerRegCredential(ctx, r.Id)
		if err != nil {
			a.log.Errorf("failed to get credential, %v", err)
			return &captenpluginspb.GetContainerRegistryResponse{
				Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to fetch container registry",
			}, err
		}
		r.RegistryAttributes = cred
		r.SecretePath = secretPath
		r.SecreteKeys = secretKeys
	}

	a.log.Infof("Found %d container registry", len(res))
	return &captenpluginspb.GetContainerRegistryResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successful",
		Registries:    res,
	}, nil

}

func (a *Agent) getContainerRegCredential(ctx context.Context, id string) (map[string]string, string, []string, error) {
	credPath := fmt.Sprintf("%s/%s/%s", credentials.GenericCredentialType, containerRegEntityName, id)
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to intialize credentials client for %s", credPath)
		a.log.Errorf("failed to get crendential for %s, %v", credPath, err)
		return nil, "", nil, err
	}

	cred, err := credAdmin.GetCredential(ctx, credentials.GenericCredentialType, containerRegEntityName, id)
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

func parseAndPrepareDockerConfigJSONContent(credMap map[string]string, server string) ([]byte, error) {
	userName := credMap["username"]
	password := credMap["password"]
	return prepareDockerConfigJSONContent(userName, password, server)
}

func prepareDockerConfigJSONContent(username, password, server string) ([]byte, error) {
	dockerConfigAuth := DockerConfigEntry{
		Username: username,
		Password: password,
		Auth:     encodeDockerConfigFieldAuth(username, password),
	}
	dockerConfigJSON := DockerConfigJSON{
		Auths: map[string]DockerConfigEntry{server: dockerConfigAuth},
	}

	return json.Marshal(dockerConfigJSON)
}

func encodeDockerConfigFieldAuth(username, password string) string {
	fieldValue := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(fieldValue))
}
