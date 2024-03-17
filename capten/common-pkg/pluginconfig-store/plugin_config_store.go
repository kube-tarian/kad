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
	insertPluginConfigByReleaseNameQuery = "INSERT INTO %s.ClusterPluginConfig(release_name) VALUES (?)"
	updatePluginConfigByReleaseNameQuery = "UPDATE %s.ClusterPluginConfig SET %s WHERE release_name = ?"
	deletePluginConfigByReleaseNameQuery = "DELETE FROM %s.ClusterPluginConfig WHERE release_name= ? "
)

const (
	storeType, chartRepo, defaultNamespace = "store_type", "chart_repo", "default_namespace"
	pluginEndpoint, capabilities, values   = "plugin_endpoint", "capabilities", "values"
	action                                 = "action"
	appName, description, category         = "app_name", "description", "category"
	chartName, repoName, repoUrl           = "chart_name", "repo_name", "repo_url"
	namespace, releaseName, version        = "namespace", "release_name", "version"
	launchUrl, launchUIDesc                = "launch_url", "launch_redirect_url"
	createNamespace, privilegedNamespace   = "create_namespace", "privileged_namespace"
	overrideValues, launchUiValues         = "override_values", "launch_ui_values"
	templateValues, defaultApp             = "template_values", "default_app"
	icon, installStatus                    = "icon", "install_status"
	updateTime                             = "update_time"
	usecase, projectUrl, status, details   = "usecase", "project_url", "status", "details"
	pluginName, pluginDescription          = "plugin_name", "plugin_description"
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

func (a *Store) UpsertPluginConfig(config *PluginConfig) error {
	if len(config.PluginName) == 0 {
		return fmt.Errorf("app release name empty")
	}

	kvPairs, isEmptyUpdate := formUpdateKvPairs(config)
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertPluginConfigByReleaseNameQuery, a.keyspace), config.PluginName)
	if !isEmptyUpdate {
		batch.Query(fmt.Sprintf(updatePluginConfigByReleaseNameQuery, a.keyspace, kvPairs), config.PluginName)
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

func (a *Store) GetPluginConfig(appReleaseName string) (*PluginConfig, error) {
	selectQuery := a.client.Session().Query(CreatePluginConfigSelectByFieldNameQuery(a.keyspace, releaseName), appReleaseName)

	config := &PluginConfig{}
	var values string

	if err := selectQuery.Scan(
		&config.StoreType, &config.PluginName, &config.Description,
		&config.Category, &config.Version, &config.Icon, &config.ChartName,
		&config.ChartRepo, &config.DefaultNamespace, &config.PrivilegedNamespace,
		&config.ApiEndpoint, &config.Capabilities, &values,
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
		&config.ApiEndpoint, &config.Capabilities, &values,
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

func formUpdateKvPairs(config *PluginConfig) (string, bool) {
	params := []string{}

	if config.Values != nil {
		encoded := base64.StdEncoding.EncodeToString(config.Values)
		params = append(params,
			fmt.Sprintf("%s = '%s'", values, encoded))
	}

	if config.PrivilegedNamespace {
		params = append(params,
			fmt.Sprintf("%s = true", privilegedNamespace))
	}

	if config.PluginName != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", appName, config.PluginName))
	}
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
	if config.ChartName != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", repoName, config.ChartName))
	}
	if config.ChartRepo != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", repoUrl, config.ChartRepo))
	}

	if config.DefaultNamespace != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", namespace, config.DefaultNamespace))
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
			fmt.Sprintf("%s = '%s'", pluginEndpoint, config.ApiEndpoint))
	}

	params = append(params,
		fmt.Sprintf("%s = '%s'", updateTime, time.Now().Format(time.RFC3339)))

	if len(params) == 0 {
		// query is empty there is nothing to update
		return "", true
	}

	params = append(params,
		fmt.Sprintf("%s = '%s'", pluginName, "helm"))

	params = append(params,
		fmt.Sprintf("%s = true", createNamespace))

	return strings.Join(params, ", "), false
}
