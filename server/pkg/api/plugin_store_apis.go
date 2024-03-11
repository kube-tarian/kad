package api

import (
	"context"

	"github.com/kube-tarian/kad/server/pkg/pb/pluginstorepb"
)

func (s *Server) ConfigurePluginStore(ctx context.Context, request *pluginstorepb.ConfigurePluginStoreRequest) (
	*pluginstorepb.ConfigurePluginStoreResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &pluginstorepb.ConfigurePluginStoreResponse{
			Status:        pluginstorepb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}
	s.log.Infof("Configure plugin store request recieved for cluster %s, [org: %s]", clusterId, orgId)

	err = s.pluginStore.ConfigureStore(clusterId, request.Config)
	if err != nil {
		s.log.Errorf("Configure plugin store request failed for cluster %s, [org: %s], %w", clusterId, orgId, err)
		return &pluginstorepb.ConfigurePluginStoreResponse{
			Status:        pluginstorepb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to configuere plugin store",
		}, err
	}

	s.log.Infof("Plugin store request processed for cluster %s, [org: %s]", clusterId, orgId)
	return &pluginstorepb.ConfigurePluginStoreResponse{
		Status: pluginstorepb.StatusCode_OK,
	}, nil
}

func (s *Server) GetPluginStoreConfig(ctx context.Context, request *pluginstorepb.GetPluginStoreConfigRequest) (
	*pluginstorepb.GetPluginStoreConfigResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &pluginstorepb.GetPluginStoreConfigResponse{
			Status:        pluginstorepb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}
	s.log.Infof("Get plugin store config request recieved for cluster %s, [org: %s]", clusterId, orgId)

	config, err := s.pluginStore.GetStoreConfig(clusterId, request.StoreType)
	if err != nil {
		s.log.Errorf("Get plugin store config request failed for cluster %s, [org: %s], %w", clusterId, orgId, err)
		return &pluginstorepb.GetPluginStoreConfigResponse{
			Status:        pluginstorepb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get plugin store config",
		}, err
	}

	s.log.Infof("Get plugin store config request processed for cluster %s, [org: %s]", clusterId, orgId)
	return &pluginstorepb.GetPluginStoreConfigResponse{
		Status: pluginstorepb.StatusCode_OK,
		Config: config,
	}, nil
}

func (s *Server) SyncPluginStore(ctx context.Context, request *pluginstorepb.SyncPluginStoreRequest) (
	*pluginstorepb.SyncPluginStoreResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &pluginstorepb.SyncPluginStoreResponse{
			Status:        pluginstorepb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}
	s.log.Infof("Sync plugin store request recieved for cluster %s, [org: %s]", clusterId, orgId)

	err = s.pluginStore.SyncPlugins(clusterId, request.StoreType)
	if err != nil {
		s.log.Errorf("Sync plugin store request failed for cluster %s, [org: %s], %w", clusterId, orgId, err)
		return &pluginstorepb.SyncPluginStoreResponse{
			Status:        pluginstorepb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: err.Error(),
		}, err
	}

	s.log.Infof("Sync plugin store request processed for cluster %s, [org: %s]", clusterId, orgId)
	return &pluginstorepb.SyncPluginStoreResponse{
		Status: pluginstorepb.StatusCode_OK,
	}, nil
}

func (s *Server) GetPlugins(ctx context.Context, request *pluginstorepb.GetPluginsRequest) (
	*pluginstorepb.GetPluginsResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &pluginstorepb.GetPluginsResponse{
			Status:        pluginstorepb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}
	s.log.Infof("Get plugins request recieved for cluster %s, [org: %s]", clusterId, orgId)

	plugins, err := s.pluginStore.GetPlugins(clusterId, request.StoreType)
	if err != nil {
		s.log.Errorf("Get plugins request failed for cluster %s, [org: %s], %w", clusterId, orgId, err)
		return &pluginstorepb.GetPluginsResponse{
			Status:        pluginstorepb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get plugins",
		}, err
	}

	s.log.Infof("Get plugins request processed for cluster %s, [org: %s]", clusterId, orgId)
	return &pluginstorepb.GetPluginsResponse{
		Status:  pluginstorepb.StatusCode_OK,
		Plugins: plugins,
	}, nil
}

func (s *Server) GetPluginValues(ctx context.Context, request *pluginstorepb.GetPluginValuesRequest) (
	*pluginstorepb.GetPluginValuesResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx, request.PluginName)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &pluginstorepb.GetPluginValuesResponse{
			Status:        pluginstorepb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}
	s.log.Infof("Get plugin values request recieved for plugin %s-%s, cluster %s, [org: %s]",
		request.PluginName, request.Version, clusterId, orgId)

	values, err := s.pluginStore.GetPluginValues(clusterId, request.StoreType, request.PluginName, request.Version)
	if err != nil {
		s.log.Errorf("Get plugins request failed for plugin %s-%s, cluster %s, [org: %s], %w",
			request.PluginName, request.Version, clusterId, orgId, err)
		return &pluginstorepb.GetPluginValuesResponse{
			Status:        pluginstorepb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to get plugins",
		}, err
	}

	s.log.Infof("Get plugin values request processed for plugin %s-%s, cluster %s, [org: %s]",
		request.PluginName, request.Version, clusterId, orgId)
	return &pluginstorepb.GetPluginValuesResponse{
		Status: pluginstorepb.StatusCode_OK,
		Values: values,
	}, nil
}

func (s *Server) DeployPlugin(ctx context.Context, request *pluginstorepb.DeployPluginRequest) (
	*pluginstorepb.DeployPluginResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx, request.PluginName)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &pluginstorepb.DeployPluginResponse{
			Status:        pluginstorepb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}
	s.log.Infof("Deploy plugin request recieved for plugin %s-%s, cluster %s, [org: %s]",
		request.PluginName, request.Version, clusterId, orgId)

	err = s.pluginStore.DeployPlugin(clusterId, request.StoreType, request.PluginName, request.Version, request.Values)
	if err != nil {
		s.log.Errorf("Deploy plugin request failed for plugin %s-%s, cluster %s, [org: %s], %w",
			request.PluginName, request.Version, clusterId, orgId, err)
		return &pluginstorepb.DeployPluginResponse{
			Status:        pluginstorepb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to deploy plugin",
		}, err
	}

	s.log.Infof("Deploy plugin request processed for plugin %s-%s, cluster %s, [org: %s]",
		request.PluginName, request.Version, clusterId, orgId)
	return &pluginstorepb.DeployPluginResponse{
		Status: pluginstorepb.StatusCode_OK,
	}, nil
}

func (s *Server) UnDeployPlugin(ctx context.Context, request *pluginstorepb.UnDeployPluginRequest) (
	*pluginstorepb.UnDeployPluginResponse, error) {
	orgId, clusterId, err := validateOrgClusterWithArgs(ctx, request.PluginName)
	if err != nil {
		s.log.Infof("request validation failed", err)
		return &pluginstorepb.UnDeployPluginResponse{
			Status:        pluginstorepb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, err
	}
	s.log.Infof("UnDeploy plugin request recieved for plugin %s, cluster %s, [org: %s]",
		request.PluginName, clusterId, orgId)

	err = s.pluginStore.UnDeployPlugin(clusterId, request.StoreType, request.PluginName)
	if err != nil {
		s.log.Errorf("UnDeploy plugin request failed for plugin %s, cluster %s, [org: %s], %w",
			request.PluginName, clusterId, orgId, err)
		return &pluginstorepb.UnDeployPluginResponse{
			Status:        pluginstorepb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "failed to undeploy plugin",
		}, err
	}

	s.log.Infof("UnDeploy plugin request processed for plugin %s, cluster %s, [org: %s]",
		request.PluginName, clusterId, orgId)
	return &pluginstorepb.UnDeployPluginResponse{
		Status: pluginstorepb.StatusCode_OK,
	}, nil
}
