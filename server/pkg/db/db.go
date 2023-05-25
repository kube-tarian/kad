package db

import (
	"fmt"

	"github.com/kube-tarian/kad/server/pkg/db/astra"
	"github.com/kube-tarian/kad/server/pkg/db/cassandra"
	"github.com/kube-tarian/kad/server/pkg/types"
)

type DB interface {
	GetAgentInfo(customerID string) (*types.AgentInfo, error)
	RegisterEndpoint(customerID, endpoint string, fileContentMap map[string]string) error
	FetchCreds(customerID, secretPlugin string) (*types.DbCreds, error)
}

func New(db string) (DB, error) {
	switch db {
	case "CASSANDRA":
		return cassandra.New()
	case "ASTRA":
		return astra.New()
	}

	return nil, fmt.Errorf("db: %s not found", db)
}
