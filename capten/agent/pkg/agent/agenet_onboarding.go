package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
)

func (a *Agent) AddOrUpdateOnboarding(ctx context.Context, request *agentpb.AddOrUpdateOnboardingRequest) (*agentpb.AddOrUpdateOnboardingResponse, error) {

	if err := a.as.AddOrUpdateOnboardingIntegration(request); err != nil {
		a.log.Errorf("Add/Update of onboarding integration failed, %v", err)
		return &agentpb.AddOrUpdateOnboardingResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Add/Update onboarding integration failed",
		}, err
	}

	a.log.Infof("Add/Update of onboarding integration successful. Project Url - %s", request.ProjectUrl)
	return &agentpb.AddOrUpdateOnboardingResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "Add/Update of onboarding integration successful",
	}, nil
}

func (a *Agent) GetOnboarding(ctx context.Context, request *agentpb.GetOnboardingRequest) (*agentpb.GetOnboardingResponse, error) {

	resp, err := a.as.GetOnboardingIntegration(request.Type, request.ProjectUrl)
	if err != nil {
		a.log.Errorf("failed to get onboarding integration, %v", err)
		return &agentpb.GetOnboardingResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get the onboarding integration",
		}, err
	}

	a.log.Infof("Successfully fetched the onboarding integration. Project Url - %s", request.ProjectUrl)
	return &agentpb.GetOnboardingResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "Successfully fetched the onboarding integration",
		Onboarding:    resp,
	}, nil
}

func (a *Agent) DeleteOnboarding(ctx context.Context, request *agentpb.DeleteOnboardingRequest) (*agentpb.DeleteOnboardingResponse, error) {

	if err := a.as.DeleteOnboardingIntegration(request.Type, request.ProjectUrl); err != nil {
		a.log.Errorf("failed to delete onboarding integration, %v", err)
		return &agentpb.DeleteOnboardingResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to delete the onboarding integration",
		}, err
	}

	a.log.Infof("Successfully deleted the onboarding integration. Project Url - %s", request.ProjectUrl)
	return &agentpb.DeleteOnboardingResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "Successfully deleted the onboarding integration",
	}, nil
}
