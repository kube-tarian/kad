package captenstore

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/kube-tarian/kad/capten/common-pkg/gerrors"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/clusterpluginspb"
	postgresdb "github.com/kube-tarian/kad/capten/common-pkg/postgres"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (a *Store) UpsertClusterPluginConfig(pluginConfig *clusterpluginspb.Plugin) error {
	if len(pluginConfig.PluginName) == 0 {
		return fmt.Errorf("plugin name empty")
	}

	plugin := &ClusterPluginConfig{}
	recordFound := true
	err := a.dbClient.Find(plugin, ClusterPluginConfig{PluginName: pluginConfig.PluginName})
	if err != nil {
		if gerrors.GetErrorType(err) != postgresdb.ObjectNotExist {
			return prepareError(err, pluginConfig.PluginName, "Fetch")
		}
		err = nil
		recordFound = false
	} else if plugin.PluginName == "" {
		recordFound = false
	}

	plugin.PluginName = pluginConfig.PluginName
	plugin.PluginStoreType = int(pluginConfig.StoreType)
	plugin.Category = pluginConfig.Category
	plugin.Capabilities = pluginConfig.Capabilities
	plugin.Description = pluginConfig.Description
	plugin.Icon = pluginConfig.Icon
	plugin.ChartName = pluginConfig.ChartName
	plugin.ChartRepo = pluginConfig.ChartRepo
	plugin.Namespace = pluginConfig.DefaultNamespace
	plugin.PrivilegedNamespace = pluginConfig.PrivilegedNamespace
	plugin.APIEndpoint = pluginConfig.ApiEndpoint
	plugin.UIEndpoint = pluginConfig.UiEndpoint
	plugin.UIModuleEndpoint = pluginConfig.UiEndpoint
	plugin.Version = pluginConfig.Version
	plugin.Values = base64.StdEncoding.EncodeToString(pluginConfig.Values)
	plugin.OverrideValues = base64.StdEncoding.EncodeToString(pluginConfig.OverrideValues)
	plugin.InstallStatus = pluginConfig.InstallStatus
	plugin.LastUpdateTime = time.Now()

	if !recordFound {
		err = a.dbClient.Create(plugin)
	} else {
		err = a.dbClient.Update(plugin, ClusterPluginConfig{PluginName: pluginConfig.PluginName})
	}
	return err
}

func (a *Store) DeleteClusterPluginConfig(pluginName string) error {
	err := a.dbClient.Delete(ClusterPluginConfig{}, ClusterPluginConfig{PluginName: pluginName})
	if err != nil {
		err = prepareError(err, pluginName, "Delete")
	}
	return err
}

func (a *Store) GetClusterPluginConfig(pluginName string) (*clusterpluginspb.Plugin, error) {
	var pluginConfig ClusterPluginConfig
	err := a.dbClient.Find(&pluginConfig, "plugin_name = ?", pluginName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	values, _ := base64.StdEncoding.DecodeString(pluginConfig.Values)
	overrideValues, _ := base64.StdEncoding.DecodeString(pluginConfig.OverrideValues)

	p := &clusterpluginspb.Plugin{
		PluginName:          pluginConfig.PluginName,
		StoreType:           clusterpluginspb.StoreType(pluginConfig.PluginStoreType),
		Category:            pluginConfig.Category,
		Capabilities:        pluginConfig.Capabilities,
		Description:         pluginConfig.Description,
		Icon:                pluginConfig.Icon,
		ChartName:           pluginConfig.ChartName,
		ChartRepo:           pluginConfig.ChartRepo,
		DefaultNamespace:    pluginConfig.Namespace,
		PrivilegedNamespace: pluginConfig.PrivilegedNamespace,
		ApiEndpoint:         pluginConfig.APIEndpoint,
		UiEndpoint:          pluginConfig.UIEndpoint,
		Version:             pluginConfig.Version,
		Values:              values,
		OverrideValues:      overrideValues,
		InstallStatus:       pluginConfig.InstallStatus,
	}

	return p, nil
}

func (a *Store) GetAllClusterPluginConfigs() ([]*clusterpluginspb.Plugin, error) {
	var plugins []ClusterPluginConfig
	err := a.dbClient.Find(&plugins, nil)
	if err != nil && gerrors.GetErrorType(err) != postgresdb.ObjectNotExist {
		return nil, fmt.Errorf("failed to fetch plugins: %v", err.Error())
	}

	var pluginConfigs []*clusterpluginspb.Plugin
	for _, p := range plugins {
		values, _ := base64.StdEncoding.DecodeString(p.Values)
		overrideValues, _ := base64.StdEncoding.DecodeString(p.OverrideValues)

		pluginConfig := &clusterpluginspb.Plugin{
			PluginName:          p.PluginName,
			StoreType:           clusterpluginspb.StoreType(p.PluginStoreType),
			Category:            p.Category,
			Capabilities:        p.Capabilities,
			Description:         p.Description,
			Icon:                p.Icon,
			ChartName:           p.ChartName,
			ChartRepo:           p.ChartRepo,
			DefaultNamespace:    p.Namespace,
			PrivilegedNamespace: p.PrivilegedNamespace,
			ApiEndpoint:         p.APIEndpoint,
			UiEndpoint:          p.UIEndpoint,
			Version:             p.Version,
			Values:              values,
			OverrideValues:      overrideValues,
			InstallStatus:       p.InstallStatus,
		}
		pluginConfigs = append(pluginConfigs, pluginConfig)
	}

	return pluginConfigs, nil
}
