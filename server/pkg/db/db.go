package db

import (
	"fmt"
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
		cassandraSession.session, err = Connect()
		if err != nil {
			log.Println("failed to connect to cassandra")
		}
	})

	return cassandraSession, err
}

func Connect() (*gocql.Session, error) {
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

func (c *cassandra) GetEndpoint(customerID string) (string, error) {
	var endpoint string
	iter := c.session.Query(`SELECT endpoint FROM endpoints WHERE customer_id = ?`, customerID).Iter()
	iter.Scan(&endpoint)
	if err := iter.Close(); err != nil {
		return "", fmt.Errorf("failed to close the db iterator %v", err)
	}

	return endpoint, nil
}

func (c *cassandra) RegisterEndpoint(customerID, endpoint string) error {
	if err := c.session.Query(`INSERT INTO endpoints (customer_id, endpoint) VALUES (?, ?, ?)`,
		customerID, endpoint).Exec(); err != nil {
		return err
	}

	return nil
}
