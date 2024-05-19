package captenstore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kube-tarian/kad/capten/common-pkg/gerrors"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/pluginstorepb"
	postgresdb "github.com/kube-tarian/kad/capten/common-pkg/postgres"
)

func (a *Store) UpsertPluginStoreConfig(config *pluginstorepb.PluginStoreConfig) error {
	pluginStoreConfig := &PluginStoreConfig{}
	recordFound := true
	err := a.dbClient.FindFirst(pluginStoreConfig, PluginStoreConfig{StoreType: int(config.StoreType)})
	if err != nil {
		if gerrors.GetErrorType(err) != postgresdb.ObjectNotExist {
			return prepareError(err, fmt.Sprintf("%d", config.StoreType), "Fetch")
		}
		err = nil
		recordFound = false
	}

	pluginStoreConfig.StoreType = int(config.StoreType)
	pluginStoreConfig.GitProjectID = uuid.MustParse(config.GitProjectId)
	pluginStoreConfig.GitProjectURL = config.GitProjectURL
	pluginStoreConfig.LastUpdateTime = time.Now()

	if !recordFound {
		err = a.dbClient.Create(pluginStoreConfig)
	} else {
		err = a.dbClient.Update(pluginStoreConfig, PluginStoreConfig{StoreType: int(config.StoreType)})
	}
	return err
}

func (a *Store) GetPluginStoreConfig(storeType pluginstorepb.StoreType) (*pluginstorepb.PluginStoreConfig, error) {
	pluginStoreConfig := &PluginStoreConfig{}
	err := a.dbClient.FindFirst(pluginStoreConfig, PluginStoreConfig{StoreType: int(storeType)})
	if err != nil {
		return nil, err
	}

	return &pluginstorepb.PluginStoreConfig{
		StoreType:     pluginstorepb.StoreType(pluginStoreConfig.StoreType),
		GitProjectId:  pluginStoreConfig.GitProjectID.String(),
		GitProjectURL: pluginStoreConfig.GitProjectURL,
	}, nil
}

func (a *Store) DeletePluginStoreConfig(storeType pluginstorepb.StoreType) error {
	return a.dbClient.Delete(PluginStoreConfig{}, PluginStoreConfig{StoreType: int(storeType)})
}
