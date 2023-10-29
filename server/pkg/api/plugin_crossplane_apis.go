package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/captenpluginspb"
)

func (s *Server) RegisterCrossplaneProject(ctx context.Context, request *captenpluginspb.RegisterCrossplaneProjectRequest) (
	*captenpluginspb.RegisterCrossplaneProjectResponse, error) {

	orgId, clusterId, err := validateOrgClusterWithArgs(ctx, request.Id)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.RegisterCrossplaneProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Register Crossplane Git project %s request for cluster %s recieved, [org: %s]",
		request.Id, clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.RegisterCrossplaneProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to initialize agent",
		}, err
	}

	_, err = agent.GetCaptenPluginsClient().RegisterCrossplaneProject(ctx, &captenpluginspb.RegisterCrossplaneProjectRequest{Id: request.Id})
	if err != nil {
		s.log.Errorf("failed to register the Crossplane, %v", err)
		return &captenpluginspb.RegisterCrossplaneProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to register the Crossplane",
		}, err
	}

	s.log.Infof("Crossplane Git project %s request for cluster %s Registered, [org: %s]",
		request.Id, clusterId, orgId)
	return &captenpluginspb.RegisterCrossplaneProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "Crossplane Registration successful",
	}, nil
}

func (s *Server) UnRegisterCrossplaneProject(ctx context.Context, request *captenpluginspb.UnRegisterCrossplaneProjectRequest) (
	*captenpluginspb.UnRegisterCrossplaneProjectResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx, request.Id)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.UnRegisterCrossplaneProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("UnRegister Crossplane Git project %s request for cluster %s recieved, [org: %s]",
		request.Id, clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.UnRegisterCrossplaneProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to initialize agent",
		}, err
	}

	_, err = agent.GetCaptenPluginsClient().UnRegisterCrossplaneProject(ctx, &captenpluginspb.UnRegisterCrossplaneProjectRequest{Id: request.Id})
	if err != nil {
		s.log.Errorf("failed to register the Crossplane, %v", err)
		return &captenpluginspb.UnRegisterCrossplaneProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to unregister the Crossplane",
		}, err
	}

	s.log.Infof("Crossplane Git project %s request for cluster %s UnRegistered, [org: %s]",
		request.Id, clusterId, orgId)
	return &captenpluginspb.UnRegisterCrossplaneProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "Crossplane UnRegistration successful",
	}, nil
}

func (s *Server) GetCrossplaneProject(ctx context.Context, request *captenpluginspb.GetCrossplaneProjectsRequest) (
	*captenpluginspb.GetCrossplaneProjectsResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.GetCrossplaneProjectsResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Get Crossplane Git projects request for cluster %s recieved, [org: %s]",
		clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.GetCrossplaneProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to initialize agent",
		}, err
	}

	CrossplaneResp, err := agent.GetCaptenPluginsClient().GetCrossplaneProject(ctx, &captenpluginspb.GetCrossplaneProjectsRequest{})
	if err != nil {
		s.log.Errorf("failed to register the Crossplane, %v", err)
		return &captenpluginspb.GetCrossplaneProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Crossplane project",
		}, err
	}

	s.log.Infof("Fetched Crossplane Git project, id: %v for cluster %s, [org: %s]",
		CrossplaneResp.Project.Id, clusterId, orgId)
	return &captenpluginspb.GetCrossplaneProjectsResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "Crossplane Get successful",
		Project:       CrossplaneResp.Project,
	}, nil
}
