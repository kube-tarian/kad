// Package cassandra contains ...
package cassandra

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
)

func Create(log logging.Logger) (Store, error) {
	dbConf, err := GetDbConfig()
	if err != nil {
		log.Errorf("failed to parse db config", err)
		return nil, err
	}

	dbStore, err := NewCassandraStore(log, dbConf.DbAddresses, dbConf.DbAdminUsername, dbConf.DbAdminPassword)
	if err != nil {
		log.Errorf("failed to connect to ca")
		return nil, err
	}

	dbConfigurator := NewDbConfigurator(log, dbStore)
	log.Info("Start DB configuration")
	if err := dbConfigurator.ConfigureDb(*dbConf); err != nil {
		log.Errorf("Could not configure db properly err: %s", err)
		return nil, err
	}

	return dbStore, nil
}

func GetDbConfig() (*DBConfig, error) {
	dbConf := &DBConfig{}
	if err := envconfig.Process("", dbConf); err != nil {
		return nil, err
	}

	return dbConf, nil
}