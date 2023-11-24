package api

import (
	"context"
	"fmt"

	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/agent/internal/capten-store"
	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/common-pkg/plugins/argocd"
	"github.com/kube-tarian/kad/capten/model"
)

const (
	argoCDRepositoryType    string = "git"
	argoCDRepositoryProject string = "Default"
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
		a.log.Infof("failed to get argocd project %s, %v", request.Id, err)
		return &captenpluginspb.RegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	accessToken, userID, err := a.getGitProjectCredential(ctx, argoCDProject.GitProjectId)
	if err != nil {
		a.log.Errorf("failed to get credential, %v", err)
		return &captenpluginspb.RegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "Error occured while adding Repository",
		}, nil
	}

	if err := a.addProjectToArgoCD(ctx, argoCDProject.GitProjectUrl, userID, accessToken); err != nil {
		a.log.Errorf("failed to add ArgoCD Repository: %v ", err)
		return &captenpluginspb.RegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Error occured while adding Repository",
		}, err
	}

	argoCDProject.Status = string(model.ArgoCDProjectConfigured)
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
		if !captenstore.IsObjectNotFound(err) {
			a.log.Infof("faile to get argocd project %s, %v", request.Id, err)
			return &captenpluginspb.UnRegisterArgoCDProjectResponse{
				Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
				StatusMessage: "request validation failed",
			}, nil
		}
	}

	if err := a.deleteProjectFromArgoCD(ctx, argoCDProject.GitProjectUrl); err != nil {
		a.log.Errorf("failed to delete ArgoCD Repository: %v ", err)
		return &captenpluginspb.UnRegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Error occured while deleting Repository",
		}, err
	}

	argoCDProject.Status = string(model.ArgoCDProjectAvailable)
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

	argocdClient, err := argocd.NewClient(&logging.Logging{})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Started fetching kubeconfig")
	err = argocdClient.GetConfig(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("finished fetching kubeconfig")

	argocdProjects := []*captenpluginspb.ArgoCDProject{}
	for _, project := range projects {
		argocdProject := &captenpluginspb.ArgoCDProject{
			Id:             project.Id,
			ProjectUrl:     project.GitProjectUrl,
			Status:         project.Status,
			LastUpdateTime: project.LastUpdateTime,
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

func (a *Agent) addProjectToArgoCD(ctx context.Context, projectUrl, userID, accessToken string) error {
	argocdClient, err := argocd.NewClient(&logging.Logging{})
	if err != nil {
		return err
	}

	repo := &argocd.Repository{
		Project:               argoCDRepositoryProject,
		Repo:                  projectUrl,
		Username:              userID,
		Password:              accessToken,
		Type:                  argoCDRepositoryType,
		Insecure:              false,
		EnableLFS:             false,
		InsecureIgnoreHostKey: false,
		ConnectionState: argocd.ConnectionState{
			Status:  "Connected",
			Message: "Repository is connected",
		},
	}

	_, err = argocdClient.CreateRepository(ctx, repo)
	if err != nil {
		return err
	}
	return nil
}

func (a *Agent) deleteProjectFromArgoCD(ctx context.Context, projectUrl string) error {
	argocdClient, err := argocd.NewClient(&logging.Logging{})
	if err != nil {
		return err
	}
	_, err = argocdClient.DeleteRepository(ctx, projectUrl)
	if err != nil {
		return err
	}
	return nil
}

func (a *Agent) isProjectRegisteredWithArgoCD(ctx context.Context, projectUrl string) (bool, error) {
	argocdClient, err := argocd.NewClient(&logging.Logging{})
	if err != nil {
		return false, err
	}
	_, err = argocdClient.GetRepository(ctx, projectUrl)
	if err != nil && fmt.Sprintf("rpc error: code = NotFound desc = repo '%s' not found", projectUrl) == err.Error() {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
