package captenstore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kube-tarian/kad/capten/common-pkg/gerrors"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/pluginstorepb"
	postgresdb "github.com/kube-tarian/kad/capten/common-pkg/postgres"
	"gorm.io/gorm"
)

func (a *Store) UpsertPluginStoreData(gitProjectID string, pluginData *pluginstorepb.PluginData) error {
	pluginStoreData := &PluginStoreData{}
	recordFound := true
	err := a.dbClient.Find(pluginStoreData, PluginStoreData{
		StoreType:    int(pluginData.StoreType),
		GitProjectID: uuid.MustParse(gitProjectID),
		PluginName:   pluginData.PluginName})
	if err != nil {
		if gerrors.GetErrorType(err) != postgresdb.ObjectNotExist {
			return prepareError(err, fmt.Sprintf("%d/%s/%s",
				pluginData.StoreType, gitProjectID, pluginData.PluginName), "Fetch")
		}
		err = nil
		recordFound = false
	} else if pluginStoreData.StoreType == 0 {
		recordFound = false
	}

	pluginStoreData.StoreType = int(pluginData.StoreType)
	pluginStoreData.GitProjectID = uuid.MustParse(gitProjectID)
	pluginStoreData.PluginName = pluginData.PluginName
	pluginStoreData.Category = pluginData.Category
	pluginStoreData.Versions = pluginData.Versions
	pluginStoreData.Icon = pluginData.Icon
	pluginStoreData.Description = pluginData.Description
	pluginStoreData.LastUpdateTime = time.Now()

	if !recordFound {
		err = a.dbClient.Create(pluginStoreData)
	} else {
		err = a.dbClient.Update(pluginStoreData, PluginStoreData{
			StoreType:    int(pluginData.StoreType),
			GitProjectID: uuid.MustParse(gitProjectID),
			PluginName:   pluginData.PluginName})
	}
	return err
}

func (a *Store) GetPluginStoreData(storeType pluginstorepb.StoreType, gitProjectId, pluginName string) (*pluginstorepb.PluginData, error) {
	pluginStoreData := &PluginStoreData{}
	err := a.dbClient.Find(pluginStoreData, PluginStoreData{
		StoreType:    int(storeType),
		GitProjectID: uuid.MustParse(gitProjectId),
		PluginName:   pluginName,
	})
	if err != nil {
		return nil, prepareError(err, fmt.Sprintf("%s/%s/%s", gitProjectId, gitProjectId, pluginName), "Fetch")
	} else if pluginStoreData.StoreType == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	pluginData := &pluginstorepb.PluginData{
		StoreType:   pluginstorepb.StoreType(pluginStoreData.StoreType),
		PluginName:  pluginStoreData.PluginName,
		Category:    pluginStoreData.Category,
		Versions:    pluginStoreData.Versions,
		Icon:        pluginStoreData.Icon,
		Description: pluginStoreData.Description,
	}
	return pluginData, nil
}

func (a *Store) GetPlugins(gitProjectId string) ([]*pluginstorepb.Plugin, error) {
	pluginDataSet := []PluginStoreData{}
	err := a.dbClient.Find(&pluginDataSet, PluginStoreData{GitProjectID: uuid.MustParse(gitProjectId)})
	if err != nil {
		return nil, prepareError(err, gitProjectId, "Fetch")
	}

	plugins := []*pluginstorepb.Plugin{}
	for _, pluginData := range pluginDataSet {
		plugins = append(plugins, &pluginstorepb.Plugin{
			StoreType:   pluginstorepb.StoreType(pluginData.StoreType),
			PluginName:  pluginData.PluginName,
			Category:    pluginData.Category,
			Versions:    pluginData.Versions,
			Icon:        pluginData.Icon,
			Description: pluginData.Description,
		})
	}
	return plugins, nil
}
func (a *Store) DeletePluginStoreData(storeType pluginstorepb.StoreType, gitProjectId, pluginName string) error {
	err := a.dbClient.Delete(PluginStoreData{}, PluginStoreData{
		StoreType:    int(storeType),
		GitProjectID: uuid.MustParse(gitProjectId),
		PluginName:   pluginName,
	})
	return err
}
