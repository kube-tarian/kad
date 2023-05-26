package db

import (
	"fmt"

	"github.com/kube-tarian/kad/server/pkg/db/astra"
	"github.com/kube-tarian/kad/server/pkg/db/cassandra"
	"github.com/kube-tarian/kad/server/pkg/types"
)

type DB interface {
	GetAgentInfo(customerID string) (*types.AgentInfo, error)
	RegisterEndpoint(customerID, endpoint string) error
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
