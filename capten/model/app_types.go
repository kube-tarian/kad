package model

type WorkFlowStatus string

type AppStatus string

const (
	AppIntallingStatus      AppStatus = "Installing"
	AppIntalledStatus       AppStatus = "Installed"
	AppUpgradingStatus      AppStatus = "Upgrading"
	AppUpgradedStatus       AppStatus = "Upgraded"
	AppIntallFailedStatus   AppStatus = "Installion Failed"
	AppUpgradeFaileddStatus AppStatus = "Upgrade Failed"
	AppUnInstalledStatus    AppStatus = "UnInstalled"
	AppUnInstallingStatus   AppStatus = "UnInstalling"
)

type AppDeployAction string

const (
	AppInstallAction   AppDeployAction = "install"
	AppUnInstallAction AppDeployAction = "delete"
	AppUpgradeAction   AppDeployAction = "upgrade"
)

type AppConfig struct {
	AppName             string `json:"AppName,omitempty"`
	Version             string `json:"Version,omitempty"`
	Category            string `json:"Category,omitempty"`
	Description         string `json:"Description,omitempty"`
	ChartName           string `json:"ChartName,omitempty"`
	RepoName            string `json:"RepoName,omitempty"`
	ReleaseName         string `json:"ReleaseName,omitempty"`
	RepoURL             string `json:"RepoURL,omitempty"`
	Namespace           string `json:"Namespace,omitempty"`
	CreateNamespace     bool   `json:"CreateNamespace"`
	PrivilegedNamespace bool   `json:"PrivilegedNamespace"`
	Icon                string `json:"Icon,omitempty"`
	LaunchURL           string `json:"LaunchURL,omitempty"`
	LaunchUIDescription string `json:"LaunchUIDescription,omitempty"`
}

type ApplicationInstallRequest struct {
	PluginName     string `json:"PluginName,omitempty"`
	RepoName       string `json:"RepoName,omitempty"`
	RepoURL        string `json:"RepoURL,omitempty"`
	ChartName      string `json:"ChartName,omitempty"`
	Namespace      string `json:"Namespace,omitempty"`
	ReleaseName    string `json:"ReleaseName,omitempty"`
	Timeout        uint32 `json:"Timeout,omitempty"`
	Version        string `json:"Version,omitempty"`
	ClusterName    string `json:"ClusterName,omitempty"`
	OverrideValues string `json:"OverrideValues,omitempty"`
}

type ApplicationDeleteRequest struct {
	PluginName  string `json:"plugin_name,omitempty"`
	Namespace   string `json:"namespace,omitempty"`
	ReleaseName string `json:"release_name,omitempty"`
	Timeout     uint32 `json:"timeout,omitempty"`
	ClusterName string `json:"cluster_name,omitempty"`
}

type ClusterGitoptsConfig struct {
	Usecase    string `json:"usecase,omitempty"`
	ProjectUrl string `json:"project_url,omitempty"`
	Status     string `json:"status,omitempty"`
}
