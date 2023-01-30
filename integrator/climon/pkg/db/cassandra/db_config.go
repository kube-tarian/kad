// Package cassandra contains ...
package cassandra

import "github.com/kube-tarian/kad/integrator/model"

type DBConfig struct {
	DbAddresses         []string `envconfig:"CASSANDRA_SERVICE_URL" required:"true"`
	DbAdminUsername     string   `envconfig:"CASSANDRA_USERNAME" required:"true"`
	DbServiceUsername   string   `envconfig:"DB_SERVICE_USERNAME" default:"user"`
	DbName              string   `envconfig:"CASSANDRA_KEYSPACE_NAME" required:"true"`
	DbReplicationFactor string   `envconfig:"DB_REPLICATION_FACTOR" default:"1"`
	DbAdminPassword     string   `envconfig:"CASSANDRA_PASSWORD" required:"true"`
	DbServicePassword   string   `envconfig:"DB_SERVICE_PASSWD" default:"password"`
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
