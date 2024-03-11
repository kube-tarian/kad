package pluginstore

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
