package captenstore

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"
	"github.com/pkg/errors"
)

const (
	insertManagedCluster                 = "INSERT INTO %s.ManagedClusters(id, cluster_name, cluster_endpoint, cluster_deploy_status, app_deploy_status, last_update_time) VALUES (?,?,?,?,?,?)"
	insertManagedClusterId               = "INSERT INTO %s.ManagedClusters(id) VALUES (?)"
	updateManagedClusterById             = "UPDATE %s.ManagedClusters SET %s WHERE id=?"
	deleteManagedClusterById             = "DELETE FROM %s.ManagedClusters WHERE id= ?"
	selectAllManagedClusters             = "SELECT id, cluster_name, cluster_endpoint, cluster_deploy_status, app_deploy_status, last_update_time FROM %s.ManagedClusters"
	selectAllManagedClustersByLabels     = "SELECT id, cluster_name, cluster_endpoint, cluster_deploy_status, app_deploy_status, last_update_time FROM %s.ManagedClusters WHERE %s"
	selectGetManagedClusterById          = "SELECT id, cluster_name, cluster_endpoint, cluster_deploy_status, app_deploy_status, last_update_time FROM %s.ManagedClusters WHERE id=%s;"
	selectGetManagedClusterByClusterName = "SELECT id, cluster_name, cluster_endpoint, cluster_deploy_status, app_deploy_status, last_update_time FROM %s.ManagedClusters WHERE cluster_name=%s ALLOW FILTERING;"
)

func (a *Store) UpsertManagedCluster(config *captenpluginspb.ManagedCluster) error {
	config.LastUpdateTime = time.Now().Format(time.RFC3339)
	batch := a.client.Session().NewBatch(gocql.LoggedBatch)

	query := fmt.Sprintf(selectGetManagedClusterByClusterName, a.keyspace, config.ClusterName)
	clusters, err := a.executeManagedClustersSelectQuery(query)
	if err != nil {
		batch.Query(fmt.Sprintf(insertManagedCluster, a.keyspace), config.Id, config.ClusterName, config.ClusterEndpoint, config.ClusterDeployStatus, config.AppDeployStatus, config.LastUpdateTime)
	} else if len(clusters) > 0 {
		updatePlaceholders, values := formUpdateKvPairsForManagedCluster(config)
		if updatePlaceholders == "" {
			return fmt.Errorf("empty values found")
		}
		query := fmt.Sprintf(updateManagedClusterById, a.keyspace, updatePlaceholders)
		args := append(values, config.Id)
		batch.Query(query, args...)
	}

	if err := a.client.Session().ExecuteBatch(batch); err != nil {
		return err
	}

	return nil
}

func (a *Store) DeleteManagedClusterById(id string) error {
	deleteAction := a.client.Session().Query(fmt.Sprintf(deleteManagedClusterById,
		a.keyspace), id)
	err := deleteAction.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (a *Store) GetManagedClusterForID(id string) (*captenpluginspb.ManagedCluster, error) {
	query := fmt.Sprintf(selectGetManagedClusterById, a.keyspace, id)
	clusters, err := a.executeManagedClustersSelectQuery(query)
	if err != nil {
		return nil, err
	}

	if len(clusters) == 0 {
		return nil, fmt.Errorf("cluster not found")
	}
	return clusters[0], nil
}

func (a *Store) GetManagedClusters() ([]*captenpluginspb.ManagedCluster, error) {
	query := fmt.Sprintf(selectAllManagedClusters, a.keyspace)
	return a.executeManagedClustersSelectQuery(query)
}

func (a *Store) executeManagedClustersSelectQuery(query string) ([]*captenpluginspb.ManagedCluster, error) {
	selectQuery := a.client.Session().Query(query)
	iter := selectQuery.Iter()

	cluster := captenpluginspb.ManagedCluster{}

	ret := make([]*captenpluginspb.ManagedCluster, 0)
	for iter.Scan(
		&cluster.Id, &cluster.ClusterName, &cluster.ClusterEndpoint,
		&cluster.ClusterDeployStatus, &cluster.AppDeployStatus, &cluster.LastUpdateTime,
	) {
		ManagedCluster := &captenpluginspb.ManagedCluster{
			Id:                  cluster.Id,
			ClusterName:         cluster.ClusterName,
			ClusterEndpoint:     cluster.ClusterEndpoint,
			ClusterDeployStatus: cluster.ClusterDeployStatus,
			AppDeployStatus:     cluster.AppDeployStatus,
			LastUpdateTime:      cluster.LastUpdateTime,
		}
		ret = append(ret, ManagedCluster)
	}

	if err := iter.Close(); err != nil {
		return nil, errors.WithMessage(err, "executeManagedClustersSelectQuery: failed to iterate through results:")
	}

	return ret, nil
}

func formUpdateKvPairsForManagedCluster(config *captenpluginspb.ManagedCluster) (updatePlaceholders string, values []interface{}) {
	params := []string{}

	if config.ClusterName != "" {
		params = append(params, "cluster_name = ?")
		values = append(values, config.ClusterName)
	}

	if config.ClusterEndpoint != "" {
		params = append(params, "cluster_endpoint = ?")
		values = append(values, config.ClusterEndpoint)
	}

	if config.ClusterDeployStatus != "" {
		params = append(params, "cluster_deploy_status = ?")
		values = append(values, config.ClusterDeployStatus)
	}

	if config.AppDeployStatus != "" {
		params = append(params, "app_deploy_status = ?")
		values = append(values, config.AppDeployStatus)
	}

	if config.LastUpdateTime != "" {
		params = append(params, "last_update_time = ?")
		values = append(values, config.LastUpdateTime)
	}

	if len(params) == 0 {
		return "", nil
	}
	return strings.Join(params, ", "), values
}
