package cassandraclient

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gocql/gocql"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/server/pkg/credential"
	"github.com/pkg/errors"
)

type Config struct {
	DBAddrs               []string `envconfig:"CASSANDRA_ADDRESSES" required:"true"`
	ServiceUsername       string   `envconfig:"CASSANDRA_SERVICE_USERNAME" required:"true"`
	EntityName            string   `envconfig:"CASSANDRA_ENTITY_NAME" required:"true"`
	Keyspace              string   `envconfig:"CASSANDRA_DB_NAME" required:"true"`
	ClusterTimeout        int      `envconfig:"CLUSTER_TIMEOUT_IN_SEC" default:"20"`
	ClusterConnectTimeout int      `envconfig:"CLUSTER_CONNECT_TIMEOUT_IN_SEC" default:"20"`
	ClusterConistency     uint16   `envconfig:"CLUSTER_CONSISTENCY" default:"6"`
	MaxRetryCount         int      `envconfig:"MAX_RETRY_COUNT" default:"3"`
	MaxConnectionCount    int      `envconfig:"MAX_CLUSTER_CONNECTION_COUNT" default:"5"`
	EnableCassandraTrace  bool     `envconfig:"ENABLE_CASSANDRA_TRACE" default:"false"`
}

type Client struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session
	conf    *Config
}

func NewClient() (*Client, error) {
	conf := &Config{}
	if err := envconfig.Process("", conf); err != nil {
		return nil, fmt.Errorf("cassandra config read faile, %v", err)
	}
	if len(conf.DBAddrs) == 0 {
		return nil, errors.New("cassandra DB addresses are empty")
	}

	cluster := gocql.NewCluster(conf.DBAddrs...)
	cluster.Consistency = gocql.Consistency(conf.ClusterConistency)
	cluster.Timeout = time.Duration(conf.ClusterTimeout) * time.Second
	cluster.ConnectTimeout = time.Duration(conf.ClusterConnectTimeout) * time.Second
	cluster.Keyspace = conf.Keyspace
	cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{NumRetries: conf.MaxRetryCount}
	cluster.NumConns = conf.MaxConnectionCount

	serviceCredential, err := credential.GetServiceUserCredential(context.Background(),
		conf.EntityName, conf.ServiceUsername)
	if err != nil {
		return nil, err
	}

	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: serviceCredential.UserName,
		Password: serviceCredential.Password,
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the DB")
	}

	if conf.EnableCassandraTrace {
		session.SetTrace(gocql.NewTraceWriter(session, os.Stdout))
	}
	store := &Client{
		cluster: cluster,
		session: session,
		conf:    conf,
	}
	store.session.SetConsistency(gocql.Consistency(conf.ClusterConistency))
	return store, nil
}

func (c *Client) Session() *gocql.Session {
	return c.session
}

func (c *Client) Keyspace() string {
	return c.conf.Keyspace
}

func (c *Client) Close() {
	c.session.Close()
}
