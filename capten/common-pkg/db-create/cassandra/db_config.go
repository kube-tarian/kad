// Package cassandra contains ...
package cassandra

import "github.com/gocql/gocql"

type DBConfig struct {
	DbAddresses         []string `envconfig:"DB_ADDRESSES" required:"true"`
	DbAdminUsername     string   `envconfig:"DB_ADMIN_USERNAME"`
	DbReplicationFactor string   `envconfig:"DB_REPLICATION_FACTOR" required:"true" default:"1"`
	DbAdminPassword     string   `envconfig:"DB_ADMIN_PASSWD"`

	DbName            string `envconfig:"DB_NAME" required:"true"`
	DbServiceUsername string `envconfig:"DB_SERVICE_USERNAME" required:"true"`
	DbServicePassword string `envconfig:"DB_SERVICE_PASSWD"`
}

type Store interface {
	Connect(dbAddrs []string, dbAdminUsername string, dbAdminPassword string) (err error)
	Session() *gocql.Session
	Close()
	CreateDbUser(serviceUsername string, servicePassword string) (err error)
	GrantPermission(serviceUsername string, dbName string) (err error)
	CreateDb(dbName string, replicationFactor string) (err error)
	CreateLockSchemaDb(replicationFactor string) (err error)
}
