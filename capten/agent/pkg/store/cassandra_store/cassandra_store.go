package cassandra_store

import (
	"fmt"
	"strings"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/agent/pkg/store"
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
		FROM %s.%s WHERE name = ? LIMIT 1`

	updateAppConfigByNameQuery string = `
		UPDATE %s.%s SET %s WHERE name = ?
		`
)

var _ store.StoreIface = CassandraStore{}

type CassandraStore struct {
	cassandra.Store
}

func New(cassandraStore cassandra.Store) *CassandraStore {
	return &CassandraStore{Store: cassandraStore}
}

func (a CassandraStore) GetAppConfigByName(name string) (config types.AppConfig, err error) {

	selectQuery := a.Store.Session().Query(fmt.Sprintf(selectAppConfigByNameQuery, AppKeyspace, AppTable), name)

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

	if overrideYaml != "" {
		err = yaml.Unmarshal([]byte(overrideYaml), &config.Override)
	}

	return
}

func (a CassandraStore) InsertAppConfig(config types.AppConfig) error {
	query := a.Store.Session().Query(fmt.Sprintf(insertAppConfigQuery, AppKeyspace, AppTable))
	var override []byte
	var err error
	if config.Override != nil {
		override, err = yaml.Marshal(config.Override)
		if err != nil {
			return err
		}
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
		config.CreateNamespace,
		config.PrivilegedNamespace,
	).Exec()
}

func (a CassandraStore) UpsertAppConfig(config types.AppConfig) error {

	if config.Name == "" {
		return fmt.Errorf("missing AppConfig.name")
	}

	kvPairs, err := formUpdateKvPairs(&config)
	if err != nil {
		return err
	}

	var insertNameOnlyAppConfigQuery string = `
	INSERT INTO %s.%s(
		name
	) VALUES (?)`

	batch := a.Store.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertNameOnlyAppConfigQuery, AppKeyspace, AppTable), config.Name)
	batch.Query(fmt.Sprintf(updateAppConfigByNameQuery, AppKeyspace, AppTable, kvPairs), config.Name)
	return a.Store.Session().ExecuteBatch(batch)

}

func (a CassandraStore) CreateDB(name string) error {
	return a.Store.CreateDb(name, "1")
}

func (a CassandraStore) DropDB(name string) error {
	return a.Store.Session().Query(fmt.Sprintf("DROP KEYSPACE IF EXISTS %s", name)).Exec()
}

func (a CassandraStore) CloseSession() {
	a.Store.Session().Close()
}

func formUpdateKvPairs(config *types.AppConfig) (string, error) {

	params := []string{}

	if config.Override != nil {
		marshaled, err := yaml.Marshal(config.Override)
		if err != nil {
			return "", err
		}
		params = append(params, fmt.Sprintf("override = '%s'", marshaled))
	}

	if config.CreateNamespace != nil {
		val := "false"
		if *config.CreateNamespace {
			val = "true"
		}
		params = append(params, fmt.Sprintf("create_namespace = '%s'", val))
	}

	if config.PrivilegedNamespace != nil {
		val := "false"
		if *config.PrivilegedNamespace {
			val = "true"
		}
		params = append(params, fmt.Sprintf("privileged_namespace = '%s'", val))
	}

	if config.ChartName != "" {
		params = append(params, fmt.Sprintf("chart_name = '%s'", config.ChartName))
	}

	if config.Namespace != "" {
		params = append(params, fmt.Sprintf("namespace = '%s'", config.Namespace))
	}

	if config.ReleaseName != "" {
		params = append(params, fmt.Sprintf("release_name = '%s'", config.ReleaseName))
	}

	if config.RepoName != "" {
		params = append(params, fmt.Sprintf("repo_name = '%s'", config.RepoName))
	}

	if config.RepoURL != "" {
		params = append(params, fmt.Sprintf("repo_url = '%s'", config.RepoURL))
	}

	if config.Version != "" {
		params = append(params, fmt.Sprintf("version = '%s'", config.Version))
	}

	return strings.Join(params, ", "), nil

}
