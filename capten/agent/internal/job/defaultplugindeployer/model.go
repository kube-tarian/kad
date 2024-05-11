package defaultplugindeployer

const (
	tmpGitProjectCloneStr          = "clone*"
	gitProjectAccessTokenAttribute = "accessToken"
	gitProjectUserId               = "userID"
)

var (
	supporttedCapabilities = map[string]bool{
		"ui-sso-oauth":    true,
		"capten-sdk":      true,
		"postgress-store": true,
		"vault-store":     true}
)

type Config struct {
	PluginsStoreProjectMount     string `envconfig:"PLUGIN_STORE_PROJECT_MOUNT" default:"/plugin-store-clone"`
	PluginsStorePath             string `envconfig:"PLUGIN_STORE_PATH" default:"/plugin-store"`
	PluginsFileName              string `envconfig:"PLUGIN_LIST_FILE" default:"default-plugin-list.yaml"`
	PluginFileName               string `envconfig:"PLUGIN_FILE" default:"plugin.yaml"`
	PluginConfigFileName         string `envconfig:"PLUGIN_CONFIG_FILE" default:"plugin-config.yaml"`
	PluginStoreProjectURL        string `envconfig:"PLUGIN_STORE_PROJECT_URL" default:"https://github.com/intelops/capten-plugins"`
	PluginStoreProjectAccess     string `envconfig:"PLUGIN_STORE_PROJECT_ACCESS" default:""`
	PluginStoreProjectID         string `envconfig:"PLUGIN_STORE_PROJECT_ID" default:"1cf5201d-5f35-4d5b-afe0-4b9d0e0d4cd2"`
	GitVaultEntityName           string `envconfig:"GIT_VAULT_ENTITY_NAME" default:"git-project"`
	DefaultPluginsGitAccessToken string `envconfig:"DEFAULT_PLUGINS_GIT_ACCESS_TOKEN" required:"true"`
}

type PluginListData struct {
	Plugins []string `yaml:"plugins"`
}

type Deployment struct {
	ControlplaneCluster *DeploymentConfig `yaml:"controlplaneCluster"`
	BussinessCluster    *DeploymentConfig `yaml:"bussinessCluster"`
}

type DeploymentConfig struct {
	Version             string `yaml:"version"`
	ChartName           string `yaml:"chartName"`
	ChartRepo           string `yaml:"chartRepo"`
	DefaultNamespace    string `yaml:"defaultNamespace"`
	PrivilegedNamespace bool   `yaml:"privilegedNamespace"`
	ValuesFile          string `yaml:"valuesFile"`
}

type PluginConfig struct {
	Deployment         Deployment `yaml:"deployment"`
	ApiEndpoint        string     `yaml:"apiEndpoint"`
	UIEndpoint         string     `yaml:"uiEndpoint"`
	UIModulePackageURL string     `yaml:"uiModulePackageURL"`
	Capabilities       []string   `yaml:"capabilities"`
}

type Plugin struct {
	PluginName  string   `yaml:"pluginName"`
	Description string   `yaml:"description"`
	Category    string   `yaml:"category"`
	Icon        string   `yaml:"icon"`
	Versions    []string `yaml:"versions"`
}
