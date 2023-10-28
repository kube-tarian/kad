package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/captenpluginspb"
)

func (s *Server) AddCloudProvider(ctx context.Context, request *captenpluginspb.AddCloudProviderRequest) (
	*captenpluginspb.AddCloudProviderResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx, request.CloudType)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.AddCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	s.log.Infof("Add Cloud Provider %s request for cluster %s recieved, [org: %s]",
		request.CloudType, clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent for cluster %s, %v", clusterId, err)
		return &captenpluginspb.AddCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add the Cluster CloudProvider",
		}, nil
	}

	response, err := agent.GetCaptenPluginsClient().AddCloudProvider(context.Background(), request)
	if err != nil {
		s.log.Errorf("failed to add the Cluster CloudProvider for cluster %s, %v", clusterId, err)
		return &captenpluginspb.AddCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add the Cluster CloudProvider",
		}, nil
	}

	if response.Status != captenpluginspb.StatusCode_OK {
		s.log.Errorf("failed to add the ClusterProject for cluster %s, %s, %s", response.Status, response.StatusMessage)
		return &captenpluginspb.AddCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add the Cluster CloudProvider",
		}, nil
	}

	s.log.Infof("Cloud Provider %s request for cluster %s added, [org: %s]",
		request.CloudType, clusterId, orgId)
	return &captenpluginspb.AddCloudProviderResponse{
		Id:            response.GetId(),
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil

}

func (s *Server) UpdateCloudProviders(ctx context.Context, request *captenpluginspb.UpdateCloudProviderRequest) (
	*captenpluginspb.UpdateCloudProviderResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx, request.Id, request.CloudType)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.UpdateCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	s.log.Infof("Update Cloud Provider %s, %s request for cluster %s recieved, [org: %s]",
		request.CloudType, request.Id, clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent for cluster %s, %v", clusterId, err)
		return &captenpluginspb.UpdateCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add the Cluster CloudProvider",
		}, nil
	}

	response, err := agent.GetCaptenPluginsClient().UpdateCloudProvider(context.Background(), request)
	if err != nil {
		s.log.Errorf("failed to add the Cluster CloudProvider for cluster %s, %v", clusterId, err)
		return &captenpluginspb.UpdateCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add the Cluster CloudProvider",
		}, nil
	}

	if response.Status != captenpluginspb.StatusCode_OK {
		s.log.Errorf("failed to update the ClusterProject for cluster %s, %s, %s", response.Status, response.StatusMessage)
		return &captenpluginspb.UpdateCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to update the Cluster CloudProvider",
		}, nil
	}

	s.log.Infof("Cloud Provider %s, %s request for cluster %s updated, [org: %s]",
		request.CloudType, request.Id, clusterId, orgId)
	return &captenpluginspb.UpdateCloudProviderResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (s *Server) DeleteCloudProvider(ctx context.Context, request *captenpluginspb.DeleteCloudProviderRequest) (
	*captenpluginspb.DeleteCloudProviderResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx, request.Id)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.DeleteCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	s.log.Infof("Delete Cloud Provider %s request for cluster %s recieved, [org: %s]",
		request.Id, clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent for cluster %s, %v", clusterId, err)
		return &captenpluginspb.DeleteCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete the Cluster CloudProvider",
		}, nil
	}

	response, err := agent.GetCaptenPluginsClient().DeleteCloudProvider(context.Background(), request)
	if err != nil {
		s.log.Errorf("failed to delete the Cluster CloudProvider for cluster %s, %v", clusterId, err)
		return &captenpluginspb.DeleteCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete the Cluster CloudProvider",
		}, nil
	}

	if response.Status != captenpluginspb.StatusCode_OK {
		s.log.Errorf("failed to update the ClusterProject for cluster %s, %s, %s", response.Status, response.StatusMessage)
		return &captenpluginspb.DeleteCloudProviderResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete the Cluster CloudProvider",
		}, nil
	}

	return &captenpluginspb.DeleteCloudProviderResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (s *Server) GetCloudProviders(ctx context.Context, request *captenpluginspb.GetCloudProvidersRequest) (
	*captenpluginspb.GetCloudProvidersResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.GetCloudProvidersResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Get Cloud Providers request for cluster %s recieved, [org: %s]",
		clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.GetCloudProvidersResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster CloudProvider",
		}, nil
	}

	response, err := agent.GetCaptenPluginsClient().GetCloudProviders(context.Background(), request)
	if err != nil {
		s.log.Errorf("failed to get the Cluster CloudProvider, %v", err)
		return &captenpluginspb.GetCloudProvidersResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster CloudProvider",
		}, nil
	}

	if response.Status != captenpluginspb.StatusCode_OK {
		s.log.Errorf("failed to get the ClusterProject")
		return &captenpluginspb.GetCloudProvidersResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster CloudProvider",
		}, nil
	}

	s.log.Infof("Fetched %d Cloud Providers request for cluster %s processed, [org: %s]",
		len(response.GetCloudProviders()), clusterId, orgId)
	return &captenpluginspb.GetCloudProvidersResponse{
		CloudProviders: response.GetCloudProviders(),
		Status:         captenpluginspb.StatusCode_OK,
		StatusMessage:  "ok",
	}, nil
}

func (s *Server) GetCloudProvidersForLabels(ctx context.Context, request *captenpluginspb.GetCloudProvidersWithFilterRequest) (
	*captenpluginspb.GetCloudProvidersWithFilterResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &captenpluginspb.GetCloudProvidersWithFilterResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Get Cloud Providers request with lables %v for cluster %s recieved, [org: %s]",
		request.Labels, clusterId, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, clusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &captenpluginspb.GetCloudProvidersWithFilterResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster CloudProvider",
		}, nil
	}

	response, err := agent.GetCaptenPluginsClient().GetCloudProvidersWithFilter(context.Background(), request)
	if err != nil {
		s.log.Errorf("failed to get the Cluster CloudProvider with lables, %v", err)
		return &captenpluginspb.GetCloudProvidersWithFilterResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster CloudProvider",
		}, nil
	}

	if response.Status != captenpluginspb.StatusCode_OK {
		s.log.Errorf("failed to get the ClusterProject with lables")
		return &captenpluginspb.GetCloudProvidersWithFilterResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to get the Cluster CloudProvider",
		}, nil
	}

	s.log.Infof("Fetched %d Cloud Providers request with lables %v for cluster %s recieved, [org: %s]",
		request.Labels, len(response.GetCloudProviders()), clusterId, orgId)
	return &captenpluginspb.GetCloudProvidersWithFilterResponse{
		CloudProviders: response.GetCloudProviders(),
		Status:         captenpluginspb.StatusCode_OK,
		StatusMessage:  "ok",
	}, nil
}
