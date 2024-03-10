package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/pluginstorepb"
	pluginstore "github.com/kube-tarian/kad/server/pkg/plugin-store"
)

func (s *Server) ConfigPluginStore(ctx context.Context, request *pluginstorepb.ConfigPluginStoreRequest) (
	*pluginstorepb.ConfigPluginStoreResponse, error) {
	return &pluginstorepb.ConfigPluginStoreResponse{
		Status: pluginstorepb.StatusCode_OK,
	}, nil
}

func (s *Server) SyncPluginStore(ctx context.Context, request *pluginstorepb.SyncPluginStoreRequest) (
	*pluginstorepb.SyncPluginStoreResponse, error) {
	err := pluginstore.SyncPluginApps(s.log, s.serverStore)
	if err != nil {
		return &pluginstorepb.SyncPluginStoreResponse{
			Status:        pluginstorepb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: err.Error(),
		}, err
	}
	return &pluginstorepb.SyncPluginStoreResponse{
		Status: pluginstorepb.StatusCode_OK,
	}, nil
}

func (s *Server) GetPlugins(ctx context.Context, request *pluginstorepb.GetPluginsRequest) (
	*pluginstorepb.GetPluginsResponse, error) {
	return &pluginstorepb.GetPluginsResponse{
		Status: pluginstorepb.StatusCode_OK,
	}, nil
}

func (s *Server) GetPluginValues(ctx context.Context, request *pluginstorepb.GetPluginValuesRequest) (
	*pluginstorepb.GetPluginValuesResponse, error) {
	return &pluginstorepb.GetPluginValuesResponse{
		Status: pluginstorepb.StatusCode_OK,
	}, nil
}

func (s *Server) DeployPlugin(ctx context.Context, request *pluginstorepb.DeployPluginRequest) (
	*pluginstorepb.DeployPluginResponse, error) {
	return &pluginstorepb.DeployPluginResponse{
		Status: pluginstorepb.StatusCode_OK,
	}, nil
}

func (s *Server) UnDeployPlugin(ctx context.Context, request *pluginstorepb.UnDeployPluginRequest) (
	*pluginstorepb.UnDeployPluginResponse, error) {
	return &pluginstorepb.UnDeployPluginResponse{
		Status: pluginstorepb.StatusCode_OK,
	}, nil
}
