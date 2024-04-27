package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/common-pkg/plugins/argocd"
)

const (
	argoCDRepositoryType    string = "git"
	argoCDRepositoryProject string = "Default"
)

func (a *Agent) RegisterArgoCDProject(ctx context.Context, request *captenpluginspb.RegisterArgoCDProjectRequest) (
	*captenpluginspb.RegisterArgoCDProjectResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (a *Agent) GetArgoCDProjects(ctx context.Context, request *captenpluginspb.GetArgoCDProjectsRequest) (
	*captenpluginspb.GetArgoCDProjectsResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (a *Agent) addProjectToArgoCD(ctx context.Context, projectUrl, userID, accessToken string) error {
	argocdClient, err := argocd.NewClient(&logging.Logging{})
	if err != nil {
		return err
	}

	if !strings.HasSuffix(projectUrl, ".git") {
		projectUrl = projectUrl + ".git"
	}

	repo := &argocd.Repository{
		Project:               argoCDRepositoryProject,
		Repo:                  projectUrl,
		Username:              userID,
		Password:              accessToken,
		Type:                  argoCDRepositoryType,
		Insecure:              false,
		EnableLFS:             false,
		InsecureIgnoreHostKey: false,
		Upsert:                true,
		ConnectionState: argocd.ConnectionState{
			Status:  "Connected",
			Message: "Repository is connected",
		},
	}

	_, err = argocdClient.CreateRepository(ctx, repo)
	if err != nil {
		return err
	}
	return nil
}

func (a *Agent) deleteProjectFromArgoCD(ctx context.Context, projectUrl string) error {
	argocdClient, err := argocd.NewClient(&logging.Logging{})
	if err != nil {
		return err
	}
	_, err = argocdClient.DeleteRepository(ctx, projectUrl)
	if err != nil {
		return err
	}
	return nil
}

func (a *Agent) isProjectRegisteredWithArgoCD(ctx context.Context, projectUrl string) (bool, error) {
	argocdClient, err := argocd.NewClient(&logging.Logging{})
	if err != nil {
		return false, err
	}
	_, err = argocdClient.GetRepository(ctx, projectUrl)
	if err != nil && fmt.Sprintf("rpc error: code = NotFound desc = repo '%s' not found", projectUrl) == err.Error() {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
