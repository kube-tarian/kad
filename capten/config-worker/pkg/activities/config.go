package activities

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
)

type tektonPluginDS map[string]string
type crossplanePluginDS map[string]string
type crossplaneProviderPluginDS map[string]string

type Config struct {
	GitDefaultCommiterName         string `envconfig:"GIT_COMMIT_NAME" default:"capten-bot"`
	GitDefaultCommiterEmail        string `envconfig:"GIT_COMMIT_EMAIL" default:"capten-bot@intelops.dev"`
	VaultEntityName                string `envconfig:"VAULT_ENTITY_NAME" default:"gitproject"`
	GitCLoneDir                    string `envconfig:"GIT_CLONE_DIR" default:"/gitCloneDir"`
	TektonPluginConfig             string `envconfig:"TEKTON_PLUGIN_CONFIG_FILE" default:"/tekton_plugin_config.json"`
	CrossPlanePluginConfig         string `envconfig:"CROSSPLANE_PLUGIN_CONFIG_FILE" default:"/crossplane_plugin_config.json"`
	CrossPlaneProviderPluginConfig string `envconfig:"CROSSPLANE_PLUGIN_CONFIG_FILE" default:"/crossplane_provider_plugin_config.json"`
}

func GetConfig() (*Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	return &cfg, err
}

func ReadTektonPluginConfig(pluginFile string) (tektonPluginDS, error) {
	data, err := os.ReadFile(filepath.Clean(pluginFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read pluginConfig File: %s, err: %w", pluginFile, err)
	}

	var pluginData tektonPluginDS
	err = json.Unmarshal(data, &pluginData)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return pluginData, nil
}

func ReadCrossPlanePluginConfig(pluginFile string) (crossplanePluginDS, error) {
	data, err := os.ReadFile(filepath.Clean(pluginFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read pluginConfig File: %s, err: %w", pluginFile, err)
	}

	var pluginData crossplanePluginDS
	err = json.Unmarshal(data, &pluginData)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return pluginData, nil
}

func ReadCrossPlaneProviderPluginConfig(pluginFile string) (crossplaneProviderPluginDS, error) {
	data, err := os.ReadFile(filepath.Clean(pluginFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read pluginConfig File: %s, err: %w", pluginFile, err)
	}

	var pluginData crossplaneProviderPluginDS
	err = json.Unmarshal(data, &pluginData)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return pluginData, nil
}
