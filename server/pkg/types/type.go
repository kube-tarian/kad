package types

import "time"

const (
	ClientCertChainFileName = "cert-chain.pem"
	ClientCertFileName      = "client.crt"
	ClientKeyFileName       = "client.key"
	AgentPortCfgKey         = "agent.port"
	AgentTlsEnabledCfgKey   = "agent.tlsEnabled"
	ServerDbCfgKey          = "server.db"
)

type AgentInfo struct {
	Endpoint string
	CaPem    string
	Cert     string
	Key      string
}

type ClusterDetails struct {
	OrgID       string
	ClusterID   string
	ClusterName string
	Endpoint    string
}

type StoreAppConfig struct {
	AppName             string `json:"appName,omitempty"`
	Version             string `json:"version,omitempty"`
	Category            string `json:"category,omitempty"`
	Description         string `json:"description,omitempty"`
	ChartName           string `json:"chartName,omitempty"`
	RepoName            string `json:"repoName,omitempty"`
	ReleaseName         string `json:"releaseName,omitempty"`
	RepoURL             string `json:"repoURL,omitempty"`
	Namespace           string `json:"namespace,omitempty"`
	CreateNamespace     bool   `json:"createNamespace"`
	PrivilegedNamespace bool   `json:"privilegedNamespace"`
	Icon                string `json:"icon,omitempty"`
	LaunchURL           string `json:"launchURL,omitempty"`
	LaunchRedirectURL   string `json:"launchRedirectURL,omitempty"`
	OverrideValues      string `json:"overrideValues,omitempty"`
	LaunchUIValues      string `json:"launchUIValues,omitempty"`
}

type AppConfig struct {
	ID                  int64     `cql:"id" json:"id,omitempty"`
	CreatedTime         time.Time `cql:"created_time" json:"created_time,omitempty"`
	LastUpdatedTime     time.Time `cql:"last_updated_time" json:"last_updated_time,omitempty"`
	LastUpdatedUser     string    `cql:"last_updated_user" json:"last_updated_user,omitempty"`
	Name                string    `cql:"name" json:"name"`
	ChartName           string    `cql:"chart_name" json:"chart_name"`
	RepoName            string    `cql:"repo_name" json:"repo_name"`
	ReleaseName         string    `cql:"release_name" json:"release_name"`
	RepoURL             string    `cql:"repo_url" json:"repo_url"`
	Namespace           string    `cql:"namespace" json:"namespace"`
	Version             string    `cql:"version" json:"version"`
	CreateNamespace     bool      `cql:"create_namespace" json:"create_namespace"`
	PrivilegedNamespace bool      `cql:"privileged_namespace" json:"privileged_namespace"`
	LaunchUIURL         string    `cql:"launch_ui_url" json:"launch_ui_url"`
	LaunchUIRedirectURL string    `cql:"launch_ui_redirect_url" json:"launch_ui_redirect_url"`
	Category            string    `cql:"category" json:"category"`
	Icon                string    `cql:"icon" json:"icon"`
	Description         string    `cql:"description" json:"description"`
	LaunchUIValues      string    `cql:"launch_ui_values" json:"launch_ui_values"`
	OverrideValues      string    `cql:"override_values" json:"override_values"`
}
