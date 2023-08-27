package model

type CreteRequestPayload struct {
	RepoName  string `json:"repo_name" required:"true"`
	RepoURL   string `json:"repo_url" required:"true"`
	ChartName string `json:"chart_name" required:"true"`

	Namespace   string `json:"namespace" required:"true"`
	ReleaseName string `json:"release_name" required:"true"`
	Timeout     int    `json:"timeout" default:"5"`
	Version     string `json:"version"`

	ClusterName string `json:"cluster_name" required:"false"`
	ValuesYaml  string `json:"values_yaml" required:"false"`

	// CreateNamespace bool `json:"createNamespace"`
}

type DeleteRequestPayload struct {
	Namespace   string `json:"namespace" required:"true"`
	ReleaseName string `json:"release_name" required:"true"`
	Timeout     int    `json:"timeout" default:"5"`

	ClusterName string `json:"cluster_name" required:"false"`
}

type ListRequestPayload struct {
	RepoName  string `json:"repo_name" required:"true"`
	Namespace string `json:"namespace" required:"true"`
	Timeout   int    `json:"timeout" default:"5"`

	ClusterName string `json:"cluster_name" required:"false"`
}
