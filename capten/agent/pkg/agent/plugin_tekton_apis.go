package agent

import (
	"context"
	"fmt"
	"os"

	"github.com/intelops/go-common/credentials"
	"github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
)

const (
	tekton string = "tekton"
)

func (a *Agent) RegisterTektonProject(ctx context.Context, request *captenpluginspb.RegisterTektonProjectRequest) (
	*captenpluginspb.RegisterTektonProjectResponse, error) {
	// get the corressponding git url from the DB.
	projectUrl := os.Getenv("ProjectURL")
	accessToken := os.Getenv("accessToken")
	regTekton := &model.RegisterTekton{Id: request.Id, ProjectUrl: projectUrl, Status: "in-progress"}
	if err := a.as.AddTektonProject(regTekton); err != nil {
		a.log.Errorf("failed to Set Cluster Gitopts Project, %v", err)
		return &captenpluginspb.RegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "inserting data to tekton db got failed",
		}, err
	}
	credPath := fmt.Sprintf("%s/%s/%s", credentials.GenericCredentialType, tekton, tekton)
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to intialize credentails client for %s", credPath)
		a.log.Errorf("failed to store credentail for %s, %v", credPath, err)
		return &captenpluginspb.RegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "",
		}, err
	}

	err = credAdmin.PutCredential(ctx, credentials.GenericCredentialType, tekton,
		tekton, map[string]string{"accessToken": accessToken})
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to store credentail for %s", credPath)
		a.log.Errorf("failed to store credentail for %s, %v", credPath, err)
		return &captenpluginspb.RegisterTektonProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "",
		}, err
	}
	a.log.Audit("security", "storecred", "success", "system", "credentail stored for %s", credPath)
	a.log.Infof("stored credentail for entity %s", credPath)

	// start the config-worker routine
	go a.configureGitRepo(regTekton, tekton)

	return &captenpluginspb.RegisterTektonProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successfully registered tekton",
	}, nil
}

func (a *Agent) UnRegisterTektonProject(ctx context.Context, request *captenpluginspb.UnRegisterTektonProjectRequest) (
	*captenpluginspb.UnRegisterTektonProjectResponse, error) {
	return &captenpluginspb.UnRegisterTektonProjectResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}

func (a *Agent) GetTektonProjects(ctx context.Context, request *captenpluginspb.GetTektonProjectsRequest) (
	*captenpluginspb.GetTektonProjectsResponse, error) {
	return &captenpluginspb.GetTektonProjectsResponse{
		Status:        captenpluginspb.StatusCode_NOT_FOUND,
		StatusMessage: "not implemented",
	}, fmt.Errorf("not implemented")
}
