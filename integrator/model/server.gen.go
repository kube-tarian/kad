package model

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

// DeployerPostRequest defines model for DeployerPostRequest.
type DeployerPostRequest struct {
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
