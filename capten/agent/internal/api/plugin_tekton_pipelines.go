package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/agent/internal/workers"
	"github.com/kube-tarian/kad/capten/model"
	captenmodel "github.com/kube-tarian/kad/capten/model"
	"github.com/pkg/errors"
)

var (
	pipelineSuffix = "-pipeline"
)

func (a *Agent) CreateTektonPipeline(ctx context.Context, request *captenpluginspb.CreateTektonPipelineRequest) (
	*captenpluginspb.CreateTektonPipelineResponse, error) {
	if err := validateArgs(request.PipelineName, request.GitOrgId, request.ContainerRegistryIds,
		request.ManagedClusterId, request.CrossPlaneGitProjectId); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.CreateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}

	if len(request.ContainerRegistryIds) != 1 {
		a.log.Infof("currently single container registry supported")
		return &captenpluginspb.CreateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "currently single container registry supported",
		}, errors.New("kindly provide only one item in container registry")
	}

	if !strings.HasSuffix(request.PipelineName, pipelineSuffix) {
		a.log.Infof("the pipeline should have the suffix %s in the name", pipelineSuffix)
		return &captenpluginspb.CreateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "the pipeline should have the suffix -pipeline in the name",
		}, errors.New("the pipeline should have the suffix -pipeline in the name")
	}

	tektonAvailable, err := a.as.GetTektonProjectForID(request.GitOrgId)
	if err != nil {
		a.log.Infof("failed to get git project %s, %v", request.GitOrgId, err)
		return &captenpluginspb.CreateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}

	_, err = a.as.GetContainerRegistryForID(request.ContainerRegistryIds[0])
	if err != nil {
		a.log.Infof("failed to get container registry %s, %v", request.ContainerRegistryIds[0], err)
		return &captenpluginspb.CreateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "failed to get container registry",
		}, err
	}

	crossplane, err := a.as.GetCrossplaneProjectForID(request.CrossPlaneGitProjectId)
	if err != nil {
		a.log.Infof("failed to get crossplane git project %s, %v", request.CrossPlaneGitProjectId, err)
		return &captenpluginspb.CreateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "failed to get crossplane git project",
		}, nil
	}

	_, err = a.as.GetManagedClusterForID(request.ManagedClusterId)
	if err != nil {
		a.log.Infof("failed to get managedCluster id %s, %v", request.ManagedClusterId, err)
		return &captenpluginspb.CreateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "failed to get managedCluster id",
		}, err
	}

	a.log.Infof("Add Create Tekton Pipeline registry %s request received", request.PipelineName)

	id := uuid.New()

	TektonPipeline := model.TektonPipeline{
		Id:                     id.String(),
		PipelineName:           request.PipelineName,
		GitProjectId:           tektonAvailable.GitProjectId,
		ContainerRegId:         request.ContainerRegistryIds,
		ManagedClusterId:       request.ManagedClusterId,
		CrossplaneGitProjectId: crossplane.GitProjectId,
	}

	if err := a.as.UpsertTektonPipelines(&TektonPipeline); err != nil {
		a.log.Errorf("failed to store create pipeline req %s to DB, %v", request.PipelineName, err)
		return &captenpluginspb.CreateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to create pipeline  reqin db",
		}, err
	}

	_, oErr := a.configureTektonPipelinesGitRepo(&TektonPipeline, model.TektonPipelineCreate, true)
	if oErr != nil {
		TektonPipeline.Status = string(model.TektonPipelineConfigurationFailed)
		TektonPipeline.WorkflowId = "NA"
		if err := a.as.UpsertTektonPipelines(&TektonPipeline); err != nil {
			a.log.Errorf("failed to configure tekton pipelines: %s, for Gitopts Project, %v", request.PipelineName, err)
		}
		a.log.Errorf("failed to configure tekton pipelines for the req: %s, %v", request.PipelineName, err)
		return &captenpluginspb.CreateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to configure tekton pipelines",
		}, oErr
	}

	a.log.Infof("create pipelines %s added with id %s", request.PipelineName, id.String())
	return &captenpluginspb.CreateTektonPipelineResponse{
		Id:            id.String(),
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) UpdateTektonPipeline(ctx context.Context, request *captenpluginspb.UpdateTektonPipelineRequest) (
	*captenpluginspb.UpdateTektonPipelineResponse, error) {
	if err := validateArgs(request.GitOrgId, request.Id, request.ContainerRegistryIds,
		request.ManagedClusterId, request.CrossPlaneGitProjectId); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.UpdateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}

	_, err := a.as.GetContainerRegistryForID(request.ContainerRegistryIds[0])
	if err != nil {
		a.log.Infof("failed to get container registry %s, %v", request.ContainerRegistryIds[0], err)
		return &captenpluginspb.UpdateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "failed to get container registry",
		}, err
	}

	_, err = a.as.GetCrossplaneProjectForID(request.CrossPlaneGitProjectId)
	if err != nil {
		a.log.Infof("failed to get crossplane git project %s, %v", request.CrossPlaneGitProjectId, err)
		return &captenpluginspb.UpdateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "failed to get crossplane git project",
		}, err
	}

	_, err = a.as.GetManagedClusterForID(request.ManagedClusterId)
	if err != nil {
		a.log.Infof("failed to get managedCluster id %s, %v", request.ManagedClusterId, err)
		return &captenpluginspb.UpdateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "failed to get managedCluster id",
		}, err
	}

	a.log.Infof("Update tekton pipelines project, %s request recieved", request.Id)

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

	pipeline.ContainerRegId = request.ContainerRegistryIds
	pipeline.GitProjectId = request.GitOrgId

	if err := a.as.UpsertTektonPipelines(pipeline); err != nil {
		a.log.Errorf("failed to update TektonPipeline: %s in db, %v", request.Id, err)
		return &captenpluginspb.UpdateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to update TektonPipeline in db",
		}, err
	}

	if _, err := a.configureTektonPipelinesGitRepo(pipeline, model.TektonPipelineSync, true); err != nil {
		a.log.Errorf("failed to configure updates for TektonPipeline: %s, %v", request.Id, err)
		return &captenpluginspb.UpdateTektonPipelineResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to configure update for TektonPipeline",
		}, err
	}

	a.log.Infof("TektonPipeline, %s updated", request.Id)
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
			WebhookURL: r.WebhookURL, Status: r.Status, GitOrgId: r.GitProjectId,
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

func (a *Agent) configureTektonPipelinesGitRepo(req *model.TektonPipeline, action string, triggerMonitor bool) (string, error) {
	a.log.Infof("configuring tekton pipeline for the git repo %s", req.GitProjectUrl)
	tektonProject := model.TektonPipeline{}
	proj, err := a.as.GetGitProjectForID(req.GitProjectId)
	if err != nil {
		a.log.Errorf("failed to send event to workflow to configure, %v", err)
		return "", fmt.Errorf("failed to send event to workflow to configure %s, %v", req.GitProjectId, err)
	}

	containerReg, err := a.as.GetContainerRegistryForID(req.ContainerRegId[0])
	if err != nil {
		a.log.Errorf("failed to send event to workflow to configure, %v", err)
		return "", fmt.Errorf("failed to send event to workflow to configure %s, %v", req.GitProjectId, err)
	}

	extraGitProject, err := a.as.GetCrossplaneProjectForID(req.CrossplaneGitProjectId)
	if err != nil {
		a.log.Infof("failed to get crossplane git project %s, %v", req.CrossplaneGitProjectId, err)
		return "", fmt.Errorf("failed to get crossplane git project %s, %v", req.CrossplaneGitProjectId, err)
	}

	managedCluster, err := a.as.GetManagedClusterForID(req.ManagedClusterId)
	if err != nil {
		a.log.Infof("failed to get managed clsuter %s, %v", req.ManagedClusterId, err)
		return "", fmt.Errorf("failed to get managed clsuter  %s, %v", req.ManagedClusterId, err)
	}

	containerRegURLIdMap := make(map[string]string)
	containerRegURLIdMap[containerReg.Id] = containerReg.RegistryUrl

	ci := captenmodel.TektonPipelineUseCase{Type: model.TektonPipelineConfigUseCase,
		PipelineName: req.PipelineName, RepoURL: proj.ProjectUrl,
		CredentialIdentifiers: map[captenmodel.Identifiers]captenmodel.CredentialIdentifier{
			captenmodel.Git: {Identifier: gitProjectEntityName, Id: req.GitProjectId},
			captenmodel.Container: {Identifier: containerRegEntityName,
				Id: req.ContainerRegId[0], Url: containerReg.RegistryUrl},
			captenmodel.ManagedCluster:  {Identifier: ManagedClusterEntityName, Id: req.ManagedClusterId, Url: managedCluster.ClusterName},
			captenmodel.ExtraGitProject: {Identifier: gitProjectEntityName, Id: req.CrossplaneGitProjectId, Url: extraGitProject.GitProjectUrl},
		}}
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
