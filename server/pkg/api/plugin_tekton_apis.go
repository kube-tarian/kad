package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/captenpluginspb"
)

func (s *Server) RegisterTektonProject(ctx context.Context, request *captenpluginspb.RegisterTektonProjectRequest) (
	*captenpluginspb.RegisterTektonProjectResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx, request.Id)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.RegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.RegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to initialize agent",
		}, err
	}

	_, err = agent.GetCaptenPluginsClient().RegisterTektonProject(ctx, &captenpluginspb.RegisterTektonProjectRequest{Id: request.Id})
	if err != nil {
		s.log.Errorf("failed to register the tekton, %v", err)
		return &captenpluginspb.RegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to register the tekton",
		}, err
	}

	return &captenpluginspb.RegisterTektonProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "Tekton Registration successful",
	}, nil
}

func (s *Server) UnRegisterTektonProject(ctx context.Context, request *captenpluginspb.UnRegisterTektonProjectRequest) (
	*captenpluginspb.UnRegisterTektonProjectResponse, error) {
	orgId, err := validateOrgWithArgs(ctx, request.ClusterID)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.UnRegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}

	agent, err := s.agentHandeler.GetAgent(orgId, request.ClusterID)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.UnRegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to initialize agent",
		}, err
	}

	_, err = agent.GetCaptenPluginsClient().UnRegisterTektonProject(ctx, &captenpluginspb.UnRegisterTektonProjectRequest{Id: request.Id})
	if err != nil {
		s.log.Errorf("failed to register the tekton, %v", err)
		return &captenpluginspb.UnRegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to register the tekton",
		}, err
	}

	return &captenpluginspb.UnRegisterTektonProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "Tekton UnRegistration successful",
	}, nil
}

func (s *Server) GetTektonProjects(ctx context.Context, request *captenpluginspb.GetTektonProjectsRequest) (
	*captenpluginspb.GetTektonProjectsResponse, error) {
	orgId, err := validateOrgWithArgs(ctx, request.ClusterID)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.GetTektonProjectsResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}

	agent, err := s.agentHandeler.GetAgent(orgId, request.ClusterID)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.GetTektonProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to initialize agent",
		}, err
	}

	tektonResp, err := agent.GetCaptenPluginsClient().GetTektonProjects(ctx, &captenpluginspb.GetTektonProjectsRequest{})
	if err != nil {
		s.log.Errorf("failed to register the tekton, %v", err)
		return &captenpluginspb.GetTektonProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to register the tekton",
		}, err
	}

	return &captenpluginspb.GetTektonProjectsResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "Tekton Get successful",
		Projects:      tektonResp.Projects,
	}, nil
}
