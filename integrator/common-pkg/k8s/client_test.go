package k8s

//+local

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/util/homedir"
)

const (
	DEFUALT = "default"
)

func TestClient(t *testing.T) {
	log := logging.NewLogger()
	home := homedir.HomeDir()
	os.Setenv("KUBECONFIG_PATH", filepath.Join(home, ".kube", "config"))

	client, err := NewK8SClient(log)
	assert.Nilf(t, err, "K8S client should be initilized successful")
	assert.Nilf(t, client, "k8s Client should be initiliazed")

	if client != nil {
		pods, err := client.ListPods(DEFUALT)
		assert.Nilf(t, err, "List Pods should be success")
		t.Logf("Pods: %+v", pods)
	}
}
