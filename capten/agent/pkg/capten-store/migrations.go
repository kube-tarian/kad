package captenstore

import (
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/db-create/cassandra"
	dbmigration "github.com/kube-tarian/kad/capten/common-pkg/db-migration"
	cassandramigrate "github.com/kube-tarian/kad/capten/common-pkg/db-migration/cassandra"
)

func CreateDb() {}

func Migrate(log logging.Logger) error {

	// CASSANDRA_DB_NAME
	// DB_SERVICE_USERNAME
	// DB_SERVICE_PASSWD
	if err := cassandra.Create(log); err != nil {
		return err
	}

	// SOURCE_URI
	mig, err := cassandramigrate.NewCassandraMigrate(log)
	if err != nil {
		return err
	}

	return mig.Run("AppConfig", dbmigration.UP)

}
