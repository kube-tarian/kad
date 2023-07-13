package store

import (
	"fmt"

	"github.com/kube-tarian/kad/capten/agent/pkg/types"
	"github.com/kube-tarian/kad/capten/common-pkg/db-create/cassandra"
	"gopkg.in/yaml.v2"
)

const (
	AppKeyspace string = "Apps"
	AppTable    string = "Configs"

	insertAppConfigQuery string = `
		INSERT INTO %s.%s(
			name, chart_name, repo_name, repo_url,
			namespace, release_name, version, override,
			create_namespace, privileged_namespace
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	selectAppConfigByNameQuery string = `
		SELECT name, chart_name, repo_name, repo_url, 
				namespace, release_name, version, override,
				create_namespace, privileged_namespace
		FROM %s.%s WHERE name = ? LIMIT 1;`
)

// get update delete
type Store struct {
	cassandraStore cassandra.Store
}

func New(cassandraStore cassandra.Store) *Store {
	return &Store{cassandraStore: cassandraStore}
}

func (a Store) GetAppConfigByName(name string) (config types.AppConfig, err error) {

	selectQuery := a.cassandraStore.Session().Query(fmt.Sprintf(selectAppConfigByNameQuery, AppKeyspace, AppTable), name)

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

func (a Store) InsertAppConfig(config types.AppConfig) error {
	query := a.cassandraStore.Session().Query(fmt.Sprintf(insertAppConfigQuery, AppKeyspace, AppTable))
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

func (a Store) CreateKeyspace(name, replicationFactor string) error {
	return a.cassandraStore.CreateDb(name, replicationFactor)
}

func (a Store) DropKeyspace(name string) error {
	return a.cassandraStore.Session().Query(fmt.Sprintf("DROP KEYSPACE IF EXISTS %s", name)).Exec()
}

// Todo: remove this from here, not needed, only being used for migration purposes
func (a Store) CreateLockSchemaDb(replicationFactor string) error {
	return a.cassandraStore.CreateLockSchemaDb(replicationFactor)
}
