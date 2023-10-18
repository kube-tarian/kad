package agent

import (
	"context"
	"fmt"

	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
)

func (a *Agent) AddGitProject(ctx context.Context, request *captenpluginspb.AddGitProjectRequest) (
	*captenpluginspb.AddGitProjectResponse, error) {
	return &captenpluginspb.AddGitProjectResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}

func (a *Agent) UpdateGitProject(ctx context.Context, request *captenpluginspb.UpdateGitProjectRequest) (
	*captenpluginspb.UpdateGitProjectResponse, error) {
	return &captenpluginspb.UpdateGitProjectResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}

func (a *Agent) DeleteGitProject(ctx context.Context, request *captenpluginspb.DeleteGitProjectRequest) (
	*captenpluginspb.DeleteGitProjectResponse, error) {
	return &captenpluginspb.DeleteGitProjectResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}

func (a *Agent) GetGitProjects(ctx context.Context, request *captenpluginspb.GetGitProjectsRequest) (
	*captenpluginspb.GetGitProjectsResponse, error) {
	return &captenpluginspb.GetGitProjectsResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}

func (a *Agent) GetGitProjectsForLabels(ctx context.Context, request *captenpluginspb.GetGitProjectsForLabelsRequest) (
	*captenpluginspb.GetGitProjectsForLabelsResponse, error) {
	return &captenpluginspb.GetGitProjectsForLabelsResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}
