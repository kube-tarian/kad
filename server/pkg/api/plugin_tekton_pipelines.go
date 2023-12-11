package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/captenpluginspb"
)

func (s *Server) GetDefaultTektonPipelines(ctx context.Context, request *captenpluginspb.GetDefaultTektonPipelinesRequest) (*captenpluginspb.GetDefaultTektonPipelinesResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.GetDefaultTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Get DefaultTektonPipelines request for cluster %s recieved, [org: %s]",
		clusterId, orgId)

	// check where to fetch the config file.

	s.log.Infof("get DefaultTektonPipelines %s request for cluster %s successful, [org: %s]",
		clusterId, orgId)
	return &captenpluginspb.GetDefaultTektonPipelinesResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "default pipeline fetch successful",
	}, nil
}

func (s *Server) GetConfiguredTektonPipelines(ctx context.Context, request *captenpluginspb.GetConfiguredTektonPipelinesRequest) (
	*captenpluginspb.GetConfiguredTektonPipelinesResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.GetConfiguredTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("GetConfiguredTektonPipelines request for cluster %s recieved, [org: %s]",
		clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.GetConfiguredTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to initialize agent",
		}, err
	}

	pieplines, err := agent.GetCaptenPluginsClient().GetConfiguredTektonPipelines(ctx, request)
	if err != nil {
		s.log.Errorf("failed to get configured tekton pipelines, %v", err)
		return &captenpluginspb.GetConfiguredTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get configured tekton pipelines",
		}, err
	}

	s.log.Infof("get configured tekton pipelines for cluster %s successful, [org: %s]",
		clusterId, orgId)
	return &captenpluginspb.GetConfiguredTektonPipelinesResponse{
		Pipelines:     pieplines.Pipelines,
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "get configured tekton pipelines successful",
	}, nil
}

func (s *Server) CreateTektonPipelines(ctx context.Context, request *captenpluginspb.CreateTektonPipelinesRequest) (
	*captenpluginspb.CreateTektonPipelinesResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.CreateTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Create TektonPipelines for cluster %s recieved, [org: %s]",
		clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.CreateTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to initialize agent",
		}, err
	}

	tektonResp, err := agent.GetCaptenPluginsClient().CreateTektonPipelines(ctx, request)
	if err != nil {
		s.log.Errorf("failed to create TektonPipelines , %v", err)
		return &captenpluginspb.CreateTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to create TektonPipelines",
		}, err
	}

	s.log.Infof("Created the TektonPipelines for cluster %s, [org: %s]",
		clusterId, orgId)
	return &captenpluginspb.CreateTektonPipelinesResponse{
		Status:        tektonResp.Status,
		StatusMessage: "Creation of TektonPipelines successful",
		PipelineName:  tektonResp.PipelineName,
	}, nil
}

func (s *Server) UpdateTektonPipelines(ctx context.Context, request *captenpluginspb.UpdateTektonPipelinesRequest) (*captenpluginspb.UpdateTektonPipelinesResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.UpdateTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("update TektonPipelines for cluster %s recieved, [org: %s]",
		clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.UpdateTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to initialize agent",
		}, err
	}

	_, err = agent.GetCaptenPluginsClient().UpdateTektonPipelines(ctx, request)
	if err != nil {
		s.log.Errorf("failed to update TektonPipelines , %v", err)
		return &captenpluginspb.UpdateTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to update TektonPipelines",
		}, err
	}

	s.log.Infof("updated the TektonPipelines for cluster %s, [org: %s]",
		clusterId, orgId)
	return &captenpluginspb.UpdateTektonPipelinesResponse{
		StatusMessage: "update of TektonPipelines successful",
		Status:        captenpluginspb.StatusCode_OK,
	}, nil
}
