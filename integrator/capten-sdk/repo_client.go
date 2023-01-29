package captensdk

import "github.com/kube-tarian/kad/integrator/common-pkg/logging"

type RepoClient struct {
	log logging.Logger
}

func (c *Client) NewRepoClient() (*RepoClient, error) {
	return &RepoClient{log: c.log}, nil
}

type RepositoryRequest struct {
}

func (r *RepoClient) Create(req *RepositoryRequest) error {
	return nil
}

func (r *RepoClient) Delete(req *RepositoryRequest) error {
	return nil
}

func (r *RepoClient) Update(req *RepositoryRequest) error {
	return nil
}
