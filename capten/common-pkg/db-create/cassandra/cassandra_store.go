// Package cassandra contains ...
package cassandra

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/intelops/go-common/logging"
)

const (
	createKeyspaceSchemaChangeCQL         = `CREATE KEYSPACE IF NOT EXISTS schema_change WITH REPLICATION = { 'class' : 'NetworkTopologyStrategy', 'datacenter1' : %s } AND DURABLE_WRITES = true`
	createTableKeyspaceLockCQL            = "CREATE TABLE IF NOT EXISTS schema_change.lock(keyspace_to_lock text, started_at timestamp, PRIMARY KEY(keyspace_to_lock)) WITH default_time_to_live = 300"
	createKeyspaceCQL                     = `CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = { 'class' : 'NetworkTopologyStrategy', 'datacenter1' : %s } AND DURABLE_WRITES = true`
	createUser                            = "CREATE USER %s WITH PASSWORD '%s' NOSUPERUSER;"
	alterUser                             = "ALTER USER %s WITH PASSWORD '%s' NOSUPERUSER;"
	grantPermission                       = "GRANT ALL PERMISSIONS ON KEYSPACE %s TO %s ;"
	grantSchemaChangeLockSelectPermission = "GRANT SELECT ON TABLE schema_change.lock TO %s ;"
	grantSchemaChangeLockModifyPermission = "GRANT MODIFY ON TABLE schema_change.lock TO %s ;"
)

type CassandraStore struct {
	logg    logging.Logger
	session *gocql.Session
}

func NewCassandraStore(logger logging.Logger, sess *gocql.Session) (store *CassandraStore) {
	store = &CassandraStore{
		logg:    logger,
		session: sess,
	}

	return
}

func (c *CassandraStore) Connect(dbAddrs []string, dbAdminUsername string, dbAdminPassword string) (err error) {
	c.logg.Info("Creating new db cluster configuration")
	cluster, err := configureClusterConfig(dbAddrs, dbAdminUsername, dbAdminPassword)
	if err != nil {
		c.logg.Error("Error creating/configuring new db store", err)
		return
	}

	c.logg.Info("Creating new db session")
	c.session, err = createDbSession(cluster)
	if err != nil {
		return
	}

	return
}

func (c *CassandraStore) Session() *gocql.Session {
	return c.session
}

func (c *CassandraStore) Close() {
	c.logg.Info("closing cassandra session")
	c.session.Close()
}
func (c *CassandraStore) CreateDbUser(serviceUsername string, servicePassword string) (err error) {
	// Create database user for service usage
	err = c.session.Query(fmt.Sprintf(createUser, serviceUsername, servicePassword)).Exec()
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return c.updateDbUser(serviceUsername, servicePassword)
		} else {
			c.logg.Error("Unable to create service user", err)
			return
		}
	}
	return
}

func (c *CassandraStore) GrantPermission(serviceUsername string, dbName string) (err error) {
	err = c.session.Query(fmt.Sprintf(grantSchemaChangeLockSelectPermission, serviceUsername)).Exec()
	if err != nil {
		c.logg.Error("Unable to grant select permission to service user on schema_change.lock table", err)
		return
	}

	err = c.session.Query(fmt.Sprintf(grantSchemaChangeLockModifyPermission, serviceUsername)).Exec()
	if err != nil {
		c.logg.Error("Unable to grant modify permission to service user on schema_change.lock table", err)
		return
	}

	err = c.session.Query(fmt.Sprintf(grantPermission, dbName, serviceUsername)).Exec()
	if err != nil {
		c.logg.Error("Unable to grant permission to service user", err)
		return
	}

	return
}

func (c *CassandraStore) CreateDb(dbName string, replicationFactor string) (err error) {
	// Create keyspace only if it does not already exist
	err = c.session.Query(fmt.Sprintf(createKeyspaceCQL, dbName, replicationFactor)).Exec()
	if err != nil {
		c.logg.Error("Unable to create the keyspace", err)
		return
	}

	return
}

func (c *CassandraStore) CreateLockSchemaDb(replicationFactor string) (err error) {
	// Create keyspace only if it does not already exist
	err = c.session.Query(fmt.Sprintf(createKeyspaceSchemaChangeCQL, replicationFactor)).Exec()
	if err != nil {
		c.logg.Error("Unable to create the schema_change keyspace", err)
		return
	}

	// Create table only if it does not already exist
	err = retry(3, 2*time.Second, func() (err error) {
		err = c.session.Query(createTableKeyspaceLockCQL).Exec()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		c.logg.Error("Unable to create the schema_change.lock table", err)
		return
	}

	return
}

var configureClusterConfig = func(addrs []string, adminUsername string, adminPassword string) (cluster *gocql.ClusterConfig, err error) {
	if len(addrs) == 0 {
		err = errors.New("you must specify a Cassandra address to connect to")
		return
	}

	cluster = gocql.NewCluster(addrs...)
	cluster.Consistency = gocql.All
	cluster.Timeout = 20 * time.Second
	cluster.ConnectTimeout = 20 * time.Second

	if adminUsername != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: adminUsername,
			Password: adminPassword,
		}
	}

	return
}

var createDbSession = func(cluster *gocql.ClusterConfig) (session *gocql.Session, err error) {
	session, err = cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return
}

func (c *CassandraStore) updateDbUser(serviceUsername string, servicePassword string) (err error) {
	// alter  database user for service usage
	err = c.session.Query(fmt.Sprintf(alterUser, serviceUsername, servicePassword)).Exec()
	if err != nil {
		c.logg.Error("Unable to update service user, failed with error : ", err)
		return
	}
	return
}

func retry(attempts int, sleep time.Duration, f func() error) (err error) {
	for i := 0; ; i++ {
		err = f()
		if err == nil {
			return
		}

		if i >= (attempts - 1) {
			break
		}
		time.Sleep(sleep)
	}
	return
}
