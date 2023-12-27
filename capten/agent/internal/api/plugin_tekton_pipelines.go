package api

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"
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

	TektonPipeline := captenpluginspb.TektonPipelines{
		Id:                  id.String(),
		PipelineName:        request.PipelineName,
		GitOrgId:            request.GitOrgId,
		ContainerRegistryId: request.ContainerRegistryId,
	}
	if err := a.as.UpsertTektonPipelines(&TektonPipeline); err != nil {
		a.log.Errorf("failed to store create pipelines req to DB, %v", err)
		return &captenpluginspb.CreateTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to create pipelines in db",
		}, nil
	}

	a.log.Infof("create pipelines %s added with id %s", request.PipelineName, id.String())
	return &captenpluginspb.CreateTektonPipelinesResponse{
		Id:            id.String(),
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) UpdateTektonPipeline(ctx context.Context, request *captenpluginspb.UpdateTektonPipelinesRequest) (
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

	TektonPipeline := captenpluginspb.TektonPipelines{
		Id:                  id.String(),
		ContainerRegistryId: request.ContainerRegistryId,
		GitOrgId:            request.GitOrgId,
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

func (a *Agent) GetTektonPipeline(ctx context.Context, request *captenpluginspb.GetTektonPipelinesRequest) (
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

	for _, r := range res {
		// ## TODO Replace the host with tekton static name
		r.WebhookURL = "" + r.PipelineName
	}

	a.log.Infof("Found %d tekton pipelines", len(res))
	return &captenpluginspb.GetTektonPipelinesResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successful",
		Pipelines:     res,
	}, nil

}
