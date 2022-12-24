package model

import (
	"encoding/json"
	"fmt"
)

type RequestPayload struct {
	PluginName string          `json:"plugin_name" required:"true"`
	Action     string          `json:"action" required:"true"`
	Data       json.RawMessage `json:"data" required:"true"` // TODO: This will be enhanced along with plugin implementation
}

func (r *RequestPayload) ToString() string {
	return fmt.Sprintf("plugin_name: %s, action: %s, data: %s", r.PluginName, r.Action, string(r.Data))
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
	Timeout     int    `json:"timeout" default:"5"`
	Version     string `json:"version"`
}
