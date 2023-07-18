package captenstore

import (
	"fmt"
	"strings"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/agent/pkg/types"
	"gopkg.in/yaml.v2"
)

const (
	AppConfigTable string = "AppConfig"

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
		UPDATE %s.%s SET %s WHERE name = ?`
)

func (a *Store) AddAppConfig(config types.AppConfig) error {
	query := a.client.Session().Query(fmt.Sprintf(insertAppConfigQuery, a.conf.Keyspace, AppConfigTable))
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

func (a *Store) UpdateAppConfig(config types.AppConfig) error {
	if config.Name == "" {
		return fmt.Errorf("invalid attribute")
	}

	kvPairs, err := formUpdateKvPairs(&config)
	if err != nil {
		return err
	}

	var insertNameOnlyAppConfigQuery string = `
	INSERT INTO %s.%s(
		name
	) VALUES (?)`

	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertNameOnlyAppConfigQuery, a.conf.Keyspace, AppConfigTable), config.Name)
	batch.Query(fmt.Sprintf(updateAppConfigByNameQuery, a.conf.Keyspace, AppConfigTable, kvPairs), config.Name)
	return a.client.Session().ExecuteBatch(batch)

}

func (a *Store) GetAppConfig(appName string) (config types.AppConfig, err error) {
	selectQuery := a.client.Session().Query(fmt.Sprintf(selectAppConfigByNameQuery, a.conf.Keyspace, AppConfigTable), appName)
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
