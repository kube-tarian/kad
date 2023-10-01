package github

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/google/go-github/v55/github"
)

type Operations interface {
	CreatePR(ctx context.Context, prRepo, prRepoOwner, prSubject, commitBranch,
		repoBranch, prBranch, prDescription string) (string, error)
}

type Operation struct {
	client *github.Client
}

func NewOperation(token string) *Operation {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	client := github.NewClient(httpClient).WithAuthToken(token)
	return &Operation{client: client}
}

func (op *Operation) CreatePR(ctx context.Context, prRepo, prRepoOwner, prSubject,
	commitBranch, repoBranch, prDescription string) (string, error) {
	newPR := &github.NewPullRequest{
		Title:               &prSubject,
		Head:                &commitBranch,
		HeadRepo:            &repoBranch,
		Base:                &repoBranch,
		Body:                &prDescription,
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := op.client.PullRequests.Create(ctx, prRepoOwner, prRepo, newPR)
	if err != nil {
		return "", err
	}

	return pr.GetHTMLURL(), nil
}
