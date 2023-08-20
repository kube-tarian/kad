package model

import (
	"encoding/json"
	"fmt"
)

type RequestPayload struct {
	PluginName string          `json:"plugin_name"`
	Action     string          `json:"action"`
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

type AppConfig struct {
	AppName             string `json:"AppName,omitempty"`
	Version             string `json:"Version,omitempty"`
	Category            string `json:"Category,omitempty"`
	Description         string `json:"Description,omitempty"`
	ChartName           string `json:"ChartName,omitempty"`
	RepoName            string `json:"RepoName,omitempty"`
	ReleaseName         string `json:"ReleaseName,omitempty"`
	RepoURL             string `json:"RepoURL,omitempty"`
	Namespace           string `json:"Namespace,omitempty"`
	CreateNamespace     bool   `json:"CreateNamespace"`
	PrivilegedNamespace bool   `json:"PrivilegedNamespace"`
	Icon                string `json:"Icon,omitempty"`
	LaunchURL           string `json:"LaunchURL,omitempty"`
	LaunchUIDescription string `json:"LaunchUIDescription,omitempty"`
}
