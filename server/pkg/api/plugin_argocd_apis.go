package api

import (
	"context"
	"fmt"

	"github.com/kube-tarian/kad/server/pkg/pb/captenpluginspb"
)

func (s *Server) RegisterArgoCDProject(ctx context.Context, request *captenpluginspb.RegisterArgoCDProjectRequest) (
	*captenpluginspb.RegisterArgoCDProjectResponse, error) {

	metadataMap := metadataContextToMap(ctx)
	orgId, clusterId := metadataMap[organizationIDAttribute], metadataMap[clusterIDAttribute]
	if orgId == "" || clusterId == "" {
		s.log.Errorf("organizationid or clusterid is missing in the request")
		return &captenpluginspb.RegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "Organization Id or Cluster Id is missing",
		}, nil
	}

	if request.Id == "" {
		return &captenpluginspb.RegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Github Project Id is required",
		}, fmt.Errorf("Github Project Id is required")
	}

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.RegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to Register the ArgoCD Project",
		}, err
	}
	response, err := agent.GetCaptenPluginsClient().RegisterArgoCDProject(context.Background(), &captenpluginspb.RegisterArgoCDProjectRequest{Id: request.Id})
	if err != nil {
		s.log.Errorf("failed to Register the ArgoCD Project, %v", err)
		return &captenpluginspb.RegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to Register the ArgoCD Project",
		}, err
	}

	return &captenpluginspb.RegisterArgoCDProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: response.StatusMessage,
	}, nil
}

func (s *Server) UnRegisterArgoCDProject(ctx context.Context, request *captenpluginspb.UnRegisterArgoCDProjectRequest) (
	*captenpluginspb.UnRegisterArgoCDProjectResponse, error) {

	metadataMap := metadataContextToMap(ctx)
	orgId, clusterId := metadataMap[organizationIDAttribute], metadataMap[clusterIDAttribute]
	if orgId == "" || clusterId == "" {
		s.log.Errorf("organizationid or clusterid is missing in the request")
		return &captenpluginspb.UnRegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "Organization Id or Cluster Id is missing",
		}, nil
	}

	if request.Id == "" {
		return &captenpluginspb.UnRegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Github Project Id is required",
		}, fmt.Errorf("Github Project Id is required")
	}

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.UnRegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to Unregister the ArgoCD Project",
		}, err
	}
	response, err := agent.GetCaptenPluginsClient().UnRegisterArgoCDProject(context.Background(), &captenpluginspb.UnRegisterArgoCDProjectRequest{Id: request.Id})
	if err != nil {
		s.log.Errorf("failed to Unregister the ArgoCD Project, %v", err)
		return &captenpluginspb.UnRegisterArgoCDProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to Unregister the ArgoCD Project",
		}, err
	}

	return &captenpluginspb.UnRegisterArgoCDProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: response.StatusMessage,
	}, nil
}

func (s *Server) GetArgoCDProjects(ctx context.Context, request *captenpluginspb.GetArgoCDProjectsRequest) (
	*captenpluginspb.GetArgoCDProjectsResponse, error) {

	metadataMap := metadataContextToMap(ctx)
	orgId, clusterId := metadataMap[organizationIDAttribute], metadataMap[clusterIDAttribute]
	if orgId == "" || clusterId == "" {
		s.log.Errorf("organizationid or clusterid is missing in the request")
		return &captenpluginspb.GetArgoCDProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "Organization Id or Cluster Id is missing",
		}, nil
	}

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.GetArgoCDProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get ArgoCD Project",
		}, err
	}

	response, err := agent.GetCaptenPluginsClient().GetArgoCDProjects(context.Background(), &captenpluginspb.GetArgoCDProjectsRequest{})
	if err != nil {
		s.log.Errorf("failed to fetch ArgoCD projects, %v", err)
		return &captenpluginspb.GetArgoCDProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch ArgoCD projects",
		}, err
	}

	return &captenpluginspb.GetArgoCDProjectsResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: response.StatusMessage,
		Projects:      response.Projects,
	}, nil
}
