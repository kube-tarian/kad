package captensdk

import "github.com/kube-tarian/kad/integrator/common-pkg/logging"

type ApplicationClient struct {
	log logging.Logger
}

func (c *Client) NewApplicationClient() (*ApplicationClient, error) {
	return &ApplicationClient{log: c.log}, nil
}

type ApplicationRequest struct {
	RepoName  string `json:"repo_name" required:"true"`
	RepoURL   string `json:"repo_url" required:"true"`
	ChartName string `json:"chart_name" required:"true"`

	Namespace   string `json:"namespace" required:"true"`
	ReleaseName string `json:"release_name" required:"true"`
	Timeout     int    `json:"timeout" default:"5"`
	Version     string `json:"version"`

	KubeConfig string `json:"kube_config" required:"true"`
}

func (a *ApplicationClient) Create(req *ApplicationRequest) error {
	return nil
}

func (a *ApplicationClient) Delete(req *ApplicationRequest) error {
	return nil
}

func (a *ApplicationClient) Update(req *ApplicationRequest) error {
	return nil
}
