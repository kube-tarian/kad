package model

// DeployerDeleteRequest defines model for DeployerDeleteRequest.
type DeployerDeleteRequest struct {
	// ClusterName Cluster in which to be deleted, default in-build cluster
	ClusterName *string `json:"cluster_name,omitempty"`

	// Namespace Namespace chart to be installed
	Namespace string `json:"namespace"`

	// PluginName Plugin name
	PluginName string `json:"plugin_name"`

	// ReleaseName Release name to be used for install
	ReleaseName string `json:"release_name"`

	// Timeout Timeout for the application installation
	Timeout int `json:"timeout"`
}

type ApplicationDeployRequest struct {
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
