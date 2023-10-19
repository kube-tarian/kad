package api

import (
	"context"
	"fmt"

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
	return &captenpluginspb.UnRegisterTektonProjectResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}

func (s *Server) GetTektonProjects(ctx context.Context, request *captenpluginspb.GetTektonProjectsRequest) (
	*captenpluginspb.GetTektonProjectsResponse, error) {
	return &captenpluginspb.GetTektonProjectsResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}
