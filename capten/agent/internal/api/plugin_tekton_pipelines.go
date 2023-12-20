package api

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kube-tarian/kad/capten/agent/internal/workers"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/model"
	captenmodel "github.com/kube-tarian/kad/capten/model"
)

const (
	tektonPipelineConfigUseCase string = "tektonpipelines"
)

func (a *Agent) CreateTektonPipelines(ctx context.Context, request *captenpluginspb.CreateTektonPipelinesRequest) (
	*captenpluginspb.CreateTektonPipelinesResponse, error) {
	if err := validateArgs(request.PipelineName, request.GitOrgId, request.ContainerRegistryId); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.CreateTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	a.log.Infof("Add Create Tekton Pipeline registry %s request received", request.PipelineName)

	id := uuid.New()

	TektonPipeline := model.TektonPipeline{
		Id:             id.String(),
		PipelineName:   request.PipelineName,
		GitProjectId:   request.GitOrgId,
		ContainerRegId: request.ContainerRegistryId,
	}
	if err := a.as.UpsertTektonPipelines(&TektonPipeline); err != nil {
		a.log.Errorf("failed to store create pipelines req to DB, %v", err)
		return &captenpluginspb.CreateTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to create pipelines in db",
		}, nil
	}

	go a.configureTektonPipelinesGitRepo(&TektonPipeline)

	a.log.Infof("create pipelines %s added with id %s", request.PipelineName, id.String())
	return &captenpluginspb.CreateTektonPipelinesResponse{
		Id:            id.String(),
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) UpdateTektonPipelines(ctx context.Context, request *captenpluginspb.UpdateTektonPipelinesRequest) (
	*captenpluginspb.UpdateTektonPipelinesResponse, error) {
	if err := validateArgs(request.GitOrgId, request.Id, request.ContainerRegistryId); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.UpdateTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	a.log.Infof("Update tekton pipelines project, %s request recieved", request.Id)

	id, err := uuid.Parse(request.Id)
	if err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.UpdateTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: fmt.Sprintf("invalid uuid: %s", request.Id),
		}, nil
	}

	TektonPipeline := model.TektonPipeline{
		Id:             id.String(),
		ContainerRegId: request.ContainerRegistryId,
		GitProjectId:   request.GitOrgId,
	}

	if err := a.as.UpsertTektonPipelines(&TektonPipeline); err != nil {
		a.log.Errorf("failed to update TektonPipeline in db, %v", err)
		return &captenpluginspb.UpdateTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to update TektonPipeline in db",
		}, nil
	}

	a.log.Infof("TektonPipeline, %s updated", request.Id)
	return &captenpluginspb.UpdateTektonPipelinesResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) GetTektonPipelines(ctx context.Context, request *captenpluginspb.GetTektonPipelinesRequest) (
	*captenpluginspb.GetTektonPipelinesResponse, error) {
	a.log.Infof("Get tekton pipeline request recieved")
	res, err := a.as.GetTektonPipeliness()
	if err != nil {
		a.log.Errorf("failed to get TektonPipeline from db, %v", err)
		return &captenpluginspb.GetTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch TektonPipelines",
		}, nil
	}

	pipeline := make([]*captenpluginspb.TektonPipelines, len(res))

	for _, r := range res {
		// ## TODO Replace the host with tekton static name
		r.WebhookURL = "" + r.PipelineName
		p := &captenpluginspb.TektonPipelines{Id: r.Id, PipelineName: r.PipelineName,
			WebhookURL: r.WebhookURL, Status: r.Status, GitOrgId: r.GitProjectId,
			ContainerRegistryId: r.ContainerRegId, LastUpdateTime: r.LastUpdateTime}
		pipeline = append(pipeline, p)
	}

	a.log.Infof("Found %d tekton pipelines", len(res))
	return &captenpluginspb.GetTektonPipelinesResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successful",
		Pipelines:     pipeline,
	}, nil

}

func (a *Agent) configureTektonPipelinesGitRepo(req *model.TektonPipeline) error {
	a.log.Infof("configuring tekton pipeline for the git repo %s", req.GitProjectUrl)
	tektonProject := model.TektonPipeline{}
	proj, err := a.as.GetGitProjectForID(req.GitProjectId)
	if err != nil {
		tektonProject.Status = string(model.TektonPipelineConfigurationFailed)
		tektonProject.WorkflowId = "NA"
		if err := a.as.UpsertTektonPipelines(req); err != nil {
			a.log.Errorf("failed to configure tekton pipelines for Gitopts Project, %v", err)
			return fmt.Errorf("failed to configure tekton pipelines for Gitopts Project, %v", err)
		}
		a.log.Errorf("failed to send event to workflow to configure, %v", err)
		return fmt.Errorf("failed to send event to workflow to configure %s, %v", req.GitProjectId, err)
	}

	containerRegs, err := a.as.GetContainerRegistrys()
	if err != nil {
		tektonProject.Status = string(model.TektonPipelineConfigurationFailed)
		tektonProject.WorkflowId = "NA"
		if err := a.as.UpsertTektonPipelines(req); err != nil {
			a.log.Errorf("failed to onfigure tekton pipelines for Gitopts Project %v", err)
			return fmt.Errorf("failed to configure tekton pipelines for Gitopts Project, %v", err)
		}
		a.log.Errorf("failed to send event to workflow to configure, %v", err)
		return fmt.Errorf("failed to send event to workflow to configure %s, %v", req.GitProjectId, err)
	}
	containerRegURLIdMap := make(map[string]string)
	for _, containerReg := range containerRegs {
		containerRegURLIdMap[containerReg.Id] = containerReg.RegistryUrl
	}

	ci := captenmodel.TektonPipelineUseCase{Type: tektonPipelineConfigUseCase,
		PipelineName: req.PipelineName, RepoURL: proj.ProjectUrl,
		GitCredId: req.GitProjectId, ContainerRegUrlIdMap: containerRegURLIdMap,
		ContainerRegCredIdentifier: containerRegEntityName, GitCredIdentifier: gitProjectEntityName}
	wd := workers.NewConfig(a.tc, a.log)

	wkfId, err := wd.SendAsyncEvent(context.TODO(),
		&captenmodel.ConfigureParameters{Resource: tektonPipelineConfigUseCase, Action: model.TektonPipelineSync}, ci)
	if err != nil {
		tektonProject.Status = string(model.TektonPipelineConfigurationFailed)
		tektonProject.WorkflowId = "NA"
		if err := a.as.UpsertTektonPipelines(req); err != nil {
			a.log.Errorf("failed to tekton pipelines for Gitopts Project, %v", err)
			return fmt.Errorf("failed to tekton pipelines for Gitopts Project, %v", err)
		}
		a.log.Errorf("failed to send event to workflow to configure, %v", err)
		return fmt.Errorf("failed to send event to workflow to configure %s, %v", proj.ProjectUrl, err)
	}

	a.log.Infof("tekton pipelines for Git project %s config workflow %s created", proj.ProjectUrl, wkfId)

	tektonProject.Status = string(model.TektonPipelineConfigured)
	tektonProject.WorkflowId = wkfId
	tektonProject.WorkflowStatus = string(model.WorkFlowStatusStarted)
	if err := a.as.UpsertTektonPipelines(req); err != nil {
		a.log.Errorf("failed to update tekton pipelines for Gitopts Project, %v", err)
		return nil
	}

	go a.monitorTektonPipelineWorkflow(req, wkfId)
	a.log.Infof("tekton pipelines for Git project %s registration workflow monitor started", tektonProject.GitProjectUrl)
	return nil
}

func (a *Agent) monitorTektonPipelineWorkflow(req *model.TektonPipeline, wkfId string) {
	// during system reboot start monitoring, add it in map or somewhere.
	wd := workers.NewConfig(a.tc, a.log)
	wkfResp, err := wd.GetWorkflowInformation(context.TODO(), wkfId)
	if err != nil {
		req.Status = string(model.TektonPipelineConfigurationFailed)
		req.WorkflowStatus = string(model.WorkFlowStatusFailed)
		if err := a.as.UpsertTektonPipelines(req); err != nil {
			a.log.Errorf("failed to configure tekton pipelines for Gitopts Project, %v", err)
			return
		}
		a.log.Errorf("failed to send event to workflow to configure %s, %v", req.GitProjectUrl, err)
		return
	}

	req.Status = string(model.TektonPipelineConfigured)
	req.WorkflowStatus = wkfResp.Status
	if err := a.as.UpsertTektonPipelines(req); err != nil {
		a.log.Errorf("failed to configure tekton pipelines for Gitopts Project, %v", err)
		return
	}
	a.log.Infof("tekton pipelines for Git project %s config workflow %s completed", req.GitProjectUrl, wkfId)
}
