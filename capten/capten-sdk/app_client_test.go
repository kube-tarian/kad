package captensdk

import (
	"testing"

	"github.com/kube-tarian/kad/integrator/capten-sdk/agentpb"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/stretchr/testify/assert"
)

var log = logging.NewLogger()

func TestApplicationCreate(t *testing.T) {
	appClient, err := appSetup(t)
	if err != nil {
		return
	}

	req := &agentpb.ApplicationInstallRequest{
		PluginName:  "argocd",
		RepoName:    "argocd-example",
		RepoUrl:     "https://gitlab.privatecloud.sk/vladoportos/argo-cd-example.git",
		ChartName:   "hello-worldhello-world",
		Namespace:   "default",
		ReleaseName: "hello-world",
		Timeout:     5,
	}
	_, err = appClient.Create(req)
	assert.Nilf(t, err, "application create should be success")
	t.Logf("error: %+v", err)
}

func TestApplicationDelete(t *testing.T) {
	appClient, err := appSetup(t)
	if err != nil {
		return
	}

	req := &agentpb.ApplicationDeleteRequest{
		PluginName:  "argocd",
		Namespace:   "default",
		ReleaseName: "hello-world",
		Timeout:     5,
	}
	_, err = appClient.Delete(req)
	assert.Nilf(t, err, "application create should be success")
}

func appSetup(t *testing.T) (*ApplicationClient, error) {
	client, err := NewClient(log)
	assert.Nilf(t, err, "New client should be initialized")
	if err != nil {
		return nil, err
	}

	appClient, err := client.NewApplicationClient(&TransportSSLOptions{IsSSLEnabled: false})
	assert.Nilf(t, err, "New application client should be initialized")
	if err != nil {
		return nil, err
	}
	return appClient, err
}
