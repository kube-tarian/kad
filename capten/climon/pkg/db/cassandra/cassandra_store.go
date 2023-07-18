// Package cassandra contains ...
package cassandra

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/kube-tarian/kad/capten/climon/pkg/types"

	"github.com/kube-tarian/kad/capten/model"

	"github.com/gocql/gocql"
	"github.com/intelops/go-common/logging"
)

const (
	createKeyspaceSchemaChangeCQL         = `CREATE KEYSPACE IF NOT EXISTS schema_change WITH REPLICATION = { 'class' : 'NetworkTopologyStrategy', %s } AND DURABLE_WRITES = true`
	createTableKeyspaceLockCQL            = "CREATE TABLE IF NOT EXISTS schema_change.lock(keyspace_to_lock text, started_at timestamp, PRIMARY KEY(keyspace_to_lock)) WITH default_time_to_live = 300"
	createKeyspaceCQL                     = "CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = { 'class' : 'NetworkTopologyStrategy', 'datacenter1' : %s } AND DURABLE_WRITES = true"
	createUser                            = "CREATE USER %s WITH PASSWORD '%s' NOSUPERUSER;"
	alterUser                             = "ALTER USER %s WITH PASSWORD '%s' NOSUPERUSER;"
	grantPermission                       = "GRANT ALL PERMISSIONS ON KEYSPACE %s TO %s ;"
	grantSchemaChangeLockSelectPermission = "GRANT SELECT ON TABLE schema_change.lock TO %s ;"
	grantSchemaChangeLockModifyPermission = "GRANT MODIFY ON TABLE schema_change.lock TO %s ;"
	createToolsTableCQL                   = "CREATE TABLE IF NOT EXISTS %s.tools ( name text, repo_name text, repo_url text, chart_name text, namespace text, release_name text, version text, PRIMARY KEY (name))"
	captenKeyspace                        = "capten"
	insertToolsCQL                        = `INSERT INTO capten.tools (name, repo_name, repo_url, chart_name, namespace, release_name, version) VALUES (?, ?, ?, ?, ?, ?, ?)`
	deleteToolsCQL                        = "DELETE FROM capten.tools where name='%s' if exists"
)

type cassandraStore struct {
	log     logging.Logger
	session *gocql.Session
}

var (
	cassandraStoreObj *cassandraStore
	once              sync.Once
)

func NewCassandraStore(dbAddress []string, username, password string) (Store, error) {
	var err error
	once.Do(func() {
		cassandraStoreObj = &cassandraStore{}
		cassandraStoreObj.session, err = Connect(dbAddress, username, password)
	})

	return cassandraStoreObj, err
}

func Connect(dbAddress []string, dbAdminUsername string, dbAdminPassword string) (*gocql.Session, error) {
	cluster, err := configureClusterConfig(dbAddress, dbAdminUsername, dbAdminPassword)
	if err != nil {
		return nil, err
	}

	return createDbSession(cluster)
}

func GetStore() *cassandraStore {
	return cassandraStoreObj
}

func (c *cassandraStore) Close() {
	c.session.Close()
}
func (c *cassandraStore) CreateDbUser(serviceUsername string, servicePassword string) (err error) {
	// Create database user for service usage
	err = c.session.Query(fmt.Sprintf(createUser, serviceUsername, servicePassword)).Exec()
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return c.updateDbUser(serviceUsername, servicePassword)
		} else {
			c.log.Error("Unable to create service user", err)
			return
		}
	}
	return
}

func (c *cassandraStore) GrantPermission(serviceUsername string, dbName string) (err error) {
	err = c.session.Query(fmt.Sprintf(grantSchemaChangeLockSelectPermission, serviceUsername)).Exec()
	if err != nil {
		c.log.Error("Unable to grant select permission to service user on schema_change.lock table", err)
		return
	}

	err = c.session.Query(fmt.Sprintf(grantSchemaChangeLockModifyPermission, serviceUsername)).Exec()
	if err != nil {
		c.log.Error("Unable to grant modify permission to service user on schema_change.lock table", err)
		return
	}

	err = c.session.Query(fmt.Sprintf(grantPermission, dbName, serviceUsername)).Exec()
	if err != nil {
		c.log.Error("Unable to grant permission to service user", err)
		return
	}

	return
}

func (c *cassandraStore) CreateDb(keyspace, dbName string, replicationFactor string) error {
	if err := c.session.Query(fmt.Sprintf(createKeyspaceCQL, keyspace, replicationFactor)).Exec(); err != nil {
		c.log.Error("Unable to create the keyspace", err)
		return err
	}

	if err := c.session.Query(fmt.Sprintf(createToolsTableCQL, keyspace)).Exec(); err != nil {
		c.log.Error("Unable to create the tools table", err)
		return err
	}

	return nil
}

func (c *cassandraStore) CreateLockSchemaDb(replicationFactor string) (err error) {
	// Create keyspace only if it does not already exist
	err = c.session.Query(fmt.Sprintf(createKeyspaceSchemaChangeCQL, replicationFactor)).Exec()
	if err != nil {
		c.log.Error("Unable to create the schema_change keyspace", err)
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
		c.log.Error("Unable to create the schema_change.lock table", err)
		return
	}

	return
}

func configureClusterConfig(addrs []string, adminUsername string, adminPassword string) (cluster *gocql.ClusterConfig, err error) {
	if len(addrs) == 0 {
		err = errors.New("you must specify a Cassandra address to connect to")
		return
	}

	cluster = gocql.NewCluster(addrs...)
	cluster.Consistency = gocql.One
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

func (c *cassandraStore) updateDbUser(serviceUsername string, servicePassword string) (err error) {
	// alter  database user for service usage
	err = c.session.Query(fmt.Sprintf(alterUser, serviceUsername, servicePassword)).Exec()
	if err != nil {
		c.log.Error("Unable to update service user, failed with error : ", err)
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

func (c *cassandraStore) InsertToolsDb(data *model.ClimonPostRequest) error {
	return c.session.Query(insertToolsCQL,
		data.ReleaseName,
		data.RepoName,
		data.RepoUrl,
		data.ChartName,
		data.Namespace,
		data.ReleaseName,
		data.Version).Exec()
}

func (c *cassandraStore) InsertApps(apps []types.App) error {
	for _, app := range apps {
		err := c.session.Query(insertToolsCQL,
			app.ReleaseName,
			app.RepoName,
			app.RepoURL,
			app.ChartName,
			app.Namespace,
			app.ReleaseName,
			app.Version).Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *cassandraStore) DeleteToolsDbEntry(data *model.ClimonDeleteRequest) error {
	return c.session.Query(fmt.Sprintf(deleteToolsCQL, data.ReleaseName)).Exec()
}
