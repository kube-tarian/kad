package store

import (
	"fmt"

	"github.com/kube-tarian/kad/server/pkg/types"

	"github.com/kube-tarian/kad/server/pkg/store/astra"
)

type ServerStore interface {
	InitializeDb() error
	GetClusterDetails(orgID, clusterID string) (*types.ClusterDetails, error)
	GetClusterForOrg(orgID string) (*types.ClusterDetails, error)
	GetClusters(orgID string) ([]types.ClusterDetails, error)
	AddCluster(orgID, clusterID, clusterName, endpoint string) error
	UpdateCluster(orgID, clusterID, clusterName, endpoint string) error
	DeleteCluster(orgID, clusterID string) error
	AddOrUpdateStoreApp(config *types.StoreAppConfig) error
	DeleteAppInStore(name, version string) error
	GetAppFromStore(name, version string) (*types.AppConfig, error)
	GetAppsFromStore() (*[]types.AppConfig, error)
	GetStoreAppValues(name, version string) (*types.AppConfig, error)
}

func NewStore(db string) (ServerStore, error) {
	switch db {
	case "astra":
		return astra.NewStore()
	}
	return nil, fmt.Errorf("db: %s not found", db)
}
