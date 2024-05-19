package captenstore

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/kube-tarian/kad/capten/common-pkg/gerrors"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/agentpb"
	postgresdb "github.com/kube-tarian/kad/capten/common-pkg/postgres"
)

func (a *Store) UpsertAppConfig(appData *agentpb.SyncAppData) error {
	if len(appData.Config.ReleaseName) == 0 {
		return fmt.Errorf("app release name empty")
	}

	appConfig := &ClusterAppConfig{}
	recordFound := true
	err := a.dbClient.Find(appConfig, ClusterAppConfig{ReleaseName: appData.Config.ReleaseName})
	if err != nil {
		if gerrors.GetErrorType(err) != postgresdb.ObjectNotExist {
			return prepareError(err, appData.Config.ReleaseName, "Fetch")
		}
		err = nil
		recordFound = false
	} else if appConfig.ReleaseName == "" {
		recordFound = false
	}

	appConfig.ReleaseName = appData.Config.ReleaseName
	appConfig.PluginName = appData.Config.PluginName
	appConfig.PluginStoreType = int(appData.Config.PluginStoreType)
	appConfig.Category = appData.Config.Category
	appConfig.Description = appData.Config.Description
	appConfig.Icon = appData.Config.Icon
	appConfig.AppName = appData.Config.ChartName
	appConfig.RepoURL = appData.Config.RepoURL
	appConfig.Namespace = appData.Config.Namespace
	appConfig.PrivilegedNamespace = appData.Config.PrivilegedNamespace
	appConfig.APIEndpoint = appData.Config.ApiEndpoint
	appConfig.UIEndpoint = appData.Config.UiEndpoint
	appConfig.UIModuleEndpoint = appData.Config.UiModuleEndpoint
	appConfig.Version = appData.Config.Version
	appConfig.OverrideValues = base64.StdEncoding.EncodeToString(appData.Values.OverrideValues)
	appConfig.TemplateValues = base64.StdEncoding.EncodeToString(appData.Values.TemplateValues)
	appConfig.LaunchUIValues = base64.StdEncoding.EncodeToString(appData.Values.LaunchUIValues)
	appConfig.InstallStatus = appData.Config.InstallStatus
	appConfig.DefaultApp = appData.Config.DefualtApp
	appConfig.LastUpdateTime = time.Now()

	if !recordFound {
		err = a.dbClient.Create(appConfig)
	} else {
		err = a.dbClient.Update(appConfig, ClusterAppConfig{ReleaseName: appData.Config.ReleaseName})
	}
	return err
}

func (a *Store) GetAppConfig(appReleaseName string) (*agentpb.SyncAppData, error) {
	appConfig := &ClusterAppConfig{}
	err := a.dbClient.Find(appConfig, ClusterAppConfig{ReleaseName: appReleaseName})
	if err != nil {
		err = prepareError(err, appReleaseName, "Fetch")
		return nil, err
	}

	overrideValues, err := base64.StdEncoding.DecodeString(appConfig.OverrideValues)
	if err != nil {
		return nil, err
	}
	launchUIValues, err := base64.StdEncoding.DecodeString(appConfig.LaunchUIValues)
	if err != nil {
		return nil, err
	}
	templateValues, err := base64.StdEncoding.DecodeString(appConfig.TemplateValues)
	if err != nil {
		return nil, err
	}

	return &agentpb.SyncAppData{
		Config: &agentpb.AppConfig{
			ReleaseName:         appConfig.ReleaseName,
			PluginName:          appConfig.PluginName,
			PluginStoreType:     agentpb.PluginStoreType(appConfig.PluginStoreType),
			Category:            appConfig.Category,
			Description:         appConfig.Description,
			Icon:                appConfig.Icon,
			ChartName:           appConfig.AppName,
			RepoURL:             appConfig.RepoURL,
			Namespace:           appConfig.Namespace,
			PrivilegedNamespace: appConfig.PrivilegedNamespace,
			ApiEndpoint:         appConfig.APIEndpoint,
			UiEndpoint:          appConfig.UIEndpoint,
			UiModuleEndpoint:    appConfig.UIModuleEndpoint,
			Version:             appConfig.Version,
			InstallStatus:       appConfig.InstallStatus,
			DefualtApp:          appConfig.DefaultApp,
			LastUpdateTime:      appConfig.LastUpdateTime.Format(time.RFC3339),
		},
		Values: &agentpb.AppValues{
			OverrideValues: overrideValues,
			LaunchUIValues: launchUIValues,
			TemplateValues: templateValues,
		},
	}, nil
}

func (a *Store) GetAllApps() ([]*agentpb.SyncAppData, error) {
	var appConfigs []ClusterAppConfig
	err := a.dbClient.Find(&appConfigs, nil)
	if err != nil && gerrors.GetErrorType(err) != postgresdb.ObjectNotExist {
		return nil, fmt.Errorf("Unable to fetch apps: %v", err.Error())
	}

	var appData []*agentpb.SyncAppData
	for _, ac := range appConfigs {
		overrideValues, err := base64.StdEncoding.DecodeString(ac.OverrideValues)
		if err != nil {
			return nil, err
		}
		launchUIValues, err := base64.StdEncoding.DecodeString(ac.LaunchUIValues)
		if err != nil {
			return nil, err
		}
		templateValues, err := base64.StdEncoding.DecodeString(ac.TemplateValues)
		if err != nil {
			return nil, err
		}

		appData = append(appData, &agentpb.SyncAppData{
			Config: &agentpb.AppConfig{
				ReleaseName:         ac.ReleaseName,
				PluginName:          ac.PluginName,
				PluginStoreType:     agentpb.PluginStoreType(ac.PluginStoreType),
				Category:            ac.Category,
				Description:         ac.Description,
				Icon:                ac.Icon,
				ChartName:           ac.AppName,
				RepoURL:             ac.RepoURL,
				Namespace:           ac.Namespace,
				PrivilegedNamespace: ac.PrivilegedNamespace,
				ApiEndpoint:         ac.APIEndpoint,
				UiEndpoint:          ac.UIEndpoint,
				UiModuleEndpoint:    ac.UIModuleEndpoint,
				Version:             ac.Version,
				InstallStatus:       ac.InstallStatus,
				DefualtApp:          ac.DefaultApp,
				LastUpdateTime:      ac.LastUpdateTime.Format(time.RFC3339),
			},
			Values: &agentpb.AppValues{
				OverrideValues: overrideValues,
				LaunchUIValues: launchUIValues,
				TemplateValues: templateValues,
			},
		})
	}

	return appData, nil
}

func (a *Store) DeleteAppConfigByReleaseName(releaseName string) error {
	err := a.dbClient.Delete(ClusterAppConfig{}, ClusterAppConfig{ReleaseName: releaseName})
	if err != nil {
		err = prepareError(err, releaseName, "Delete")
	}
	return err
}
