package captensdk

import "github.com/kube-tarian/kad/integrator/common-pkg/logging"

type ClusterClient struct {
	log logging.Logger
}

func (c *Client) NewClusterClient() (*ClusterClient, error) {
	return &ClusterClient{log: c.log}, nil
}

type ClusterRequest struct {
}

func (cc *ClusterClient) Create(req *ClusterRequest) error {
	return nil
}

func (cc *ClusterClient) Delete(req *ClusterRequest) error {
	return nil
}

func (cc *ClusterClient) Update(req *ClusterRequest) error {
	return nil
}
