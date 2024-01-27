package api

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/agent/internal/workers"
	"github.com/kube-tarian/kad/capten/model"
	captenmodel "github.com/kube-tarian/kad/capten/model"
)

const (
	tektonConfigUseCase string = "tekton"
)

func (a *Agent) RegisterTektonProject(ctx context.Context, request *captenpluginspb.RegisterTektonProjectRequest) (
	*captenpluginspb.RegisterTektonProjectResponse, error) {
	a.log.Infof("registering the tekton project")
	if err := validateArgs(request.Id); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.RegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}

	res, err := a.as.GetTektonPipeliness()
	if err != nil {
		a.log.Errorf("failed to get TektonPipeline from db, %v", err)
		return &captenpluginspb.RegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch TektonPipelines",
		}, err
	}

	pipelines := make([]*captenpluginspb.TektonPipelines, len(res))
	for index, r := range res {
		pipelines[index] = &captenpluginspb.TektonPipelines{
			Id: r.Id, PipelineName: r.PipelineName,
			WebhookURL: r.WebhookURL, Status: r.Status, GitOrgId: r.GitProjectId,
			ContainerRegistryIds: r.ContainerRegId, LastUpdateTime: r.LastUpdateTime,
		}
		if _, err := a.configureTektonPipelinesGitRepo(r, model.TektonPipelineSync, true); err != nil {
			a.log.Errorf("failed to trigger the tekton pipeline sync for %s, %v", r.PipelineName, err)
			pipelines[index].Status = "failed to trigger the sync"
			continue
		}
	}

	a.log.Infof("Successfully registered the project")
	return &captenpluginspb.RegisterTektonProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "Successfully registered the project",
	}, nil
}

func (a *Agent) UnRegisterTektonProject(ctx context.Context, request *captenpluginspb.UnRegisterTektonProjectRequest) (
	*captenpluginspb.UnRegisterTektonProjectResponse, error) {
	if err := validateArgs(request.Id); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.UnRegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}
	a.log.Infof("UnRegister Tekton Git project %s request recieved", request.Id)

	tektonProject, err := a.as.GetTektonProjectForID(request.Id)
	if err != nil {
		a.log.Infof("faile to get git project %s, %v", request.Id, err)
		return &captenpluginspb.UnRegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	tektonProject.Status = string(model.TektonProjectAvailable)
	if err := a.as.UpsertTektonProject(tektonProject); err != nil {
		a.log.Errorf("failed to Set Cluster Gitopts Project, %v", err)
		return &captenpluginspb.UnRegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "inserting data to tekton db got failed",
		}, err
	}

	a.log.Infof("UnRegister Tekton Git project %s request processed", request.Id)
	return &captenpluginspb.UnRegisterTektonProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successfully delete the tekton project",
	}, nil
}

func (a *Agent) GetTektonProject(ctx context.Context, request *captenpluginspb.GetTektonProjectRequest) (
	*captenpluginspb.GetTektonProjectResponse, error) {
	a.log.Infof("Get Tekton Git projects request recieved")

	projects, err := a.as.GetTektonProjects()
	if err != nil {
		a.log.Errorf("failed to get tekton Project, %v", err)
		return &captenpluginspb.GetTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get tekton Project",
		}, err
	}

	tekTonProject := &captenpluginspb.TektonProject{}
	for _, project := range projects {
		tekTonProject = &captenpluginspb.TektonProject{
			Id:             project.Id,
			GitProjectUrl:  project.GitProjectUrl,
			Status:         project.Status,
			LastUpdateTime: project.LastUpdateTime,
		}
	}

	a.log.Infof("Fetched Tekton Git projects")
	return &captenpluginspb.GetTektonProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successfully fetched the tekton projects",
		Project:       tekTonProject,
	}, nil
}

func (a *Agent) configureTektonGitRepo(req *model.TektonProject) {
	ci := captenmodel.UseCase{Type: tektonConfigUseCase, RepoURL: req.GitProjectUrl,
		VaultCredIdentifier: req.Id, PushToDefaultBranch: !a.createPr}
	wd := workers.NewConfig(a.tc, a.log)

	wkfId, err := wd.SendAsyncEvent(context.TODO(), &captenmodel.ConfigureParameters{Resource: tektonConfigUseCase}, ci)
	if err != nil {
		req.Status = string(model.TektonProjectConfigurationFailed)
		req.WorkflowId = "NA"
		if err := a.as.UpsertTektonProject(req); err != nil {
			a.log.Errorf("failed to update Cluster Gitopts Project, %v", err)
			return
		}
		a.log.Errorf("failed to send event to workflow to configure %s, %v", req.GitProjectUrl, err)
		return
	}

	a.log.Infof("Tekton Git project %s config workflow event %s created", wkfId)

	req.Status = string(model.TektonProjectConfigured)
	req.WorkflowId = wkfId
	req.WorkflowStatus = string(model.WorkFlowStatusStarted)
	if err := a.as.UpsertTektonProject(req); err != nil {
		a.log.Errorf("failed to update Cluster Gitopts Project, %v", err)
		return
	}

	go a.monitorWorkflow(req, wkfId)
	a.log.Infof("Tekton Git project %s registration completed", req.Id)
}

func (a *Agent) monitorWorkflow(req *model.TektonProject, wkfId string) {
	// during system reboot start monitoring, add it in map or somewhere.
	wd := workers.NewConfig(a.tc, a.log)
	wkfResp, err := wd.GetWorkflowInformation(context.TODO(), wkfId)
	if err != nil {
		req.Status = string(model.TektonProjectConfigurationFailed)
		req.WorkflowStatus = string(model.WorkFlowStatusFailed)
		if err := a.as.UpsertTektonProject(req); err != nil {
			a.log.Errorf("failed to update Cluster Gitopts Project, %v", err)
			return
		}
		a.log.Errorf("failed to send event to workflow to configure %s, %v", req.GitProjectUrl, err)
		return
	}

	a.log.Infof("Monitoring Tekton Git project %s config workflow event %s created", wkfId)

	req.Status = string(model.TektonProjectConfigured)
	req.WorkflowStatus = wkfResp.Status
	if err := a.as.UpsertTektonProject(req); err != nil {
		a.log.Errorf("failed to update Cluster Gitopts Project, %v", err)
		return
	}

	a.log.Infof("Tekton Git project %s monitoring completed", req.Id)
}
