package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/captenpluginspb"
)

func (s *Server) GetTektonPipelines(ctx context.Context, request *captenpluginspb.GetTektonPipelinesRequest) (
	*captenpluginspb.GetTektonPipelinesResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.GetTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("GetTektonPipelines request for cluster %s recieved, [org: %s]",
		clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.GetTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to initialize agent",
		}, err
	}

	pieplines, err := agent.GetCaptenPluginsClient().GetTektonPipelines(ctx, request)
	if err != nil {
		s.log.Errorf("failed to get  tekton pipelines, %v", err)
		return &captenpluginspb.GetTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get  tekton pipelines",
		}, err
	}

	s.log.Infof("get tekton pipelines for cluster %s successful, [org: %s]",
		clusterId, orgId)
	return &captenpluginspb.GetTektonPipelinesResponse{
		Pipelines:     pieplines.Pipelines,
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "get  tekton pipelines successful",
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

func (s *Server) SyncTektonPipelines(ctx context.Context, request *captenpluginspb.SyncTektonPipelinesRequest) (*captenpluginspb.SyncTektonPipelinesResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.SyncTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Sync TektonPipelines for cluster %s recieved, [org: %s]",
		clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.SyncTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to initialize agent",
		}, err
	}

	resp, err := agent.GetCaptenPluginsClient().SyncTektonPipelines(ctx, request)
	if err != nil {
		s.log.Errorf("failed to sync TektonPipelines , %v", err)
		return &captenpluginspb.SyncTektonPipelinesResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to sync TektonPipelines",
		}, err
	}

	s.log.Infof("sync the TektonPipelines for cluster %s, [org: %s]",
		clusterId, orgId)
	return &captenpluginspb.SyncTektonPipelinesResponse{
		Pipelines:     resp.Pipelines,
		StatusMessage: "sync of TektonPipelines successful",
		Status:        captenpluginspb.StatusCode_OK,
	}, nil
}
