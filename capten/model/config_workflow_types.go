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

type CrossplaneClusterUpdate struct {
	Name             string `json:"name,omitempty"`
	GitProjectId     string `json:"gitProjectId,omitempty"`
	ManagedClusterId string `json:"managedClusterId,omitempty"`
	RepoURL          string `json:"repoURL,omitempty"`
}
