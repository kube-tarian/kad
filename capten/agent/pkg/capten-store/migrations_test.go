package captenstore

import (
	"os"
	"testing"

	"github.com/intelops/go-common/logging"
	"github.com/stretchr/testify/require"
)

func TestMigrations(t *testing.T) {
	assert := require.New(t)

	setEnvVars()

	logger := logging.NewLogger()
	err := Migrate(logger)
	assert.Nil(err)

}

func setEnvVars() {

	os.Setenv("CASSANDRA_DB_NAME", "apps")
	os.Setenv("DB_SERVICE_USERNAME", "apps_user")
	os.Setenv("DB_SERVICE_PASSWD", "apps_password")
	os.Setenv("SOURCE_URI", "file://test_migrations")

}
