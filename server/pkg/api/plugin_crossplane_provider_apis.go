package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/captenpluginspb"
)

func (s *Server) AddCrossplanProvider(ctx context.Context, request *captenpluginspb.AddCrossplanProviderRequest) (
	*captenpluginspb.AddCrossplanProviderResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx, request.CloudProviderId, request.CloudType)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.AddCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	s.log.Infof("Add Crosssplane Provider cloudId: %s request for cluster %s recieved, [org: %s]",
		request.CloudProviderId, clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent for cluster %s, %v", clusterId, err)
		return &captenpluginspb.AddCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add the Cluster CrossplanProvider",
		}, nil
	}

	response, err := agent.GetCaptenPluginsClient().AddCrossplanProvider(context.Background(), request)
	if err != nil {
		s.log.Errorf("failed to add the Cluster CrossplanProvider for cluster %s, %v", clusterId, err)
		return &captenpluginspb.AddCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add the Cluster CrossplanProvider",
		}, nil
	}

	if response.Status != captenpluginspb.StatusCode_OK {
		s.log.Errorf("failed to add the CrossplanProvider for cluster %s, %s, %s", response.Status, response.StatusMessage)
		return &captenpluginspb.AddCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add the Cluster CrossplanProvider",
		}, nil
	}

	s.log.Infof("Crossplane Provider cloudId: %s request for cluster %s added, [org: %s]",
		request.CloudProviderId, clusterId, orgId)
	return &captenpluginspb.AddCrossplanProviderResponse{
		Id:            response.GetId(),
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil

}

func (s *Server) UpdateCrossplanProvider(ctx context.Context, request *captenpluginspb.UpdateCrossplanProviderRequest) (
	*captenpluginspb.UpdateCrossplanProviderResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx, request.Id, request.CloudType, request.CloudProviderId)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.UpdateCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	s.log.Infof("Update Crossplane provider cloudId: %s, %s request for cluster %s recieved, [org: %s]",
		request.CloudProviderId, request.CloudType, clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent for cluster %s, %v", clusterId, err)
		return &captenpluginspb.UpdateCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add the Cluster CrossplanProvider",
		}, nil
	}

	response, err := agent.GetCaptenPluginsClient().UpdateCrossplanProvider(context.Background(), request)
	if err != nil {
		s.log.Errorf("failed to add the Cluster CrossplanProvider for cluster %s, %v", clusterId, err)
		return &captenpluginspb.UpdateCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add the Cluster CrossplanProvider",
		}, nil
	}

	if response.Status != captenpluginspb.StatusCode_OK {
		s.log.Errorf("failed to update the CrossplanProvider for cluster %s, %s, %s", response.Status, response.StatusMessage)
		return &captenpluginspb.UpdateCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to update the Cluster CrossplanProvider",
		}, nil
	}

	s.log.Infof("Crossplane Provider cloudId: %s, %s request for cluster %s updated, [org: %s]",
		request.CloudProviderId, request.CloudType, clusterId, orgId)
	return &captenpluginspb.UpdateCrossplanProviderResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (s *Server) DeleteCrossplanProvider(ctx context.Context, request *captenpluginspb.DeleteCrossplanProviderRequest) (
	*captenpluginspb.DeleteCrossplanProviderResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx, request.Id)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.DeleteCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	s.log.Infof("Delete Crossplane Provider %s request for cluster %s recieved, [org: %s]",
		request.Id, clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent for cluster %s, %v", clusterId, err)
		return &captenpluginspb.DeleteCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete the Cluster CrossplanProvider",
		}, nil
	}

	response, err := agent.GetCaptenPluginsClient().DeleteCrossplanProvider(context.Background(), request)
	if err != nil {
		s.log.Errorf("failed to delete the Cluster CrossplanProvider for cluster %s, %v", clusterId, err)
		return &captenpluginspb.DeleteCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete the Cluster CrossplanProvider",
		}, nil
	}

	if response.Status != captenpluginspb.StatusCode_OK {
		s.log.Errorf("failed to update the CrossplanProvider for cluster %s, %s, %s", response.Status, response.StatusMessage)
		return &captenpluginspb.DeleteCrossplanProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete the Cluster CrossplanProvider",
		}, nil
	}

	return &captenpluginspb.DeleteCrossplanProviderResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (s *Server) GetCrossplanProviders(ctx context.Context, request *captenpluginspb.GetCrossplanProvidersRequest) (
	*captenpluginspb.GetCrossplanProvidersResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.GetCrossplanProvidersResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Get Crossplan Providers request for cluster %s recieved, [org: %s]",
		clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.GetCrossplanProvidersResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster CrossplanProvider",
		}, nil
	}

	response, err := agent.GetCaptenPluginsClient().GetCrossplanProviders(context.Background(), request)
	if err != nil {
		s.log.Errorf("failed to get the Cluster Crossplane Providers, %v", err)
		return &captenpluginspb.GetCrossplanProvidersResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster CrossplanProvider",
		}, nil
	}

	if response.Status == captenpluginspb.StatusCode_NOT_FOUND {
		response.Providers = []*captenpluginspb.CrossplaneProvider{}

	} else if response.Status != captenpluginspb.StatusCode_OK {
		s.log.Errorf("failed to get the Crossplane Providers")
		return &captenpluginspb.GetCrossplanProvidersResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster CrossplanProvider",
		}, nil
	}

	s.log.Infof("Fetched %d Crossplane Providers for cluster %s processed, [org: %s]",
		len(response.Providers), clusterId, orgId)
	return &captenpluginspb.GetCrossplanProvidersResponse{
		Providers:     response.Providers,
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}
