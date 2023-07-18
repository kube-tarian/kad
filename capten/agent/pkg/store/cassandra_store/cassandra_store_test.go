package cassandra_store

import (
	"os"
	"testing"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/pkg/store"
	"github.com/kube-tarian/kad/capten/agent/pkg/types"
	"github.com/kube-tarian/kad/capten/common-pkg/db-create/cassandra"
	dbmigration "github.com/kube-tarian/kad/capten/common-pkg/db-migration"
	migrator "github.com/kube-tarian/kad/capten/common-pkg/db-migration/cassandra"
	"github.com/stretchr/testify/suite"
)

type StoreSuite struct {
	suite.Suite

	logger logging.Logger
	cs     store.StoreIface
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
	err = cdb.CreateLockSchemaDb("1")
	if err != nil {
		t.Fatalf("createLockSchemaDb failed, err: %v", err)
	}
	ss.cs = New(cdb)

	suite.Run(t, ss)

}

func (suite *StoreSuite) SetupSuite() {

	// create keyspace
	err := suite.cs.CreateDB("apps")
	if err != nil {
		suite.FailNowf("error creating db/keyspace", "err: %v", err)
	}

	// setEnvVars
	os.Setenv("SOURCE_URI", "file://migrations")
	os.Setenv("CASSANDRA_DB_NAME", "apps") // dbName/keyspace

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
	if err := suite.cs.DropDB("apps"); err != nil {
		suite.logger.Errorf("Error dropping keyspace apps, err: %v", err)
	}
	suite.cs.CloseSession()
}

func (suite *StoreSuite) TestInsertAndGetAppConfigs() {

	for _, config := range configs {
		err := suite.cs.InsertAppConfig(config)
		suite.Nil(err)
	}

	for _, config := range upsertConfigs {
		err := suite.cs.UpsertAppConfig(config)
		suite.Nil(err)
	}

	cfg, err := suite.cs.GetAppConfigByName("App1")
	suite.Nil(err)

	suite.Equal("App1", cfg.Name)
	suite.Equal("App1ChartName", cfg.ChartName)
	suite.Empty(cfg.RepoName, "non empty repoName")
	suite.Equal("2", cfg.Version)

}

var configs = []types.AppConfig{
	{
		Name:      "App1",
		ChartName: "App1ChartName",
		Version:   "1",
		Override:  &types.Override{LaunchUIValues: map[string]any{}, Values: map[string]any{}},
	},

	{
		Name:      "App2",
		ChartName: "App2ChartName",
		Override:  &types.Override{LaunchUIValues: map[string]any{}, Values: map[string]any{}},
	},
}

var upsertConfigs = []types.AppConfig{
	{
		Name:    "App1",
		Version: "2",
	},
}
