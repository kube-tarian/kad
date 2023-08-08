package dbmigrate

import (
	"os"
	"testing"

	"github.com/intelops/go-common/logging"
	dbinit "github.com/kube-tarian/kad/capten/common-pkg/cassandra/db-init"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	log := logging.NewLogger()
	setEnvConfig()

	err := dbinit.CreatedDatabase(log)
	assert.Nil(t, err)

	err = RunMigrations(log, UP)
	assert.Nil(t, err)

	err = RunMigrations(log, PURGE)
	assert.Nil(t, err)
}

func setEnvConfig() {
	os.Setenv("DB_ADDRESSES", "127.0.0.1:9042")
	os.Setenv("DB_ADMIN_USERNAME", "cassandra")
	os.Setenv("DB_SERVICE_USERNAME", "agent")
	os.Setenv("DB_NAME", "integrator")
	os.Setenv("DB_REPLICATION_FACTOR", `'datacenter1': 1`)
	os.Setenv("DB_ADMIN_PASSWD", "cassandra")
	os.Setenv("DB_SERVICE_PASSWD", "agent")
	os.Setenv("SOURCE_URI", "file://tests/migrations/cassandra")
}
