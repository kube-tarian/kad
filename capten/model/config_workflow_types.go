package model

type Identifiers string

const (
	Container            Identifiers = "container"
	GitOrg               Identifiers = "gitOrg"
	ManagedCluster       Identifiers = "managedCluster"
	CrossplaneGitProject Identifiers = "crossplaneGitProject"
	TektonGitProject     Identifiers = "tektonGitProject"
)

var IdentifiersList = []Identifiers{Container, GitOrg, ManagedCluster, CrossplaneGitProject, TektonGitProject}

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
	ManagedClusterName string `json:"managedClusterName,omitempty"`
	ManagedClusterId   string `json:"managedClusterId,omitempty"`
	GitProjectId       string `json:"gitProjectId,omitempty"`
	RepoURL            string `json:"repoURL,omitempty"`
}

type CredentialIdentifier struct {
	Identifier string `json:"Identifier,omitempty"`
	Id         string `json:"Id,omitempty"`
	Url        string `json:"Url,omitempty"`
}

type CrossplaneProviderUpdate struct {
	ProviderId   string `json:"providerId,omitempty"`
	ProviderName string `json:"providerName,omitempty"`
	CloudType    string `json:"cloudType,omitempty"`
	GitProjectId string `json:"gitProjectId,omitempty"`
	RepoURL      string `json:"repoURL,omitempty"`
}

type TektonPipelineUseCase struct {
	Type                  string                               `json:"Type,omitempty"`
	RepoURL               string                               `json:"RepoURL,omitempty"`
	PipelineName          string                               `json:"PipelineName,omitempty"`
	Timeout               uint32                               `json:"Timeout,omitempty"`
	CredentialIdentifiers map[Identifiers]CredentialIdentifier `json:"CredentialIdentifiers,omitempty"`
}
