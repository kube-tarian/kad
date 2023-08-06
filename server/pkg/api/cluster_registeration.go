package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/agent"
	"github.com/kube-tarian/kad/server/pkg/credential"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) NewClusterRegistration(ctx context.Context, request *serverpb.NewClusterRegistrationRequest) (
	*serverpb.NewClusterRegistrationResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if len(orgId) == 0 {
		s.log.Error("organizationID is missing in the request")
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "organizationID is missing",
		}, nil
	}

	s.log.Infof("[%s] New cluster registration request for cluster %s recieved", orgId, request.ClusterName)
	agentConfig := &agent.Config{
		Address: request.AgentEndpoint,
		CaCert:  request.ClientCAChainData,
		Cert:    request.ClientKeyData,
		Key:     request.ClientCertData,
	}
	if err := s.agentHandeler.AddAgent(orgId, request.ClusterName, agentConfig); err != nil {
		s.log.Errorf("[%s] failed to connect to agent on cluster %s, %v", orgId, request.ClusterName, err)
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to connect to agent",
		}, nil
	}

	err := credential.PutClusterCerts(ctx, orgId, request.ClusterName,
		request.ClientCAChainData, request.ClientKeyData, request.ClientCertData)
	if err != nil {
		s.log.Errorf("[%s] failed to store cert in vault for cluster %s, %v", orgId, request.ClusterName, err)
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed register cluster",
		}, nil
	}

	clusterID, err := s.serverStore.AddCluster(orgId, request.ClusterName, request.AgentEndpoint)
	if err != nil {
		s.log.Errorf("[%s] failed to store cluster %s to db, %v", orgId, request.ClusterName, err)
		return &serverpb.NewClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed register cluster",
		}, nil
	}

	s.log.Infof("[%s] New cluster registration successful for %s cluster", orgId, request.ClusterName)
	return &serverpb.NewClusterRegistrationResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "register cluster success",
		ClusterID:     clusterID,
	}, nil
}

func (s *Server) UpdateClusterRegistration(ctx context.Context, request *serverpb.UpdateClusterRegistrationRequest) (
	*serverpb.UpdateClusterRegistrationResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if len(orgId) == 0 {
		s.log.Error("organizationID is missing in the request")
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "organizationID is missing",
		}, nil
	}

	s.log.Infof("[%s] Update cluster registration request for cluster %s recieved", orgId, request.ClusterName)
	agentConfig := &agent.Config{
		Address: request.AgentEndpoint,
		CaCert:  request.ClientCAChainData,
		Cert:    request.ClientKeyData,
		Key:     request.ClientCertData,
	}

	if err := s.agentHandeler.UpdateAgent(orgId, request.ClusterID, agentConfig); err != nil {
		s.log.Errorf("[%s] failed to connect to agent on cluster %s, %v", orgId, request.ClusterName, err)
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to connect to agent",
		}, nil
	}

	err := credential.PutClusterCerts(ctx, orgId, request.ClusterName,
		request.ClientCAChainData, request.ClientKeyData, request.ClientCertData)
	if err != nil {
		s.log.Errorf("[%s] failed to update cert in vault for cluster %s, %v", orgId, request.ClusterName, err)
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed update register cluster",
		}, nil
	}

	err = s.serverStore.UpdateCluster(orgId, request.ClusterID, request.ClusterName, request.AgentEndpoint)
	if err != nil {
		s.log.Errorf("[%s] failed to update cluster %s in db, %v", orgId, request.ClusterName, err)
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed update register cluster",
		}, nil
	}

	s.log.Infof("[%s] Update cluster registration successful for %s cluster", orgId, request.ClusterName)
	return &serverpb.UpdateClusterRegistrationResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "cluster register update success",
	}, nil
}

func (s *Server) DeleteClusterRegistration(ctx context.Context, request *serverpb.DeleteClusterRegistrationRequest) (
	*serverpb.DeleteClusterRegistrationResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if len(orgId) == 0 {
		s.log.Error("organizationID is missing in the request")
		return &serverpb.DeleteClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "organizationID is missing",
		}, nil
	}

	s.log.Infof("[%s] Delete cluster registration request for cluster %s recieved", orgId, request.ClusterID)
	s.agentHandeler.RemoveAgent(orgId, request.ClusterID)
	err := credential.DeleteClusterCerts(ctx, orgId, request.ClusterID)
	if err != nil {
		s.log.Errorf("[%s] failed to delete cert in vault for cluster %s, %v", orgId, request.ClusterID, err)
		return &serverpb.DeleteClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed delete register cluster",
		}, nil
	}

	err = s.serverStore.DeleteCluster(orgId, request.ClusterID)
	if err != nil {
		s.log.Errorf("[%s] failed to delete cluster %s from db, %v", orgId, request.ClusterID, err)
		return &serverpb.DeleteClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed delete register cluster",
		}, nil
	}

	s.log.Infof("[%s] Delete cluster registration request for cluster %s successful", orgId, request.ClusterID)
	return &serverpb.DeleteClusterRegistrationResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "cluster deletion success",
	}, nil
}

func (s *Server) GetClusters(ctx context.Context, request *serverpb.GetClustersRequest) (
	*serverpb.GetClustersResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if len(orgId) == 0 {
		s.log.Error("organizationID is missing in the request")
		return &serverpb.GetClustersResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "organizationID is missing",
		}, nil
	}

	s.log.Infof("[%s] GetClusters request recieved", orgId)
	clusterDetails, err := s.serverStore.GetClusters(orgId)
	if err != nil {
		s.log.Errorf("[%s] failed to get clusters, %v", orgId, err)
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

	s.log.Infof("[%s] Found %d clusters", orgId, len(data))
	return &serverpb.GetClustersResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "get cluster details success",
		Data:          data,
	}, nil
}
