package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/agent"
	"github.com/kube-tarian/kad/server/pkg/credential"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) UpdateClusterRegistration(ctx context.Context, request *serverpb.UpdateClusterRegistrationRequest) (
	*serverpb.UpdateClusterRegistrationResponse, error) {
	orgId, err := validateOrgWithArgs(ctx, request.ClusterID)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Update cluster registration request for cluster %s recieved, [org: %s]", request.ClusterName, orgId)

	caData, caDataErr := getBase64DecodedString(request.ClientCAChainData)
	clientKey, clientKeyErr := getBase64DecodedString(request.ClientKeyData)
	clientCrt, clientCrtErr := getBase64DecodedString(request.ClientCertData)
	if caDataErr != nil || clientKeyErr != nil || clientCrtErr != nil {
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "only base64 encoded certificates are allowed",
		}, nil
	}

	agentConfig := &agent.Config{
		ClusterName: request.ClusterName,
		Address:     request.AgentEndpoint,
		CaCert:      caData,
		Key:         clientKey,
		Cert:        clientCrt,
	}

	if err := s.agentHandeler.UpdateAgent(request.ClusterID, agentConfig); err != nil {
		s.log.Errorf("[org: %s] failed to connect to agent on cluster %s, %v", orgId, request.ClusterName, err)
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to connect to agent",
		}, nil
	}

	err = credential.PutClusterCerts(ctx, request.ClusterID, caData, clientKey, clientCrt)
	if err != nil {
		s.log.Errorf("[org: %s] failed to update cert in vault for cluster %s, %v", orgId, request.ClusterID, err)
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed update register cluster",
		}, nil
	}

	err = s.serverStore.UpdateCluster(orgId, request.ClusterID, request.ClusterName, request.AgentEndpoint)
	if err != nil {
		s.log.Errorf("[org: %s] failed to update cluster %s in db, %v", orgId, request.ClusterName, err)
		return &serverpb.UpdateClusterRegistrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed update register cluster",
		}, nil
	}

	s.log.Infof("Update cluster registration request for cluster %s successful, [org: %s]", request.ClusterName, orgId)
	return &serverpb.UpdateClusterRegistrationResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "cluster register update success",
	}, nil
}
