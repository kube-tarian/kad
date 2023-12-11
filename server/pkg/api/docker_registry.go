package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/captenpluginspb"
)

func (s *Server) AddDockerRegistry(ctx context.Context, request *captenpluginspb.AddDockerRegistryRequest) (
	*captenpluginspb.AddDockerRegistryResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.AddDockerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("add docker registry for cluster %s recieved, [org: %s]",
		clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.AddDockerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to initialize agent",
		}, err
	}

	resp, err := agent.GetCaptenPluginsClient().AddDockerRegistry(ctx, request)
	if err != nil {
		s.log.Errorf("failed to add docker registry, %v", err)
		return &captenpluginspb.AddDockerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add docker registry",
		}, err
	}

	s.log.Infof("docker regsitry creattion for cluster %s successful, [org: %s]",
		clusterId, orgId)
	return &captenpluginspb.AddDockerRegistryResponse{
		Status:        resp.Status,
		StatusMessage: "add docker registry successful",
	}, nil
}

func (s *Server) UpdateDockerRegistry(ctx context.Context, request *captenpluginspb.UpdateDockerRegistryRequest) (
	*captenpluginspb.UpdateDockerRegistryResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.UpdateDockerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("update docker registry for cluster %s recieved, [org: %s]",
		clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.UpdateDockerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to initialize agent",
		}, err
	}

	resp, err := agent.GetCaptenPluginsClient().UpdateDockerRegistry(ctx, request)
	if err != nil {
		s.log.Errorf("failed to update docker registry, %v", err)
		return &captenpluginspb.UpdateDockerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to update docker registry",
		}, err
	}

	s.log.Infof("docker regsitry update for cluster %s successful, [org: %s]",
		clusterId, orgId)
	return &captenpluginspb.UpdateDockerRegistryResponse{
		Status:        resp.Status,
		StatusMessage: "update docker registry successful",
	}, nil
}
func (s *Server) DeleteDockerRegistry(ctx context.Context, request *captenpluginspb.DeleteDockerRegistryRequest) (
	*captenpluginspb.DeleteDockerRegistryResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.DeleteDockerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("delete docker registry for cluster %s recieved, [org: %s]",
		clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.DeleteDockerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to initialize agent",
		}, err
	}

	resp, err := agent.GetCaptenPluginsClient().DeleteDockerRegistry(ctx, request)
	if err != nil {
		s.log.Errorf("failed to delete docker registry, %v", err)
		return &captenpluginspb.DeleteDockerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete docker registry",
		}, err
	}

	s.log.Infof("docker regsitry delete for cluster %s successful, [org: %s]",
		clusterId, orgId)
	return &captenpluginspb.DeleteDockerRegistryResponse{
		Status:        resp.Status,
		StatusMessage: "delete docker registry successful",
	}, nil

}
func (s *Server) GetDockerRegistry(ctx context.Context, request *captenpluginspb.GetDockerRegistryRequest) (
	*captenpluginspb.GetDockerRegistryResponse, error) {

	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.GetDockerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("get docker registry for cluster %s recieved, [org: %s]",
		clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.GetDockerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to initialize agent",
		}, err
	}

	resp, err := agent.GetCaptenPluginsClient().GetDockerRegistry(ctx, request)
	if err != nil {
		s.log.Errorf("failed to get docker registry, %v", err)
		return &captenpluginspb.GetDockerRegistryResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to update docker registry",
		}, err
	}

	s.log.Infof("docker regsitry get for cluster %s successful, [org: %s]",
		clusterId, orgId)
	return &captenpluginspb.GetDockerRegistryResponse{
		Registries:    resp.Registries,
		Status:        resp.Status,
		StatusMessage: "update docker registry successful",
	}, nil
}
