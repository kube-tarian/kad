package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/captenpluginspb"
)

func (s *Server) RegisterArgoCDProject(ctx context.Context, request *captenpluginspb.RegisterArgoCDProjectRequest) (
	*captenpluginspb.RegisterArgoCDProjectResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx, request.Id)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.RegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Register ArgoCD Git project %s request for cluster %s recieved, [org: %s]",
		request.Id, clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.RegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to Register the ArgoCD Project",
		}, err
	}

	response, err := agent.GetCaptenPluginsClient().RegisterArgoCDProject(context.Background(),
		&captenpluginspb.RegisterArgoCDProjectRequest{Id: request.Id})
	if err != nil {
		s.log.Errorf("failed to Register the ArgoCD Project, %v", err)
		return &captenpluginspb.RegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to Register the ArgoCD Project",
		}, err
	}

	s.log.Infof("ArgoCD Git project %s request for cluster %s Registered, [org: %s]",
		request.Id, clusterId, orgId)
	return &captenpluginspb.RegisterArgoCDProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: response.StatusMessage,
	}, nil
}

func (s *Server) UnRegisterArgoCDProject(ctx context.Context, request *captenpluginspb.UnRegisterArgoCDProjectRequest) (
	*captenpluginspb.UnRegisterArgoCDProjectResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx, request.Id)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.UnRegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("UnRegister ArgoCD Git project %s request for cluster %s recieved, [org: %s]",
		request.Id, clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.UnRegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to Unregister the ArgoCD Project",
		}, err
	}

	response, err := agent.GetCaptenPluginsClient().UnRegisterArgoCDProject(context.Background(),
		&captenpluginspb.UnRegisterArgoCDProjectRequest{Id: request.Id})
	if err != nil {
		s.log.Errorf("failed to Unregister the ArgoCD Project, %v", err)
		return &captenpluginspb.UnRegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to Unregister the ArgoCD Project",
		}, err
	}

	s.log.Infof("ArgoCD Git project %s request for cluster %s UnRegistered, [org: %s]",
		request.Id, clusterId, orgId)
	return &captenpluginspb.UnRegisterArgoCDProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: response.StatusMessage,
	}, nil
}

func (s *Server) GetArgoCDProjects(ctx context.Context, request *captenpluginspb.GetArgoCDProjectsRequest) (
	*captenpluginspb.GetArgoCDProjectsResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.GetArgoCDProjectsResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Get ArgoCD Git projects request for cluster %s recieved, [org: %s]",
		clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.GetArgoCDProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get ArgoCD Project",
		}, err
	}

	response, err := agent.GetCaptenPluginsClient().GetArgoCDProjects(context.Background(),
		&captenpluginspb.GetArgoCDProjectsRequest{})
	if err != nil {
		s.log.Errorf("failed to fetch ArgoCD projects, %v", err)
		return &captenpluginspb.GetArgoCDProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch ArgoCD projects",
		}, err
	}

	s.log.Infof("Fetch %d ArgoCD Git projects for cluster %s, [org: %s]",
		len(response.Projects), clusterId, orgId)
	return &captenpluginspb.GetArgoCDProjectsResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: response.StatusMessage,
		Projects:      response.Projects,
	}, nil
}
