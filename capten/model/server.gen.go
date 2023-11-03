package model

import "github.com/kube-tarian/kad/capten/agent/pkg/model"

// AgentRequest defines model for AgentRequest.
type AgentRequest struct {
	CustomerId string `json:"customer_id"`
	Endpoint   string `json:"endpoint"`
}

// ClimonDeleteRequest defines model for ClimonDeleteRequest.
type ClimonDeleteRequest struct {
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

// ClimonPostRequest defines model for ClimonPostRequest.
type ClimonPostRequest struct {
	// ChartName Chart name in Repository
	ChartName string `json:"chart_name"`

	// ClusterName Cluster in which to be installed, default in-build cluster
	ClusterName *string `json:"cluster_name,omitempty"`

	// Namespace Namespace chart to be installed
	Namespace string `json:"namespace"`

	// PluginName Plugin name
	PluginName string `json:"plugin_name"`

	// ReleaseName Release name to be used for install
	ReleaseName string `json:"release_name"`

	// RepoName Repository name
	RepoName string `json:"repo_name"`

	// RepoUrl Repository URL
	RepoUrl string `json:"repo_url"`

	// Timeout Timeout for the application installation
	Timeout int `json:"timeout"`

	// Version Version of the chart
	Version *string `json:"version,omitempty"`
}

// ClusterRequest defines model for ClusterRequest.
type ClusterRequest struct {
	ClusterName string `json:"cluster_name"`
	PluginName  string `json:"plugin_name"`
}

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

type UseCase struct {
	Type                string            `json:"Type,omitempty"`
	RepoURL             string            `json:"RepoURL,omitempty"`
	VaultCredIdentifier string            `json:"VaultCredIdentifier,omitempty"`
	Timeout             uint32            `json:"Timeout,omitempty"`
	OverrideValues      map[string]string `json:"OverrideValues,omitempty"`
	PushToDefaultBranch bool              `json:"PushToDefaultBranch,omitempty"`
}

type CrossplaneUseCase struct {
	Type                string                     `json:"Type,omitempty"`
	RepoURL             string                     `json:"RepoURL,omitempty"`
	VaultCredIdentifier string                     `json:"VaultCredIdentifier,omitempty"`
	Timeout             uint32                     `json:"Timeout,omitempty"`
	OverrideValues      map[string]string          `json:"OverrideValues,omitempty"`
	PushToDefaultBranch bool                       `json:"PushToDefaultBranch,omitempty"`
	CrossplaneProviders []model.CrossplaneProvider `json:"ProviderInfo,omitempty"`
}

// ProjectDeleteRequest defines model for ProjectDeleteRequest.
type ProjectDeleteRequest struct {
	// PluginName Plugin name
	PluginName string `json:"plugin_name"`

	// ProjectName Project name to be created in plugin
	ProjectName string `json:"project_name"`
}

// ProjectPostRequest defines model for ProjectPostRequest.
type ProjectPostRequest struct {
	// PluginName Plugin name
	PluginName string `json:"plugin_name"`

	// ProjectName Project name to be created in plugin
	ProjectName string `json:"project_name"`
}

// RepositoryDeleteRequest defines model for RepositoryDeleteRequest.
type RepositoryDeleteRequest struct {
	// PluginName Plugin name
	PluginName string `json:"plugin_name"`

	// RepoName Repository to added to plugin
	RepoName string `json:"repo_name"`
}

// RepositoryPostRequest defines model for RepositoryPostRequest.
type RepositoryPostRequest struct {
	// PluginName Plugin name
	PluginName string `json:"plugin_name"`

	// RepoName Repository to added to plugin
	RepoName string `json:"repo_name"`

	// RepoUrl Repository URL
	RepoUrl string `json:"repo_url"`
}

// Response Configuration request response
type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}
