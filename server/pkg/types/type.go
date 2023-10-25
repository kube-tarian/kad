package types

const (
	ClientCertChainFileName = "cert-chain.pem"
	ClientCertFileName      = "client.crt"
	ClientKeyFileName       = "client.key"
	AgentPortCfgKey         = "agent.port"
	AgentTlsEnabledCfgKey   = "agent.tlsEnabled"
	ServerDbCfgKey          = "server.db"
)

type AgentInfo struct {
	Endpoint string
	CaPem    string
	Cert     string
	Key      string
}

type ClusterDetails struct {
	ClusterID   string
	Endpoint    string
	OrgID       string
	ClusterName string
}

type StoreAppConfig struct {
	AppName             string `json:"appName,omitempty"`
	Version             string `json:"version,omitempty"`
	Category            string `json:"category,omitempty"`
	Description         string `json:"description,omitempty"`
	ChartName           string `json:"chartName,omitempty"`
	RepoName            string `json:"repoName,omitempty"`
	ReleaseName         string `json:"releaseName,omitempty"`
	RepoURL             string `json:"repoURL,omitempty"`
	Namespace           string `json:"namespace,omitempty"`
	CreateNamespace     bool   `json:"createNamespace"`
	PrivilegedNamespace bool   `json:"privilegedNamespace"`
	Icon                string `json:"icon,omitempty"`
	LaunchURL           string `yaml:"LaunchURL,omitempty"`
	LaunchUIDescription string `yaml:"LaunchUIDescription,omitempty"`
	OverrideValues      string `json:"overrideValues,omitempty"`
	LaunchUIValues      string `json:"launchUIValues,omitempty"`
	TemplateValues      string `json:"templateValues,omitempty"`
	PluginName          string `yaml:"PluginName"`
	PluginDescription   string `yaml:"PluginDescription"`
}

type AppConfig struct {
	Name                string `yaml:"Name"`
	ChartName           string `yaml:"ChartName"`
	Category            string `yaml:"Category"`
	RepoName            string `yaml:"RepoName"`
	RepoURL             string `yaml:"RepoURL"`
	Namespace           string `yaml:"Namespace"`
	ReleaseName         string `yaml:"ReleaseName"`
	Version             string `yaml:"Version"`
	Description         string `yaml:"Description"`
	LaunchURL           string `yaml:"LaunchURL"`
	LaunchUIDescription string `yaml:"LaunchUIDescription"`
	LaunchUIIcon        string `yaml:"LaunchUIIcon"`
	LaunchUIValues      []byte `yaml:"LaunchUIValues"`
	OverrideValues      []byte `yaml:"OverrideValues"`
	CreateNamespace     bool   `yaml:"CreateNamespace"`
	PrivilegedNamespace bool   `yaml:"PrivilegedNamespace"`
	TemplateValues      []byte `yaml:"TemplateValues"`
	Icon                string `yaml:"Icon"`
	PluginName          string `yaml:"PluginName"`
	PluginDescription   string `yaml:"PluginDescription"`
}

type AppInstallRequest struct {
	PluginName  string `json:"plugin_name"`
	RepoName    string `json:"repo_name"`
	RepoUrl     string `json:"repo_url"`
	ChartName   string `json:"chart_name"`
	Namespace   string `json:"namespace"`
	ReleaseName string `json:"release_name"`
	Timeout     int    `json:"timeout"`
}
