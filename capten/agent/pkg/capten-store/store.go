package captenstore

import (
	"github.com/intelops/go-common/logging"
	dbclient "github.com/kube-tarian/kad/capten/common-pkg/cassandra/db-client"
)

type Store struct {
	client   *dbclient.Client
	log      logging.Logger
	keyspace string
}

func NewStore(log logging.Logger) (*Store, error) {
	client, err := dbclient.NewClient()
	if err != nil {
		return nil, err
	}
	return &Store{log: log, client: client, keyspace: client.Keyspace()}, nil
}
