package astra

const (
	createCaptenClusterTableQuery = "CREATE TABLE IF NOT EXISTS %s.capten_clusters (cluster_id uuid, org_id uuid, cluster_name text, endpoint text, PRIMARY KEY (org_id, cluster_id));"
	createAppConfigTableQuery     = "CREATE TABLE IF NOT EXISTS %s.store_app_config(id TEXT, created_time timestamp, last_updated_time timestamp, last_updated_user TEXT, name TEXT, chart_name TEXT, repo_name TEXT, release_name TEXT, repo_url TEXT, namespace TEXT, version TEXT, create_namespace BOOLEAN, privileged_namespace BOOLEAN, launch_ui_url TEXT, launch_ui_redirect_url TEXT, category TEXT, icon TEXT, description TEXT, launch_ui_values TEXT, override_values TEXT, template_values TEXT, plugin_name TEXT, plugin_description TEXT, PRIMARY KEY (name, version));"
	dropCaptenClusterTableQuery   = "DROP TABLE IF EXISTS %s.capten_clusters;"
	dropAppConfigTableQuery       = "DROP TABLE IF EXISTS %s.store_app_config;"
)

type StoreAppConfig struct {
	Name                string                 `yaml:"Name"`
	ChartName           string                 `yaml:"ChartName"`
	Category            string                 `yaml:"Category"`
	Description         string                 `yaml:"Description"`
	RepoName            string                 `yaml:"RepoName"`
	RepoURL             string                 `yaml:"RepoURL"`
	Namespace           string                 `yaml:"Namespace"`
	ReleaseName         string                 `yaml:"ReleaseName"`
	Version             string                 `yaml:"Version"`
	Icon                string                 `yaml:"Icon"`
	CreateNamespace     bool                   `yaml:"CreateNamespace"`
	PrivilegedNamespace bool                   `yaml:"PrivilegedNamespace"`
	OverrideValues      map[string]interface{} `yaml:"OverrideValues"`
	LaunchUIValues      map[string]interface{} `yaml:"LaunchUIValues"`
	LaunchURL           string                 `yaml:"LaunchURL"`
	LaunchUIDescription string                 `yaml:"LaunchUIDescription"`
}
