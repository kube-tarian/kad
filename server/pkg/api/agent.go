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

	return &serverpb.StoreCredentialResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "Credential store success",
	}, nil
}

func (s *Server) AddOrUpdateOnboarding(ctx context.Context, request *serverpb.AddOrUpdateOnboardingRequest) (
	*serverpb.AddOrUpdateOnboardingResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if orgId == "" {
		s.log.Errorf("organization ID is missing in the request")
		return &serverpb.AddOrUpdateOnboardingResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Organization Id is missing",
		}, nil
	}

	agent, err := s.agentHandeler.GetAgent(orgId, request.ClusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &serverpb.AddOrUpdateOnboardingResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Onboarding Add/Update failed",
		}, nil
	}

	response, err := agent.GetClient().AddOrUpdateOnboarding(context.Background(), &agentpb.AddOrUpdateOnboardingRequest{
		Type:       request.Type,
		ProjectUrl: request.ProjectUrl,
		Status:     "started",
		Details:    request.Details,
	})
	if err != nil {
		s.log.Errorf("failed to App/Update onboarding, %v", err)
		return &serverpb.AddOrUpdateOnboardingResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Onboarding Add/Update failed",
		}, nil
	}

	if response.Status != agentpb.StatusCode_OK {
		s.log.Errorf("failed to Add/Update onboarding")
		return &serverpb.AddOrUpdateOnboardingResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Onboarding Add/Update failed",
		}, nil
	}

	return &serverpb.AddOrUpdateOnboardingResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "Add/Update of onboarding success",
	}, nil
}

func (s *Server) GetOnboarding(ctx context.Context, request *serverpb.GetOnboardingRequest) (
	*serverpb.GetOnboardingResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if orgId == "" {
		s.log.Errorf("organization ID is missing in the request")
		return &serverpb.GetOnboardingResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Organization Id is missing",
		}, nil
	}

	agent, err := s.agentHandeler.GetAgent(orgId, request.ClusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &serverpb.GetOnboardingResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get the onboarding",
		}, nil
	}

	response, err := agent.GetClient().GetOnboarding(context.Background(), &agentpb.GetOnboardingRequest{
		Type:       request.Type,
		ProjectUrl: request.ProjectUrl,
	})
	if err != nil {
		s.log.Errorf("failed to get the onboarding, %v", err)
		return &serverpb.GetOnboardingResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get the onboarding",
		}, nil
	}

	if response.Status != agentpb.StatusCode_OK {
		s.log.Errorf("failed to get the onboarding")
		return &serverpb.GetOnboardingResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get the onboarding",
		}, nil
	}

	return &serverpb.GetOnboardingResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "Successfully fetched the onboarding",
	}, nil
}

func (s *Server) DeleteOnboarding(ctx context.Context, request *serverpb.DeleteOnboardingRequest) (
	*serverpb.DeleteOnboardingResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if orgId == "" {
		s.log.Errorf("organization ID is missing in the request")
		return &serverpb.DeleteOnboardingResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Organization Id is missing",
		}, nil
	}

	agent, err := s.agentHandeler.GetAgent(orgId, request.ClusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &serverpb.DeleteOnboardingResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Onboarding delation failed",
		}, nil
	}

	response, err := agent.GetClient().DeleteOnboarding(context.Background(), &agentpb.DeleteOnboardingRequest{
		Type:       request.Type,
		ProjectUrl: request.ProjectUrl,
	})
	if err != nil {
		s.log.Errorf("failed to deletion onboarding, %v", err)
		return &serverpb.DeleteOnboardingResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Onboarding deletion failed",
		}, nil
	}

	if response.Status != agentpb.StatusCode_OK {
		s.log.Errorf("failed to delete onboarding")
		return &serverpb.DeleteOnboardingResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Onboarding delete failed",
		}, nil
	}

	return &serverpb.DeleteOnboardingResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "Successfully deleted the onboarding",
	}, nil
}
