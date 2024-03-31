package astra

import (
	"fmt"
	"strings"
	"time"

	"github.com/kube-tarian/kad/server/pkg/pb/pluginstorepb"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
)

const (
	insertStoreConfig           = `INSERT INTO %s.plugin_store_config (cluster_id, store_type, git_project_id, git_project_url, last_updated_time) VALUES (%s, %d, '%s', '%s', '%s') IF NOT EXISTS`
	updateStoreConfig           = `UPDATE %s.plugin_store_config SET git_project_id = '%s', git_project_url = '%s', last_updated_time = '%s' WHERE cluster_id = %s and store_type = %d`
	readStoreConfigForStoreType = `SELECT git_project_id, git_project_url, last_updated_time FROM %s.plugin_store_config WHERE cluster_id = %s and store_type = %d`

	insertPluginData            = `INSERT INTO %s.plugin_data (git_project_id, plugin_name, last_updated_time, store_type, description, category, icon, versions) VALUES (%s, '%s', '%s', %d, '%s', '%s', '%s', %v) IF NOT EXISTS`
	updatePluginData            = `UPDATE %s.plugin_data SET last_updated_time = '%s', store_type = %d, description = '%s', category = '%s', icon = '%s', versions = %v WHERE git_project_id = %s and plugin_name = '%s'`
	readPlugins                 = `SELECT plugin_name, last_updated_time, store_type, description, category, icon, versions FROM %s.plugin_data WHERE git_project_id = %s`
	readPluginDataForPluginName = `SELECT plugin_name, last_updated_time, store_type, description, category, icon, versions FROM %s.plugin_data WHERE git_project_id = %s and plugin_name = '%s'`
	deletePluginData            = "DELETE FROM %s.plugin_data WHERE git_project_id = %s and plugin_name = '%s'"
)

func (a *AstraServerStore) WritePluginStoreConfig(clusterId string, config *pluginstorepb.PluginStoreConfig) error {
	query := &pb.Query{
		Cql: fmt.Sprintf(insertStoreConfig,
			a.keyspace, clusterId, config.StoreType, config.GitProjectId, config.GitProjectURL,
			time.Now().Format(time.RFC3339)),
	}

	resp, err := a.c.Session().ExecuteQuery(query)
	if err != nil {
		return fmt.Errorf("failed to insert/update the store config, %s, %w", query.Cql, err)
	}

	applied, err := client.ToBoolean(resp.GetResultSet().Rows[0].Values[0])
	if err != nil {
		return fmt.Errorf("failed to check execution status, %w", err)
	}
	if applied {
		a.log.Debug("inserted store config")
		return nil
	}

	query = &pb.Query{
		Cql: fmt.Sprintf(updateStoreConfig,
			a.keyspace, config.GitProjectId, config.GitProjectURL,
			time.Now().Format(time.RFC3339), clusterId, config.StoreType),
	}

	_, err = a.c.Session().ExecuteQuery(query)
	if err != nil {
		return fmt.Errorf("failed to insert/update the store config, %s, %w", query.Cql, err)
	}
	a.log.Debug("update store config")
	return nil
}

func (a *AstraServerStore) ReadPluginStoreConfig(clusterId string, storeType pluginstorepb.StoreType) (*pluginstorepb.PluginStoreConfig, error) {
	selectQuery := &pb.Query{
		Cql: fmt.Sprintf(readStoreConfigForStoreType, a.keyspace, clusterId, storeType),
	}

	response, err := a.c.Session().ExecuteQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to read store config, %s, %w", selectQuery.Cql, err)
	}

	result := response.GetResultSet()
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("not found")
	}

	config, err := toPluginStoreConfig(result.Rows[0])
	if err != nil {
		return nil, err
	}
	config.StoreType = storeType
	return config, nil
}

func toPluginStoreConfig(row *pb.Row) (*pluginstorepb.PluginStoreConfig, error) {
	gitProjectId, err := client.ToString(row.Values[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get chart name: %w", err)
	}
	gitProjectURL, err := client.ToString(row.Values[1])
	if err != nil {
		return nil, fmt.Errorf("failed to get chart name: %w", err)
	}

	return &pluginstorepb.PluginStoreConfig{
		GitProjectId:  gitProjectId,
		GitProjectURL: gitProjectURL,
	}, nil
}

func (a *AstraServerStore) WritePluginData(gitProjectId string, pluginData *pluginstorepb.PluginData) error {
	query := &pb.Query{
		Cql: fmt.Sprintf(insertPluginData,
			a.keyspace, gitProjectId, pluginData.PluginName, time.Now().Format(time.RFC3339),
			pluginData.StoreType, pluginData.Description, pluginData.Category, pluginData.Icon, getSQLStringArray(pluginData.Versions)),
	}

	resp, err := a.c.Session().ExecuteQuery(query)
	if err != nil {
		return fmt.Errorf("failed to insert the plugin data, %s, %w", query.Cql, err)
	}

	applied, err := client.ToBoolean(resp.GetResultSet().Rows[0].Values[0])
	if err != nil {
		return fmt.Errorf("failed to check execution status, %w", err)
	}
	if applied {
		a.log.Debugf("inserted store plugin %s", pluginData.PluginName)
		return nil
	}

	query = &pb.Query{
		Cql: fmt.Sprintf(updatePluginData,
			a.keyspace, time.Now().Format(time.RFC3339),
			pluginData.StoreType, pluginData.Description, pluginData.Category, pluginData.Icon,
			getSQLStringArray(pluginData.Versions), gitProjectId, pluginData.PluginName),
	}
	_, err = a.c.Session().ExecuteQuery(query)
	if err != nil {
		return fmt.Errorf("failed to update the plugin data, %s, %w", query.Cql, err)
	}
	a.log.Debugf("updated store plugin %s", pluginData.PluginName)
	return nil
}

func (a *AstraServerStore) ReadPluginData(gitProjectId string, pluginName string) (*pluginstorepb.PluginData, error) {
	selectQuery := &pb.Query{
		Cql: fmt.Sprintf(readPluginDataForPluginName, a.keyspace, gitProjectId, pluginName),
	}

	response, err := a.c.Session().ExecuteQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin data, %w", err)
	}

	result := response.GetResultSet()
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("not found")
	}

	pluginData, err := toPluginData(result.Rows[0], result.Columns)
	if err != nil {
		return nil, err
	}

	return pluginData, nil
}

func (a *AstraServerStore) ReadPlugins(gitProjectId string) ([]*pluginstorepb.Plugin, error) {
	selectQuery := &pb.Query{
		Cql: fmt.Sprintf(readPlugins, a.keyspace, gitProjectId),
	}

	response, err := a.c.Session().ExecuteQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin data, %w", err)
	}

	result := response.GetResultSet()
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("not found")
	}

	plugins := []*pluginstorepb.Plugin{}
	for _, row := range result.Rows {
		plugin, err := toPlugin(row, result.Columns)
		if err != nil {
			return nil, err
		}
		plugins = append(plugins, plugin)
	}

	return plugins, nil
}

func (a *AstraServerStore) DeletePlugin(gitProjectId, pluginName string) error {
	deleteQuery := &pb.Query{
		Cql: fmt.Sprintf(deletePluginData, a.keyspace, gitProjectId, pluginName),
	}

	_, err := a.c.Session().ExecuteQuery(deleteQuery)
	if err != nil {
		return fmt.Errorf("failed to delete plugin, %w", err)
	}

	return nil
}

func toPluginData(row *pb.Row, columns []*pb.ColumnSpec) (*pluginstorepb.PluginData, error) {
	pluginName, err := client.ToString(row.Values[0])
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin_name: %w", err)
	}
	/*_, err = client.ToDate(row.Values[1])
	if err != nil {
		return nil, fmt.Errorf("failed to read last_updated_time: %w", err)
	}*/
	storeType, err := client.ToInt(row.Values[2])
	if err != nil {
		return nil, fmt.Errorf("failed to read storeType: %w", err)
	}
	description, err := client.ToString(row.Values[3])
	if err != nil {
		return nil, fmt.Errorf("failed to read description: %w", err)
	}
	category, err := client.ToString(row.Values[4])
	if err != nil {
		return nil, fmt.Errorf("failed to read category: %w", err)
	}
	icon, err := client.ToString(row.Values[5])
	if err != nil {
		return nil, fmt.Errorf("failed to read icon: %w", err)
	}
	versions, err := client.ToList(row.Values[6], columns[6].Type)
	if err != nil {
		return nil, fmt.Errorf("failed to read versions: %w", err)
	}
	pluginVersions, err := convertToSlice(versions)
	if err != nil {
		return nil, fmt.Errorf("failed to convert versions: %w", err)
	}

	return &pluginstorepb.PluginData{
		StoreType:  pluginstorepb.StoreType(storeType),
		PluginName: pluginName, Description: description,
		Category: category, Icon: []byte(icon),
		Versions: pluginVersions,
	}, nil
}

func toPlugin(row *pb.Row, columns []*pb.ColumnSpec) (*pluginstorepb.Plugin, error) {
	pluginName, err := client.ToString(row.Values[0])
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin_name: %w", err)
	}
	/*_, err = client.ToDate(row.Values[1])
	if err != nil {
		return nil, fmt.Errorf("failed to read last_updated_time: %w", err)
	}*/
	storeType, err := client.ToInt(row.Values[2])
	if err != nil {
		return nil, fmt.Errorf("failed to read storeType: %w", err)
	}
	description, err := client.ToString(row.Values[3])
	if err != nil {
		return nil, fmt.Errorf("failed to read description: %w", err)
	}
	category, err := client.ToString(row.Values[4])
	if err != nil {
		return nil, fmt.Errorf("failed to read category: %w", err)
	}
	icon, err := client.ToString(row.Values[5])
	if err != nil {
		return nil, fmt.Errorf("failed to read icon: %w", err)
	}
	versions, err := client.ToList(row.Values[6], columns[6].Type)
	if err != nil {
		return nil, fmt.Errorf("failed to read versions: %w", err)
	}
	pluginVersions, err := convertToSlice(versions)
	if err != nil {
		return nil, fmt.Errorf("failed to convert versions: %w", err)
	}

	return &pluginstorepb.Plugin{
		StoreType:  pluginstorepb.StoreType(storeType),
		PluginName: pluginName, Description: description,
		Category: category, Icon: []byte(icon),
		Versions: pluginVersions,
	}, nil
}

func convertToSlice(input interface{}) ([]string, error) {
	switch v := input.(type) {
	case []string:
		return v, nil
	case []interface{}:
		result := make([]string, len(v))
		for i, item := range v {
			strVal, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("unable to convert element at index %d to string", i)
			}
			result[i] = strVal
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", input)
	}
}

func getSQLStringArray(val []string) string {
	return "['" + strings.Join(val, "', '") + "']"
}
