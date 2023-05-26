package captensdk

import "github.com/kube-tarian/kad/integrator/common-pkg/logging"

type ProjectClient struct {
	log logging.Logger
}

func (c *Client) NewProjectClient() (*ProjectClient, error) {
	return &ProjectClient{log: c.log}, nil
}

type ProjectRequest struct {
}

func (r *ProjectClient) Create(req *ProjectRequest) error {
	return nil
}

func (r *ProjectClient) Delete(req *ProjectRequest) error {
	return nil
}

func (r *ProjectClient) Update(req *ProjectRequest) error {
	return nil
}
