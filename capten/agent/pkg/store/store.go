package store

import "github.com/kube-tarian/kad/capten/agent/pkg/types"

type StoreIface interface {
	InsertAppConfig(config types.AppConfig) error
	UpsertAppConfig(config types.AppConfig) error
	GetAppConfigByName(name string) (config types.AppConfig, err error)

	DropDB(dbName string) error
	CreateDB(dbName string) error

	CloseSession()
}
