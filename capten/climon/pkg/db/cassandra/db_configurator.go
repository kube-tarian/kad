// Package cassandra contains ...
package cassandra

import (
	"github.com/kube-tarian/kad/capten/common-pkg/logging"
)

type DbConfigurator struct {
	logg  logging.Logger
	store Store
}

func NewDbConfigurator(logger logging.Logger, store Store) (handler *DbConfigurator) {
	handler = &DbConfigurator{
		logg:  logger,
		store: store,
	}

	return
}

func (dbConf *DbConfigurator) ConfigureDb(conf DBConfig) (err error) {
	//dbConf.logg.Debug("Creating new db cluster configuration")
	//err = dbConf.store.Connect(conf.DbAddresses, conf.DbAdminUsername, conf.DbAdminPassword)
	//if err != nil {
	//	return
	//}

	dbConf.logg.Infof("Creating new db %s with %s", conf.DbName, conf.DbReplicationFactor)
	err = dbConf.store.CreateDb(captenKeyspace, conf.DbName, conf.DbReplicationFactor)
	if err != nil {
		return
	}
	return
}
