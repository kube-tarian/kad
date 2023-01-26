package captensdk

import "github.com/kube-tarian/kad/integrator/common-pkg/logging"

type CaptenClient interface {
	NewApplicationClient() (*ApplicationClient, error)
	NewClusterClient() (*ClusterClient, error)
	NewRepoClient() (*RepoClient, error)
	NewProjectClient() (*ProjectClient, error)
}

type Client struct {
	log logging.Logger
}

func NewClient(log logging.Logger) (*Client, error) {
	return &Client{
		log: log,
	}, nil
}
