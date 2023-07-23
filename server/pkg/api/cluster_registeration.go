package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/agent"
	"github.com/kube-tarian/kad/server/pkg/credential"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) NewClusterRegistration(ctx context.Context, request *serverpb.NewClusterRegistrationRequest) (
	*serverpb.NewClusterRegistrationResponse, error) {
	orgId, ok := ctx.Value("organizationID").(string)
	if !ok || orgId == "" {
		s.log.Error("organizationID is missing in the request")
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "organizationID is missing",
		}, nil
	}

	agentConfig := &agent.Config{
		Address: request.AgentEndpoint,
		CaCert:  request.ClientCAChainData,
		Cert:    request.ClientKeyData,
		Key:     request.ClientCertData,
	}
	if err := s.agentHandeler.AddAgent(orgId, request.ClusterName, agentConfig); err != nil {
		s.log.Errorf("failed to connect to agent, %s", err)
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to connect to agent",
		}, nil
	}

	err := credential.PutClusterCerts(ctx, orgId, request.ClusterName,
		request.ClientCAChainData, request.ClientKeyData, request.ClientCertData)
	if err != nil {
		s.log.Errorf("failed to store cert in vault, %v", err)
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed register cluster",
		}, nil
	}

	err = s.serverStore.AddCluster(orgId, request.ClusterName, request.AgentEndpoint)
	if err != nil {
		s.log.Errorf("failed to get db session, %v", err)
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed register cluster",
		}, nil
	}

	return &serverpb.NewClusterRegistrationResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "register cluster success",
	}, nil
}

func (s *Server) UpdateClusterRegistration(ctx context.Context, request *serverpb.UpdateClusterRegistrationRequest) (
	*serverpb.UpdateClusterRegistrationResponse, error) {
	orgId, ok := ctx.Value("organizationID").(string)
	if !ok || orgId == "" {
		s.log.Error("organizationID is missing in the request")
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "organizationID is missing",
		}, nil
	}

	agentConfig := &agent.Config{
		Address: request.AgentEndpoint,
		CaCert:  request.ClientCAChainData,
		Cert:    request.ClientKeyData,
		Key:     request.ClientCertData,
	}

	if err := s.agentHandeler.UpdateAgent(orgId, request.ClusterName, agentConfig); err != nil {
		s.log.Errorf("failed to connect to agent, %s", err)
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to connect to agent",
		}, nil
	}

	err := credential.PutClusterCerts(ctx, orgId, request.ClusterName,
		request.ClientCAChainData, request.ClientKeyData, request.ClientCertData)
	if err != nil {
		s.log.Errorf("failed to store cert in vault, %v", err)
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed update register cluster",
		}, nil
	}

	err = s.serverStore.UpdateCluster(orgId, request.ClusterName, request.AgentEndpoint)
	if err != nil {
		s.log.Errorf("failed to get db session, %v", err)
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed update register cluster",
		}, nil
	}

	return &serverpb.UpdateClusterRegistrationResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "cluster register update success",
	}, nil
}

func (s *Server) DeleteClusterRegistration(ctx context.Context, request *serverpb.DeleteClusterRegistrationRequest) (
	*serverpb.DeleteClusterRegistrationResponse, error) {
	orgId, ok := ctx.Value("organizationID").(string)
	if !ok || orgId == "" {
		s.log.Error("organizationID is missing in the request")
		return &serverpb.DeleteClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "organizationID is missing",
		}, nil
	}

	s.agentHandeler.RemoveAgent(orgId, request.ClusterName)
	err := credential.DeleteClusterCerts(ctx, orgId, request.ClusterName)
	if err != nil {
		s.log.Errorf("failed to delete cert in vault, %v", err)
		return &serverpb.DeleteClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed delete register cluster",
		}, nil
	}

	err = s.serverStore.DeleteCluster(orgId, request.ClusterName)
	if err != nil {
		s.log.Errorf("failed to get db session, %v", err)
		return &serverpb.DeleteClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed delete register cluster",
		}, nil
	}

	return &serverpb.DeleteClusterRegistrationResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "cluster deletion success",
	}, nil
}

func (s *Server) GetClusters(ctx context.Context, request *serverpb.GetClustersRequest) (
	*serverpb.GetClustersResponse, error) {
	orgId, ok := ctx.Value("organizationID").(string)
	if !ok || orgId == "" {
		s.log.Error("organizationID is missing in the request")
		return &serverpb.GetClustersResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "organizationID is missing",
		}, nil
	}

	clusterDetails, err := s.serverStore.GetClusters(orgId)
	if err != nil {
		s.log.Errorf("failed to get cluster details, %v", err)
		return &serverpb.GetClustersResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed get cluster details",
		}, err
	}

	var data []*serverpb.ClusterInfo
	for _, cluster := range clusterDetails {
		data = append(data, &serverpb.ClusterInfo{
			ClusterName:   cluster.ClusterName,
			AgentEndpoint: cluster.Endpoint,
		})
	}

	return &serverpb.GetClustersResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "get cluster details success",
		Data:          data,
	}, nil
}
