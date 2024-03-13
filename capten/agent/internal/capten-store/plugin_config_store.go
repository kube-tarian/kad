package captenstore

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/agent/internal/pb/agentpb"
	"github.com/kube-tarian/kad/capten/common-pkg/cluster-plugins/clusterpluginspb"
	"github.com/pkg/errors"
)

const (
	insertPluginConfigByReleaseNameQuery = "INSERT INTO %s.ClusterPluginConfig(release_name) VALUES (?)"
	updatePluginConfigByReleaseNameQuery = "UPDATE %s.ClusterPluginConfig SET %s WHERE release_name = ?"
	deletePluginConfigByReleaseNameQuery = "DELETE FROM %s.ClusterPluginConfig WHERE release_name= ? "
)

const (
	storeType, chartRepo, defaultNamespace = "store_type", "chart_repo", "default_namespace"
	pluginEndpoint, capabilities, values   = "plugin_endpoint", "capabilities", "values"
	action                                 = "action"
)

var (
	pluginConfigfields = []string{
		storeType, pluginName, description, category,
		version, icon, chartName, chartRepo,
		defaultNamespace, privilegedNamespace, pluginEndpoint,
		capabilities, values, action,
	}
)

func CreatePluginConfigSelectByFieldNameQuery(keyspace, field string) string {
	return CreatePluginConfigSelectAllQuery(keyspace) + fmt.Sprintf(" WHERE %s = ?", field)
}

func CreatePluginConfigSelectAllQuery(keyspace string) string {
	return fmt.Sprintf("SELECT %s FROM %s.ClusterPluginConfig", strings.Join(pluginConfigfields, ", "), keyspace)
}

func (a *Store) UpsertPluginConfig(config *agentpb.SyncAppData) error {
	if len(config.Config.ReleaseName) == 0 {
		return fmt.Errorf("app release name empty")
	}

	kvPairs, isEmptyUpdate := formUpdateKvPairs(config)
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertPluginConfigByReleaseNameQuery, a.keyspace), config.Config.ReleaseName)
	if !isEmptyUpdate {
		batch.Query(fmt.Sprintf(updatePluginConfigByReleaseNameQuery, a.keyspace, kvPairs), config.Config.ReleaseName)
	}
	return a.client.Session().ExecuteBatch(batch)
}
func (a *Store) DeletePluginConfigByReleaseName(releaseName string) error {

	deleteQuery := a.client.Session().Query(fmt.Sprintf(deletePluginConfigByReleaseNameQuery,
		a.keyspace), releaseName)

	err := deleteQuery.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (a *Store) GetPluginConfig(appReleaseName string) (*clusterpluginspb.Plugin, error) {
	selectQuery := a.client.Session().Query(CreatePluginConfigSelectByFieldNameQuery(a.keyspace, releaseName), appReleaseName)

	config := &clusterpluginspb.Plugin{}
	var values string

	if err := selectQuery.Scan(
		&config.StoreType, &config.PluginName, &config.Description,
		&config.Category, &config.Version, &config.Icon, &config.ChartName,
		&config.ChartRepo, &config.DefaultNamespace, &config.PrivilegedNamespace,
		&config.PluginEndpoint, &config.Capabilities, &values,
	); err != nil {
		return nil, err
	}

	config.Values, _ = base64.StdEncoding.DecodeString(values)

	return config, nil
}

func (a *Store) GetAllPlugins() ([]*clusterpluginspb.Plugin, error) {
	selectAllQuery := a.client.Session().Query(CreatePluginConfigSelectAllQuery(a.keyspace))
	iter := selectAllQuery.Iter()

	config := clusterpluginspb.Plugin{}
	var values string

	ret := make([]*clusterpluginspb.Plugin, 0)
	for iter.Scan(
		&config.StoreType, &config.PluginName, &config.Description,
		&config.Category, &config.Version, &config.Icon, &config.ChartName,
		&config.ChartRepo, &config.DefaultNamespace, &config.PrivilegedNamespace,
		&config.PluginEndpoint, &config.Capabilities, &values,
	) {
		configCopy := config
		configCopy.Values, _ = base64.StdEncoding.DecodeString(values)
		ret = append(ret, &configCopy)
	}

	if err := iter.Close(); err != nil {
		return nil, errors.WithMessage(err, "failed to iterate through results:")
	}
	return ret, nil
}
