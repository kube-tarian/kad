package api

import (
	"context"

	"github.com/kube-tarian/kad/capten/common-pkg/cluster-plugins/clusterpluginspb"
)

func (a *Agent) GetClusterPlugins(ctx context.Context, request *clusterpluginspb.GetClusterPluginsRequest) (
	*clusterpluginspb.GetClusterPluginsResponse, error) {
	return &clusterpluginspb.GetClusterPluginsResponse{
		Status: clusterpluginspb.StatusCode_OK}, nil
}

func (a *Agent) DeployClusterPlugin(ctx context.Context, request *clusterpluginspb.DeployClusterPluginRequest) (
	*clusterpluginspb.DeployClusterPluginResponse, error) {
	return &clusterpluginspb.DeployClusterPluginResponse{
		Status: clusterpluginspb.StatusCode_OK}, nil
}

func (a *Agent) UnDeployClusterPlugin(ctx context.Context, request *clusterpluginspb.UnDeployClusterPluginRequest) (
	*clusterpluginspb.UnDeployClusterPluginResponse, error) {
	return &clusterpluginspb.UnDeployClusterPluginResponse{
		Status: clusterpluginspb.StatusCode_OK}, nil
}
