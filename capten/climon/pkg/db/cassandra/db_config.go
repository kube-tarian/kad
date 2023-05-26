// Package cassandra contains ...
package cassandra

import (
	"github.com/kube-tarian/kad/capten/climon/pkg/types"
	"github.com/kube-tarian/kad/capten/model"
)

type DBConfig struct {
	DbAddresses         []string `envconfig:"CASSANDRA_SERVICE_URL" default:"localhost:9042"`
	DbAdminUsername     string   `envconfig:"CASSANDRA_USERNAME" default:"user"`
	DbServiceUsername   string   `envconfig:"DB_SERVICE_USERNAME" default:"user"`
	DbName              string   `envconfig:"CASSANDRA_KEYSPACE_NAME" default:"capten"`
	DbReplicationFactor string   `envconfig:"DB_REPLICATION_FACTOR" default:"1"`
	DbAdminPassword     string   `envconfig:"CASSANDRA_PASSWORD" default:"password"`
	DbServicePassword   string   `envconfig:"DB_SERVICE_PASSWD" default:"password"`
}

type Store interface {
	Close()
	InsertToolsDb(data *model.ClimonPostRequest) error
	InsertApps(apps []types.App) error
	DeleteToolsDbEntry(data *model.ClimonDeleteRequest) error
	CreateDbUser(serviceUsername string, servicePassword string) (err error)
	GrantPermission(serviceUsername string, dbName string) (err error)
	CreateDb(keyspace, dbName string, replicationFactor string) (err error)
	CreateLockSchemaDb(replicationFactor string) (err error)
}
