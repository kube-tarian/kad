package pluginconfigstore

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/common-pkg/cluster-plugins/clusterpluginspb"
	"github.com/pkg/errors"
)

const (
	insertPluginConfigByReleaseNameQuery = "INSERT INTO %s.ClusterPluginConfig(plugin_name) VALUES (?)"
	updatePluginConfigByReleaseNameQuery = "UPDATE %s.ClusterPluginConfig SET %s WHERE plugin_name = ?"
	deletePluginConfigByReleaseNameQuery = "DELETE FROM %s.ClusterPluginConfig WHERE plugin_name= ? "
)

const (
	storeType, chartRepo, defaultNamespace = "store_type", "chart_repo", "default_namespace"
	capabilities, values                   = "capabilities", "values"
	action, uiEndpoint                     = "action", "ui_endpoint"
)

var (
	pluginConfigfields = []string{
		storeType, pluginName, description, category,
		version, icon, chartName, chartRepo,
		defaultNamespace, privilegedNamespace, apiEndpoint,
		uiEndpoint, capabilities, values, overrideValues,
		installStatus, updateTime,
	}
)

func CreatePluginConfigSelectByFieldNameQuery(keyspace, field string) string {
	return CreatePluginConfigSelectAllQuery(keyspace) + fmt.Sprintf(" WHERE %s = ?", field)
}

func CreatePluginConfigSelectAllQuery(keyspace string) string {
	return fmt.Sprintf("SELECT %s FROM %s.ClusterPluginConfig", strings.Join(pluginConfigfields, ", "), keyspace)
}

func (a *Store) UpsertPluginConfig(config *PluginConfig) error {
	if len(config.PluginName) == 0 {
		return fmt.Errorf("app release name empty")
	}

	kvPairs, isEmptyUpdate := formPluginUpdateKvPairs(config)
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertPluginConfigByReleaseNameQuery, a.keyspace), config.PluginName)
	if !isEmptyUpdate {
		batch.Query(fmt.Sprintf(updatePluginConfigByReleaseNameQuery, a.keyspace, kvPairs), config.PluginName)
	}
	return a.client.Session().ExecuteBatch(batch)
}

func (a *Store) DeletePluginConfigByPluginName(releaseName string) error {

	deleteQuery := a.client.Session().Query(fmt.Sprintf(deletePluginConfigByReleaseNameQuery,
		a.keyspace), releaseName)

	err := deleteQuery.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (a *Store) GetPluginConfig(pluginNameKey string) (*PluginConfig, error) {
	a.log.Debugf("Select query: %v", CreatePluginConfigSelectByFieldNameQuery(a.keyspace, pluginName))
	selectQuery := a.client.Session().Query(CreatePluginConfigSelectByFieldNameQuery(a.keyspace, pluginName), pluginNameKey)

	config := &PluginConfig{}
	var valuesValue, overrideValuesValue, capabilitiesValue string

	if err := selectQuery.Scan(
		&config.StoreType, &config.PluginName, &config.Description,
		&config.Category, &config.Version, &config.Icon, &config.ChartName,
		&config.ChartRepo, &config.DefaultNamespace, &config.PrivilegedNamespace,
		&config.ApiEndpoint, &config.UiEndpoint, &capabilitiesValue, &valuesValue,
		&overrideValuesValue, &config.InstallStatus, &config.LastUpdateTime,
	); err != nil {
		return nil, err
	}

	config.Values, _ = base64.StdEncoding.DecodeString(valuesValue)
	config.OverrideValues, _ = base64.StdEncoding.DecodeString(overrideValuesValue)
	capabilityList := strings.Split(capabilitiesValue, ",")
	config.Capabilities = capabilityList
	// config.StoreType = clusterpluginspb.StoreType(storeTypeValue)

	return config, nil
}

func (a *Store) GetAllPlugins() ([]*clusterpluginspb.Plugin, error) {
	selectAllQuery := a.client.Session().Query(CreatePluginConfigSelectAllQuery(a.keyspace))
	iter := selectAllQuery.Iter()

	config := &PluginConfig{}
	var valuesValue, overrideValuesValue, capabilitiesValue string

	ret := make([]*clusterpluginspb.Plugin, 0)
	for iter.Scan(
		&config.StoreType, &config.PluginName, &config.Description,
		&config.Category, &config.Version, &config.Icon, &config.ChartName,
		&config.ChartRepo, &config.DefaultNamespace, &config.PrivilegedNamespace,
		&config.ApiEndpoint, &config.UiEndpoint, &capabilitiesValue, &valuesValue,
		&overrideValuesValue, &config.InstallStatus, &config.LastUpdateTime,
	) {
		configCopy := config.Plugin
		configCopy.Values, _ = base64.StdEncoding.DecodeString(values)
		configCopy.Capabilities = strings.Split(capabilitiesValue, ",")
		ret = append(ret, &configCopy)
	}

	if err := iter.Close(); err != nil {
		return nil, errors.WithMessage(err, "failed to iterate through results:")
	}
	return ret, nil
}

func formPluginUpdateKvPairs(config *PluginConfig) (string, bool) {
	params := []string{}

	if config.Values != nil {
		encoded := base64.StdEncoding.EncodeToString(config.Values)
		params = append(params,
			fmt.Sprintf("%s = '%s'", values, encoded))
	}
	if config.OverrideValues != nil {
		encoded := base64.StdEncoding.EncodeToString(config.OverrideValues)
		params = append(params,
			fmt.Sprintf("%s = '%s'", overrideValues, encoded))
	}

	if config.PrivilegedNamespace {
		params = append(params,
			fmt.Sprintf("%s = true", privilegedNamespace))
	}

	params = append(params,
		fmt.Sprintf("%s = %d", storeType, config.StoreType))

	if config.Description != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", description, config.Description))
	}
	if config.Category != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", category, config.Category))
	}

	if config.ChartName != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", chartName, config.ChartName))
	}
	if config.ChartRepo != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", chartRepo, config.ChartRepo))
	}

	if config.DefaultNamespace != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", defaultNamespace, config.DefaultNamespace))
	}
	if config.Version != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", version, config.Version))
	}

	if config.Icon != nil && len(config.Icon) > 0 {
		params = append(params,
			fmt.Sprintf("%s = 0x%s", icon, hex.EncodeToString(config.Icon)))
	}
	if len(config.InstallStatus) > 0 {
		params = append(params,
			fmt.Sprintf("%s = '%s'", installStatus, config.InstallStatus))
	}
	if config.ApiEndpoint != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", apiEndpoint, config.ApiEndpoint))
	}
	if config.UiEndpoint != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", uiEndpoint, config.UiEndpoint))
	}
	if len(config.Capabilities) > 0 {
		capabilitiesList := strings.Join(config.Capabilities, ",")
		params = append(params,
			fmt.Sprintf("%s = '%s'", capabilities, capabilitiesList))
	}

	params = append(params,
		fmt.Sprintf("%s = '%s'", updateTime, time.Now().Format(time.RFC3339)))

	return strings.Join(params, ", "), false
}
