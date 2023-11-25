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
	Name           string `json:"name,omitempty"`
	ValuesPath     string `json:"valuesPath,omitempty"`
	RepoURL        string `json:"repoURL,omitempty"`
	Namespace      string `json:"namespace,omitempty"`
	Chart          string `json:"chart,omitempty"`
	TargetRevision string `json:"targetRevision,omitempty"`
}

type Cluster struct {
	Name    string        `json:"name,omitempty"`
	Server  string        `json:"server,omitempty"`
	DefApps []DefaultApps `json:"defaultApps,omitempty"`
}

type ArgoCDAppValue struct {
	Project      string       `json:"project,omitempty"`
	Global       GlobalValues `json:"global,omitempty"`
	Src          Source       `json:"source,omitempty"`
	Destination  Dest         `json:"destination,omitempty"`
	SyncPolicy   interface{}  `json:"syncPolicy,omitempty"`
	Compositions interface{}  `json:"compositions,omitempty"`
	Clusters     *[]Cluster   `json:"clusters,omitempty"`
}
