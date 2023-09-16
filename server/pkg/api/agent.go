package api

import (
	"context"

	"github.com/intelops/go-common/credentials"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) StoreCredential(ctx context.Context, request *serverpb.StoreCredentialRequest) (
	*serverpb.StoreCredentialResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if orgId == "" {
		s.log.Errorf("organization ID is missing in the request")
		return &serverpb.StoreCredentialResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Organization Id is missing",
		}, nil
	}

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
		CredIdentifier: request.CredIdentifier,
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

	return &serverpb.StoreCredentialResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "Credential store success",
	}, nil
}
