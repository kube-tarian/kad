package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kube-tarian/kad/server/pkg/pb/captenpluginspb"
)

func (s *Server) RegisterArgoCDProject(ctx context.Context, request *captenpluginspb.RegisterArgoCDProjectRequest) (
	*captenpluginspb.RegisterArgoCDProjectResponse, error) {
	return &captenpluginspb.RegisterArgoCDProjectResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}

func (s *Server) UnRegisterArgoCDProject(ctx context.Context, request *captenpluginspb.UnRegisterArgoCDProjectRequest) (
	*captenpluginspb.UnRegisterArgoCDProjectResponse, error) {
	return &captenpluginspb.UnRegisterArgoCDProjectResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
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
			StatusMessage: "failed to add the Cluster GitProject",
		}, err
	}

	response, err := agent.GetCaptenPluginsClient().GetArgoCDProjects(context.Background(), &captenpluginspb.GetArgoCDProjectsRequest{})
	if err != nil {
		s.log.Errorf("failed to fetch Argocd projects, %v", err)
		return &captenpluginspb.GetArgoCDProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch ArgoCD projects",
		}, err
	}
	v, _ := json.Marshal(response)
	fmt.Println("Response =>" + string(v))

	return &captenpluginspb.GetArgoCDProjectsResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
		Projects:      nil,
	}, nil
}
