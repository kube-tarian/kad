package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kube-tarian/kad/capten/agent/internal/pb/captensdkpb"
)

const (
	ResourceTypeGit               = "git"
	ResourceTypeCloudProvider     = "cloud_provider"
	ResourceTypeContainerRegistry = "container_registry"
)

func (a *Agent) GetGitProjectById(ctx context.Context, request *captensdkpb.GetGitProjectByIdRequest) (
	*captensdkpb.GetGitProjectByIdResponse, error) {

	if request.Id == "" {
		a.log.Error("Project Id is not provided")
		return &captensdkpb.GetGitProjectByIdResponse{
			Status:        captensdkpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "project Id is not provided",
		}, nil
	}

	a.log.Infof("Get Git project By Id request recieved for Id: %s", request.Id)

	res, err := a.as.GetGitProjectForID(request.Id)
	if err != nil {
		a.log.Errorf("failed to get gitProject from db for project Id: %s, %v", request.Id, err)
		return &captensdkpb.GetGitProjectByIdResponse{
			Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch git project for " + request.Id,
		}, nil
	}

	accessToken, _, err := a.getGitProjectCredential(ctx, res.Id)
	if err != nil {
		a.log.Errorf("failed to get git credential for project Id: %s, %v", request.Id, err)
		return &captensdkpb.GetGitProjectByIdResponse{
			Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch git project for " + request.Id,
		}, nil
	}

	project := &captensdkpb.GitProject{
		Id:             res.Id,
		ProjectUrl:     res.ProjectUrl,
		AccessToken:    accessToken,
		Labels:         res.Labels,
		LastUpdateTime: res.LastUpdateTime,
	}

	a.log.Infof("Fetched %s git project", res.Id)

	return &captensdkpb.GetGitProjectByIdResponse{
		Project:       project,
		Status:        captensdkpb.StatusCode_OK,
		StatusMessage: "successfully fetched git project for " + request.Id,
	}, nil
}

func (a *Agent) GetContainerRegistryById(ctx context.Context, request *captensdkpb.GetContainerRegistryByIdRequest) (
	*captensdkpb.GetContainerRegistryByIdResponse, error) {

	if request.Id == "" {
		a.log.Error("Container registry Id is not provided")
		return &captensdkpb.GetContainerRegistryByIdResponse{
			Status:        captensdkpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "container registry Id is not provided",
		}, nil
	}

	a.log.Infof("Get Container registry By Id request recieved for Id: %s", request.Id)

	res, err := a.as.GetContainerRegistryForID(request.Id)
	if err != nil {
		a.log.Errorf("failed to get ContainerRegistry from db, %v", err)
		return &captensdkpb.GetContainerRegistryByIdResponse{
			Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch container registry for " + request.Id,
		}, nil
	}

	registry := &captensdkpb.ContainerRegistry{
		Id:             res.Id,
		RegistryUrl:    res.RegistryUrl,
		Labels:         res.Labels,
		LastUpdateTime: res.LastUpdateTime,
		RegistryType:   res.RegistryType,
	}

	cred, err := a.getContainerRegCredential(ctx, res.Id)
	if err != nil {
		a.log.Errorf("failed to get container registry credential for %s, %v", request.Id, err)
		return &captensdkpb.GetContainerRegistryByIdResponse{
			Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch container registry for " + request.Id,
		}, nil
	}
	registry.RegistryAttributes = cred

	a.log.Infof("Fetched %s container registry", request.Id)
	return &captensdkpb.GetContainerRegistryByIdResponse{
		Registry:      registry,
		Status:        captensdkpb.StatusCode_OK,
		StatusMessage: "successfully fetched container registry for " + request.Id,
	}, nil

}

func (a *Agent) RegisterResourceUsage(ctx context.Context, request *captensdkpb.RegisterResourceUsageRequest) (
	*captensdkpb.RegisterResourceUsageResponse, error) {

	if request.ResourceId == "" {
		a.log.Error("Resouce Id is not provided")
		return &captensdkpb.RegisterResourceUsageResponse{
			Status:        captensdkpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "Resouce Id is not provided",
		}, nil
	}

	if request.ResourceType == "" {
		a.log.Error("Resouce type is not provided")
		return &captensdkpb.RegisterResourceUsageResponse{
			Status:        captensdkpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "Resouce type is not provided",
		}, nil
	}

	if request.PluginName == "" {
		a.log.Error("Plugin name is not provided")
		return &captensdkpb.RegisterResourceUsageResponse{
			Status:        captensdkpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "Plugin name is not provided",
		}, nil
	}

	if !(request.PluginName == "tekton" || request.PluginName == "crossplane") {
		a.log.Error("Invalid plugin name is provided")
		return &captensdkpb.RegisterResourceUsageResponse{
			Status:        captensdkpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "Invalid plugin name is provided",
		}, nil
	}

	a.log.Infof("Adding Plugin name %s of type %s for Id : %s", request.PluginName, request.ResourceType, request.ResourceId)

	switch request.ResourceType {
	case ResourceTypeGit:
		fmt.Println("git resouce type")
		res, err := a.as.GetGitProjectForID(request.ResourceId)
		if err != nil {
			a.log.Errorf("failed to get git project from db, %v", err)
			return &captensdkpb.RegisterResourceUsageResponse{
				Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to fetch git repo for " + request.ResourceId,
			}, nil
		}
		res.UsedPlugins = append(res.UsedPlugins, request.PluginName)
		x, _ := json.Marshal(res)
		fmt.Println(string(x))
		if err := a.as.UpsertGitProject(res); err != nil {
			a.log.Errorf("failed to Upsert Git repo from db, %v", err)
			return &captensdkpb.RegisterResourceUsageResponse{
				Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to Upsert Git repo for " + request.ResourceId,
			}, nil
		}
		fmt.Println("done")
	case ResourceTypeCloudProvider:
		res, err := a.as.GetCloudProviderForID(request.ResourceId)
		if err != nil {
			a.log.Errorf("failed to get Cloud provider from db, %v", err)
			return &captensdkpb.RegisterResourceUsageResponse{
				Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to fetch Cloud provider for " + request.ResourceId,
			}, nil
		}
		res.UsedPlugins = append(res.UsedPlugins, request.PluginName)
		if err := a.as.UpsertCloudProvider(res); err != nil {
			a.log.Errorf("failed to Upsert Cloud provider from db, %v", err)
			return &captensdkpb.RegisterResourceUsageResponse{
				Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to Upsert Cloud provider for " + request.ResourceId,
			}, nil
		}

	case ResourceTypeContainerRegistry:
		res, err := a.as.GetContainerRegistryForID(request.ResourceId)
		if err != nil {
			a.log.Errorf("failed to get Container registry from db, %v", err)
			return &captensdkpb.RegisterResourceUsageResponse{
				Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to fetch Container registry for " + request.ResourceId,
			}, nil
		}
		res.UsedPlugins = append(res.UsedPlugins, request.PluginName)
		if err := a.as.UpsertContainerRegistry(res); err != nil {
			a.log.Errorf("failed to Upsert Container registry from db, %v", err)
			return &captensdkpb.RegisterResourceUsageResponse{
				Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to Upsert Container registry for " + request.ResourceId,
			}, nil
		}
	default:
		return &captensdkpb.RegisterResourceUsageResponse{
			Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "Invalid Resouce type for " + request.ResourceId,
		}, nil
	}

	a.log.Infof("Successfully added Plugin name %s of type %s for Id : %s", request.PluginName, request.ResourceType, request.ResourceId)

	return &captensdkpb.RegisterResourceUsageResponse{
		Status:        captensdkpb.StatusCode_OK,
		StatusMessage: fmt.Sprintf("Successfully added Plugin name %s of type %s for Id : %s", request.PluginName, request.ResourceType, request.ResourceId),
	}, nil
}

func (a *Agent) DeRegisterResourceUsage(ctx context.Context, request *captensdkpb.DeRegisterResourceUsageRequest) (
	*captensdkpb.DeRegisterResourceUsageResponse, error) {

	if request.ResourceId == "" {
		a.log.Error("Resouce Id is not provided")
		return &captensdkpb.DeRegisterResourceUsageResponse{
			Status:        captensdkpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "Resouce Id is not provided",
		}, nil
	}

	if request.ResourceType == "" {
		a.log.Error("Resouce type is not provided")
		return &captensdkpb.DeRegisterResourceUsageResponse{
			Status:        captensdkpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "Resouce type is not provided",
		}, nil
	}

	if request.PluginName == "" {
		a.log.Error("Plugin name is not provided")
		return &captensdkpb.DeRegisterResourceUsageResponse{
			Status:        captensdkpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "Plugin name is not provided",
		}, nil
	}

	if !(request.PluginName == "tekton" || request.PluginName == "crossplane") {
		a.log.Error("Invalid plugin name is provided")
		return &captensdkpb.DeRegisterResourceUsageResponse{
			Status:        captensdkpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "Invalid plugin name is provided",
		}, nil
	}

	a.log.Infof("Removing Plugin name %s of type %s for Id : %s", request.PluginName, request.ResourceType, request.ResourceId)

	switch request.ResourceType {
	case ResourceTypeGit:
		res, err := a.as.GetGitProjectForID(request.ResourceId)
		if err != nil {
			a.log.Errorf("failed to get Git repo from db, %v", err)
			return &captensdkpb.DeRegisterResourceUsageResponse{
				Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to fetch Git repo for " + request.ResourceId,
			}, nil
		}
		res.UsedPlugins = removeUsedPlugin(request.PluginName, res.UsedPlugins)
		if err := a.as.UpsertGitProject(res); err != nil {
			a.log.Errorf("failed to Upsert Git repo from db, %v", err)
			return &captensdkpb.DeRegisterResourceUsageResponse{
				Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to Upsert Git repo for " + request.ResourceId,
			}, nil
		}
	case ResourceTypeCloudProvider:
		res, err := a.as.GetCloudProviderForID(request.ResourceId)
		if err != nil {
			a.log.Errorf("failed to get Cloud provider from db, %v", err)
			return &captensdkpb.DeRegisterResourceUsageResponse{
				Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to fetch Cloud provider for " + request.ResourceId,
			}, nil
		}
		res.UsedPlugins = removeUsedPlugin(request.PluginName, res.UsedPlugins)
		if err := a.as.UpsertCloudProvider(res); err != nil {
			a.log.Errorf("failed to Upsert Cloud provider from db, %v", err)
			return &captensdkpb.DeRegisterResourceUsageResponse{
				Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to Upsert Cloud provider for " + request.ResourceId,
			}, nil
		}

	case ResourceTypeContainerRegistry:
		res, err := a.as.GetContainerRegistryForID(request.ResourceId)
		if err != nil {
			a.log.Errorf("failed to get Container registry from db, %v", err)
			return &captensdkpb.DeRegisterResourceUsageResponse{
				Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to fetch Container registry for " + request.ResourceId,
			}, nil
		}
		res.UsedPlugins = removeUsedPlugin(request.PluginName, res.UsedPlugins)
		if err := a.as.UpsertContainerRegistry(res); err != nil {
			a.log.Errorf("failed to Upsert Container registry from db, %v", err)
			return &captensdkpb.DeRegisterResourceUsageResponse{
				Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to Upsert Container registry for " + request.ResourceId,
			}, nil
		}
	default:
		return &captensdkpb.DeRegisterResourceUsageResponse{
			Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "Invalid Resouce type for " + request.ResourceId,
		}, nil
	}

	a.log.Infof("Successfully removed Plugin name %s of type %s for Id : %s", request.PluginName, request.ResourceType, request.ResourceId)

	return &captensdkpb.DeRegisterResourceUsageResponse{
		Status:        captensdkpb.StatusCode_OK,
		StatusMessage: fmt.Sprintf("Successfully removed Plugin name %s of type %s for Id : %s", request.PluginName, request.ResourceType, request.ResourceId),
	}, nil
}

func removeUsedPlugin(pluginName string, usedPlugins []string) []string {
	plugins := []string{}
	for _, v := range usedPlugins {
		if v != pluginName {
			plugins = append(plugins, v)
		}
	}

	return plugins
}
