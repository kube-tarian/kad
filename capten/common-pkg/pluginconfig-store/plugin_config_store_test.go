package pluginconfigstore

import (
	"os"
	"testing"

	dbinit "github.com/kube-tarian/kad/capten/common-pkg/cassandra/db-init"
	dbmigrate "github.com/kube-tarian/kad/capten/common-pkg/cassandra/db-migrate"
	"github.com/kube-tarian/kad/capten/common-pkg/cluster-plugins/clusterpluginspb"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	log := logging.NewLogger()
	setEnvConfig()

	// DB init
	err := dbinit.CreatedDatabase(log)
	assert.Nil(t, err)

	// run migrations
	err = dbmigrate.RunMigrations(log, dbmigrate.UP)
	assert.Nil(t, err)

	// Store
	pcStore, err := NewStore(log)
	assert.Nil(t, err)

	log.Info("%v", pcStore)
	err = pcStore.UpsertPluginConfig(&PluginConfig{
		Plugin: &clusterpluginspb.Plugin{
			PluginName:          "test",
			StoreType:           clusterpluginspb.StoreType_LOCAL_CAPTEN_STORE,
			Description:         "test",
			Category:            "ci/cd",
			Version:             "0.0.1",
			Icon:                []byte{},
			ChartName:           "test",
			ChartRepo:           "http://test.com",
			DefaultNamespace:    "default",
			PrivilegedNamespace: true,
			ApiEndpoint:         "http://test.{{domain}}",
			UiEndpoint:          "http://testui.{{domain}}",
			Capabilities:        []string{"vaultstore", "captensdk", "postgresstore"},
			Values:              []byte{},
			OverrideValues:      []byte{},
			InstallStatus:       "",
		},
	})
	assert.Nil(t, err)

	plugConfigTest, err := pcStore.GetPluginConfig("test")
	assert.Nil(t, err)
	t.Logf(" get plugin: %+v", plugConfigTest)

	plugins, err := pcStore.GetAllPlugins()
	assert.Nil(t, err)
	t.Logf("All plugins: %+v", plugins)
}

func setEnvConfig() {
	os.Setenv("DB_ADDRESSES", "127.0.0.1:9042")
	os.Setenv("DB_ADMIN_USERNAME", "cassandra")
	os.Setenv("DB_ENTITY_NAME", "cassandra")
	os.Setenv("DB_SERVICE_USERNAME", "agent")
	os.Setenv("DB_NAME", "integrator")
	os.Setenv("DB_ADMIN_PASSWD", "cassandra")
	os.Setenv("DB_SERVICE_PASSWORD", "agent")
	os.Setenv("SOURCE_URI", "file://./tests/migrations")
	// DB_ENTITY_NAME
}
