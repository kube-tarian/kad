package agent

import (
	"context"
	"fmt"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/common-pkg/plugins/argocd"
)

func (a *Agent) RegisterArgoCDProject(ctx context.Context, request *captenpluginspb.RegisterArgoCDProjectRequest) (
	*captenpluginspb.RegisterArgoCDProjectResponse, error) {

	if request.Id == "" {
		return &captenpluginspb.RegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Id is required",
		}, fmt.Errorf("Id is required")
	}

	argocdClient, err := argocd.NewClient(&logging.Logging{})
	if err != nil {
		a.log.Errorf("failed to get ArgoCD client: %v ", err)
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
		Repo:          "",
	}

	resp, err := argocdClient.CreateRepository(ctx, repo)
	if err != nil {
		a.log.Errorf("failed to register ArgoCD Repository: %v ", err)
		return &captenpluginspb.RegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Error occured while registering ArgoCD Repository",
		}, err
	}

	if err := a.as.AddArgoCDProjectsData(request.Id, resp.ConnectionState.Status); err != nil {
		a.log.Errorf("failed to add ArgoCD Project Data: %v ", err)
		return &captenpluginspb.RegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Error occured while adding ArgoCD project Data",
		}, err

	}

	return &captenpluginspb.RegisterArgoCDProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "Sucessfully registered ArgoCD Repository",
	}, nil
}

func (a *Agent) UnRegisterArgoCDProject(ctx context.Context, request *captenpluginspb.UnRegisterArgoCDProjectRequest) (
	*captenpluginspb.UnRegisterArgoCDProjectResponse, error) {

	if request.Id == "" {
		return &captenpluginspb.UnRegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Id is required",
		}, fmt.Errorf("Id is required")
	}

	argocdClient, err := argocd.NewClient(&logging.Logging{})
	if err != nil {
		a.log.Errorf("failed to get ArgoCD client: %v ", err)
		return &captenpluginspb.UnRegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Error occured while getting ArgoCD client",
		}, err
	}

	_, err = argocdClient.DeleteRepository(ctx, request.Id)
	if err != nil {
		a.log.Errorf("failed to delete ArgoCD Repository: %v ", err)
		return &captenpluginspb.UnRegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Error occured while deleting Repository",
		}, err
	}

	if err := a.as.DeleteArgoCDProjectsData(request.Id); err != nil {
		a.log.Errorf("failed to delete ArgoCD Project Data: %v ", err)
		return &captenpluginspb.UnRegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Error occured while deleting ArgoCD project Data",
		}, err

	}

	return &captenpluginspb.UnRegisterArgoCDProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "Successfully deleted ArgoCD Repository",
	}, nil
}

func (a *Agent) GetArgoCDProjects(ctx context.Context, request *captenpluginspb.GetArgoCDProjectsRequest) (
	*captenpluginspb.GetArgoCDProjectsResponse, error) {

	argocdClient, err := argocd.NewClient(&logging.Logging{})
	if err != nil {
		a.log.Errorf("failed to get ArgoCD client: %v ", err)
		return &captenpluginspb.GetArgoCDProjectsResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Error occured while getting ArgoCD client",
		}, err
	}

	list, err := argocdClient.ListRepositories(ctx)
	if err != nil {
		a.log.Errorf("failed to get Repository list: %v ", err)
		return &captenpluginspb.GetArgoCDProjectsResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Error occured while fetching Repositories",
		}, err
	}

	projects := []*captenpluginspb.ArgoCDProject{}
	for _, v := range list.Items {
		projects = append(projects, &captenpluginspb.ArgoCDProject{
			ProjectUrl: v.Repo,
			Status:     v.ConnectionState.Status,
		})
	}

	return &captenpluginspb.GetArgoCDProjectsResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "Successfully fetched the Repositories",
		Projects:      projects,
	}, nil
}
