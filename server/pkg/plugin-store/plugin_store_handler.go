package pluginstore

import (
	"fmt"
	"os"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/server/pkg/store"
	"github.com/kube-tarian/kad/server/pkg/types"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	PluginsStorePath string `envconfig:"PLUGIN_APP_CONFIG_PATH" default:"/data/app-store/data"`
	PluginsFileName  string `envconfig:"APP_STORE_CONFIG_FILE" default:"plugins.yaml"`
}

type PluginStoreConfig struct {
	Plugins []string `yaml:"plugins"`
}

type DeploymentConfig struct {
	Versions            []string `yaml:"versions"`
	ChartName           string   `yaml:"chartName"`
	ChartRepo           string   `yaml:"chartRepo"`
	DefaultNamespace    string   `yaml:"defaultNamespace"`
	PrivilegedNamespace bool     `yaml:"privilegedNamespace"`
}

type PluginConfig struct {
	Endpoint     string   `yaml:"Endpoint"`
	Capabilities []string `yaml:"capabilities"`
}

type Plugin struct {
	PluginName       string           `yaml:"pluginName"`
	Description      string           `yaml:"description"`
	Category         string           `yaml:"category"`
	Icon             string           `yaml:"icon"`
	DeploymentConfig DeploymentConfig `yaml:"deploymentConfig"`
	PluginConfig     PluginConfig     `yaml:"pluginConfig"`
}

func SyncPluginApps(log logging.Logger, appStore store.ServerStore) error {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return err
	}

	appListData, err := os.ReadFile(cfg.PluginsStorePath + "/" + cfg.PluginsFileName)
	if err != nil {
		return errors.WithMessage(err, "failed to read store config file")
	}

	var config PluginStoreConfig
	if err := yaml.Unmarshal(appListData, &config); err != nil {
		return errors.WithMessage(err, "failed to unmarshall store config file")
	}

	for _, pluginName := range config.Plugins {
		err := addPluginApp(pluginName, cfg, appStore)
		if err != nil {
			log.Errorf("%v", err)
		}
	}
	return nil
}

func addPluginApp(pluginName string, cfg *Config, appStore store.ServerStore) error {
	appData, err := os.ReadFile(cfg.PluginsStorePath + "/" + pluginName + "/plugin.yaml")
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

	if err := appStore.AddOrUpdatePlugin(plugin); err != nil {
		return errors.WithMessagef(err, "failed to store plugin %s", pluginName)
	}
	return nil
}
