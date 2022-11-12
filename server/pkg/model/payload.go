package model

import (
	"encoding/json"
	"fmt"
)

type DeployRequestPayload struct {
	PluginName string          `json:"plugin_name" required:"true"`
	Action     string          `json:"action" required:"true"`
	Data       json.RawMessage `json:"data" required:"true"`
}

type DeployResponsePayload struct {
	Status  string          `json:"status"`
	Message json.RawMessage `json:"message,omitempty"` // TODO: This will be enhanced along with plugin implementation
}

func (rsp *DeployResponsePayload) ToString() string {
	return fmt.Sprintf("Status: %s, Message: %s", rsp.Status, string(rsp.Message))
}

type DeployRequestData struct {
	RepoName  string `json:"repo_name" required:"true"`
	RepoURL   string `json:"repo_url" required:"true"`
	ChartName string `json:"chart_name" required:"true"`

	Namespace   string `json:"namespace" required:"true"`
	ReleaseName string `json:"release_name" required:"true"`
	Timeout     int    `json:"timeout" default:"5"`
	Version     string `json:"version"`
}
