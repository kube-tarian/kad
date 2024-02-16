package api

import (
	"context"
	"fmt"

	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/agent/internal/pb/captensdkpb"
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

func (a *Agent) AddPluginUsage(ctx context.Context, request *captensdkpb.AddPluginUsageRequest) (
	*captensdkpb.AddPluginUsageResponse, error) {

	if request.UsedPlugin == "" {
		a.log.Error("Used Plugin name is not provided")
		return &captensdkpb.AddPluginUsageResponse{
			Status:        captensdkpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "Used plugin name is is not provided",
		}, nil
	}

	if request.Id == "" {
		a.log.Error("Git repo Id is not provided")
		return &captensdkpb.AddPluginUsageResponse{
			Status:        captensdkpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "Git repo Id is not provided",
		}, nil
	}

	a.log.Infof("Adding Used plugin %s for git repo Id: %s", request.UsedPlugin, request.Id)

	res, err := a.as.GetGitProjectForID(request.Id)
	if err != nil {
		a.log.Errorf("failed to get git project from db, %v", err)
		return &captensdkpb.AddPluginUsageResponse{
			Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch git repo for " + request.Id,
		}, nil
	}

	gitProject := captenpluginspb.GitProject{
		Id:          res.Id,
		ProjectUrl:  res.ProjectUrl,
		Labels:      res.Labels,
		UsedPlugins: append(res.UsedPlugins, request.UsedPlugin),
	}
	if err := a.as.UpsertGitProject(&gitProject); err != nil {
		a.log.Errorf("failed to Upsert Git repo from db, %v", err)
		return &captensdkpb.AddPluginUsageResponse{
			Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to Upsert Git repo for " + request.Id,
		}, nil
	}

	a.log.Infof("Successfully adding Used plugin %s for git repo Id: %s", request.UsedPlugin, request.Id)

	return &captensdkpb.AddPluginUsageResponse{
		Status:        captensdkpb.StatusCode_OK,
		StatusMessage: fmt.Sprintf("Successfully adding Used plugin %s for git repo Id: %s", request.UsedPlugin, request.Id),
	}, nil
}

func (a *Agent) RemovePluginUsage(ctx context.Context, request *captensdkpb.RemovePluginUsageRequest) (
	*captensdkpb.RemovePluginUsageResponse, error) {

	if request.UsedPlugin == "" {
		a.log.Error("Used Plugin name is not provided")
		return &captensdkpb.RemovePluginUsageResponse{
			Status:        captensdkpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "Used plugin name is is not provided",
		}, nil
	}

	if request.Id == "" {
		a.log.Error("Git repo Id is not provided")
		return &captensdkpb.RemovePluginUsageResponse{
			Status:        captensdkpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "Git repo Id is not provided",
		}, nil
	}

	a.log.Infof("Removing Used plugin %s for git repo Id: %s", request.UsedPlugin, request.Id)

	res, err := a.as.GetGitProjectForID(request.Id)
	if err != nil {
		a.log.Errorf("failed to get Git project from db, %v", err)
		return &captensdkpb.RemovePluginUsageResponse{
			Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch Git repo for " + request.Id,
		}, nil
	}

	usedPlugins := []string{}
	for _, plugin := range res.UsedPlugins {
		if plugin != request.UsedPlugin {
			usedPlugins = append(usedPlugins, plugin)
		}
	}

	gitProject := captenpluginspb.GitProject{
		Id:          res.Id,
		ProjectUrl:  res.ProjectUrl,
		Labels:      res.Labels,
		UsedPlugins: usedPlugins,
	}
	if err := a.as.UpsertGitProject(&gitProject); err != nil {
		a.log.Errorf("failed to Upsert Git repo from db, %v", err)
		return &captensdkpb.RemovePluginUsageResponse{
			Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to Upsert Git repo for " + request.Id,
		}, nil
	}

	a.log.Infof("Successfully removed Used plugin %s for git repo Id: %s", request.UsedPlugin, request.Id)

	return &captensdkpb.RemovePluginUsageResponse{
		Status:        captensdkpb.StatusCode_OK,
		StatusMessage: fmt.Sprintf("Successfully removed Used plugin %s for git repo Id: %s", request.UsedPlugin, request.Id),
	}, nil
}
