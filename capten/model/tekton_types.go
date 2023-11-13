package model

type TektonProjectStatus string

const (
	TektonProjectAvailable            TektonProjectStatus = "available"
	TektonProjectConfigured           TektonProjectStatus = "configured"
	TektonProjectConfigurationOngoing TektonProjectStatus = "configuration-ongoing"
	TektonProjectConfigurationFailed  TektonProjectStatus = "configuration-failed"
)

type TektonProject struct {
	Id             string `json:"id,omitempty"`
	GitProjectId   string `json:"git_project_id,omitempty"`
	GitProjectUrl  string `json:"git_project_url,omitempty"`
	Status         string `json:"status,omitempty"`
	LastUpdateTime string `json:"last_update_time,omitempty"`
	WorkflowId     string `json:"workflow_id,omitempty"`
	WorkflowStatus string `json:"workflow_status,omitempty"`
}
