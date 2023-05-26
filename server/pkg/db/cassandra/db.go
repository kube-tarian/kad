package cassandra

import (
	"fmt"
	"github.com/kube-tarian/kad/server/pkg/types"
	"log"
	"os"
	"sync"

	"github.com/gocql/gocql"
)

const (
	keyspace = "capten"
)

type cassandra struct {
	session *gocql.Session
}

var (
	cassandraSession *cassandra
	once             sync.Once
)

func New() (*cassandra, error) {
	var err error
	once.Do(func() {
		cassandraSession = &cassandra{}
		cassandraSession.session, err = connect()
		if err != nil {
			log.Println("failed to connect to cassandra")
		}
	})

	return cassandraSession, err
}

func connect() (*gocql.Session, error) {
	if os.Getenv("CASSANDRA_HOST") == "" {
		return nil, fmt.Errorf("CASSANDRA_HOST env empty")
	}

	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_HOST"))
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: os.Getenv("CASSANDRA_USERNAME"),
		Password: os.Getenv("CASSANDRA_PASSWORD"),
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create c")
	}

	return session, nil
}

func (c *cassandra) GetAgentInfo(customerID string) (*types.AgentInfo, error) {
	agentInfo := types.AgentInfo{}
	iter := c.session.Query(`SELECT endpoint, ca_pem, client_crt, client_key FROM endpoints WHERE customer_id = ?`, customerID).Iter()
	iter.Scan(&agentInfo.Endpoint, &agentInfo.CaPem, &agentInfo.Cert, &agentInfo.Key)
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close the db iterator %v", err)
	}

	return &agentInfo, nil
}

func (c *cassandra) RegisterEndpoint(customerID, endpoint string) error {
	return c.session.Query(
		`INSERT INTO endpoints (customer_id, endpoint) VALUES (?, ?)`,
		customerID,
		endpoint).Exec()
}
