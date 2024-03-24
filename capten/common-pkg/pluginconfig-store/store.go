package pluginconfigstore

import (
	"strings"

	"github.com/intelops/go-common/logging"
	dbclient "github.com/kube-tarian/kad/capten/common-pkg/cassandra/db-client"
	"github.com/kube-tarian/kad/capten/common-pkg/cluster-plugins/clusterpluginspb"
)

const (
	objectNotFoundErrorMessage = "object not found"
)

type Store struct {
	client   *dbclient.Client
	log      logging.Logger
	keyspace string
}

type PluginConfig struct {
	clusterpluginspb.Plugin
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

func (p *PluginConfig) String() string {
	return p.Plugin.String()
}
