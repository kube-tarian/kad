package model

import (
	"encoding/json"
	"fmt"
)

type RequestPayload struct {
	PluginName string          `json:"plugin_name"`
	Action     string          `json:"sub_action"`
	Data       json.RawMessage `json:"data,omitempty"` // TODO: This will be enhanced along with plugin implementation
}

type ResponsePayload struct {
	Status  string          `json:"status"`
	Message json.RawMessage `json:"message,omitempty"` // TODO: This will be enhanced along with plugin implementation
}

func (rsp *ResponsePayload) ToString() string {
	return fmt.Sprintf("Status: %s, Message: %s", rsp.Status, string(rsp.Message))
}

type Request struct {
	RepoName  string `json:"repo_name" required:"true"`
	RepoURL   string `json:"repo_url" required:"true"`
	ChartName string `json:"chart_name" required:"true"`

	Namespace   string `json:"namespace" required:"true"`
	ReleaseName string `json:"release_name" required:"true"`
	Timeout     int    `json:"timeout" required:"true"`
	Version     string `json:"version" required:"true"`
}
