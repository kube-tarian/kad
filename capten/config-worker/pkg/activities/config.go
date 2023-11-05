package activities

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
)

type tektonPluginConfig map[string]string
type crossplanePluginConfig map[string]string

type Config struct {
	GitDefaultCommiterName  string `envconfig:"GIT_COMMIT_NAME" default:"capten-bot"`
	GitDefaultCommiterEmail string `envconfig:"GIT_COMMIT_EMAIL" default:"capten-bot@intelops.dev"`
	VaultEntityName         string `envconfig:"VAULT_ENTITY_NAME" default:"gitproject"`
	GitCloneDir             string `envconfig:"GIT_CLONE_DIR" default:"/gitCloneDir"`
	TektonPluginConfig      string `envconfig:"TEKTON_PLUGIN_CONFIG_FILE" default:"/tekton_plugin_config.json"`
	CrossPlanePluginConfig  string `envconfig:"CROSSPLANE_PLUGIN_CONFIG_FILE" default:"/crossplane_plugin_config.json"`
}

func GetConfig() (*Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	return &cfg, err
}

func readTektonPluginConfig(pluginFile string) (tektonPluginConfig, error) {
	data, err := os.ReadFile(filepath.Clean(pluginFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read pluginConfig File: %s, err: %w", pluginFile, err)
	}

	var pluginData tektonPluginConfig
	err = json.Unmarshal(data, &pluginData)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return pluginData, nil
}

func readCrossPlanePluginConfig(pluginFile string) (crossplanePluginConfig, error) {
	data, err := os.ReadFile(filepath.Clean(pluginFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read pluginConfig File: %s, err: %w", pluginFile, err)
	}

	var pluginData crossplanePluginConfig
	err = json.Unmarshal(data, &pluginData)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return pluginData, nil
}
