package model

const (
	WorkFlowStatusStarted    WorkFlowStatus = "started"
	WorkFlowStatusCompleted  WorkFlowStatus = "completed"
	WorkFlowStatusInProgress WorkFlowStatus = "in-progress"
	WorkFlowStatusFailed     WorkFlowStatus = "failed"
)

type ConfigureParameters struct {
	Resource string
	Action   string
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
	Type                string               `json:"Type,omitempty"`
	RepoURL             string               `json:"RepoURL,omitempty"`
	VaultCredIdentifier string               `json:"VaultCredIdentifier,omitempty"`
	Timeout             uint32               `json:"Timeout,omitempty"`
	OverrideValues      map[string]string    `json:"OverrideValues,omitempty"`
	CrossplaneProviders []CrossplaneProvider `json:"ProviderInfo,omitempty"`
}

type CrossplaneClusterEndpoint struct {
	Name       string `json:"name,omitempty"`
	Endpoint   string `json:"endpoint,omitempty"`
	Kubeconfig string `json:"kubeconfig,omitempty"`
	Id         string `json:"id,omitempty"`
	RepoURL    string `json:"repoURL,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
}

type Source struct {
	RepoURL        string `json:"repoURL,omitempty"`
	TargetRevision string `json:"targetRevision,omitempty"`
}

type Dest struct {
	Server    string `json:"server,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type DefaultApps struct {
	Name           string `json:"name,omitempty"`
	AppConfigPath  string `json:"appConfigPath,omitempty"`
	RepoURL        string `json:"repoURL,omitempty"`
	Namespace      string `json:"namespace,omitempty"`
	Chart          string `json:"chart,omitempty"`
	TargetRevision string `json:"targetRevision,omitempty"`
}
type Cluster struct {
	Name       string        `json:"name,omitempty"`
	ConfigPath string        `json:"configPath,omitempty"`
	Server     string        `json:"server,omitempty"`
	DefApps    []DefaultApps `json:"defaultApps,omitempty"`
}

type ArgoCDAppValue struct {
	Project      string      `json:"project,omitempty"`
	Src          Source      `json:"source,omitempty"`
	Destination  Dest        `json:"destination,omitempty"`
	SyncPolicy   interface{} `json:"syncPolicy,omitempty"`
	Compositions interface{} `json:"compositions,omitempty"`
	Clusters     []Cluster   `json:"clusters,omitempty"`
}
