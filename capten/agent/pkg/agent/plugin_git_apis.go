package agent

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/intelops/go-common/credentials"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
)

const gitProjectEntityName = "gitproject"

func (a *Agent) AddGitProject(ctx context.Context, request *captenpluginspb.AddGitProjectRequest) (
	*captenpluginspb.AddGitProjectResponse, error) {

	// create new id
	id := uuid.New()

	// save to vault
	if err := a.storeAccesToken(ctx, id.String(), request.AccessToken); err != nil {
		return &captenpluginspb.AddGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: err.Error(),
		}, nil
	}

	// save to db
	gitProject := captenpluginspb.GitProject{
		Id:             id.String(),
		ProjectUrl:     request.ProjectUrl,
		Labels:         request.Labels,
		LastUpdateTime: request.LastUpdateTime,
	}
	if err := a.as.UpsertGitProject(&gitProject); err != nil {
		a.log.Errorf("failed to add gitProject in db, %v", err)
		return &captenpluginspb.AddGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to add gitProject in db",
		}, nil
	}

	return &captenpluginspb.AddGitProjectResponse{
		Id:            id.String(),
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) UpdateGitProject(ctx context.Context, request *captenpluginspb.UpdateGitProjectRequest) (
	*captenpluginspb.UpdateGitProjectResponse, error) {

	// validate id
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return &captenpluginspb.UpdateGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: fmt.Sprintf("invalid uuid: %s", request.Id),
		}, nil
	}

	if request.AccessToken != "" {
		// save to vault
		if err := a.storeAccesToken(ctx, request.Id, request.AccessToken); err != nil {
			return &captenpluginspb.UpdateGitProjectResponse{
				Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
				StatusMessage: err.Error(),
			}, nil
		}
	}

	// save to db
	gitProject := captenpluginspb.GitProject{
		Id:             id.String(),
		ProjectUrl:     request.ProjectUrl,
		Labels:         request.Labels,
		LastUpdateTime: request.LastUpdateTime,
	}
	if err := a.as.UpsertGitProject(&gitProject); err != nil {
		a.log.Errorf("failed to update gitProject in db, %v", err)
		return &captenpluginspb.UpdateGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to update gitProject in db",
		}, nil
	}

	return &captenpluginspb.UpdateGitProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) DeleteGitProject(ctx context.Context, request *captenpluginspb.DeleteGitProjectRequest) (
	*captenpluginspb.DeleteGitProjectResponse, error) {
	if err := a.as.DeleteGitProjectById(request.Id); err != nil {
		a.log.Errorf("failed to delete gitProject from db, %v", err)
		return &captenpluginspb.DeleteGitProjectResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to delete gitProject from db",
		}, nil
	}

	return &captenpluginspb.DeleteGitProjectResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "ok",
	}, nil
}

func (a *Agent) GetGitProjects(ctx context.Context, request *captenpluginspb.GetGitProjectsRequest) (
	*captenpluginspb.GetGitProjectsResponse, error) {

	res, err := a.as.GetGitProjects()
	if err != nil {
		a.log.Errorf("failed to get gitProjects from db, %v", err)
		return &captenpluginspb.GetGitProjectsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch git projects",
		}, nil
	}

	for _, r := range res {
		accessToken, err := a.getAccesToken(ctx, r.Id)
		if err != nil {
			a.log.Errorf("failed to get credential, %v", err)
			return &captenpluginspb.GetGitProjectsResponse{
				Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to fetch git projects",
			}, nil
		}
		r.AccessToken = accessToken
	}

	a.log.Infof("Found %d gitProjects", len(res))
	return &captenpluginspb.GetGitProjectsResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successful",
		GitProjects:   res,
	}, nil

}

func (a *Agent) GetGitProjectsForLabels(ctx context.Context, request *captenpluginspb.GetGitProjectsForLabelsRequest) (
	*captenpluginspb.GetGitProjectsForLabelsResponse, error) {

	res, err := a.as.GetGitProjectsByLabels(request.Labels)
	if err != nil {
		a.log.Errorf("failed to get gitProjects for labels from db, %v", err)
		return &captenpluginspb.GetGitProjectsForLabelsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch git projects",
		}, nil
	}

	for _, r := range res {
		// fetch cred
		accessToken, err := a.getAccesToken(ctx, r.Id)
		if err != nil {
			a.log.Errorf("failed to get credential, %v", err)
			return &captenpluginspb.GetGitProjectsForLabelsResponse{
				Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to fetch git projects",
			}, nil
		}
		r.AccessToken = accessToken
	}

	return &captenpluginspb.GetGitProjectsForLabelsResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successful",
		Projects:      res,
	}, nil
}

func (a *Agent) getAccesToken(ctx context.Context, id string) (string, error) {
	credPath := fmt.Sprintf("%s/%s/%s", credentials.GenericCredentialType, gitProjectEntityName, id)
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to intialize credentials client for %s", credPath)
		a.log.Errorf("failed to get crendential for %s, %v", credPath, err)
		return "", err
	}

	cred, err := credAdmin.GetCredential(ctx, credentials.GenericCredentialType, gitProjectEntityName, id)
	if err != nil {
		a.log.Errorf("failed to get credential for %s, %v", credPath, err)
		return "", err
	}
	return cred["accessToken"], nil
}

func (a *Agent) storeAccesToken(ctx context.Context, id string, accessToken string) error {
	credPath := fmt.Sprintf("%s/%s/%s", credentials.GenericCredentialType, gitProjectEntityName, id)
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to intialize credentials client for %s", credPath)
		a.log.Errorf("failed to store credential for %s, %v", credPath, err)
		return err
	}

	credentialMap := map[string]string{
		"accessToken": accessToken,
	}
	err = credAdmin.PutCredential(ctx, credentials.GenericCredentialType, gitProjectEntityName,
		id, credentialMap)

	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to store crendential for %s", credPath)
		a.log.Errorf("failed to store credential for %s, %v", credPath, err)
		return err
	}
	a.log.Audit("security", "storecred", "success", "system", "credential stored for %s", credPath)
	a.log.Infof("stored credential for entity %s", credPath)
	return nil
}
