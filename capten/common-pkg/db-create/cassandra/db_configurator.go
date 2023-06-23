// Package cassandra contains ...
package cassandra

import (
	"github.com/intelops/go-common/logging"
)

type DbConfigurator struct {
	logg  logging.Logger
	store Store
}

func NewDbConfigurator(logger logging.Logger, stor Store) (handler *DbConfigurator) {
	handler = &DbConfigurator{
		logg:  logger,
		store: stor,
	}

	return
}

func (dbConf *DbConfigurator) ConfigureDb(conf DBConfig) (err error) {
	dbConf.logg.Debug("Creating new db cluster configuration")
	err = dbConf.store.Connect(conf.DbAddresses, conf.DbAdminUsername, conf.DbAdminPassword)
	if err != nil {
		return
	}

	dbConf.logg.Info("Creating lock schema change db")
	err = dbConf.store.CreateLockSchemaDb(conf.DbReplicationFactor)
	if err != nil {
		return
	}

	dbConf.logg.Infof("Creating new db %s with %s", conf.DbName, conf.DbReplicationFactor)
	err = dbConf.store.CreateDb(conf.DbName, conf.DbReplicationFactor)
	if err != nil {
		return
	}

	dbConf.logg.Info("Creating new service users")
	err = dbConf.store.CreateDbUser(conf.DbServiceUsername, conf.DbServicePassword)
	if err != nil {
		return
	}

	dbConf.logg.Info("Grant permission to service user")
	err = dbConf.store.GrantPermission(conf.DbServiceUsername, conf.DbName)
	if err != nil {
		return
	}

	dbConf.store.Close()

	return
}
