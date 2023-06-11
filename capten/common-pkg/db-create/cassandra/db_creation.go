// Package cassandra contains ...
package cassandra

import (
	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
)

func Create(log logging.Logger) error {
	dbconf := &DBConfig{}
	if err := envconfig.Process("", dbconf); err != nil {
		log.Errorf("Could not parse service config, Usage: %v ", err)
		return err
	}
	dbStore := NewCassandraStore(log, nil)
	dbConfigurator := NewDbConfigurator(log, dbStore)
	log.Info("Start DB configuration")
	err := dbConfigurator.ConfigureDb(*dbconf)
	if err != nil {
		log.Errorf("Could not configure db properly err: %s", err)
		return err
	}
	return nil
}
