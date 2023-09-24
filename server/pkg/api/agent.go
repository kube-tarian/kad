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

func (s *Server) SetClusterGitoptsProject(ctx context.Context, request *serverpb.SetClusterGitoptsProjectRequest) (
	*serverpb.SetClusterGitoptsProjectResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if orgId == "" {
		s.log.Errorf("organization ID is missing in the request")
		return &serverpb.SetClusterGitoptsProjectResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Organization Id is missing",
		}, nil
	}

	agent, err := s.agentHandeler.GetAgent(orgId, request.ClusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &serverpb.SetClusterGitoptsProjectResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Cluster Gitopts Project Set failed",
		}, nil
	}

	response, err := agent.GetClient().SetClusterGitoptsProject(context.Background(), &agentpb.SetClusterGitoptsProjectRequest{
		GitoptsUsecase: request.GitoptsUsecase,
		ProjectUrl:     request.ProjectUrl,
		AccessToken:    request.AccessToken,
	})
	if err != nil {
		s.log.Errorf("failed to set Cluster Gitopts Project, %v", err)
		return &serverpb.SetClusterGitoptsProjectResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Cluster Gitopts Project Set failed",
		}, nil
	}

	if response.Status != agentpb.StatusCode_OK {
		s.log.Errorf("Cluster Gitopts Project Set failed")
		return &serverpb.SetClusterGitoptsProjectResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Cluster Gitopts Project Set failed",
		}, nil
	}

	return &serverpb.SetClusterGitoptsProjectResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "Successfully Set Cluster Gitopts Project",
	}, nil
}

func (s *Server) GetClusterGitoptsProject(ctx context.Context, request *serverpb.GetClusterGitoptsProjectRequest) (
	*serverpb.GetClusterGitoptsProjectResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if orgId == "" {
		s.log.Errorf("organization ID is missing in the request")
		return &serverpb.GetClusterGitoptsProjectResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Organization Id is missing",
		}, nil
	}

	agent, err := s.agentHandeler.GetAgent(orgId, request.ClusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &serverpb.GetClusterGitoptsProjectResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get the Cluster Gitopts Project",
		}, nil
	}

	response, err := agent.GetClient().GetClusterGitoptsProject(context.Background(), &agentpb.GetClusterGitoptsProjectRequest{
		Usecase: request.Usecase,
	})
	if err != nil {
		s.log.Errorf("failed to get the Cluster Gitopts Project, %v", err)
		return &serverpb.GetClusterGitoptsProjectResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get the Cluster Gitopts Project",
		}, nil
	}

	if response.Status != agentpb.StatusCode_OK {
		s.log.Errorf("failed to get the Cluster Gitopts Project")
		return &serverpb.GetClusterGitoptsProjectResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get the Cluster Gitopts Project",
		}, nil
	}

	return &serverpb.GetClusterGitoptsProjectResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "Successfully fetched the Cluster Gitopts Project",
		ClusterGitoptsConfig: &serverpb.ClusterGitoptsConfig{
			Usecase:     response.ClusterGitoptsConfig.Usecase,
			ProjectUrl:  response.ClusterGitoptsConfig.ProjectUrl,
			AccessToken: response.ClusterGitoptsConfig.AccessToken,
			Status:      response.ClusterGitoptsConfig.Status,
		},
	}, nil
}

func (s *Server) DeleteClusterGitoptsProject(ctx context.Context, request *serverpb.DeleteClusterGitoptsProjectRequest) (
	*serverpb.DeleteClusterGitoptsProjectResponse, error) {
	metadataMap := metadataContextToMap(ctx)
	orgId := metadataMap[organizationIDAttribute]
	if orgId == "" {
		s.log.Errorf("organization ID is missing in the request")
		return &serverpb.DeleteClusterGitoptsProjectResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Organization Id is missing",
		}, nil
	}

	agent, err := s.agentHandeler.GetAgent(orgId, request.ClusterId)
	if err != nil {
		s.log.Errorf("failed to initialize agent, %v", err)
		return &serverpb.DeleteClusterGitoptsProjectResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Delete Cluster Gitopts Project failed",
		}, nil
	}

	response, err := agent.GetClient().DeleteClusterGitoptsProject(context.Background(), &agentpb.DeleteClusterGitoptsProjectRequest{
		Usecase:    request.Usecase,
		ProjectUrl: request.ProjectUrl,
	})
	if err != nil {
		s.log.Errorf("failed to Delete Cluster Gitopts Project, %v", err)
		return &serverpb.DeleteClusterGitoptsProjectResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Delete Cluster Gitopts Project failed",
		}, nil
	}

	if response.Status != agentpb.StatusCode_OK {
		s.log.Errorf("failed to Delete Cluster Gitopts Project")
		return &serverpb.DeleteClusterGitoptsProjectResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Delete Cluster Gitopts Project failed",
		}, nil
	}

	return &serverpb.DeleteClusterGitoptsProjectResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "Successfully Deleted Cluster Gitopts Project",
	}, nil
}
