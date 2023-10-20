package agent

import (
	"context"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/common-pkg/plugins/argocd"
)

const (
	argoCDProjectAvailable           = "available"
	argoCDProjectConfigured          = "configured"
	argoCDProjectConfigurationFailed = "configuration-failed"
)

func (a *Agent) RegisterArgoCDProject(ctx context.Context, request *captenpluginspb.RegisterArgoCDProjectRequest) (
	*captenpluginspb.RegisterArgoCDProjectResponse, error) {
	if err := validateArgs(request.Id); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.RegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	a.log.Infof("Register ArgoCD Git project %s request recieved", request.Id)

	argoCDProject, err := a.as.GetArgoCDProjectForID(request.Id)
	if err != nil {
		a.log.Infof("faile to get argocd project %s, %v", request.Id, err)
		return &captenpluginspb.RegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	argocdClient, err := argocd.NewClient(&logging.Logging{})
	if err != nil {
		a.log.Errorf("failed to get ArgoCD client, %v ", err)
		return &captenpluginspb.RegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Error occured while getting ArgoCD client",
		}, err
	}

	// fetch Github project details with ID and fill those details to repository
	repo := &argocd.Repository{
		Project:       "Default",
		SSHPrivateKey: "",
		Type:          "git",
		Repo:          argoCDProject.GitProjectUrl,
	}

	_, err = argocdClient.CreateRepository(ctx, repo)
	if err != nil {
		a.log.Errorf("failed to configure git Project %s to argoCD, %v ", argoCDProject.GitProjectUrl, err)
		return &captenpluginspb.RegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Error occured while registering ArgoCD Repository",
		}, err
	}

	argoCDProject.Status = argoCDProjectConfigured
	if err := a.as.UpsertArgoCDProject(argoCDProject); err != nil {
		a.log.Errorf("failed to store argoCD git Project %s, %v ", argoCDProject.GitProjectUrl, err)
		return &captenpluginspb.RegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Error occured while adding ArgoCD project Data",
		}, err
	}

	a.log.Infof("ArgoCD Git project %s. %s Registered", request.Id, argoCDProject.GitProjectUrl)
	return &captenpluginspb.RegisterArgoCDProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "Sucessfully registered ArgoCD Repository",
	}, nil
}

func (a *Agent) UnRegisterArgoCDProject(ctx context.Context, request *captenpluginspb.UnRegisterArgoCDProjectRequest) (
	*captenpluginspb.UnRegisterArgoCDProjectResponse, error) {
	if err := validateArgs(request.Id); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.UnRegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	a.log.Infof("UnRegister ArgoCD Git project %s request recieved", request.Id)

	argoCDProject, err := a.as.GetArgoCDProjectForID(request.Id)
	if err != nil {
		a.log.Infof("faile to get argocd project %s, %v", request.Id, err)
		return &captenpluginspb.UnRegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	argocdClient, err := argocd.NewClient(&logging.Logging{})
	if err != nil {
		a.log.Errorf("failed to get ArgoCD client, %v ", err)
		return &captenpluginspb.UnRegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Error occured while getting ArgoCD client",
		}, err
	}
	_, err = argocdClient.DeleteRepository(ctx, argoCDProject.GitProjectUrl)
	if err != nil {
		a.log.Errorf("failed to delete ArgoCD Repository: %v ", err)
		return &captenpluginspb.UnRegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Error occured while deleting Repository",
		}, err
	}

	argoCDProject.Status = argoCDProjectAvailable
	if err := a.as.UpsertArgoCDProject(argoCDProject); err != nil {
		a.log.Errorf("failed to store argoCD git Project %s, %v ", argoCDProject.GitProjectUrl, err)
		return &captenpluginspb.UnRegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Error occured while adding ArgoCD project Data",
		}, err
	}

	a.log.Infof("ArgoCD Git project %s. %s UnRegistered", request.Id, argoCDProject.GitProjectUrl)
	return &captenpluginspb.UnRegisterArgoCDProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "Successfully unregisterted ArgoCD Repository",
	}, nil
}

func (a *Agent) GetArgoCDProjects(ctx context.Context, request *captenpluginspb.GetArgoCDProjectsRequest) (
	*captenpluginspb.GetArgoCDProjectsResponse, error) {
	a.log.Infof("Get ArgoCD Git projects request recieved")

	projects, err := a.as.GetArgoCDProjects()
	if err != nil {
		a.log.Errorf("failed to get argocd Project, %v", err)
		return &captenpluginspb.GetArgoCDProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get argocd Project",
		}, err
	}

	argocdProjects := []*captenpluginspb.ArgoCDProject{}
	for _, project := range projects {
		argocdProject := &captenpluginspb.ArgoCDProject{
			Id:         project.Id,
			ProjectUrl: project.GitProjectUrl,
			Status:     project.Status,
		}
		argocdProjects = append(argocdProjects, argocdProject)
	}

	a.log.Infof("Fetched %d ArgoCD Git projects", len(argocdProjects))
	return &captenpluginspb.GetArgoCDProjectsResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "Successfully fetched the Repositories",
		Projects:      argocdProjects,
	}, nil
}
