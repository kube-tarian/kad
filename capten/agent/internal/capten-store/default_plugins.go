package captenstore

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"
)

const (
	insertPluginData            = `INSERT INTO %s.PluginData (git_project_id, plugin_name, last_updated_time, description, category, icon, versions) VALUES (%s, '%s', '%s', '%s', '%s', '%s', %v) IF NOT EXISTS`
	updatePluginData            = `UPDATE %s.PluginData SET last_updated_time = '%s', description = '%s', category = '%s', icon = '%s', versions = %v WHERE git_project_id = %s and plugin_name = '%s'`
	readPlugins                 = `SELECT plugin_name, last_updated_time, description, category, icon, versions FROM %s.PluginData WHERE git_project_id = %s`
	readPluginDataForPluginName = `SELECT plugin_name, last_updated_time, store_type, description, category, icon, versions FROM %s.PluginData WHERE git_project_id = %s and plugin_name = '%s'`
	deletePluginData            = "DELETE FROM %s.PluginData WHERE git_project_id = %s and plugin_name = '%s'"
)

type PluginData struct {
	PluginName     string
	Description    string
	Category       string
	Versions       []string
	Icon           []byte
	LastUpdateTime string
}

func (a *Store) UpsertPluginData(gitProjectId string, pluginData *PluginData) error {
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertPluginData,
		a.keyspace, gitProjectId, pluginData.PluginName, time.Now().Format(time.RFC3339),
		pluginData.Description, pluginData.Category, pluginData.Icon,
		getSQLStringArray(pluginData.Versions)))
	err := a.client.Session().ExecuteBatch(batch)
	if err != nil {
		batch.Query(fmt.Sprintf(updatePluginData,
			a.keyspace, time.Now().Format(time.RFC3339),
			pluginData.Description, pluginData.Category, pluginData.Icon,
			getSQLStringArray(pluginData.Versions), gitProjectId, pluginData.PluginName))
		batch = a.client.Session().NewBatch(gocql.LoggedBatch)
		err = a.client.Session().ExecuteBatch(batch)
	}
	return err
}

func (a *Store) ReadPluginData(gitProjectId string, pluginName string) (*PluginData, error) {
	query := fmt.Sprintf(readPluginDataForPluginName, a.keyspace, gitProjectId, pluginName)
	plugins, err := a.executePluginDataSelectQuery(query)
	if err != nil {
		return nil, err
	}

	if len(plugins) == 0 {
		return nil, fmt.Errorf("project not found")
	}
	return plugins[0], nil
}

func (a *Store) ReadPlugins(gitProjectId string) ([]*PluginData, error) {
	query := fmt.Sprintf(readPlugins, a.keyspace, gitProjectId)
	return a.executePluginDataSelectQuery(query)
}

func (a *Store) DeletePlugin(gitProjectId, pluginName string) error {
	deleteAction := a.client.Session().Query(fmt.Sprintf(deletePluginData, a.keyspace, gitProjectId, pluginName))
	return deleteAction.Exec()
}

func getSQLStringArray(val []string) string {
	return "['" + strings.Join(val, "', '") + "']"
}

func (a *Store) executePluginDataSelectQuery(query string) ([]*PluginData, error) {
	selectQuery := a.client.Session().Query(query)
	iter := selectQuery.Iter()

	tempPluginData := PluginData{}

	pluginDataList := make([]*PluginData, 0)
	for iter.Scan(
		&tempPluginData.PluginName,
		&tempPluginData.LastUpdateTime,
		&tempPluginData.Description,
		&tempPluginData.Category,
		&tempPluginData.Icon,
		&tempPluginData.Versions,
	) {
		pluginData := &PluginData{
			PluginName:     tempPluginData.PluginName,
			Description:    tempPluginData.Description,
			Category:       tempPluginData.Category,
			Versions:       tempPluginData.Versions,
			Icon:           tempPluginData.Icon,
			LastUpdateTime: tempPluginData.LastUpdateTime,
		}

		pluginDataList = append(pluginDataList, pluginData)
	}

	if err := iter.Close(); err != nil {
		return nil, errors.WithMessage(err, "failed to iterate through results:")
	}

	return pluginDataList, nil
}
