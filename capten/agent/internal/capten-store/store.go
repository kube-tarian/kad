package captenstore

import (
	"strings"

	"github.com/intelops/go-common/logging"
	dbclient "github.com/kube-tarian/kad/capten/common-pkg/cassandra/db-client"
)

const (
	objectNotFoundErrorMessage = "object not found"
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

func IsObjectNotFound(err error) bool {
	if err == nil {
		return false
	}

	if strings.Contains(err.Error(), objectNotFoundErrorMessage) {
		return true
	}
	return false
}

func removePlugin(targetPlugin string, usedPlugins []string) []string {
	result := []string{}
	for _, v := range usedPlugins {
		if v != targetPlugin {
			result = append(result, v)
		}
	}
	return result
}
