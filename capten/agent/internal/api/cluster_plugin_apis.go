package api

import (
	"context"

	"github.com/kube-tarian/kad/capten/common-pkg/cluster-plugins/clusterpluginspb"
)

func (a *Agent) GetClusterPlugins(ctx context.Context, request *clusterpluginspb.GetClusterPluginsRequest) (
	*clusterpluginspb.GetClusterPluginsResponse, error) {
	a.log.Infof("Recieved get cluster plugins request")
	return &clusterpluginspb.GetClusterPluginsResponse{
		Status: clusterpluginspb.StatusCode_OK}, nil
}

func (a *Agent) DeployClusterPlugin(ctx context.Context, request *clusterpluginspb.DeployClusterPluginRequest) (
	*clusterpluginspb.DeployClusterPluginResponse, error) {
	a.log.Infof("Recieved deploy cluster plugin request, %+v", request.Plugin)
	return &clusterpluginspb.DeployClusterPluginResponse{
		Status: clusterpluginspb.StatusCode_OK}, nil
}

func (a *Agent) UnDeployClusterPlugin(ctx context.Context, request *clusterpluginspb.UnDeployClusterPluginRequest) (
	*clusterpluginspb.UnDeployClusterPluginResponse, error) {
	a.log.Infof("Recieved undeploy cluster plugin request, %s", request.PluginName)
	return &clusterpluginspb.UnDeployClusterPluginResponse{
		Status: clusterpluginspb.StatusCode_OK}, nil
}
