package argocd

import (
	"context"
	"net/url"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient/repository"
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/argoproj/argo-cd/v2/util/io"
)

func (a *ArgoCDClient) CreateRepository(ctx context.Context, repo *Repository) (*v1alpha1.Repository, error) {
	conn, appClient, err := a.client.NewRepoClient()
	if err != nil {
		return nil, err
	}
	defer io.Close(conn)

	var repoUpsert bool
	existingRepo, err := a.GetRepository(ctx, repo.Repo)
	if existingRepo != nil && err == nil {
		repoUpsert = true
	}

	resp, err := appClient.CreateRepository(ctx, &repository.RepoCreateRequest{
		Repo: &v1alpha1.Repository{
			Project:               repo.Project,
			Repo:                  repo.Repo,
			Username:              repo.Username,
			Password:              repo.Password,
			Type:                  repo.Type,
			Insecure:              repo.Insecure,
			EnableLFS:             repo.EnableLFS,
			InsecureIgnoreHostKey: repo.InsecureIgnoreHostKey,
			ConnectionState: v1alpha1.ConnectionState{
				Status:  repo.ConnectionState.Status,
				Message: repo.ConnectionState.Message,
			},
		},
		Upsert: repoUpsert,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *ArgoCDClient) DeleteRepository(ctx context.Context, repo string) (*repository.RepoResponse, error) {
	conn, appClient, err := a.client.NewRepoClient()
	if err != nil {
		return nil, err
	}
	defer io.Close(conn)

	encodedRepo := url.QueryEscape(repo)

	resp, err := appClient.DeleteRepository(ctx, &repository.RepoQuery{Repo: encodedRepo})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *ArgoCDClient) GetRepository(ctx context.Context, repo string) (*v1alpha1.Repository, error) {
	conn, appClient, err := a.client.NewRepoClient()
	if err != nil {
		return nil, err
	}
	defer io.Close(conn)

	encodedRepo := url.QueryEscape(repo)

	repository, err := appClient.Get(ctx, &repository.RepoQuery{Repo: encodedRepo})
	if err != nil {
		return nil, err
	}

	return repository, nil
}

func (a *ArgoCDClient) ListRepositories(ctx context.Context) (*v1alpha1.RepositoryList, error) {
	conn, appClient, err := a.client.NewRepoClient()
	if err != nil {
		return nil, err
	}
	defer io.Close(conn)

	list, err := appClient.ListRepositories(ctx, &repository.RepoQuery{})
	if err != nil {
		return nil, err
	}

	return list, nil
}
