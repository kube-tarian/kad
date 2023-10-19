package agent

import (
	"context"
	"fmt"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/common-pkg/plugins/argocd"
)

func (a *Agent) RegisterArgoCDProject(ctx context.Context, request *captenpluginspb.RegisterArgoCDProjectRequest) (
	*captenpluginspb.RegisterArgoCDProjectResponse, error) {
	return &captenpluginspb.RegisterArgoCDProjectResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}

func (a *Agent) UnRegisterArgoCDProject(ctx context.Context, request *captenpluginspb.UnRegisterArgoCDProjectRequest) (
	*captenpluginspb.UnRegisterArgoCDProjectResponse, error) {
	return &captenpluginspb.UnRegisterArgoCDProjectResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}

func (a *Agent) GetArgoCDProjects(ctx context.Context, request *captenpluginspb.GetArgoCDProjectsRequest) (
	*captenpluginspb.GetArgoCDProjectsResponse, error) {

	argocdClient, err := argocd.NewClient(&logging.Logging{})
	if err != nil {
		fmt.Println("error occured Error clinet -> " + err.Error())
		return &captenpluginspb.GetArgoCDProjectsResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Error occured while argocd client",
		}, err
	}
	list, err := argocdClient.ListRepositories(ctx)
	if err != nil {
		fmt.Println("error occured Error fewtch -> " + err.Error())
		return &captenpluginspb.GetArgoCDProjectsResponse{
			Status:        captenpluginspb.StatusCode_NOT_FOUND,
			StatusMessage: "Error occured while fetching repositories",
		}, err
	}

	fmt.Println(list)
	return &captenpluginspb.GetArgoCDProjectsResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}
