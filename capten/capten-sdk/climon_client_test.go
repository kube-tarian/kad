package captensdk

import (
	"testing"

	"github.com/kube-tarian/kad/integrator/capten-sdk/agentpb"
	"github.com/stretchr/testify/assert"
)

func TestClimonCreate(t *testing.T) {
	appClient, err := climonSetup(t)
	if err != nil {
		return
	}

	req := &agentpb.ClimonInstallRequest{
		PluginName:  "helm",
		RepoName:    "argo",
		RepoUrl:     "https://argoproj.github.io/argo-helm",
		ChartName:   "argo-cd",
		Namespace:   "default",
		ReleaseName: "argocd",
		Timeout:     5,
	}
	_, err = appClient.Create(req)
	assert.Nilf(t, err, "application create should be success")
	t.Logf("error: %+v", err)
}

func TestClimonDelete(t *testing.T) {
	appClient, err := climonSetup(t)
	if err != nil {
		return
	}

	req := &agentpb.ClimonDeleteRequest{
		PluginName:  "helm",
		Namespace:   "default",
		ReleaseName: "argocd",
		Timeout:     5,
	}
	_, err = appClient.Delete(req)
	assert.Nilf(t, err, "application create should be success")
}

func climonSetup(t *testing.T) (*ClimonClient, error) {
	client, err := NewClient(log)
	assert.Nilf(t, err, "New client should be initialized")
	if err != nil {
		return nil, err
	}

	climonClient, err := client.NewClimonClient(&TransportSSLOptions{IsSSLEnabled: false})
	assert.Nilf(t, err, "New application client should be initialized")
	if err != nil {
		return nil, err
	}
	return climonClient, err
}
