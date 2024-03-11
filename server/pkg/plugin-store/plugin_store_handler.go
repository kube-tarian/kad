package pluginstore

import (
	"fmt"
	"os"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/server/pkg/pb/pluginstorepb"
	"github.com/kube-tarian/kad/server/pkg/store"
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
	return p.dbStore.WritePluginStoreConfig(clusterId, config)
}

func (p *PluginStore) GetStoreConfig(clusterId string, storeType pluginstorepb.StoreType) (*pluginstorepb.PluginStoreConfig, error) {
	if storeType == pluginstorepb.StoreType_LOCAL_STORE {
		return p.dbStore.ReadPluginStoreConfig(clusterId)
	} else if storeType == pluginstorepb.StoreType_CENTRAL_STORE {
		return &pluginstorepb.PluginStoreConfig{
			StoreType:     pluginstorepb.StoreType_CENTRAL_STORE,
			GitProjectId:  p.cfg.PluginStoreProjectURL,
			GitProjectURL: p.cfg.PluginStoreProjectID,
		}, nil
	} else {
		return nil, fmt.Errorf("not supported store type")
	}
}

func (p *PluginStore) SyncPlugins(clusterId string, storeType pluginstorepb.StoreType) error {
	config, err := p.GetStoreConfig(clusterId, storeType)
	if err != nil {
		return err
	}

	pluginStoreDir, err := p.clonePluginStoreProject(config.GitProjectId, config.GitProjectURL)
	if err != nil {
		return err
	}
	defer os.RemoveAll(pluginStoreDir)

	p.log.Infof("Loading plugin data from project %s clone dir %s", config.GitProjectURL, pluginStoreDir)
	pluginListData, err := os.ReadFile(p.cfg.PluginsStorePath + "/" + p.cfg.PluginsFileName)
	if err != nil {
		return errors.WithMessage(err, "failed to read store config file")
	}

	var plugins PluginListData
	if err := yaml.Unmarshal(pluginListData, &plugins); err != nil {
		return errors.WithMessage(err, "failed to unmarshall store config file")
	}

	for _, pluginName := range plugins.Plugins {
		err := p.addPluginApp(config.GitProjectId, pluginName)
		if err != nil {
			p.log.Errorf("%v", err)
			continue
		}
		p.log.Infof("stored plugin data for plugin %s for cluster %s", pluginName, clusterId)
	}
	return nil
}

func (p *PluginStore) clonePluginStoreProject(projectURL, _ string) (pluginStoreDir string, err error) {
	pluginStoreDir, err = os.MkdirTemp(p.cfg.PluginsStoreProjectMount, tmpGitProjectCloneStr)
	if err != nil {
		err = fmt.Errorf("failed to create template tmp dir, err: %v", err)
		return
	}

	gitClient := NewGitClient()
	if err = gitClient.Clone(pluginStoreDir, projectURL, ""); err != nil {
		os.RemoveAll(pluginStoreDir)
		err = fmt.Errorf("failed to Clone template repo, err: %v", err)
		return
	}
	return
}

func (p *PluginStore) addPluginApp(gitProjectId, pluginName string) error {
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

	plugin := &pluginstorepb.PluginData{
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

	if err := p.dbStore.WritePluginData(gitProjectId, plugin); err != nil {
		return errors.WithMessagef(err, "failed to store plugin %s", pluginName)
	}
	return nil
}

func (p *PluginStore) GetPlugins(clusterId string, storeType pluginstorepb.StoreType) ([]*pluginstorepb.Plugin, error) {
	config, err := p.GetStoreConfig(clusterId, storeType)
	if err != nil {
		return nil, err
	}

	return p.dbStore.ReadPlugins(config.GitProjectId)
}

func (p *PluginStore) GetPluginData(clusterId string, storeType pluginstorepb.StoreType, pluginName string) (*pluginstorepb.PluginData, error) {
	config, err := p.GetStoreConfig(clusterId, storeType)
	if err != nil {
		return nil, err
	}

	return p.dbStore.ReadPluginData(config.GitProjectId, pluginName)
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
