package captenstore

import (
	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	cassandraclient "github.com/kube-tarian/kad/capten/common-pkg/cassandra-client"
)

type StoreConfig struct {
	Keyspace string `envconfig:"DB_NAME" required:"true"`
}

func GetStoreConfig() (StoreConfig, error) {
	conf := StoreConfig{}
	err := envconfig.Process("", &conf)
	return conf, err
}

type Store struct {
	client *cassandraclient.Client
	conf   StoreConfig
	log    logging.Logger
}

func NewStore(log logging.Logger) (*Store, error) {
	conf, err := GetStoreConfig()
	if err != nil {
		return nil, err
	}
	client, err := cassandraclient.NewClient()
	if err != nil {
		return nil, err
	}
	return &Store{log: log, conf: conf, client: client}, nil
}
