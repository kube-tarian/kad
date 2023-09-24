package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/model"
)

func (a *Agent) SetClusterGitoptsProject(ctx context.Context, request *agentpb.SetClusterGitoptsProjectRequest) (*agentpb.SetClusterGitoptsProjectResponse, error) {

	confi := &model.ClusterGitoptsConfig{
		Usecase:     request.GitoptsUsecase,
		ProjectUrl:  request.ProjectUrl,
		AccessToken: request.AccessToken,
		Status:      "started",
	}

	if err := a.as.AddOrUpdateOnboardingIntegration(confi); err != nil {
		a.log.Errorf("failed to Set Cluster Gitopts Project, %v", err)
		return &agentpb.SetClusterGitoptsProjectResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Cluster Gitopts Project Set failed",
		}, err
	}

	a.log.Infof("Set Cluster Gitopts Project successful. Project Url - %s", request.ProjectUrl)
	return &agentpb.SetClusterGitoptsProjectResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "Set Cluster Gitopts Project successful",
	}, nil
}

func (a *Agent) GetClusterGitoptsProject(ctx context.Context, request *agentpb.GetClusterGitoptsProjectRequest) (*agentpb.GetClusterGitoptsProjectResponse, error) {

	resp, err := a.as.GetOnboardingIntegration(request.Usecase)
	if err != nil {
		a.log.Errorf("failed to get the Cluster Gitopts Project, %v", err)
		return &agentpb.GetClusterGitoptsProjectResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get the Cluster Gitopts Project",
		}, err
	}

	a.log.Infof("Successfully fetched the the Cluster Gitopts Project. Project Url - %s", request.Usecase)
	return &agentpb.GetClusterGitoptsProjectResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "Successfully fetched the onboarding integration",
		ClusterGitoptsConfig: &agentpb.ClusterGitoptsConfig{
			Usecase:     resp.Usecase,
			ProjectUrl:  resp.ProjectUrl,
			AccessToken: resp.AccessToken,
			Status:      resp.Status,
		},
	}, nil
}

func (a *Agent) DeleteClusterGitoptsProject(ctx context.Context, request *agentpb.DeleteClusterGitoptsProjectRequest) (*agentpb.DeleteClusterGitoptsProjectResponse, error) {

	if err := a.as.DeleteOnboardingIntegration(request.Usecase, request.ProjectUrl); err != nil {
		a.log.Errorf("failed to delete onboarding integration, %v", err)
		return &agentpb.DeleteClusterGitoptsProjectResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to delete the onboarding integration",
		}, err
	}

	a.log.Infof("Successfully deleted the onboarding integration. Project Url - %s", request.ProjectUrl)
	return &agentpb.DeleteClusterGitoptsProjectResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "Successfully deleted the onboarding integration",
	}, nil
}
