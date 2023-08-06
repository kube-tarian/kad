package captenstore

import (
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/db-create/cassandra"
	dbmigration "github.com/kube-tarian/kad/capten/common-pkg/db-migration"
	cassandramigrate "github.com/kube-tarian/kad/capten/common-pkg/db-migration/cassandra"
)

func Migrate(log logging.Logger) error {
	if err := cassandra.Create(log); err != nil {
		return err
	}

	mig, err := cassandramigrate.NewCassandraMigrate(log)
	if err != nil {
		return err
	}

	return mig.Run("AppConfig", dbmigration.UP)
}

func MigratePurge(log logging.Logger) error {
	mig, err := cassandramigrate.NewCassandraMigrate(log)
	if err != nil {
		return err
	}

	return mig.Run("AppConfig", dbmigration.PURGE)
}
