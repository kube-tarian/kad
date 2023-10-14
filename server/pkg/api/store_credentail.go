package api

import (
	"context"

	"github.com/intelops/go-common/credentials"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) StoreCredential(ctx context.Context, request *serverpb.StoreCredentialRequest) (
	*serverpb.StoreCredentialResponse, error) {
	orgId, err := validateRequest(ctx, request.ClusterID)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.StoreCredentialResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Store Credential request recieved for cluster %s, [org: %s]", request.ClusterID, orgId)

	agent, err := s.agentHandeler.GetAgent(orgId, request.ClusterID)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &serverpb.StoreCredentialResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Credential store failed",
		}, nil
	}

	response, err := agent.GetClient().StoreCredential(context.Background(), &agentpb.StoreCredentialRequest{
		CredentialType: credentials.GenericCredentialType,
		CredEntityName: request.CredentialEntiryName,
		CredIdentifier: request.CredentialIdentifier,
		Credential:     request.Credential,
	})
	if err != nil {
		s.log.Errorf("failed to store credentials, %v", err)
		return &serverpb.StoreCredentialResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Credential store failed",
		}, nil
	}

	if response.Status != agentpb.StatusCode_OK {
		s.log.Errorf("failed to store credentials")
		return &serverpb.StoreCredentialResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Credential store failed",
		}, nil
	}

	s.log.Infof("Store Credential request for cluster %s successful, [org: %s]", request.ClusterID, orgId)
	return &serverpb.StoreCredentialResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "Credential store success",
	}, nil
}
