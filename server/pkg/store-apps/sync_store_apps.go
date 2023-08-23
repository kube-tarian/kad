package storeapps

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/server/pkg/store"
	"github.com/kube-tarian/kad/server/pkg/types"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	AppStoreAppConfigPath string `envconfig:"APP_STORE_APP_CONFIG_PATH" default:"/data/store-apps/conf"`
	SyncAppStore          bool   `envconfig:"SYNC_APP_STORE" default:"false"`
	AppStoreConfigFile    string `envconfig:"APP_STORE_CONFIG_FILE" default:"/data/store-apps/app_list.yaml"`
}

type AppStoreConfig struct {
	EnabledApps  []string `yaml:"enabledApps"`
	DisabledApps []string `yaml:"disabledApps"`
}

func SyncStoreApps(log logging.Logger, appStore store.ServerStore) error {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return err
	}

	if !cfg.SyncAppStore {
		log.Info("app store config synch disabled")
		return nil
	}

	configData, err := os.ReadFile(cfg.AppStoreConfigFile)
	if err != nil {
		return errors.WithMessage(err, "failed to read store config file")
	}

	var config AppStoreConfig
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return errors.WithMessage(err, "failed to unmarshall store config file")
	}

	for _, appName := range config.EnabledApps {
		appData, err := os.ReadFile(cfg.AppStoreAppConfigPath + "/" + appName + ".yaml")
		if err != nil {
			return errors.WithMessagef(err, "failed to read store app config for %s", appName)
		}

		var appConfig types.AppConfig
		if err := yaml.Unmarshal(appData, &appConfig); err != nil {
			return errors.WithMessagef(err, "failed to unmarshall store app config for %s", appName)
		}

		if appConfig.Name == "" || appConfig.Version == "" {
			return fmt.Errorf("app name/version is missing for %s", appName)
		}

		storeAppConfig := &types.StoreAppConfig{
			AppName:             appConfig.Name,
			Version:             appConfig.Version,
			Category:            appConfig.Category,
			Description:         appConfig.Description,
			ChartName:           appConfig.ChartName,
			RepoName:            appConfig.RepoName,
			ReleaseName:         appConfig.ReleaseName,
			RepoURL:             appConfig.RepoURL,
			Namespace:           appConfig.Namespace,
			CreateNamespace:     appConfig.CreateNamespace,
			PrivilegedNamespace: appConfig.PrivilegedNamespace,
			Icon:                appConfig.Icon,
			LaunchURL:           appConfig.LaunchUIURL,
			LaunchUIDescription: appConfig.LaunchUIDescription,
		}

		overrideValuesJSON, err := json.Marshal(appConfig.OverrideValues)
		if err != nil {
			return errors.WithMessagef(err, "failed to unmarshall store app config values for %s", appName)
		}

		storeAppConfig.OverrideValues = string(overrideValuesJSON)
		launchUIValues, err := json.Marshal(appConfig.LaunchUIValues)
		if err != nil {
			return errors.WithMessagef(err, "failed to unmarshall store app config UI values for %s", appName)
		}
		storeAppConfig.LaunchUIValues = string(launchUIValues)

		if err := appStore.AddOrUpdateApp(storeAppConfig); err != nil {
			return errors.WithMessagef(err, "failed to store app config for %s", appName)
		}
	}
	return nil
}
