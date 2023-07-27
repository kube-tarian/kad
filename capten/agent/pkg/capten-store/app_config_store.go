package captenstore

import (
	"fmt"
	"strings"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
)

func init() {
	p := make([]string, len(fields))
	for i := 0; i < len(fields); i++ {
		p[i] = "?"
	}

	placeholder := strings.Join(p, ", ")

	InsertAllAppConfigQuery = "INSERT INTO apps.AppConfig( " +
		strings.Join(fields, ", ") +
		fmt.Sprintf(" ) VALUES (%s)", placeholder)

}

func CreateSelectByFieldNameQuery(field string) string {
	return "SELECT" +
		strings.Join(fields, ", ") +
		"FROM apps.AppConfig" + fmt.Sprintf(" WHERE %s = ?", field)
}

const (
	appName, description, category       = "app_name", "description", "category"
	chartName, repoName, repoUrl         = "chart_name", "repo_name", "repo_url"
	namespace, releaseName, version      = "namespace", "release_name", "version"
	launchUrl, launchRedirectUrl         = "launch_url", "launch_redirect_url"
	createNamespace, privilegedNamespace = "create_namespace", "privileged_namespace"
	overrideValues, launchUiValues       = "override_values", "launch_ui_values"
)

var (
	fields []string = []string{
		appName, description, category,
		chartName, repoName, repoUrl,
		namespace, releaseName, version,
		launchUrl, launchRedirectUrl,
		createNamespace, privilegedNamespace,
		overrideValues, launchUiValues,
	}

	InsertAllAppConfigQuery           string
	UpdateAppConfigByReleaseNameQuery = "UPDATE apps.AppConfig SET %s WHERE release_name = ?"
)

func (a *Store) AddAppConfig(config *agentpb.SyncAppData) error {
	if config.Config == nil {
		return fmt.Errorf("config value nil")
	}
	query := a.client.Session().Query(InsertAllAppConfigQuery)

	override, launchUiValues := "", ""
	if config.Values != nil {
		override, launchUiValues =
			string(config.Values.OverrideValues), string(config.Values.LaunchUIValues)
	}
	c := config.Config
	return query.Bind(
		c.AppName, c.Description, c.Category,

		c.ChartName, c.RepoName, c.RepoURL,

		c.Namespace, c.ReleaseName, c.Version,

		c.LaunchURL, c.LaunchRedirectURL,

		c.CreateNamespace, c.PrivilegedNamespace,

		override, launchUiValues,
	).Exec()
}

func (a *Store) UpdateAppConfig(config *agentpb.SyncAppData) error {

	if config.Config != nil && config.Config.ReleaseName == "" {
		return fmt.Errorf("no release name")
	}

	kvPairs, err := formUpdateKvPairs(config)
	if err != nil {
		return err
	}

	var insertReleaseNameOnlyAppConfigQuery string = `
	INSERT INTO apps.AppConfig(
		release_name
	) VALUES (?)`

	batch := a.client.Session().NewBatch(gocql.LoggedBatch)
	batch.Query(insertReleaseNameOnlyAppConfigQuery, config.Config.ReleaseName)
	batch.Query(fmt.Sprintf(UpdateAppConfigByReleaseNameQuery, kvPairs), config.Config.ReleaseName)
	return a.client.Session().ExecuteBatch(batch)

}

func (a *Store) GetAppConfig(column, value string) (*agentpb.SyncAppData, error) {
	selectQuery := a.client.Session().Query(CreateSelectByFieldNameQuery(column))

	config := new(agentpb.AppConfig)

	var overrideValues, launchUiValues string
	if err := selectQuery.Scan(

		config.AppName, config.Description, config.Category,

		config.ChartName, config.RepoName, config.RepoURL,

		config.Namespace, config.ReleaseName, config.Version,

		config.LaunchURL, config.LaunchRedirectURL,

		config.CreateNamespace, config.PrivilegedNamespace,

		&overrideValues, &launchUiValues,
	); err != nil {
		return nil, err
	}

	return &agentpb.SyncAppData{
		Config: config,
		Values: &agentpb.AppValues{
			OverrideValues: []byte(overrideValues),
			LaunchUIValues: []byte(launchUiValues)},
	}, nil

}

func formUpdateKvPairs(config *agentpb.SyncAppData) (string, error) {
	params := []string{}

	if config.Values != nil && len(config.Values.OverrideValues) > 0 {
		params = append(params,
			fmt.Sprintf("%s = '%s'", overrideValues, string(config.Values.OverrideValues)))
	}

	if config.Values != nil && len(config.Values.LaunchUIValues) > 0 {
		params = append(params,
			fmt.Sprintf("%s = '%s'", launchUiValues, string(config.Values.LaunchUIValues)))

	}

	// -- namespace
	{

		if config.Config.CreateNamespace {
			params = append(params,
				fmt.Sprintf("%s = 'true'", createNamespace))
		}
		if config.Config.PrivilegedNamespace {
			params = append(params,
				fmt.Sprintf("%s = 'true'", privilegedNamespace))
		}

	}

	// -- launchUrl, launchRedirecturl
	{

		if config.Config.LaunchURL != "" {
			params = append(params,
				fmt.Sprintf("%s = '%s'", launchUrl, config.Config.LaunchURL))
		}
		if config.Config.LaunchRedirectURL != "" {
			params = append(params,
				fmt.Sprintf("%s = '%s'", launchRedirectUrl, config.Config.LaunchRedirectURL))
		}
	}

	// -- appName, description, category,
	{
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
	}

	// -- chartName, repoName, repoUrl,
	{
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

	}

	// -- namespace, releaseName, version,
	{
		if config.Config.Namespace != "" {
			params = append(params,
				fmt.Sprintf("%s = '%s'", namespace, config.Config.Namespace))
		}
		if config.Config.ReleaseName != "" {
			params = append(params,
				fmt.Sprintf("%s = '%s'", releaseName, config.Config.ReleaseName))
		}
		if config.Config.Version != "" {
			params = append(params,
				fmt.Sprintf("%s = '%s'", version, config.Config.Version))
		}
	}

	return strings.Join(params, ", "), nil

}
