// Package cassandra contains ...
package cassandra

import "github.com/kube-tarian/kad/integrator/model"

type DBConfig struct {
	DbAddresses         []string `envconfig:"DB_ADDRESSES" required:"true"`
	DbAdminUsername     string   `envconfig:"DB_ADMIN_USERNAME" required:"true"`
	DbServiceUsername   string   `envconfig:"DB_SERVICE_USERNAME" required:"false"`
	DbName              string   `envconfig:"CASSANDRA_DB_NAME" required:"true"`
	DbReplicationFactor string   `envconfig:"DB_REPLICATION_FACTOR" required:"true"`
	DbAdminPassword     string   `envconfig:"DB_ADMIN_PASSWD" required:"true"`
	DbServicePassword   string   `envconfig:"DB_SERVICE_PASSWD" required:"false"`
}

type Store interface {
	Close()
	InsertToolsDb(data *model.Request) error
	DeleteToolsDbEntry(data *model.Request) error
	CreateDbUser(serviceUsername string, servicePassword string) (err error)
	GrantPermission(serviceUsername string, dbName string) (err error)
	CreateDb(keyspace, dbName string, replicationFactor string) (err error)
	CreateLockSchemaDb(replicationFactor string) (err error)
}
