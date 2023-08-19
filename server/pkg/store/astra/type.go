package astra

import (
	"encoding/json"

	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
)

const (
	createKeyspaceQuery             = "CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1};"
	createClusterEndpointTableQuery = "CREATE TABLE IF NOT EXISTS %s.capten_clusters (cluster_id uuid, org_id uuid, cluster_name text, endpoint text, PRIMARY KEY (cluster_id, org_id));"
	createAppConfigTableQuery       = "CREATE TABLE IF NOT EXISTS %s.app_config(id TEXT, created_time timestamp, last_updated_time timestamp, last_updated_user TEXT, name TEXT, chart_name TEXT, repo_name TEXT, release_name TEXT, repo_url TEXT, namespace TEXT, version TEXT, create_namespace BOOLEAN, privileged_namespace BOOLEAN, launch_ui_url TEXT, launch_ui_redirect_url TEXT, category TEXT, icon TEXT, description TEXT, launch_ui_values TEXT, override_values TEXT, PRIMARY KEY (name, version));"
)

const (
	appStoreConfigFileName string = "storeconfig.yaml"
)

var (
	UuidSetSpec = &pb.TypeSpec{
		Spec: &pb.TypeSpec_Set_{
			Set: &pb.TypeSpec_Set{
				Element: &pb.TypeSpec{
					Spec: &pb.TypeSpec_Basic_{
						Basic: pb.TypeSpec_UUID,
					},
				},
			},
		},
	}
)

type AppStoreConfig struct {
	CreateStoreApps []string `yaml:"CreateStoreApps"`
	UpdateStoreApps []string `yaml:"UpdateStoreApps"`
}

type AppConfig struct {
	Name                string          `yaml:"Name"`
	ChartName           string          `yaml:"ChartName"`
	Category            string          `yaml:"Category"`
	Description         string          `yaml:"Description"`
	RepoName            string          `yaml:"RepoName"`
	RepoURL             string          `yaml:"RepoURL"`
	Namespace           string          `yaml:"Namespace"`
	ReleaseName         string          `yaml:"ReleaseName"`
	Version             string          `yaml:"Version"`
	Icon                string          `yaml:"Icon"`
	CreateNamespace     bool            `yaml:"CreateNamespace"`
	PrivilegedNamespace bool            `yaml:"PrivilegedNamespace"`
	OverrideValues      json.RawMessage `yaml:"OverrideValues"`
	LaunchUIValues      json.RawMessage `yaml:"LaunchUIValues"`
	LaunchURL           string          `yaml:"LaunchURL"`
	LaunchRedirectURL   string          `yaml:"LaunchRedirectURL"`
}
