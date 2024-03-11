package pluginstore

import (
	"fmt"
	"os"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/server/pkg/pb/pluginstorepb"
	"github.com/kube-tarian/kad/server/pkg/store"
	"github.com/kube-tarian/kad/server/pkg/types"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type PluginStore struct {
	log     logging.Logger
	cfg     *Config
	dbStore store.ServerStore
}

func NewPluginStore(log logging.Logger, dbStore store.ServerStore) (*PluginStore, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}

	return &PluginStore{
		log:     log,
		cfg:     cfg,
		dbStore: dbStore,
	}, nil
}

func (p *PluginStore) ConfigureStore(clusterId string, config *pluginstorepb.PluginStoreConfig) error {
	return nil
}

func (p *PluginStore) GetStoreConfig(clusterId string, storeType pluginstorepb.StoreType) (*pluginstorepb.PluginStoreConfig, error) {
	return nil, nil
}

func (p *PluginStore) SyncPlugins(clusterId string, storeType pluginstorepb.StoreType) error {
	appListData, err := os.ReadFile(p.cfg.PluginsStorePath + "/" + p.cfg.PluginsFileName)
	if err != nil {
		return errors.WithMessage(err, "failed to read store config file")
	}

	var config PluginStoreConfig
	if err := yaml.Unmarshal(appListData, &config); err != nil {
		return errors.WithMessage(err, "failed to unmarshall store config file")
	}

	for _, pluginName := range config.Plugins {
		err := p.addPluginApp(pluginName)
		if err != nil {
			p.log.Errorf("%v", err)
		}
	}
	return nil
}

func (p *PluginStore) addPluginApp(pluginName string) error {
	appData, err := os.ReadFile(p.cfg.PluginsStorePath + "/" + pluginName + "/plugin.yaml")
	if err != nil {
		return errors.WithMessagef(err, "failed to read store plugin %s", pluginName)
	}

	var appConfig Plugin
	if err := yaml.Unmarshal(appData, &appConfig); err != nil {
		return errors.WithMessagef(err, "failed to unmarshall store plugin %s", pluginName)
	}

	if appConfig.PluginName == "" || len(appConfig.DeploymentConfig.Versions) == 0 {
		return fmt.Errorf("app name/version is missing for %s", pluginName)
	}

	plugin := &types.Plugin{
		PluginName:          appConfig.PluginName,
		Description:         appConfig.Description,
		Category:            appConfig.Category,
		ChartName:           appConfig.DeploymentConfig.ChartName,
		ChartRepo:           appConfig.DeploymentConfig.ChartRepo,
		Versions:            appConfig.DeploymentConfig.Versions,
		DefaultNamespace:    appConfig.DeploymentConfig.DefaultNamespace,
		PrivilegedNamespace: appConfig.DeploymentConfig.PrivilegedNamespace,
		PluginEndpoint:      appConfig.PluginConfig.Endpoint,
		Capabilities:        appConfig.PluginConfig.Capabilities,
	}

	if err := p.dbStore.AddOrUpdatePlugin(plugin); err != nil {
		return errors.WithMessagef(err, "failed to store plugin %s", pluginName)
	}
	return nil
}

func (p *PluginStore) GetPlugins(clusterId string, storeType pluginstorepb.StoreType) ([]*pluginstorepb.Plugin, error) {
	return nil, nil
}

func (p *PluginStore) GetPluginValues(clusterId string, storeType pluginstorepb.StoreType,
	pluginName, version string) ([]byte, error) {
	return nil, nil
}

func (p *PluginStore) DeployPlugin(clusterId string, storeType pluginstorepb.StoreType,
	pluginName, version string, values []byte) error {
	return nil
}

func (p *PluginStore) UnDeployPlugin(clusterId string, storeType pluginstorepb.StoreType, pluginName string) error {
	return nil
}
