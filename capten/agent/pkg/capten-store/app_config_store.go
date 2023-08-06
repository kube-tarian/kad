package captenstore

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
)

const (
	insertAppConfigByReleaseNameQuery = "INSERT INTO apps.AppConfig(release_name) VALUES (?)"
	updateAppConfigByReleaseNameQuery = "UPDATE apps.AppConfig SET %s WHERE release_name = ?"
)

func CreateSelectByFieldNameQuery(field string) string {
	return CreateSelectAllQuery() + fmt.Sprintf(" WHERE %s = ?", field)
}

func CreateSelectAllQuery() string {
	return "SELECT " +
		strings.Join(appConfigfields, ", ") +
		" FROM apps.AppConfig"
}

const (
	appName, description, category       = "app_name", "description", "category"
	chartName, repoName, repoUrl         = "chart_name", "repo_name", "repo_url"
	namespace, releaseName, version      = "namespace", "release_name", "version"
	launchUrl, launchRedirectUrl         = "launch_url", "launch_redirect_url"
	createNamespace, privilegedNamespace = "create_namespace", "privileged_namespace"
	overrideValues, launchUiValues       = "override_values", "launch_ui_values"
	icon, installStatus                  = "icon", "install_status"
)

var (
	appConfigfields = []string{
		appName, description, category,
		chartName, repoName, repoUrl,
		namespace, releaseName, version,
		launchUrl, launchRedirectUrl,
		createNamespace, privilegedNamespace,
		overrideValues, launchUiValues,
		icon, installStatus,
	}
)

func (a *Store) UpsertAppConfig(config *agentpb.SyncAppData) error {
	if len(config.Config.ReleaseName) == 0 {
		return fmt.Errorf("app release name empty")
	}

	kvPairs, isEmptyUpdate := formUpdateKvPairs(config)
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(insertAppConfigByReleaseNameQuery, config.Config.ReleaseName)
	if !isEmptyUpdate {
		batch.Query(fmt.Sprintf(updateAppConfigByReleaseNameQuery, kvPairs), config.Config.ReleaseName)
	}
	return a.client.Session().ExecuteBatch(batch)
}

func (a *Store) GetAppConfig(appReleaseName string) (*agentpb.SyncAppData, error) {
	selectQuery := a.client.Session().Query(CreateSelectByFieldNameQuery(releaseName), appReleaseName)

	config := agentpb.AppConfig{}
	var overrideValues, launchUiValues string

	if err := selectQuery.Scan(
		&config.AppName, &config.Description, &config.Category,
		&config.ChartName, &config.RepoName, &config.RepoURL,
		&config.Namespace, &config.ReleaseName, &config.Version,
		&config.LaunchURL, &config.LaunchRedirectURL,
		&config.CreateNamespace, &config.PrivilegedNamespace,
		&overrideValues, &launchUiValues,
		&config.Icon, &config.InstallStatus,
	); err != nil {
		return nil, err
	}

	return &agentpb.SyncAppData{
		Config: &config,
		Values: &agentpb.AppValues{
			OverrideValues: []byte(overrideValues),
			LaunchUIValues: []byte(launchUiValues)},
	}, nil
}

func (a *Store) GetAllApps() ([]*agentpb.SyncAppData, error) {
	selectAllQuery := a.client.Session().Query(CreateSelectAllQuery())
	iter := selectAllQuery.Iter()

	config := agentpb.AppConfig{}
	var overrideValues, launchUiValues string

	ret := make([]*agentpb.SyncAppData, 0)
	for iter.Scan(
		&config.AppName, &config.Description, &config.Category,
		&config.ChartName, &config.RepoName, &config.RepoURL,
		&config.Namespace, &config.ReleaseName, &config.Version,
		&config.LaunchURL, &config.LaunchRedirectURL,
		&config.CreateNamespace, &config.PrivilegedNamespace,
		&overrideValues, &launchUiValues,
		&config.Icon, &config.InstallStatus,
	) {
		a := &agentpb.SyncAppData{
			Config: &config,
			Values: &agentpb.AppValues{
				OverrideValues: []byte(overrideValues),
				LaunchUIValues: []byte(launchUiValues)},
		}
		ret = append(ret, a)
	}

	if err := iter.Close(); err != nil {
		a.log.Fatal("Failed to iterate through results:", err)
		return nil, err
	}
	return ret, nil
}

func formUpdateKvPairs(config *agentpb.SyncAppData) (string, bool) {
	params := []string{}

	if config.Values != nil && len(config.Values.OverrideValues) > 0 {
		params = append(params,
			fmt.Sprintf("%s = '%s'", overrideValues, string(config.Values.OverrideValues)))
	}

	if config.Values != nil && len(config.Values.LaunchUIValues) > 0 {
		params = append(params,
			fmt.Sprintf("%s = '%s'", launchUiValues, string(config.Values.LaunchUIValues)))
	}

	if config.Config.CreateNamespace {
		params = append(params,
			fmt.Sprintf("%s = 'true'", createNamespace))
	}
	if config.Config.PrivilegedNamespace {
		params = append(params,
			fmt.Sprintf("%s = 'true'", privilegedNamespace))
	}

	if config.Config.LaunchURL != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", launchUrl, config.Config.LaunchURL))
	}
	if config.Config.LaunchRedirectURL != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", launchRedirectUrl, config.Config.LaunchRedirectURL))
	}

	if config.Config.AppName != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", appName, config.Config.AppName))
	}
	if config.Config.Description != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", description, config.Config.Description))
	}
	if config.Config.Category != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", category, config.Config.Category))
	}

	if config.Config.ChartName != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", chartName, config.Config.ChartName))
	}
	if config.Config.RepoName != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", repoName, config.Config.RepoName))
	}
	if config.Config.RepoURL != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", repoUrl, config.Config.RepoURL))
	}

	if config.Config.Namespace != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", namespace, config.Config.Namespace))
	}

	if config.Config.Version != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", version, config.Config.Version))
	}

	if config.Config.Icon != nil && len(config.Config.Icon) > 0 {
		params = append(params,
			fmt.Sprintf("%s = 0x%s", icon, hex.EncodeToString(config.Config.Icon)))
	}
	if len(config.Config.InstallStatus) > 0 {
		params = append(params,
			fmt.Sprintf("%s = '%s'", installStatus, config.Config.InstallStatus))
	}

	if len(params) == 0 {
		// query is empty there is nothing to update
		return "", true
	}
	return strings.Join(params, ", "), false
}
