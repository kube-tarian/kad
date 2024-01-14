package api

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/agent/internal/workers"
	"github.com/kube-tarian/kad/capten/model"
	captenmodel "github.com/kube-tarian/kad/capten/model"
)

func (a *Agent) CreateTektonPipelines(ctx context.Context, request *captenpluginspb.CreateTektonPipelinesRequest) (
	*captenpluginspb.CreateTektonPipelinesResponse, error) {
	if err := validateArgs(request.PipelineName, request.GitOrgId, request.ContainerRegistryIds); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.CreateTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	if len(request.ContainerRegistryIds) != 1 {
		a.log.Infof("currently single container registry supported")
		return &captenpluginspb.CreateTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "currently single container registry supported",
		}, nil
	}

	tektonAvailable, err := a.as.GetTektonProjectForID(request.GitOrgId)
	if err != nil {
		a.log.Infof("faile to get git project %s, %v", request.GitOrgId, err)
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
		GitProjectId:   tektonAvailable.GitProjectId,
		ContainerRegId: request.ContainerRegistryIds,
	}
	if err := a.as.UpsertTektonPipelines(&TektonPipeline); err != nil {
		a.log.Errorf("failed to store create pipeline req %s to DB, %v", request.PipelineName, err)
		return &captenpluginspb.CreateTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to create pipeline  reqin db",
		}, nil
	}

	_, err = a.configureTektonPipelinesGitRepo(&TektonPipeline, model.TektonPipelineCreate, true)
	if err != nil {
		TektonPipeline.Status = string(model.TektonPipelineConfigurationFailed)
		TektonPipeline.WorkflowId = "NA"
		if err := a.as.UpsertTektonPipelines(&TektonPipeline); err != nil {
			a.log.Errorf("failed to configure tekton pipelines: %s, for Gitopts Project, %v", request.PipelineName, err)
		}
		a.log.Errorf("failed to configure tekton pipelines for the req: %s, %v", request.PipelineName, err)
		return &captenpluginspb.CreateTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to configure tekton pipelines",
		}, nil
	}

	a.log.Infof("create pipelines %s added with id %s", request.PipelineName, id.String())
	return &captenpluginspb.CreateTektonPipelinesResponse{
		Id:            id.String(),
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) UpdateTektonPipelines(ctx context.Context, request *captenpluginspb.UpdateTektonPipelinesRequest) (
	*captenpluginspb.UpdateTektonPipelinesResponse, error) {
	if err := validateArgs(request.GitOrgId, request.Id, request.ContainerRegistryIds); err != nil {
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

	pipeline, err := a.as.GetTektonPipelinesForID(id.String())
	if err != nil {
		a.log.Infof("failed to get the tekton pipeline: %s", request.Id, err)
		return &captenpluginspb.UpdateTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: fmt.Sprintf("failed to get the tekton pipeline: %s", request.Id),
		}, nil
	}

	pipeline.ContainerRegId = request.ContainerRegistryIds
	pipeline.GitProjectId = request.GitOrgId

	if err := a.as.UpsertTektonPipelines(pipeline); err != nil {
		a.log.Errorf("failed to update TektonPipeline: %s in db, %v", request.Id, err)
		return &captenpluginspb.UpdateTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to update TektonPipeline in db",
		}, nil
	}

	if _, err := a.configureTektonPipelinesGitRepo(pipeline, model.TektonPipelineSync, true); err != nil {
		a.log.Errorf("failed to configure updates for TektonPipeline: %s, %v", request.Id, err)
		return &captenpluginspb.UpdateTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to configure update for TektonPipeline",
		}, nil
	}

	a.log.Infof("TektonPipeline, %s updated", request.Id)
	return &captenpluginspb.UpdateTektonPipelinesResponse{
		Id:            id.String(),
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

	for index, r := range res {
		r.WebhookURL = "https://" + model.TektonHostName + "." + a.cfg.DomainName + "/" + r.PipelineName
		p := &captenpluginspb.TektonPipelines{Id: r.Id, PipelineName: r.PipelineName,
			WebhookURL: r.WebhookURL, Status: r.Status, GitOrgId: r.GitProjectId,
			ContainerRegistryIds: r.ContainerRegId, LastUpdateTime: r.LastUpdateTime}
		pipeline[index] = p
	}

	a.log.Infof("Found %d tekton pipelines", len(res))
	return &captenpluginspb.GetTektonPipelinesResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successful",
		Pipelines:     pipeline,
	}, nil

}

func (a *Agent) DeleteTektonPipelines(ctx context.Context, request *captenpluginspb.DeleteTektonPipelinesRequest) (
	*captenpluginspb.DeleteTektonPipelinesResponse, error) {
	a.log.Infof("Delete tekton pipeline request recieved")
	if err := validateArgs(request.Id); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.DeleteTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	pipeline, err := a.as.GetTektonPipelinesForID(request.Id)
	if err != nil {
		a.log.Errorf("failed to get TektonPipeline from db, %v", err)
		return &captenpluginspb.DeleteTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch TektonPipelines",
		}, nil
	}

	wkfID, err := a.configureTektonPipelinesGitRepo(pipeline, model.TektonPipelineDelete, false)
	if err != nil {
		a.log.Errorf("failed to initiate cleanup of  TektonPipeline from git repo, %v", err)
		return &captenpluginspb.DeleteTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to  initiate cleanup of  TektonPipeline from git repo",
		}, nil
	}
	a.monitorTektonPipelineWorkflow(pipeline, wkfID)
	if pipeline.Status != string(model.TektonPipelineConfigured) {
		a.log.Infof("failed to delete tekton pipeline %s", pipeline.PipelineName)

		return &captenpluginspb.DeleteTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed",
		}, nil
	}

	err = a.as.DeleteTektonPipelinesById(request.Id)
	if err != nil {
		a.log.Errorf("failed to delete TektonPipeline from db, %v", err)
		return &captenpluginspb.DeleteTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete TektonPipelines",
		}, nil
	}

	a.log.Infof("Deleted tekton pipeline %s", pipeline.PipelineName)

	return &captenpluginspb.DeleteTektonPipelinesResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successful",
	}, nil

}

func (a *Agent) SyncTektonPipelines(ctx context.Context, request *captenpluginspb.SyncTektonPipelinesRequest) (
	*captenpluginspb.SyncTektonPipelinesResponse, error) {
	a.log.Infof("Get tekton pipeline request recieved")
	res, err := a.as.GetTektonPipeliness()
	if err != nil {
		a.log.Errorf("failed to get TektonPipeline from db, %v", err)
		return &captenpluginspb.SyncTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch TektonPipelines",
		}, nil
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

	a.log.Infof("Successfully triggered the sync")
	return &captenpluginspb.SyncTektonPipelinesResponse{
		Pipelines:     pipelines,
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "Successfully triggered the sync",
	}, nil
}

func (a *Agent) configureTektonPipelinesGitRepo(req *model.TektonPipeline, action string, triggerMonitor bool) (string, error) {
	a.log.Infof("configuring tekton pipeline for the git repo %s", req.GitProjectUrl)
	tektonProject := model.TektonPipeline{}
	proj, err := a.as.GetGitProjectForID(req.GitProjectId)
	if err != nil {
		a.log.Errorf("failed to send event to workflow to configure, %v", err)
		return "", fmt.Errorf("failed to send event to workflow to configure %s, %v", req.GitProjectId, err)
	}

	containerRegs, err := a.as.GetContainerRegistries()
	if err != nil {
		a.log.Errorf("failed to send event to workflow to configure, %v", err)
		return "", fmt.Errorf("failed to send event to workflow to configure %s, %v", req.GitProjectId, err)
	}
	containerRegURLIdMap := make(map[string]string)
	for _, containerReg := range containerRegs {
		containerRegURLIdMap[containerReg.Id] = containerReg.RegistryUrl
	}

	ci := captenmodel.TektonPipelineUseCase{Type: model.TektonPipelineConfigUseCase,
		PipelineName: req.PipelineName, RepoURL: proj.ProjectUrl,
		GitCredId: req.GitProjectId, ContainerRegUrlIdMap: containerRegURLIdMap,
		ContainerRegCredIdentifier: containerRegEntityName, GitCredIdentifier: gitProjectEntityName}
	wd := workers.NewConfig(a.tc, a.log)

	wkfId, err := wd.SendAsyncEvent(context.TODO(),
		&captenmodel.ConfigureParameters{Resource: model.TektonPipelineConfigUseCase, Action: action}, ci)
	if err != nil {
		a.log.Errorf("failed to send event to workflow to configure, %v", err)
		return "", fmt.Errorf("failed to send event to workflow to configure %s, %v", proj.ProjectUrl, err)
	}

	a.log.Infof("tekton pipelines for Git project %s config workflow %s initiated", proj.ProjectUrl, wkfId)

	tektonProject.Status = string(model.TektonPipelineConfigurationOngoing)
	tektonProject.WorkflowId = wkfId
	tektonProject.WorkflowStatus = string(model.WorkFlowStatusStarted)
	if err := a.as.UpsertTektonPipelines(req); err != nil {
		a.log.Errorf("failed to update tekton pipelines for Gitopts Project, %v", err)
		return "", nil
	}

	if triggerMonitor {
		go a.monitorTektonPipelineWorkflow(req, wkfId)
		a.log.Infof("started monitoring the tekton pipelines for Git project %s", tektonProject.GitProjectUrl)
	}

	return wkfId, nil
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
