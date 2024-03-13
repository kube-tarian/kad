package astra

import (
	"fmt"

	astraclient "github.com/kube-tarian/kad/server/pkg/astra-client"
	pb "github.com/stargate/stargate-grpc-go-client/stargate/pkg/proto"
)

type AstraServerStore struct {
	c        *astraclient.Client
	keyspace string
}

func NewStore() (*AstraServerStore, error) {
	a := &AstraServerStore{}
	err := a.initClient()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to astra db, %w", err)
	}
	a.keyspace = a.c.Keyspace()
	return a, nil
}

func (a *AstraServerStore) initClient() error {
	var err error
	a.c, err = astraclient.NewClient()
	return err
}

func (a *AstraServerStore) InitializeDatabase() error {
	initDbQueries := []string{
		fmt.Sprintf(createCaptenClusterTableQuery, a.keyspace),
		fmt.Sprintf(createAppConfigTableQuery, a.keyspace),
		fmt.Sprintf(createPluginStoreConfigTableQuery, a.keyspace),
		fmt.Sprintf(createPluginStoreTableQuery, a.keyspace),
	}
	return a.executeDBQueries(initDbQueries)
}

func (a *AstraServerStore) CleanupDatabase() error {
	initDbQueries := []string{
		//fmt.Sprintf(dropCaptenClusterTableQuery, a.keyspace),
		//fmt.Sprintf(dropAppConfigTableQuery, a.keyspace),
		fmt.Sprintf(dropPluginStoreTableQuery, a.keyspace),
	}
	return a.executeDBQueries(initDbQueries)
}

func (a *AstraServerStore) executeDBQueries(queries []string) error {
	for _, query := range queries {
		createQuery := &pb.Query{
			Cql: query,
		}

		_, err := a.c.Session().ExecuteQuery(createQuery)
		if err != nil {
			return fmt.Errorf("failed to initialise db: %w", err)
		}
	}
	return nil
}
