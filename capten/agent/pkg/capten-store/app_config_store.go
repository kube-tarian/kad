package captenstore

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/pkg/errors"
)

const (
	insertAppConfigByReleaseNameQuery = "INSERT INTO %s.ClusterAppConfig(release_name) VALUES (?)"
	updateAppConfigByReleaseNameQuery = "UPDATE %s.ClusterAppConfig SET %s WHERE release_name = ?"
)

func CreateSelectByFieldNameQuery(keyspace, field string) string {
	return CreateSelectAllQuery(keyspace) + fmt.Sprintf(" WHERE %s = ?", field)
}

func CreateSelectAllQuery(keyspace string) string {
	return fmt.Sprintf("SELECT %s FROM %s.ClusterAppConfig", strings.Join(appConfigfields, ", "), keyspace)
}

const (
	appName, description, category       = "app_name", "description", "category"
	chartName, repoName, repoUrl         = "chart_name", "repo_name", "repo_url"
	namespace, releaseName, version      = "namespace", "release_name", "version"
	launchUrl, launchUIDesc              = "launch_url", "launch_redirect_url"
	createNamespace, privilegedNamespace = "create_namespace", "privileged_namespace"
	overrideValues, launchUiValues       = "override_values", "launch_ui_values"
	templateValues, defaultApp           = "template_values", "default_app"
	icon, installStatus                  = "icon", "install_status"
	updateTime                           = "update_time"
)

var (
	appConfigfields = []string{
		appName, description, category,
		chartName, repoName, repoUrl,
		namespace, releaseName, version,
		launchUrl, launchUIDesc,
		createNamespace, privilegedNamespace,
		overrideValues, launchUiValues,
		templateValues, defaultApp,
		icon, installStatus,
		updateTime,
	}
)

func (a *Store) UpsertAppConfig(config *agentpb.SyncAppData) error {
	if len(config.Config.ReleaseName) == 0 {
		return fmt.Errorf("app release name empty")
	}

	kvPairs, isEmptyUpdate := formUpdateKvPairs(config)
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(fmt.Sprintf(insertAppConfigByReleaseNameQuery, a.keyspace), config.Config.ReleaseName)
	if !isEmptyUpdate {
		batch.Query(fmt.Sprintf(updateAppConfigByReleaseNameQuery, a.keyspace, kvPairs), config.Config.ReleaseName)
	}
	return a.client.Session().ExecuteBatch(batch)
}

func (a *Store) GetAppConfig(appReleaseName string) (*agentpb.SyncAppData, error) {
	selectQuery := a.client.Session().Query(CreateSelectByFieldNameQuery(a.keyspace, releaseName), appReleaseName)

	config := agentpb.AppConfig{}
	var overrideValues, launchUiValues, templateValues string

	if err := selectQuery.Scan(
		&config.AppName, &config.Description, &config.Category,
		&config.ChartName, &config.RepoName, &config.RepoURL,
		&config.Namespace, &config.ReleaseName, &config.Version,
		&config.LaunchURL, &config.LaunchUIDescription,
		&config.CreateNamespace, &config.PrivilegedNamespace,
		&overrideValues, &launchUiValues,
		&templateValues, &config.DefualtApp,
		&config.Icon, &config.InstallStatus,
		&config.LastUpdateTime,
	); err != nil {
		return nil, err
	}

	overrideValuesCopy, _ := base64.StdEncoding.DecodeString(overrideValues)
	launchUiValuesCopy, _ := base64.StdEncoding.DecodeString(launchUiValues)
	templateValuesCopy, _ := base64.StdEncoding.DecodeString(templateValues)

	return &agentpb.SyncAppData{
		Config: &config,
		Values: &agentpb.AppValues{
			OverrideValues: overrideValuesCopy,
			LaunchUIValues: launchUiValuesCopy,
			TemplateValues: templateValuesCopy,
		},
	}, nil
}

func (a *Store) GetAllApps() ([]*agentpb.SyncAppData, error) {
	selectAllQuery := a.client.Session().Query(CreateSelectAllQuery(a.keyspace))
	iter := selectAllQuery.Iter()

	config := agentpb.AppConfig{}
	var overrideValues, launchUiValues, templateValues string

	ret := make([]*agentpb.SyncAppData, 0)
	for iter.Scan(
		&config.AppName, &config.Description, &config.Category,
		&config.ChartName, &config.RepoName, &config.RepoURL,
		&config.Namespace, &config.ReleaseName, &config.Version,
		&config.LaunchURL, &config.LaunchUIDescription,
		&config.CreateNamespace, &config.PrivilegedNamespace,
		&overrideValues, &launchUiValues,
		&templateValues, &config.DefualtApp,
		&config.Icon, &config.InstallStatus,
		&config.LastUpdateTime,
	) {
		configCopy := config
		overrideValuesCopy, _ := base64.StdEncoding.DecodeString(overrideValues)
		launchUiValuesCopy, _ := base64.StdEncoding.DecodeString(launchUiValues)
		templateValuesCopy, _ := base64.StdEncoding.DecodeString(templateValues)
		a := &agentpb.SyncAppData{
			Config: &configCopy,
			Values: &agentpb.AppValues{
				OverrideValues: overrideValuesCopy,
				LaunchUIValues: launchUiValuesCopy,
				TemplateValues: templateValuesCopy,
			},
		}
		ret = append(ret, a)
	}

	if err := iter.Close(); err != nil {
		return nil, errors.WithMessage(err, "failed to iterate through results:")
	}
	return ret, nil
}

func formUpdateKvPairs(config *agentpb.SyncAppData) (string, bool) {
	params := []string{}

	if config.Values != nil {
		if len(config.Values.OverrideValues) > 0 {
			encoded := base64.StdEncoding.EncodeToString(config.Values.OverrideValues)
			params = append(params,
				fmt.Sprintf("%s = '%s'", overrideValues, encoded))
		}

		if len(config.Values.LaunchUIValues) > 0 {
			encoded := base64.StdEncoding.EncodeToString(config.Values.LaunchUIValues)
			params = append(params,
				fmt.Sprintf("%s = '%s'", launchUiValues, encoded))
		}

		if len(config.Values.TemplateValues) > 0 {
			encoded := base64.StdEncoding.EncodeToString(config.Values.TemplateValues)
			params = append(params,
				fmt.Sprintf("%s = '%s'", templateValues, encoded))
		}
	}

	if config.Config.CreateNamespace {
		params = append(params,
			fmt.Sprintf("%s = true", createNamespace))
	}
	if config.Config.PrivilegedNamespace {
		params = append(params,
			fmt.Sprintf("%s = true", privilegedNamespace))
	}

	if config.Config.LaunchURL != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", launchUrl, config.Config.LaunchURL))
	}

	if config.Config.LaunchUIDescription != "" {
		params = append(params,
			fmt.Sprintf("%s = '%s'", launchUIDesc, config.Config.LaunchUIDescription))
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

	params = append(params,
		fmt.Sprintf("%s = '%s'", updateTime, time.Now().Format(time.RFC3339)))

	params = append(params,
		fmt.Sprintf("%s = %v", defaultApp, config.Config.DefualtApp))

	if len(params) == 0 {
		// query is empty there is nothing to update
		return "", true
	}
	return strings.Join(params, ", "), false
}
