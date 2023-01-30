package fetcher

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/kube-tarian/kad/integrator/common-pkg/plugins/utils"
)

const (
	FetchPluginQuery = `select name, repo_name, repo_url, chart_name, namespace, release_name, version from ? where name = ?;`
)

type PluginConfiguration struct {
	TableName string `envconfig:"CASSANDRA_TABLE_NAME" default:"tools"`
}

func FetchPluginDetails(log logging.Logger, pluginName string) (*PluginDetails, error) {
	cfg := &PluginConfiguration{}
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Errorf("Cassandra configuration detail missing, %v", err)
		return nil, err
	}

	// Fetch the plugin details from Cassandra
	store, err := utils.NewStore(log)
	if err != nil {
		log.Errorf("Store initialization failed, %v", err)
		return nil, err
	}
	defer store.Close()

	pd := &PluginDetails{}
	// name, repo_name, repo_url, chart_name, namespace, release_name, version
	query := store.GetSession().Query(FetchPluginQuery, cfg.TableName, pluginName)
	err = query.Scan(
		&pd.Name,
		&pd.RepoName,
		&pd.RepoURL,
		&pd.ChartName,
		&pd.Namespace,
		&pd.ReleaseName,
		&pd.Version,
	)

	if err != nil {
		log.Errorf("Fetch plugin details failed, %v", err)
		return nil, err
	}
	return pd, nil
}
