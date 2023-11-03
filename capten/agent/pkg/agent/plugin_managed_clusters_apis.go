package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
)

func (a *Agent) GetManagedClusters(ctx context.Context, request *captenpluginspb.GetManagedClustersRequest) (
	*captenpluginspb.GetManagedClustersResponse, error) {
	a.log.Infof("Get Managed Clusters request recieved")
	// TO DO
	a.log.Infof("Fetched %d Managed Clusters", 0)
	return &captenpluginspb.GetManagedClustersResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successfully fetched the Crossplane projects",
		Clusters:      []*captenpluginspb.ManagedCluster{},
	}, nil
}
