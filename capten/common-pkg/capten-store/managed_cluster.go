package captenstore

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kube-tarian/kad/capten/common-pkg/gerrors"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/captenpluginspb"
	postgresdb "github.com/kube-tarian/kad/capten/common-pkg/postgres"
	"gorm.io/gorm"
)

func (a *Store) UpsertManagedCluster(managedCluster *captenpluginspb.ManagedCluster) error {
	if managedCluster.Id == "" {
		cluster := ManagedCluster{
			ID:                  uuid.New(),
			ClusterName:         managedCluster.ClusterName,
			ClusterEndpoint:     managedCluster.ClusterEndpoint,
			ClusterDeployStatus: managedCluster.ClusterDeployStatus,
			AppDeployStatus:     managedCluster.AppDeployStatus,
			LastUpdateTime:      time.Now(),
		}
		return a.dbClient.Create(&cluster)
	}

	cluster := ManagedCluster{ClusterName: managedCluster.ClusterName,
		ClusterEndpoint:     managedCluster.ClusterEndpoint,
		ClusterDeployStatus: managedCluster.ClusterDeployStatus,
		AppDeployStatus:     managedCluster.AppDeployStatus,
		LastUpdateTime:      time.Now()}
	return a.dbClient.Update(&cluster, ManagedCluster{ID: uuid.MustParse(managedCluster.Id)})
}

func (a *Store) DeleteManagedClusterById(id string) error {
	err := a.dbClient.Delete(ManagedCluster{}, ManagedCluster{ID: uuid.MustParse(id)})
	if err != nil {
		err = prepareError(err, id, "Delete")
	}
	return err
}

func (a *Store) GetManagedClusterForID(id string) (*captenpluginspb.ManagedCluster, error) {
	cluster := ManagedCluster{}
	err := a.dbClient.Find(&cluster, ManagedCluster{ID: uuid.MustParse(id)})
	if err != nil {
		return nil, err
	} else if cluster.ID == uuid.Nil {
		return nil, gorm.ErrRecordNotFound
	}

	result := &captenpluginspb.ManagedCluster{
		Id:                  cluster.ID.String(),
		ClusterName:         cluster.ClusterName,
		ClusterEndpoint:     cluster.ClusterEndpoint,
		ClusterDeployStatus: cluster.ClusterDeployStatus,
		AppDeployStatus:     cluster.AppDeployStatus,
		LastUpdateTime:      cluster.LastUpdateTime.String(),
	}
	return result, err
}

func (a *Store) GetManagedClusters() ([]*captenpluginspb.ManagedCluster, error) {
	clusters := []ManagedCluster{}
	err := a.dbClient.Find(&clusters, nil)
	if err != nil && gerrors.GetErrorType(err) != postgresdb.ObjectNotExist {
		return nil, fmt.Errorf("failed to fetch clusters: %v", err.Error())
	}

	result := []*captenpluginspb.ManagedCluster{}
	for _, cluster := range clusters {
		result = append(result, &captenpluginspb.ManagedCluster{
			Id:                  cluster.ID.String(),
			ClusterName:         cluster.ClusterName,
			ClusterEndpoint:     cluster.ClusterEndpoint,
			ClusterDeployStatus: cluster.ClusterDeployStatus,
			AppDeployStatus:     cluster.AppDeployStatus,
			LastUpdateTime:      cluster.LastUpdateTime.String(),
		})
	}

	return result, err
}
