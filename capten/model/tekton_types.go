package model

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TektonProjectStatus string

const (
	TektonProjectAvailable            TektonProjectStatus = "available"
	TektonProjectConfigured           TektonProjectStatus = "configured"
	TektonProjectConfigurationOngoing TektonProjectStatus = "configuration-ongoing"
	TektonProjectConfigurationFailed  TektonProjectStatus = "configuration-failed"
)

const (
	TektonPipelineConfigUseCase = "tekton-pipelines"
	TektonHostName              = "tekton"
	TektonPipelineCreate        = "tekton-pipeline-create"
	TektonPipelineSync          = "tekton-pipeline-sync"
	TektonPipelineDelete        = "tekton-pipeline-delete"
)

type TektonPipelineStatus string

const (
	TektonPipelineAvailable            TektonPipelineStatus = "available"
	TektonPipelineConfigured           TektonPipelineStatus = "configured"
	TektonPipelineConfigurationOngoing TektonPipelineStatus = "configuration-ongoing"
	TektonPipelineConfigurationFailed  TektonPipelineStatus = "configuration-failed"
)

const (
	TektonPipelineOutofSynch    TektonPipelineStatus = "OutOfSynch"
	TektonPipelineInSynch       TektonPipelineStatus = "InSynch"
	TektonPipelineFailedToSynch TektonPipelineStatus = "FailedToSynch"
	TektonPipelineReady         TektonPipelineStatus = "Ready"
	TektonPipelineNotReady      TektonPipelineStatus = "NotReady"
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

type TektonPipeline struct {
	Id             string   `json:"id,omitempty"`
	PipelineName   string   `json:"pipeline_name,omitempty"`
	WebhookURL     string   `json:"webhook_url,omitempty"`
	GitProjectId   string   `json:"git_project_id,omitempty"`
	GitProjectUrl  string   `json:"git_project_url,omitempty"`
	ContainerRegId []string `json:"container_reg_id,omitempty"`
	Status         string   `json:"status,omitempty"`
	LastUpdateTime string   `json:"last_update_time,omitempty"`
	WorkflowId     string   `json:"workflow_id,omitempty"`
	WorkflowStatus string   `json:"workflow_status,omitempty"`
}

type EventListenerStatus struct {
	ConditionedStatus `json:",inline"`
}

type EventListener struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status EventListenerStatus `json:"status,omitempty"`
}

type EventListeners struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EventListener `json:"items"`
}
