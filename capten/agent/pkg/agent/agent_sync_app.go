package agent

import (
	"context"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"gopkg.in/yaml.v2"
)

func (a Agent) syncApp(ctx context.Context, request *agentpb.SyncAppRequest) error {

	var appConfig AppConfig
	if err := yaml.Unmarshal(request.Payload, &appConfig); err != nil {
		a.log.Errorf("could not unmarshal appConfig yaml: %v", err)
		return err
	}

	// Note: We are creating keyspace and table when needed
	if err := a.store.Session().Query(fmt.Sprintf(createKeyspaceCQL, "1")).Exec(); err != nil {
		a.log.Errorf("could not create keyspace, err: %v", err)
		return err
	}

	if err := a.store.Session().Query(createTableQuery).Exec(); err != nil {
		a.log.Errorf("could not create table, err: %v", err)
		return err
	}

	if err := insert(a.store.Session(), appConfig); err != nil {
		a.log.Errorf("could not insert, err: %v", err)
		return err
	}

	return nil
}

// Todo: These structs along with proto file should be a part of common package

type LaunchUIConfig struct {
	RedirectURL string `yaml:"RedirectURL"`
}

type Override struct {
	LaunchUIConfig LaunchUIConfig `yaml:"LaunchUIConfig"`
	LaunchUIValues map[string]any `yaml:"LaunchUIValues"`
	Values         map[string]any `yaml:"Values"`
}

type AppConfig struct {
	Name                string   `yaml:"Name"`
	ChartName           string   `yaml:"ChartName"`
	RepoName            string   `yaml:"RepoName"`
	RepoURL             string   `yaml:"RepoURL"`
	Namespace           string   `yaml:"Namespace"`
	ReleaseName         string   `yaml:"ReleaseName"`
	Version             string   `yaml:"Version"`
	Override            Override `yaml:"Override"` // this can be marshled as json and stored as string
	CreateNamespace     bool     `yaml:"CreateNamespace"`
	PrivilegedNamespace bool     `yaml:"PrivilegedNamespace"`
}

const createKeyspaceCQL = "CREATE KEYSPACE IF NOT EXISTS Apps WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : %s } AND DURABLE_WRITES = true"

const createTableQuery = `
CREATE TABLE IF NOT EXISTS Apps.Configs(
	name text PRIMARY KEY,
	chart_name text,
	repo_name text,
	repo_url text,
	namespace text,
	release_name text,
	version text,
	override text,
	create_namespace boolean,
	privileged_namespace boolean
) WITH bloom_filter_fp_chance = 0.01
AND caching = {'keys': 'ALL', 'rows_per_partition': 'NONE'}
AND comment = ''
AND compaction = {'class': 'org.apache.cassandra.db.compaction.LeveledCompactionStrategy', 'tombstone_compaction_interval': '1800', 'tombstone_threshold': '0.01', 'unchecked_tombstone_compaction': 'true'}
AND compression = {'chunk_length_in_kb': '64', 'class': 'org.apache.cassandra.io.compress.LZ4Compressor'}
AND crc_check_chance = 1.0
AND default_time_to_live = 0
AND gc_grace_seconds = 3600
AND max_index_interval = 2048
AND memtable_flush_period_in_ms = 0
AND min_index_interval = 128
AND speculative_retry = '99PERCENTILE';
`

const insertQuery = `
INSERT INTO Apps.Configs(
    name, chart_name, repo_name, repo_url, namespace, release_name, version, override, create_namespace, privileged_namespace
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

func insert(session *gocql.Session, config AppConfig) error {
	query := session.Query(insertQuery)
	override, err := yaml.Marshal(config.Override)
	if err != nil {
		return err
	}
	return query.Bind(
		config.Name,
		config.ChartName,
		config.RepoName,
		config.RepoURL,
		config.Namespace,
		config.ReleaseName,
		config.Version,
		string(override),
		// this is a complex mapping value, hence storing as string
		// while querying, we query and then unmarshal
		// Note: the value is already verified before inserting
		config.CreateNamespace,
		config.PrivilegedNamespace,
	).Exec()
}

func getAppsByName(session *gocql.Session, name string) (config AppConfig, err error) {

	selectQuery := session.Query(`
	SELECT name, chart_name, repo_name, repo_url, namespace, release_name, version, override, create_namespace, privileged_namespace
	FROM Apps.Configs WHERE name = ? LIMIT 1;
	`, name)

	var overrideYaml string

	if err = selectQuery.Scan(
		&config.Name,
		&config.ChartName,
		&config.RepoName,
		&config.RepoURL,
		&config.Namespace,
		&config.ReleaseName,
		&config.Version,
		&overrideYaml,
		&config.CreateNamespace,
		&config.PrivilegedNamespace,
	); err != nil {
		return
	}

	err = yaml.Unmarshal([]byte(overrideYaml), &config.Override)

	return

}
