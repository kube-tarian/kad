package crossplane

type Source struct {
	RepoURL        string `json:"repoURL,omitempty"`
	TargetRevision string `json:"targetRevision,omitempty"`
}

type Dest struct {
	Server    string `json:"server,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type GlobalValues struct {
	ClusterConfigPath string `json:"clusterConfigPath,omitempty"`
}

type DefaultApps struct {
	Name           string `yaml:"name" json:"name,omitempty"`
	ValuesPath     string `yaml:"valuesPath" json:"valuesPath,omitempty"`
	RepoURL        string `yaml:"repoURL" json:"repoURL,omitempty"`
	Namespace      string `yaml:"namespace" json:"namespace,omitempty"`
	Chart          string `yaml:"chart" json:"chart,omitempty"`
	TargetRevision string `yaml:"targetRevision" json:"targetRevision,omitempty"`
}

type DefaultAppList struct {
	DefaultApps []DefaultApps `yaml:"defaultApps"`
}

type Cluster struct {
	Name    string        `json:"name,omitempty"`
	Server  string        `json:"server,omitempty"`
	DefApps []DefaultApps `json:"defaultApps,omitempty"`
}

type ClusterConfigValues struct {
	Project      string       `json:"project,omitempty"`
	Global       GlobalValues `json:"global,omitempty"`
	Src          Source       `json:"source,omitempty"`
	Destination  Dest         `json:"destination,omitempty"`
	SyncPolicy   interface{}  `json:"syncPolicy,omitempty"`
	Compositions interface{}  `json:"compositions,omitempty"`
	Clusters     *[]Cluster   `json:"clusters,omitempty"`
}
