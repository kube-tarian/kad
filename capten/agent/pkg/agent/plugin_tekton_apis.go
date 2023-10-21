package agent

import (
	"context"
	"fmt"

	"github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
)

const (
	tekton     string = "tekton"
	gitProject string = "git-project"
)

func (a *Agent) RegisterTektonProject(ctx context.Context, request *captenpluginspb.RegisterTektonProjectRequest) (
	*captenpluginspb.RegisterTektonProjectResponse, error) {
	gitProject, err := a.as.GetGitProject(request.Id)
	if err != nil || len(gitProject) == 0 {
		return nil, fmt.Errorf("git project not found or error occured: %v", err)
	}

	regTekton := &model.RegisterTekton{Id: request.Id, ProjectUrl: gitProject[0].ProjectUrl, Status: "in-progress"}
	if err := a.as.AddTektonProject(regTekton); err != nil {
		a.log.Errorf("failed to Set Cluster Gitopts Project, %v", err)
		return &captenpluginspb.RegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "inserting data to tekton db got failed",
		}, err
	}

	// start the config-worker routine
	go a.configureGitRepo(regTekton, tekton)

	return &captenpluginspb.RegisterTektonProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successfully registered tekton",
	}, nil
}

func (a *Agent) UnRegisterTektonProject(ctx context.Context, request *captenpluginspb.UnRegisterTektonProjectRequest) (
	*captenpluginspb.UnRegisterTektonProjectResponse, error) {
	if err := a.as.DeleteTektonProject(request.Id); err != nil {
		return &captenpluginspb.UnRegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete the tekton project",
		}, fmt.Errorf("failed to delete the tekton project")
	}
	return &captenpluginspb.UnRegisterTektonProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successfully delete the tekton project",
	}, nil
}

func (a *Agent) GetTektonProjects(ctx context.Context, request *captenpluginspb.GetTektonProjectsRequest) (
	*captenpluginspb.GetTektonProjectsResponse, error) {
	projects, err := a.as.GetTektonProjects()
	if err != nil {
		a.log.Errorf("failed to get tekton Project, %v", err)
		return &captenpluginspb.GetTektonProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get tekton Project",
		}, err
	}

	return &captenpluginspb.GetTektonProjectsResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successfully fetched the tekton projects",
		Projects:      projects,
	}, nil
}
