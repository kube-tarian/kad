package pluginstore

const (
	tmpGitProjectCloneStr = "clone*"
)

type Config struct {
	PluginsStoreProjectMount string `envconfig:"PLUGIN_STORE_PROJECT_MOUNT" default:"/plugin-store-clone"`
	PluginsStorePath         string `envconfig:"PLUGIN_STORE_PATH" default:"/plugin-store"`
	PluginsFileName          string `envconfig:"PLUGIN_LIST_FILE" default:"plugin-list.yaml"`
	PluginStoreProjectURL    string `envconfig:"PLUGIN_STORE_PROJECT_URL" default:"https://github.com/vramk23/capten-plugins"`
	PluginStoreProjectID     string `envconfig:"PLUGIN_STORE_PROJECT_ID" default:"1cf5201d-5f35-4d5b-afe0-4b9d0e0d4cd2"`
}

type PluginListData struct {
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
