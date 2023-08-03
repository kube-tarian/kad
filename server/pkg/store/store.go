package store

import (
	"fmt"

	"github.com/kube-tarian/kad/server/pkg/store/cassandra"
	"github.com/kube-tarian/kad/server/pkg/types"

	"github.com/kube-tarian/kad/server/pkg/store/astra"
)

type ServerStore interface {
	InitializeDb() error
	GetClusterEndpoint(organizationID, clusterName string) (string, error)
	GetClusters(organizationID string) ([]types.ClusterDetails, error)
	AddCluster(organizationID, clusterName, endpoint string) error
	UpdateCluster(organizationID, clusterName, endpoint string) error
	DeleteCluster(organizationID, clusterName string) error
	AddOrUpdateApp(config *types.StoreAppConfig) error
	DeleteAppInStore(name, version string) error
	GetAppFromStore(name, version string) (*types.AppConfig, error)
	GetAppsFromStore() (*[]types.AppConfig, error)
}

func NewStore(db string) (ServerStore, error) {
	switch db {
	case "cassandra":
		return cassandra.NewStore()
	case "astra":
		return astra.NewStore()
	}
	return nil, fmt.Errorf("db: %s not found", db)
}
