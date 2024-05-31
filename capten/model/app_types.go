package model

import "fmt"

type WorkFlowStatus string

type AppStatus string

const (
	AppIntallingStatus         AppStatus = "Installing"
	AppIntalledStatus          AppStatus = "Installed"
	AppUpgradingStatus         AppStatus = "Upgrading"
	AppUpgradedStatus          AppStatus = "Upgraded"
	AppIntallFailedStatus      AppStatus = "Installion Failed"
	AppUpgradeFaileddStatus    AppStatus = "Upgrade Failed"
	AppUnInstalledStatus       AppStatus = "UnInstalled"
	AppUnUninstallFailedStatus AppStatus = "UnInstall Failed"
	AppUnInstallingStatus      AppStatus = "UnInstalling"
)

type AppDeployAction string

const (
	AppInstallAction   AppDeployAction = "install"
	AppUnInstallAction AppDeployAction = "delete"
	AppUpgradeAction   AppDeployAction = "upgrade"
	AppUpdateAction    AppDeployAction = "update"
)

type DeployRequest interface {
	String() string
}


type PluginDeployRequest struct {
    Data string
}

func (p *PluginDeployRequest) String() string {
    return p.Data
}

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

func (a *ApplicationInstallRequest) String() string {
	return fmt.Sprintf("pluginName: %s, repoName: %s, repoURL: %s, chartName: %s, Namespace: %s, releaseName: %s, timeout: %v, version: %v, clusterName: %s",
		a.PluginName, a.RepoName, a.RepoURL, a.ChartName, a.Namespace, a.ReleaseName, a.Timeout, a.Version, a.ClusterName)
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
