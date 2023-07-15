package db

import (
	"fmt"
	"github.com/kube-tarian/kad/server/pkg/db/cassandra"
	"github.com/kube-tarian/kad/server/pkg/types"

	"github.com/kube-tarian/kad/server/pkg/db/astra"
)

type DB interface {
	GetClusterEndpoint(organizationID, clusterName string) (string, error)
	GetClusters(organizationID string) ([]types.ClusterDetails, error)
	RegisterCluster(organizationID, clusterName, endpoint string) error
	UpdateCluster(organizationID, clusterName, endpoint string) error
	DeleteCluster(organizationID, clusterName string) error
}

func New(db string) (DB, error) {
	switch db {
	case "cassandra":
		return cassandra.New()
	case "astra":
		return astra.New()
	}

	return nil, fmt.Errorf("db: %s not found", db)
}
