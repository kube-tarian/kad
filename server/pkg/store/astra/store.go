package astra

import (
	"encoding/json"
	"fmt"
	"os"

	astraclient "github.com/kube-tarian/kad/server/pkg/astra-client"
	"github.com/kube-tarian/kad/server/pkg/config"
	"github.com/kube-tarian/kad/server/pkg/types"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
	"gopkg.in/yaml.v3"
)

type AstraServerStore struct {
	c        *astraclient.Client
	keyspace string
}

func NewStore() (*AstraServerStore, error) {
	a := &AstraServerStore{}
	err := a.initClient()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to astra db, %w", err)
	}
	a.keyspace = a.c.Keyspace()
	return a, nil
}

func (a *AstraServerStore) initClient() error {
	var err error
	a.c, err = astraclient.NewClient()
	return err
}

func (a *AstraServerStore) InitializeDb() error {
	initDbQueries := []string{
		fmt.Sprintf(createClusterEndpointTableQuery, a.keyspace),
		fmt.Sprintf(createAppConfigTableQuery, a.keyspace),
	}

	for _, query := range initDbQueries {
		createQuery := &pb.Query{
			Cql: query,
		}

		_, err := a.c.Session().ExecuteQuery(createQuery)
		if err != nil {
			return fmt.Errorf("failed to initialise db: %w", err)
		}
	}

	if err := appStoreConfig(a, a.c.Session()); err != nil {
		return err
	}

	return nil
}

func appStoreConfig(handler *AstraServerStore, session *client.StargateClient) error {

	cfg, err := config.GetServiceConfig()
	if err != nil {
		return fmt.Errorf("failed to load service config: %w", err)
	}

	if cfg.ReadAppStoreConfig {
		configData, err := os.ReadFile(cfg.AppStorConfig + "/" + appStoreConfigFileName)
		if err != nil {
			return fmt.Errorf("failed to read store config file: %w", err)
		}

		var config AppStoreConfig
		if err := yaml.Unmarshal(configData, &config); err != nil {
			return fmt.Errorf("failed to unmarshall store config file: %w", err)
		}

		for _, v := range append(config.CreateStoreApps, config.UpdateStoreApps...) {
			appData, err := os.ReadFile(cfg.AppStorConfig + "/" + v + ".yaml")
			if err != nil {
				return fmt.Errorf("failed to read app store config file: %w. App name - %s", err, v)
			}

			var appConfig AppConfig
			if err := yaml.Unmarshal(appData, &appConfig); err != nil {
				return fmt.Errorf("failed to unmarshall app store config file: %w. App name - %s", err, v)
			}

			if appConfig.Name == "" || appConfig.Version == "" {
				return fmt.Errorf("failed to add app config to store, %v", "App name/version is missing")
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
				LaunchURL:           appConfig.LaunchURL,
				LaunchRedirectURL:   appConfig.LaunchRedirectURL,
			}

			overrideValuesJSON, err := json.Marshal(appConfig.OverrideValues)
			if err != nil {
				return fmt.Errorf("failed to marshall app store config file: %w. App name - %s", err, v)
			}
			storeAppConfig.OverrideValues = string(overrideValuesJSON)
			launchUIValues, err := json.Marshal(appConfig.LaunchUIValues)
			if err != nil {
				return fmt.Errorf("failed to marshall app store config file: %w. App name - %s", err, v)
			}
			storeAppConfig.LaunchUIValues = string(launchUIValues)

			if err := handler.AddOrUpdateApp(storeAppConfig); err != nil {
				return fmt.Errorf("failed to add app config to store, %v", err)
			}

		}
	}

	return nil
}
