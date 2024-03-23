package postgres

import (
	"os"
	"testing"

	"github.com/intelops/go-common/logging"
	dbinit "github.com/kube-tarian/kad/capten/common-pkg/postgres/db-init"
	"github.com/kube-tarian/kad/capten/common-pkg/postgres/db-migrate/migration"
	"github.com/kube-tarian/kad/capten/common-pkg/postgres/db-migrate/source"
	"github.com/stretchr/testify/assert"
)

var migrationsPath = "tests/migrations/postgres/"

func TestCreateUsingFileSource(t *testing.T) {
	log := logging.NewLogger()
	setEnvConfig("test1")
	os.Setenv("PG_SOURCE_URI", "file://tests/migrations/postgres")

	err := dbinit.CreatedDatabase(log)
	assert.Nil(t, err)

	err = RunMigrations(migration.UP)
	assert.Nil(t, err)
}

func TestCreateUsingBindataSource(t *testing.T) {
	log := logging.NewLogger()
	setEnvConfig("test2")
	os.Setenv("PG_SOURCE_URI", "go-bindata")

	err := dbinit.CreatedDatabase(log)
	assert.Nil(t, err)

	var data = map[string][]byte{}
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		log.Fatalf("reading migration scripts from path %s failed: %v", migrationsPath, err)
	}

	for _, file := range files {
		data[file.Name()], err = getFileContent(migrationsPath, file.Name())
		if err != nil {
			log.Fatalf("reading the migrations file %s content failed, %v", file.Name(), err)
		}
	}

	binData := source.NewBinData(data)

	err = RunMigrationsBinData(GetResource(binData.FileNames, binData.Asset), migration.UP)
	assert.Nil(t, err)
}

func setEnvConfig(serviceUser string) {
	os.Setenv("PG_DB_ADDRESS", "127.0.0.1")
	os.Setenv("PG_DB_ENTITY_NAME", "postgres")
	os.Setenv("PG_DB_ADMIN_PASSWORD", "example")
	os.Setenv("PG_DB_SERVICE_USERNAME", serviceUser)
	os.Setenv("PG_DB_SERVICE_USER_PASSWORD", serviceUser)
	os.Setenv("PG_DB_NAME", serviceUser)
}

func getFileContent(migrationsPath, fileName string) ([]byte, error) {
	return os.ReadFile(migrationsPath + fileName)
}
