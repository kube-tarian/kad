package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"
)

func (s *Server) SetClusterGitoptsProject(ctx context.Context, request *serverpb.SetClusterGitoptsProjectRequest) (
	*serverpb.SetClusterGitoptsProjectResponse, error) {
	orgId, err := validateRequest(ctx, request.ClusterId)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.SetClusterGitoptsProjectResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Set cluster gitopts project request recieved for cluster %s, [org: %s]", request.ClusterId, orgId)

	if v, ok := request.Credential[credentialAccessTokenKey]; !ok || v == "" {
		s.log.Errorf("accessToken is missing in the request")
		return &serverpb.SetClusterGitoptsProjectResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "accessToken credential is missing",
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
		Usecase:    request.Usecase,
		ProjectUrl: request.ProjectUrl,
		Credential: request.Credential,
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

	s.log.Infof("Set cluster gitopts project request for cluster %s successful, [org: %s]", request.ClusterId, orgId)
	return &serverpb.SetClusterGitoptsProjectResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "Successfully Set Cluster Gitopts Project",
	}, nil
}

func (s *Server) GetClusterGitoptsProject(ctx context.Context, request *serverpb.GetClusterGitoptsProjectRequest) (
	*serverpb.GetClusterGitoptsProjectResponse, error) {
	orgId, err := validateRequest(ctx, request.ClusterId)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.GetClusterGitoptsProjectResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Get cluster gitopts project request recieved for cluster %s, [org: %s]", request.ClusterId, orgId)

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

	s.log.Infof("Get cluster gitopts project request for cluster %s successful, [org: %s]", request.ClusterId, orgId)
	return &serverpb.GetClusterGitoptsProjectResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "Successfully fetched the Cluster Gitopts Project",
		ClusterGitoptsConfig: &serverpb.ClusterGitoptsConfig{
			Usecase:    response.ClusterGitoptsConfig.Usecase,
			ProjectUrl: response.ClusterGitoptsConfig.ProjectUrl,
			Status:     response.ClusterGitoptsConfig.Status,
			Credential: response.ClusterGitoptsConfig.Credential,
		},
	}, nil
}

func (s *Server) DeleteClusterGitoptsProject(ctx context.Context, request *serverpb.DeleteClusterGitoptsProjectRequest) (
	*serverpb.DeleteClusterGitoptsProjectResponse, error) {
	orgId, err := validateRequest(ctx, request.ClusterId)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &serverpb.DeleteClusterGitoptsProjectResponse{
			Status:        serverpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}
	s.log.Infof("Delete cluster gitopts project request recieved for cluster %s, [org: %s]", request.ClusterId, orgId)

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

	s.log.Infof("Delete cluster gitopts project request for cluster %s successful, [org: %s]", request.ClusterId, orgId)
	return &serverpb.DeleteClusterGitoptsProjectResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "Successfully Deleted Cluster Gitopts Project",
	}, nil
}
