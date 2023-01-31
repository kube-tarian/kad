package fetcher

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/kube-tarian/kad/integrator/common-pkg/plugins/utils"
)

const (
	FetchClusterQuery = `select name, kubeconfig from clusters where name = ?;`
)

type ClusterStoreConfiguration struct {
	TableName string `envconfig:"CASSANDRA_CLUSTER_TABLE_NAME" default:"clusters"`
}

func FetchClusterDetails(log logging.Logger, clusterName string) (*ClusterDetails, error) {
	cfg := &ClusterStoreConfiguration{}
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

	pd := &ClusterDetails{}
	// name, kubeconfig
	query := store.GetSession().Query(FetchClusterQuery, clusterName)
	err = query.Scan(
		&pd.Name,
		&pd.Kubeconfig,
	)

	if err != nil {
		log.Errorf("Fetch plugin details failed, %v", err)
		return nil, err
	}
	return pd, nil
}
