package activities

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
)

type pluginDS map[string]map[string]string

type Config struct {
	GitDefaultCommiterName  string `envconfig:"GIT_COMMIT_NAME" default:"capten-bot"`
	GitDefaultCommiterEmail string `envconfig:"GIT_COMMIT_EMAIL" default:"capten-bot@intelops.dev"`
	VaultEntityName         string `envconfig:"VAULT_ENTITY_NAME" default:"gitproject"`
	GitCLoneDir             string `envconfig:"GIT_CLONE_DIR" default:"/gitCloneDir"`
	PluginConfig            string `envconfig:"PLUGIN_CONFIG_FILE" default:"/plugin_config.json"`
}

func GetConfig() (*Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	return &cfg, err
}

func ReadPluginConfig(pluginFile string) (pluginDS, error) {
	data, err := os.ReadFile(filepath.Clean(pluginFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read pluginConfig File: %s, err: %w", pluginFile, err)
	}

	var pluginData pluginDS
	err = json.Unmarshal(data, &pluginData)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return pluginData, nil
}
