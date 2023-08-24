package dbinit

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/intelops/go-common/logging"
	"github.com/pkg/errors"
)

const (
	createKeyspaceSchemaChangeCQL         = `CREATE KEYSPACE IF NOT EXISTS schema_change WITH REPLICATION = { 'class' : 'NetworkTopologyStrategy', 'replication_factor' : %s } AND DURABLE_WRITES = true`
	createTableKeyspaceLockCQL            = "CREATE TABLE IF NOT EXISTS schema_change.lock(keyspace_to_lock text, started_at timestamp, PRIMARY KEY(keyspace_to_lock)) WITH default_time_to_live = 300"
	createKeyspaceCQL                     = `CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = { 'class' : 'NetworkTopologyStrategy', 'replication_factor' : %s } AND DURABLE_WRITES = true`
	createUser                            = "CREATE USER %s WITH PASSWORD '%s' NOSUPERUSER;"
	alterUser                             = "ALTER USER %s WITH PASSWORD '%s' NOSUPERUSER;"
	grantPermission                       = "GRANT ALL PERMISSIONS ON KEYSPACE %s TO %s ;"
	grantSchemaChangeLockSelectPermission = "GRANT SELECT ON TABLE schema_change.lock TO %s ;"
	grantSchemaChangeLockModifyPermission = "GRANT MODIFY ON TABLE schema_change.lock TO %s ;"
)

type CassandraAdmin struct {
	log     logging.Logger
	session *gocql.Session
}

func NewCassandraAdmin(logger logging.Logger, dbAddrs []string, dbAdminUsername string, dbAdminPassword string) (*CassandraAdmin, error) {
	ca := &CassandraAdmin{
		log: logger,
	}
	err := ca.initSession(dbAddrs, dbAdminUsername, dbAdminPassword)
	if err != nil {
		return nil, err
	}
	return ca, nil
}

func (c *CassandraAdmin) initSession(dbAddrs []string, dbAdminUsername string, dbAdminPassword string) (err error) {
	cluster, err := configureClusterConfig(dbAddrs, dbAdminUsername, dbAdminPassword)
	if err != nil {
		return
	}

	c.session, err = createDbSession(cluster)
	if err != nil {
		return
	}
	return
}

func (c *CassandraAdmin) Session() *gocql.Session {
	return c.session
}

func (c *CassandraAdmin) Close() {
	c.session.Close()
}

func (c *CassandraAdmin) CreateDbUser(serviceUsername string, servicePassword string) (err error) {
	err = c.session.Query(fmt.Sprintf(createUser, serviceUsername, servicePassword)).Exec()
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return c.updateDbUser(serviceUsername, servicePassword)
		} else {
			err = errors.WithMessage(err, "failed to create service user")
			return
		}
	}
	return
}

func (c *CassandraAdmin) GrantPermission(serviceUsername string, dbName string) (err error) {
	err = c.session.Query(fmt.Sprintf(grantSchemaChangeLockSelectPermission, serviceUsername)).Exec()
	if err != nil {
		err = errors.WithMessage(err, "failed to grant select permission to service user on schema_change.lock table")
		return
	}

	err = c.session.Query(fmt.Sprintf(grantSchemaChangeLockModifyPermission, serviceUsername)).Exec()
	if err != nil {
		err = errors.WithMessage(err, "failed to grant modify permission to service user on schema_change.lock table")
		return
	}

	err = c.session.Query(fmt.Sprintf(grantPermission, dbName, serviceUsername)).Exec()
	if err != nil {
		err = errors.WithMessage(err, "failed to grant permission to service user")
		return
	}
	return
}

func (c *CassandraAdmin) CreateDb(dbName string, replicationFactor string) (err error) {
	err = c.session.Query(fmt.Sprintf(createKeyspaceCQL, dbName, replicationFactor)).Exec()
	if err != nil {
		err = errors.WithMessage(err, "failed to create the keyspace")
		return
	}
	return
}

func (c *CassandraAdmin) CreateLockSchemaDb(replicationFactor string) (err error) {
	err = c.session.Query(fmt.Sprintf(createKeyspaceSchemaChangeCQL, replicationFactor)).Exec()
	if err != nil {
		err = errors.WithMessage(err, "failed to create the schema_change keyspace")
		return
	}

	err = retry(3, 2*time.Second, func() (err error) {
		err = c.session.Query(createTableKeyspaceLockCQL).Exec()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		err = errors.WithMessage(err, "failed to create the schema_change.lock table")
		return
	}
	return
}

func configureClusterConfig(addrs []string, adminUsername string, adminPassword string) (cluster *gocql.ClusterConfig, err error) {
	if len(addrs) == 0 {
		err = errors.New("Cassandra addresses are empty")
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

func createDbSession(cluster *gocql.ClusterConfig) (session *gocql.Session, err error) {
	session, err = cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return
}

func (c *CassandraAdmin) updateDbUser(serviceUsername string, servicePassword string) (err error) {
	err = c.session.Query(fmt.Sprintf(alterUser, serviceUsername, servicePassword)).Exec()
	if err != nil {
		err = errors.WithMessage(err, "failed to update service user")
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
