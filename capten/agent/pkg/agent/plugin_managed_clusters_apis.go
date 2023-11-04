package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
)

func (a *Agent) GetManagedClusters(ctx context.Context, request *captenpluginspb.GetManagedClustersRequest) (
	*captenpluginspb.GetManagedClustersResponse, error) {
	a.log.Infof("Get Managed Clusters request recieved")
	// TODO
	a.log.Infof("Fetched %d Managed Clusters", 0)
	return &captenpluginspb.GetManagedClustersResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successfully fetched the Crossplane projects",
		Clusters:      []*captenpluginspb.ManagedCluster{},
	}, nil
}

func (a *Agent) GetManagedClusterKubeconfig(ctx context.Context, request *captenpluginspb.GetManagedClusterKubeconfigRequest) (
	*captenpluginspb.GetManagedClusterKubeconfigResponse, error) {
	if err := validateArgs(request.Id); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.GetManagedClusterKubeconfigResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	a.log.Infof("Get Managed Cluster %s kubeconfig request recieved", request.Id)
	// TODO
	a.log.Infof("Fetched %d Managed Clusters", 0)
	return &captenpluginspb.GetManagedClusterKubeconfigResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successfully fetched the Crossplane projects",
		Kubeconfig:    "",
	}, nil
}
