package cassandrastore

import "github.com/kube-tarian/kad/integrator/common-pkg/logging"

// This package creates a client to cassandra using go-cql
// CRUD operations specific to table under keyspace should be implemented by user

type Configuration struct {
	CassandraServiceURL string
	Username            string
	Password            string
}

type Client struct {
	log  logging.Logger
	conf *Configuration
}

func NewClientWithConfig(log logging.Logger, conf *Configuration) (*Client, error) {
	return &Client{
		log:  log,
		conf: conf,
	}, nil
}
