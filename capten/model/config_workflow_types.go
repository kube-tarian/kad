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

type TektonPipelineUseCase struct {
	Type                       string            `json:"Type,omitempty"`
	RepoURL                    string            `json:"RepoURL,omitempty"`
	PipelineName               string            `json:"PipelineName,omitempty"`
	GitCredIdentifier          string            `json:"GitCredIdentifier,omitempty"`
	GitCredId                  string            `json:"GitCredId,omitempty"`
	ContainerRegCredIdentifier string            `json:"ContainerRegCredIdentifier,omitempty"`
	ContainerRegUrlIdMap       map[string]string `json:"ContainerRegUrlIdMap,omitempty"`
	Timeout                    uint32            `json:"Timeout,omitempty"`
	OverrideValues             map[string]string `json:"OverrideValues,omitempty"`
}

type CrossplaneClusterUpdate struct {
	ManagedClusterName string `json:"managedClusterName,omitempty"`
	ManagedClusterId   string `json:"managedClusterId,omitempty"`
	GitProjectId       string `json:"gitProjectId,omitempty"`
	RepoURL            string `json:"repoURL,omitempty"`
}
