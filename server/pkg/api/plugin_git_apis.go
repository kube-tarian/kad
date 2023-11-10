package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/captenpluginspb"
)

func (s *Server) AddGitProject(ctx context.Context, request *captenpluginspb.AddGitProjectRequest) (
	*captenpluginspb.AddGitProjectResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx, request.ProjectUrl, request.AccessToken)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.AddGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	s.log.Infof("Add Git project %s request for cluster %s recieved, [org: %s]",
		request.ProjectUrl, clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent for cluster %s, %v", clusterId, err)
		return &captenpluginspb.AddGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add the Cluster GitProject",
		}, nil
	}

	response, err := agent.GetCaptenPluginsClient().AddGitProject(context.Background(), request)
	if err != nil {
		s.log.Errorf("failed to add the Cluster GitProject for cluster %s, %v", clusterId, err)
		return &captenpluginspb.AddGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add the Cluster GitProject",
		}, nil
	}

	if response.Status != captenpluginspb.StatusCode_OK {
		s.log.Errorf("failed to add the ClusterProject for cluster %s, %s, %s", response.Status, response.StatusMessage)
		return &captenpluginspb.AddGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add the Cluster GitProject",
		}, nil
	}

	s.log.Infof("Git project %s request for cluster %s added, [org: %s]",
		request.ProjectUrl, clusterId, orgId)
	return &captenpluginspb.AddGitProjectResponse{
		Id:            response.GetId(),
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil

}

func (s *Server) UpdateGitProject(ctx context.Context, request *captenpluginspb.UpdateGitProjectRequest) (
	*captenpluginspb.UpdateGitProjectResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx, request.Id, request.ProjectUrl, request.AccessToken)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.UpdateGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	s.log.Infof("Update Git project %s, %s request for cluster %s recieved, [org: %s]",
		request.ProjectUrl, request.Id, clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent for cluster %s, %v", clusterId, err)
		return &captenpluginspb.UpdateGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add the Cluster GitProject",
		}, nil
	}

	response, err := agent.GetCaptenPluginsClient().UpdateGitProject(context.Background(), request)
	if err != nil {
		s.log.Errorf("failed to add the Cluster GitProject for cluster %s, %v", clusterId, err)
		return &captenpluginspb.UpdateGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add the Cluster GitProject",
		}, nil
	}

	if response.Status != captenpluginspb.StatusCode_OK {
		s.log.Errorf("failed to update the ClusterProject for cluster %s, %s, %s", response.Status, response.StatusMessage)
		return &captenpluginspb.UpdateGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to update the Cluster GitProject",
		}, nil
	}

	s.log.Infof("Git project %s, %s request for cluster %s updated, [org: %s]",
		request.ProjectUrl, request.Id, clusterId, orgId)
	return &captenpluginspb.UpdateGitProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (s *Server) DeleteGitProject(ctx context.Context, request *captenpluginspb.DeleteGitProjectRequest) (
	*captenpluginspb.DeleteGitProjectResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx, request.Id)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.DeleteGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	s.log.Infof("Delete Git project %s request for cluster %s recieved, [org: %s]",
		request.Id, clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent for cluster %s, %v", clusterId, err)
		return &captenpluginspb.DeleteGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete the Cluster GitProject",
		}, nil
	}

	response, err := agent.GetCaptenPluginsClient().DeleteGitProject(context.Background(), request)
	if err != nil {
		s.log.Errorf("failed to delete the Cluster GitProject for cluster %s, %v", clusterId, err)
		return &captenpluginspb.DeleteGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete the Cluster GitProject",
		}, nil
	}

	if response.Status != captenpluginspb.StatusCode_OK {
		s.log.Errorf("failed to update the ClusterProject for cluster %s, %s, %s", response.Status, response.StatusMessage)
		return &captenpluginspb.DeleteGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete the Cluster GitProject",
		}, nil
	}

	return &captenpluginspb.DeleteGitProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (s *Server) GetGitProjects(ctx context.Context, request *captenpluginspb.GetGitProjectsRequest) (
	*captenpluginspb.GetGitProjectsResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.GetGitProjectsResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Get Git projects request for cluster %s recieved, [org: %s]",
		clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.GetGitProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster GitProject",
		}, nil
	}

	response, err := agent.GetCaptenPluginsClient().GetGitProjects(context.Background(), request)
	if err != nil {
		s.log.Errorf("failed to get the Cluster GitProject, %v", err)
		return &captenpluginspb.GetGitProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster GitProject",
		}, nil
	}

	if response.Status != captenpluginspb.StatusCode_OK {
		s.log.Errorf("failed to get the ClusterProject, %s", response.StatusMessage)
		return &captenpluginspb.GetGitProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster GitProject",
		}, nil
	}

	s.log.Infof("Fetched %d Git projects request for cluster %s processed, [org: %s]",
		len(response.Projects), clusterId, orgId)
	return &captenpluginspb.GetGitProjectsResponse{
		Projects:      response.Projects,
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (s *Server) GetGitProjectsForLabels(ctx context.Context, request *captenpluginspb.GetGitProjectsForLabelsRequest) (
	*captenpluginspb.GetGitProjectsForLabelsResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.GetGitProjectsForLabelsResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Get Git projects request with lables %v for cluster %s recieved, [org: %s]",
		request.Labels, clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.GetGitProjectsForLabelsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster GitProject",
		}, nil
	}

	response, err := agent.GetCaptenPluginsClient().GetGitProjectsForLabels(context.Background(), request)
	if err != nil {
		s.log.Errorf("failed to get the Cluster GitProject with lables, %v", err)
		return &captenpluginspb.GetGitProjectsForLabelsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster GitProject",
		}, nil
	}

	if response.Status != captenpluginspb.StatusCode_OK {
		s.log.Errorf("failed to get the ClusterProject with lables, %s", response.StatusMessage)
		return &captenpluginspb.GetGitProjectsForLabelsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster GitProject",
		}, nil
	}

	s.log.Infof("Fetched %d Git projects request with lables %v for cluster %s recieved, [org: %s]",
		request.Labels, len(response.Projects), clusterId, orgId)
	return &captenpluginspb.GetGitProjectsForLabelsResponse{
		Projects:      response.Projects,
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}
