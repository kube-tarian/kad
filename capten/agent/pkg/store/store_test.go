package store

import (
	"os"
	"testing"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/pkg/types"
	"github.com/kube-tarian/kad/capten/common-pkg/db-create/cassandra"
	dbmigration "github.com/kube-tarian/kad/capten/common-pkg/db-migration"
	migrator "github.com/kube-tarian/kad/capten/common-pkg/db-migration/cassandra"
	"github.com/stretchr/testify/suite"
)

type StoreSuite struct {
	suite.Suite

	logger logging.Logger
	store  *Store
	mig    *migrator.CassandraMigrate
}

func TestStoreTestSuite(t *testing.T) {

	ss := new(StoreSuite)
	ss.logger = logging.NewLogger()

	cdb := cassandra.NewCassandraStore(ss.logger, nil)
	err := cdb.Connect([]string{"localhost:9042"}, "cassandra", "cassandra")
	if err != nil {
		t.Fatalf("couldn't connect to cassandra, err: %v", err)
	}
	ss.store = New(cdb)

	err = ss.store.CreateLockSchemaDb("1")
	if err != nil {
		t.Fatalf("createLockSchemaDb failed, err: %v", err)
	}

	suite.Run(t, ss)

}

func (suite *StoreSuite) SetupSuite() {

	// create keyspace
	err := suite.store.CreateKeyspace("Apps", "1")
	if err != nil {
		suite.FailNowf("error creating db/keyspace", "err: %v", err)
	}
	suite.logger.Infof("Keyspace %v created!", "Apps")

	// setEnvVars
	os.Setenv("SOURCE_URI", "file://migrations")
	os.Setenv("CASSANDRA_DB_NAME", "apps") // keyspace
	// Note: apparently this keyspace is created in small case

	mig, err := migrator.NewCassandraMigrate(suite.logger)
	if err != nil {
		suite.FailNowf("migrator initialization error", "err: %v", err)
		return
	}
	suite.mig = mig

	// run migrations
	suite.mig.Run("cassandra", dbmigration.UP)
	if err != nil {
		suite.FailNowf("could not perform migrations", "err: %v", err)
	}

}

func (suite *StoreSuite) TearDownSuite() {

	if err := suite.mig.Run("cassandra", dbmigration.PURGE); err != nil {
		suite.logger.Errorf("Error migrating down, err: %v", err)
	}
	if err := suite.store.DropKeyspace("apps"); err != nil {
		suite.logger.Errorf("Error dropping keyspace apps, err: %v", err)
	}
	suite.store.cassandraStore.Session().Close()
}

func (suite *StoreSuite) TestInsertAndGetAppConfigs() {

	for _, config := range configs {
		err := suite.store.InsertAppConfig(config)
		suite.Nil(err)
	}

}

var configs = []types.AppConfig{
	{
		Name:      "App1",
		ChartName: "App1ChartName",
		Override:  types.Override{LaunchUIValues: map[string]any{}, Values: map[string]any{}},
	},

	{
		Name:      "App2",
		ChartName: "App2ChartName",
		Override:  types.Override{LaunchUIValues: map[string]any{}, Values: map[string]any{}},
	},
}
