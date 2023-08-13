package astra

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kube-tarian/kad/server/pkg/types"
	"github.com/stargate/stargate-grpc-go-client/stargate/pkg/client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
)

const (
	createAppConfigQuery         string = "INSERT INTO %s.app_config (name, chart_name, repo_name,release_name, repo_url, namespace, version, create_namespace, privileged_namespace, launch_ui_url, launch_ui_redirect_url, category, icon, description, launch_ui_values, override_values, created_time, id) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', %t, %t, '%s', '%s', '%s', '%s', '%s', '%s', '%s','%v', '%s' );"
	updateAppConfigQuery         string = "UPDATE %s.app_config SET chart_name = '%s', repo_name = '%s', repo_url = '%s', namespace = '%s', create_namespace = %t, privileged_namespace = %t, launch_ui_url = '%s', launch_ui_redirect_url = '%s', category = '%s', icon = '%s', description = '%s', launch_ui_values = '%s', override_values = '%s',last_updated_time='%v' WHERE name = '%s' AND version = '%s';"
	deleteAppConfigQuery         string = "DELETE FROM %s.app_config WHERE name='%s' AND version='%s';"
	getAppConfigQuery            string = "SELECT name,chart_name,repo_name,repo_url,namespace,version,create_namespace,privileged_namespace,launch_ui_url,launch_ui_redirect_url,category,icon,description,launch_ui_values,override_values, release_name FROM %s.app_config WHERE name='%s' AND version='%s';"
	getAllAppConfigsQuery        string = "SELECT name,chart_name,repo_name,repo_url,namespace,version,create_namespace,privileged_namespace,launch_ui_url,launch_ui_redirect_url,category,icon,description,launch_ui_values,override_values, release_name FROM %s.app_config;"
	appConfigExistanceCheckQuery string = "SELECT name, version FROM %s.app_config WHERE name='%s' AND version ='%s';"
)

func (a *AstraServerStore) AddOrUpdateApp(config *types.StoreAppConfig) error {
	appExists, err := a.isAppExistsInStore(config.AppName, config.Version)
	if err != nil {
		return fmt.Errorf("failed to check app config existance : %w", err)
	}

	fmt.Println("inside AddOrUpdateApp")
	log.Print("inside AddOrUpdateApp")

	var query *pb.Query
	if appExists {
		query = &pb.Query{
			Cql: fmt.Sprintf(updateAppConfigQuery,
				a.keyspace, config.ChartName, config.RepoName, config.RepoURL, config.Namespace, config.CreateNamespace, config.PrivilegedNamespace, config.LaunchURL, config.LaunchRedirectURL, config.Category, config.Icon, config.Description, config.LaunchUIValues, config.OverrideValues, time.Now().Format("2006-01-02 15:04:05"), config.AppName, config.Version),
		}
	} else {
		query = &pb.Query{
			Cql: fmt.Sprintf(createAppConfigQuery,
				a.keyspace, config.AppName, config.ChartName, config.RepoName, config.ReleaseName, config.RepoURL, config.Namespace, config.Version, config.CreateNamespace, config.PrivilegedNamespace, config.LaunchURL, config.LaunchRedirectURL, config.Category, config.Icon, config.Description, config.LaunchUIValues, config.OverrideValues, time.Now().Format("2006-01-02 15:04:05"), uuid.New().String()),
		}
	}

	_, err = a.c.Session().ExecuteQuery(query)
	if err != nil {
		return fmt.Errorf("failed to insert/update the app config into the app_config table : %w", err)
	}

	return nil
}

func (a *AstraServerStore) DeleteAppInStore(name, version string) error {

	deleteQuery := &pb.Query{
		Cql: fmt.Sprintf(
			deleteAppConfigQuery,
			a.keyspace, name, version),
	}

	_, err := a.c.Session().ExecuteQuery(deleteQuery)
	if err != nil {
		return fmt.Errorf("failed to initialise db: %w", err)
	}

	return nil
}

func (a *AstraServerStore) GetAppFromStore(name, version string) (*types.AppConfig, error) {

	selectQuery := &pb.Query{
		Cql: fmt.Sprintf(getAppConfigQuery,
			a.keyspace, name, version),
	}

	response, err := a.c.Session().ExecuteQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to initialise db: %w", err)
	}

	result := response.GetResultSet()

	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("app: %s not found", name)
	}

	config, err := toAppConfig(result.Rows[0])
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (a *AstraServerStore) GetAppsFromStore() (*[]types.AppConfig, error) {

	selectQuery := &pb.Query{
		Cql: fmt.Sprintf(getAllAppConfigsQuery,
			a.keyspace),
	}

	response, err := a.c.Session().ExecuteQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to initialise db: %w", err)
	}

	result := response.GetResultSet()

	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("app config's not found")
	}

	var appConfigs []types.AppConfig
	for _, row := range result.Rows {
		config, err := toAppConfig(row)
		if err != nil {
			return nil, err
		}
		appConfigs = append(appConfigs, *config)
	}

	return &appConfigs, nil
}

func toAppConfig(row *pb.Row) (*types.AppConfig, error) {

	cqlAppName, err := client.ToString(row.Values[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get app name: %w", err)
	}
	cqlChartName, err := client.ToString(row.Values[1])
	if err != nil {
		return nil, fmt.Errorf("failed to get chart name: %w", err)
	}
	cqlRepoName, err := client.ToString(row.Values[2])
	if err != nil {
		return nil, fmt.Errorf("failed to get repo name: %w", err)
	}
	cqlRepoURL, err := client.ToString(row.Values[3])
	if err != nil {
		return nil, fmt.Errorf("failed to get repo url: %w", err)
	}
	cqlNamespace, err := client.ToString(row.Values[4])
	if err != nil {
		return nil, fmt.Errorf("failed to get Namespace: %w", err)
	}
	cqlVersion, err := client.ToString(row.Values[5])
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %w", err)
	}
	cqlCreateNamespace, err := client.ToBoolean(row.Values[6])
	if err != nil {
		return nil, fmt.Errorf("failed to get Create Namespace: %w", err)
	}
	cqlPrivilegedNamespace, err := client.ToBoolean(row.Values[7])
	if err != nil {
		return nil, fmt.Errorf("failed to get Privileged Namespace: %w", err)
	}
	cqlLaunchUiUrl, err := client.ToString(row.Values[8])
	if err != nil {
		return nil, fmt.Errorf("failed to get launch ui url: %w", err)
	}
	cqlLaunchUiRedirectUrl, err := client.ToString(row.Values[9])
	if err != nil {
		return nil, fmt.Errorf("failed to get launch ui redirect url: %w", err)
	}
	cqlCategory, err := client.ToString(row.Values[10])
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	cqlIcon, err := client.ToString(row.Values[11])
	if err != nil {
		return nil, fmt.Errorf("failed to get icon: %w", err)
	}
	cqlDescription, err := client.ToString(row.Values[12])
	if err != nil {
		return nil, fmt.Errorf("failed to get launch ui redirect url: %w", err)
	}
	cqlLaunchUiValues, err := client.ToString(row.Values[13])
	if err != nil {
		return nil, fmt.Errorf("failed to get launch ui values: %w", err)
	}
	cqlOverrideValues, err := client.ToString(row.Values[14])
	if err != nil {
		return nil, fmt.Errorf("failed to get override values: %w", err)
	}
	cqlReleaseNameValues, err := client.ToString(row.Values[15])
	if err != nil {
		return nil, fmt.Errorf("failed to get override values: %w", err)
	}

	config := &types.AppConfig{
		Name:                cqlAppName,
		ChartName:           cqlChartName,
		RepoName:            cqlRepoName,
		RepoURL:             cqlRepoURL,
		Namespace:           cqlNamespace,
		Version:             cqlVersion,
		CreateNamespace:     cqlCreateNamespace,
		PrivilegedNamespace: cqlPrivilegedNamespace,
		LaunchUIURL:         cqlLaunchUiUrl,
		LaunchUIRedirectURL: cqlLaunchUiRedirectUrl,
		Category:            cqlCategory,
		Icon:                cqlIcon,
		Description:         cqlDescription,
		LaunchUIValues:      cqlLaunchUiValues,
		OverrideValues:      cqlOverrideValues,
		ReleaseName:         cqlReleaseNameValues,
	}
	return config, nil
}

func (a *AstraServerStore) isAppExistsInStore(name, version string) (bool, error) {
	selectClusterQuery := &pb.Query{
		Cql: fmt.Sprintf(appConfigExistanceCheckQuery,
			a.keyspace, name, version),
	}

	response, err := a.c.Session().ExecuteQuery(selectClusterQuery)
	if err != nil {
		return false, fmt.Errorf("failed to initialise db: %w", err)
	}

	result := response.GetResultSet()
	if len(result.Rows) > 0 {
		return true, nil
	}

	return false, nil
}

func (a *AstraServerStore) GetStoreAppValues(name, version string) (*types.AppConfig, error) {

	selectQuery := &pb.Query{
		Cql: fmt.Sprintf(getAppConfigQuery,
			a.keyspace, name, version),
	}

	response, err := a.c.Session().ExecuteQuery(selectQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch app store values: %w", err)
	}

	result := response.GetResultSet()

	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("app: %s not found", name)
	}

	config, err := toAppConfig(result.Rows[0])
	if err != nil {
		return nil, err
	}

	return config, nil
}
