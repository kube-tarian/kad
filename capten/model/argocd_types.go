package model

type ArgoCDProjectStatus string

const (
	ArgoCDProjectAvailable           ArgoCDProjectStatus = "available"
	ArgoCDProjectConfigured          ArgoCDProjectStatus = "configured"
	ArgoCDProjectConfigurationFailed ArgoCDProjectStatus = "configuration-failed"
)

type ArgoCDProject struct {
	Id             string `json:"id,omitempty"`
	GitProjectId   string `json:"git_project_id,omitempty"`
	GitProjectUrl  string `json:"git_project_url,omitempty"`
	Status         string `json:"status,omitempty"`
	LastUpdateTime string `json:"last_update_time,omitempty"`
}
