package astra

const (
	createCaptenClusterTableQuery     = "CREATE TABLE IF NOT EXISTS %s.capten_clusters (cluster_id uuid, org_id uuid, cluster_name text, endpoint text, PRIMARY KEY (org_id, cluster_id));"
	createPluginStoreTableQuery       = "CREATE TABLE IF NOT EXISTS %s.plugin_data(git_project_id uuid, plugin_name TEXT, last_updated_time timestamp, store_type INT, description TEXT, category TEXT, icon TEXT, versions LIST<TEXT>, PRIMARY KEY (git_project_id, plugin_name));"
	createPluginStoreConfigTableQuery = "CREATE TABLE IF NOT EXISTS %s.plugin_store_config(cluster_id uuid, store_type INT, git_project_id TEXT, git_project_url TEXT, last_updated_time timestamp, PRIMARY KEY (cluster_id, store_type));"
	dropCaptenClusterTableQuery       = "DROP TABLE IF EXISTS %s.capten_clusters;"
	dropPluginStoreTableQuery         = "DROP TABLE IF EXISTS %s.plugin_data;"
	dropPluginStoreConfigTableQuery   = "DROP TABLE IF EXISTS %s.plugin_store_config;"
)
