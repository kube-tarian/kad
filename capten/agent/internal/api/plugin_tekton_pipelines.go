package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/agent/internal/workers"
	"github.com/kube-tarian/kad/capten/model"
	"github.com/pkg/errors"
)

var (
	pipelineSuffix = "-pipeline"
)

func (a *Agent) CreateTektonPipeline(ctx context.Context, request *captenpluginspb.CreateTektonPipelineRequest) (
	*captenpluginspb.CreateTektonPipelineResponse, error) {
	a.log.Infof("Create Tekton Pipeline %s request received", request.PipelineName)
	if err := validateArgs(request.PipelineName, request.GitOrgId, request.ContainerRegistryIds,
		request.ManagedClusterId); err != nil {
		a.log.Infof("request validation failed for Tekton Pipeline %s", request.PipelineName, err)
		return &captenpluginspb.CreateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}

	if len(request.ContainerRegistryIds) != 1 {
		a.log.Infof("currently single container registry supported, skipping create pipeline %s", request.PipelineName)
		return &captenpluginspb.CreateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "currently single container registry supported",
		}, errors.New("kindly provide only one item in container registry")
	}

	if !strings.HasSuffix(request.PipelineName, pipelineSuffix) {
		a.log.Infof("the pipeline %s should have the suffix %s in the name", request.PipelineName, pipelineSuffix)
		return &captenpluginspb.CreateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "the pipeline should have the suffix -pipeline in the name",
		}, errors.New("the pipeline should have the suffix -pipeline in the name")
	}

	gitOrgProject, err := a.as.GetGitProjectForID(request.GitOrgId)
	if err != nil {
		a.log.Infof("failed to get git org project, skipping create pipeline %s, %v", request.PipelineName, err)
		return &captenpluginspb.CreateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "failed to get git org project",
		}, err
	}

	crossplaneProject, err := a.as.GetCrossplaneProject()
	if err != nil {
		a.log.Infof("failed to get the crossplane project, skipping pipeline %s, %v", request.PipelineName, err)
		return &captenpluginspb.CreateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: fmt.Sprintf("failed to get the crossplane project"),
		}, errors.New("crossplane project not configured")
	}

	id := uuid.New()
	tektonPipeline := model.TektonPipeline{
		Id:                     id.String(),
		PipelineName:           request.PipelineName,
		GitOrgId:               gitOrgProject.Id,
		GitOrgUrl:              gitOrgProject.ProjectUrl,
		ContainerRegId:         request.ContainerRegistryIds,
		ManagedClusterId:       request.ManagedClusterId,
		CrossplaneGitProjectId: crossplaneProject.GitProjectId,
	}

	if err := a.as.UpsertTektonPipelines(&tektonPipeline); err != nil {
		a.log.Errorf("failed to store create pipeline req %s to DB, %v", request.PipelineName, err)
		return &captenpluginspb.CreateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to create pipeline  reqin db",
		}, err
	}

	_, err = a.configureTektonPipelinesGitRepo(&tektonPipeline, model.TektonPipelineCreate, true)
	if err != nil {
		a.log.Errorf("failed to configure tekton pipeline %s, %v", request.PipelineName, err)
		tektonPipeline.Status = string(model.TektonPipelineConfigurationFailed)
		tektonPipeline.WorkflowId = "NA"
		if err := a.as.UpsertTektonPipelines(&tektonPipeline); err != nil {
			a.log.Errorf("failed to configure tekton pipeline %s, for Gitopts Project, %v", request.PipelineName, err)
		}
		return &captenpluginspb.CreateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to configure tekton pipelines",
		}, err
	}

	a.log.Infof("pipeline %s created with id %s, %+v", request.PipelineName, id.String(), tektonPipeline)
	return &captenpluginspb.CreateTektonPipelineResponse{
		Id:            id.String(),
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) UpdateTektonPipeline(ctx context.Context, request *captenpluginspb.UpdateTektonPipelineRequest) (
	*captenpluginspb.UpdateTektonPipelineResponse, error) {
	a.log.Infof("Update Tekton pipeline id %s request received", request.Id)

	if err := validateArgs(request.GitOrgId, request.Id, request.ContainerRegistryIds,
		request.ManagedClusterId); err != nil {
		a.log.Infof("request validation failed for Tekton pipeline id %s", request.Id, err)
		return &captenpluginspb.UpdateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}

	gitOrgProject, err := a.as.GetGitProjectForID(request.GitOrgId)
	if err != nil {
		a.log.Infof("failed to get git org project, %v", err)
		return &captenpluginspb.UpdateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "failed to get git org project",
		}, err
	}

	id, err := uuid.Parse(request.Id)
	if err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.UpdateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: fmt.Sprintf("invalid uuid: %s", request.Id),
		}, err
	}

	pipeline, err := a.as.GetTektonPipelinesForID(id.String())
	if err != nil {
		a.log.Infof("failed to get the tekton pipeline: %s", request.Id, err)
		return &captenpluginspb.UpdateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: fmt.Sprintf("failed to get the tekton pipeline: %s", request.Id),
		}, err
	}

	crossplaneProject, err := a.as.GetCrossplaneProject()
	if err != nil {
		a.log.Infof("failed to get the crossplane project: %s", request.Id, err)
		return &captenpluginspb.UpdateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: fmt.Sprintf("failed to get the tekton pipeline: %s", request.Id),
		}, err
	}

	pipeline.ContainerRegId = request.ContainerRegistryIds
	pipeline.GitOrgId = gitOrgProject.Id
	pipeline.GitOrgUrl = gitOrgProject.ProjectUrl
	pipeline.CrossplaneGitProjectId = crossplaneProject.GitProjectId

	if err := a.as.UpsertTektonPipelines(pipeline); err != nil {
		a.log.Errorf("failed to update TektonPipeline: %s in db, %v", request.Id, err)
		return &captenpluginspb.UpdateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to update TektonPipeline in db",
		}, err
	}

	if _, err := a.configureTektonPipelinesGitRepo(pipeline, model.CrossPlaneClusterUpdate, true); err != nil {
		a.log.Errorf("failed to configure updates for TektonPipeline %s, %v", pipeline.PipelineName, err)
		return &captenpluginspb.UpdateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to configure update for TektonPipeline",
		}, err
	}

	a.log.Infof("pipeline %s update for pipeline id %s, %+v", pipeline.PipelineName, id.String(), pipeline)
	return &captenpluginspb.UpdateTektonPipelineResponse{
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
		}, err
	}

	pipeline := make([]*captenpluginspb.TektonPipelines, len(res))
	for index, r := range res {
		r.WebhookURL = "https://" + model.TektonHostName + "." + a.cfg.DomainName + "/" + r.PipelineName
		p := &captenpluginspb.TektonPipelines{Id: r.Id, PipelineName: r.PipelineName,
			WebhookURL: r.WebhookURL, Status: r.Status, GitOrgId: r.GitOrgId,
			ManagedClusterId: r.ManagedClusterId, CrossPlaneGitProjectId: r.CrossplaneGitProjectId,
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

func (a *Agent) DeleteTektonPipeline(ctx context.Context, request *captenpluginspb.DeleteTektonPipelineRequest) (
	*captenpluginspb.DeleteTektonPipelineResponse, error) {
	a.log.Infof("Delete tekton pipeline request recieved")
	if err := validateArgs(request.Id); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.DeleteTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}

	pipeline, err := a.as.GetTektonPipelinesForID(request.Id)
	if err != nil {
		a.log.Errorf("failed to get TektonPipeline from db, %v", err)
		return &captenpluginspb.DeleteTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch TektonPipelines",
		}, err
	}

	wkfID, err := a.configureTektonPipelinesGitRepo(pipeline, model.TektonPipelineDelete, false)
	if err != nil {
		a.log.Errorf("failed to initiate cleanup of  TektonPipeline from git repo, %v", err)
		return &captenpluginspb.DeleteTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to  initiate cleanup of  TektonPipeline from git repo",
		}, err
	}
	a.monitorTektonPipelineWorkflow(pipeline, wkfID)
	if pipeline.Status != string(model.TektonPipelineConfigured) {
		a.log.Infof("failed to delete tekton pipeline %s", pipeline.PipelineName)

		return &captenpluginspb.DeleteTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed",
		}, err
	}

	err = a.as.DeleteTektonPipelinesById(request.Id)
	if err != nil {
		a.log.Errorf("failed to delete TektonPipeline from db, %v", err)
		return &captenpluginspb.DeleteTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete TektonPipelines",
		}, err
	}

	a.log.Infof("Deleted tekton pipeline %s", pipeline.PipelineName)

	return &captenpluginspb.DeleteTektonPipelineResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successful",
	}, nil

}

func (a *Agent) configureTektonPipelinesGitRepo(req *model.TektonPipeline,
	action string, triggerMonitor bool) (string, error) {
	a.log.Infof("configuring tekton pipeline %s for git org %s", req.PipelineName, req.GitOrgUrl)
	crossplaneProject, err := a.as.GetCrossplaneProject()
	if err != nil {
		return "", fmt.Errorf("failed to get crossplane git project, %v", err)
	}

	containerRegistry, err := a.as.GetContainerRegistryForID(req.ContainerRegId[0])
	if err != nil {
		return "", fmt.Errorf("failed to get container registry, %v", err)
	}

	managedCluster, err := a.as.GetManagedClusterForID(req.ManagedClusterId)
	if err != nil {
		return "", fmt.Errorf("failed to get managed clsuter  %s, %v", req.ManagedClusterId, err)
	}

	tektonProject, err := a.as.GetTektonProject()
	if err != nil {
		return "", fmt.Errorf("tekton project not available, %v", err)
	}

	containerRegURLIdMap := make(map[string]string)
	containerRegURLIdMap[containerRegistry.Id] = containerRegistry.RegistryUrl

	ci := model.TektonPipelineUseCase{
		Type:         model.TektonPipelineConfigUseCase,
		PipelineName: req.PipelineName,
		RepoURL:      tektonProject.GitProjectUrl,
		CredentialIdentifiers: map[model.Identifiers]model.CredentialIdentifier{
			model.GitOrg: {
				Identifier: gitProjectEntityName,
				Id:         req.GitOrgId},
			model.Container: {
				Identifier: containerRegEntityName,
				Id:         req.ContainerRegId[0],
				Url:        containerRegistry.RegistryUrl},
			model.ManagedCluster: {
				Identifier: ManagedClusterEntityName,
				Id:         req.ManagedClusterId,
				Url:        managedCluster.ClusterName},
			model.CrossplaneGitProject: {
				Identifier: gitProjectEntityName,
				Id:         crossplaneProject.GitProjectId,
				Url:        crossplaneProject.GitProjectUrl},
			model.TektonGitProject: {
				Identifier: gitProjectEntityName,
				Id:         tektonProject.GitProjectId,
				Url:        tektonProject.GitProjectUrl},
		}}

	wd := workers.NewConfig(a.tc, a.log)
	wkfId, err := wd.SendAsyncEvent(context.TODO(),
		&model.ConfigureParameters{Resource: model.TektonPipelineConfigUseCase, Action: action}, ci)
	if err != nil {
		a.log.Errorf("failed to send event to workflow to configure, %v", err)
		return "", fmt.Errorf("failed to send event to workflow to configure %s, %v", req.GitOrgUrl, err)
	}

	a.log.Infof("tekton pipelines for Git project %s config workflow %s initiated", req.GitOrgUrl, wkfId)

	req.Status = string(model.TektonPipelineConfigurationOngoing)
	req.WorkflowId = wkfId
	req.WorkflowStatus = string(model.WorkFlowStatusStarted)
	if err := a.as.UpsertTektonPipelines(req); err != nil {
		a.log.Errorf("failed to update tekton pipelines for Gitopts Project, %v", err)
		return "", nil
	}

	if triggerMonitor {
		go a.monitorTektonPipelineWorkflow(req, wkfId)
		a.log.Infof("started monitoring the tekton pipeline %s for Git project %s", req.PipelineName, req.GitOrgUrl)
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
			a.log.Errorf("failed to update tekton pipeline to DB, %v", err)
			return
		}
		a.log.Errorf("failed to send pipeline %s event to workflow, %v", req.PipelineName, err)
		return
	}

	req.Status = string(model.TektonPipelineConfigured)
	req.WorkflowStatus = wkfResp.Status
	if err := a.as.UpsertTektonPipelines(req); err != nil {
		a.log.Errorf("failed to configure tekton pipelines for Gitopts Project, %v", err)
		return
	}
	a.log.Infof("tekton pipeline %s for Git org %s workflow %s completed", req.PipelineName, req.GitOrgUrl, wkfId)
}
