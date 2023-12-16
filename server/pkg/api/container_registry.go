package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/captenpluginspb"
)

func (s *Server) AddContainerRegistry(ctx context.Context, request *captenpluginspb.AddContainerRegistryRequest) (
	*captenpluginspb.AddContainerRegistryResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.AddContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("add Container registry for cluster %s recieved, [org: %s]",
		clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.AddContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to initialize agent",
		}, err
	}

	resp, err := agent.GetCaptenPluginsClient().AddContainerRegistry(ctx, request)
	if err != nil {
		s.log.Errorf("failed to add Container registry, %v", err)
		return &captenpluginspb.AddContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add Container registry",
		}, err
	}

	s.log.Infof("Container regsitry creattion for cluster %s successful, [org: %s]",
		clusterId, orgId)
	return &captenpluginspb.AddContainerRegistryResponse{
		Status:        resp.Status,
		StatusMessage: "add Container registry successful",
	}, nil
}

func (s *Server) UpdateContainerRegistry(ctx context.Context, request *captenpluginspb.UpdateContainerRegistryRequest) (
	*captenpluginspb.UpdateContainerRegistryResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.UpdateContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("update Container registry for cluster %s recieved, [org: %s]",
		clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.UpdateContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to initialize agent",
		}, err
	}

	resp, err := agent.GetCaptenPluginsClient().UpdateContainerRegistry(ctx, request)
	if err != nil {
		s.log.Errorf("failed to update Container registry, %v", err)
		return &captenpluginspb.UpdateContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to update Container registry",
		}, err
	}

	s.log.Infof("Container regsitry update for cluster %s successful, [org: %s]",
		clusterId, orgId)
	return &captenpluginspb.UpdateContainerRegistryResponse{
		Status:        resp.Status,
		StatusMessage: "update Container registry successful",
	}, nil
}
func (s *Server) DeleteContainerRegistry(ctx context.Context, request *captenpluginspb.DeleteContainerRegistryRequest) (
	*captenpluginspb.DeleteContainerRegistryResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.DeleteContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("delete Container registry for cluster %s recieved, [org: %s]",
		clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.DeleteContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to initialize agent",
		}, err
	}

	resp, err := agent.GetCaptenPluginsClient().DeleteContainerRegistry(ctx, request)
	if err != nil {
		s.log.Errorf("failed to delete Container registry, %v", err)
		return &captenpluginspb.DeleteContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete Container registry",
		}, err
	}

	s.log.Infof("Container regsitry delete for cluster %s successful, [org: %s]",
		clusterId, orgId)
	return &captenpluginspb.DeleteContainerRegistryResponse{
		Status:        resp.Status,
		StatusMessage: "delete Container registry successful",
	}, nil

}
func (s *Server) GetContainerRegistry(ctx context.Context, request *captenpluginspb.GetContainerRegistryRequest) (
	*captenpluginspb.GetContainerRegistryResponse, error) {

	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.GetContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("get Container registry for cluster %s recieved, [org: %s]",
		clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.GetContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to initialize agent",
		}, err
	}

	resp, err := agent.GetCaptenPluginsClient().GetContainerRegistry(ctx, request)
	if err != nil {
		s.log.Errorf("failed to get Container registry, %v", err)
		return &captenpluginspb.GetContainerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to update Container registry",
		}, err
	}

	s.log.Infof("Container regsitry get for cluster %s successful, [org: %s]",
		clusterId, orgId)
	return &captenpluginspb.GetContainerRegistryResponse{
		Registries:    resp.Registries,
		Status:        resp.Status,
		StatusMessage: "update Container registry successful",
	}, nil
}
