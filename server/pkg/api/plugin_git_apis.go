package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/captenpluginspb"
)

func (s *Server) AddGitProject(ctx context.Context, request *captenpluginspb.AddGitProjectRequest) (
	*captenpluginspb.AddGitProjectResponse, error) {

	metadataMap := metadataContextToMap(ctx)
	orgId, clusterId := metadataMap[organizationIDAttribute], metadataMap[clusterIDAttribute]
	if orgId == "" || clusterId == "" {
		s.log.Errorf("organization or cluster ID is missing in the request")
		return &captenpluginspb.AddGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "Organization or Cluster Id is missing",
		}, nil
	}

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.AddGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add the Cluster GitProject",
		}, nil
	}

	response, err := agent.GetCaptenPluginClient().AddGitProject(context.Background(), request)
	if err != nil {
		s.log.Errorf("failed to add the Cluster GitProject, %v", err)
		return &captenpluginspb.AddGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add the Cluster GitProject",
		}, nil
	}

	if response.Status != captenpluginspb.StatusCode_OK {
		s.log.Errorf("failed to add the ClusterProject")
		return &captenpluginspb.AddGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add the Cluster GitProject",
		}, nil
	}

	return &captenpluginspb.AddGitProjectResponse{
		Id:            response.GetId(),
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil

}

func (s *Server) UpdateGitProject(ctx context.Context, request *captenpluginspb.UpdateGitProjectRequest) (
	*captenpluginspb.UpdateGitProjectResponse, error) {

	metadataMap := metadataContextToMap(ctx)
	orgId, clusterId := metadataMap[organizationIDAttribute], metadataMap[clusterIDAttribute]
	if orgId == "" || clusterId == "" {
		s.log.Errorf("organization or cluster ID is missing in the request")
		return &captenpluginspb.UpdateGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "Organization or Cluster Id is missing",
		}, nil
	}

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.UpdateGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to update the Cluster GitProject",
		}, nil
	}

	response, err := agent.GetCaptenPluginClient().UpdateGitProject(context.Background(), request)
	if err != nil {
		s.log.Errorf("failed to update the Cluster GitProject, %v", err)
		return &captenpluginspb.UpdateGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to update the Cluster GitProject",
		}, nil
	}

	if response.Status != captenpluginspb.StatusCode_OK {
		s.log.Errorf("failed to update the ClusterProject")
		return &captenpluginspb.UpdateGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to update the Cluster GitProject",
		}, nil
	}

	return &captenpluginspb.UpdateGitProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (s *Server) DeleteGitProject(ctx context.Context, request *captenpluginspb.DeleteGitProjectRequest) (
	*captenpluginspb.DeleteGitProjectResponse, error) {

	metadataMap := metadataContextToMap(ctx)
	orgId, clusterId := metadataMap[organizationIDAttribute], metadataMap[clusterIDAttribute]
	if orgId == "" || clusterId == "" {
		s.log.Errorf("organization or cluster ID is missing in the request")
		return &captenpluginspb.DeleteGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "Organization or Cluster Id is missing",
		}, nil
	}

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.DeleteGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete the Cluster GitProject",
		}, nil
	}

	response, err := agent.GetCaptenPluginClient().DeleteGitProject(context.Background(), request)
	if err != nil {
		s.log.Errorf("failed to add the Cluster GitProject, %v", err)
		return &captenpluginspb.DeleteGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete the Cluster GitProject",
		}, nil
	}

	if response.Status != captenpluginspb.StatusCode_OK {
		s.log.Errorf("failed to delete the ClusterProject")
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

	metadataMap := metadataContextToMap(ctx)
	orgId, clusterId := metadataMap[organizationIDAttribute], metadataMap[clusterIDAttribute]
	if orgId == "" || clusterId == "" {
		s.log.Errorf("organization or cluster ID is missing in the request")
		return &captenpluginspb.GetGitProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "Organization or Cluister Id is missing",
		}, nil
	}

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.GetGitProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster GitProject",
		}, nil
	}

	response, err := agent.GetCaptenPluginClient().GetGitProjects(context.Background(), request)
	if err != nil {
		s.log.Errorf("failed to get the Cluster GitProject, %v", err)
		return &captenpluginspb.GetGitProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster GitProject",
		}, nil
	}

	if response.Status != captenpluginspb.StatusCode_OK {
		s.log.Errorf("failed to get the ClusterProject")
		return &captenpluginspb.GetGitProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster GitProject",
		}, nil
	}

	return &captenpluginspb.GetGitProjectsResponse{
		GitProjects:   response.GetGitProjects(),
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (s *Server) GetGitProjectsForLabels(ctx context.Context, request *captenpluginspb.GetGitProjectsForLabelsRequest) (
	*captenpluginspb.GetGitProjectsForLabelsResponse, error) {

	metadataMap := metadataContextToMap(ctx)
	orgId, clusterId := metadataMap[organizationIDAttribute], metadataMap[clusterIDAttribute]
	if orgId == "" || clusterId == "" {
		s.log.Errorf("organization or cluster ID is missing in the request")
		return &captenpluginspb.GetGitProjectsForLabelsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "Organization or Cluister Id is missing",
		}, nil
	}

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.GetGitProjectsForLabelsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster GitProject",
		}, nil
	}

	response, err := agent.GetCaptenPluginClient().GetGitProjectsForLabels(context.Background(), request)
	if err != nil {
		s.log.Errorf("failed to get the Cluster GitProject, %v", err)
		return &captenpluginspb.GetGitProjectsForLabelsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster GitProject",
		}, nil
	}

	if response.Status != captenpluginspb.StatusCode_OK {
		s.log.Errorf("failed to get the ClusterProject")
		return &captenpluginspb.GetGitProjectsForLabelsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster GitProject",
		}, nil
	}

	return &captenpluginspb.GetGitProjectsForLabelsResponse{
		Projects:      response.GetProjects(),
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}
