package utils

import (
	"os"
	"testing"

	"github.com/kube-tarian/kad/capten/common-pkg/logging"
	"github.com/stretchr/testify/assert"
)

var log = logging.NewLogger()

func TestFetchArgoCDPluginDetails(t *testing.T) {
	os.Setenv("CASSANDRA_SERVICE_URL", "127.0.0.1:9042")
	os.Setenv("CASSANDRA_USERNAME", "user")
	os.Setenv("CASSANDRA_PASSWORD", "password")
	os.Setenv("CASSANDRA_KEYSPACE_NAME", "capten")
	os.Setenv("CASSANDRA_TABLE_NAME", "tools")

	store, err := NewStore(log)
	assert.Nilf(t, err, "Store initialization should be passed")
	if err != nil {
		return
	}
	assert.NotNilf(t, store, "store session should get initialized")

	// pd, err := store.FetchPluginDetails("argocd")
	// assert.Nilf(t, err, "argocd plugin details should be able fetch")
	// assert.NotNilf(t, pd, "argocd plugin details failed to fetch")
	// if err != nil {
	// 	return
	// }

	// t.Logf("argocd plugin details: %+v", pd)
}
