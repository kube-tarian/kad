package integrationtests

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/kube-tarian/kad/integrator/deployment-worker/pkg/application"
	"github.com/kube-tarian/kad/integrator/deployment-worker/pkg/model"
	"github.com/kube-tarian/kad/integrator/deployment-worker/pkg/workflows"
	"github.com/kube-tarian/kad/integrator/pkg/logging"
	"go.temporal.io/sdk/client"
)

var logger = logging.NewLogger()

func TestMain(m *testing.M) {
	m.Run()
}

func TestIntegrationArgocdDeploymentEvent(t *testing.T) {
	testData := setup()

	stop := startMain()

	data := &model.Request{
		RepoName:    "argocd-example",
		RepoURL:     "https://gitlab.privatecloud.sk/vladoportos/argo-cd-example.git",
		ChartName:   "hello-world",
		Namespace:   "default",
		ReleaseName: "hello-world",
		Timeout:     5,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Data marshalling failed, %v", err)
	}

	sendDeploymentEvent(t, "argocd", dataJSON, "install")
	logger.Info("Sleeping now for 5 seconds")
	time.Sleep(5 * time.Second)

	logger.Info("Starting teardown")
	tearDown(testData)
	stop <- true
}

func TestIntegrationArgocdDeleteEvent(t *testing.T) {
	testData := setup()

	stop := startMain()

	data := &model.Request{
		RepoName:    "argocd-example",
		RepoURL:     "https://gitlab.privatecloud.sk/vladoportos/argo-cd-example.git",
		ChartName:   "hello-world",
		Namespace:   "default",
		ReleaseName: "hello-world",
		Timeout:     5,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Data marshalling failed, %v", err)
	}

	sendDeploymentEvent(t, "argocd", dataJSON, "delete")
	logger.Info("Sleeping now for 5 seconds")
	time.Sleep(5 * time.Second)

	logger.Info("Starting teardown")
	tearDown(testData)
	stop <- true
}

func TestIntegrationHelmDeploymentEvent(t *testing.T) {
	testData := setup()

	stop := startMain()

	data := &model.Request{
		RepoName:    "argo",
		RepoURL:     "https://argoproj.github.io/argo-helm",
		ChartName:   "argo-cd",
		Namespace:   "default",
		ReleaseName: "argocd",
		Timeout:     5,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Data marshalling failed, %v", err)
	}

	sendDeploymentEvent(t, "helm", dataJSON, "install")
	logger.Info("Sleeping now for 5 seconds")
	time.Sleep(5 * time.Second)

	logger.Info("Starting teardown")
	tearDown(testData)
	stop <- true
}

func TestIntegrationHelmDeleteEvent(t *testing.T) {
	testData := setup()

	stop := startMain()

	data := &model.Request{
		RepoName:    "argo",
		RepoURL:     "https://argoproj.github.io/argo-helm",
		ChartName:   "argo-cd",
		Namespace:   "default",
		ReleaseName: "argocd",
		Timeout:     5,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Data marshalling failed, %v", err)
	}

	sendDeploymentEvent(t, "helm", dataJSON, "delete")
	logger.Info("Sleeping now for 5 seconds")
	time.Sleep(5 * time.Second)

	logger.Info("Starting teardown")
	tearDown(testData)
	stop <- true
}

func sendDeploymentEvent(t *testing.T, pluginName string, dataJSON json.RawMessage, action string) {
	// The client is a heavyweight object that should be created once per process.
	temporalAddress := os.Getenv("TEMPORAL_SERVICE_URL")
	if len(temporalAddress) == 0 {
		temporalAddress = "127.0.0.1:7233"
	}
	c, err := client.Dial(client.Options{
		HostPort: temporalAddress,
		Logger:   logger,
	})
	if err != nil {
		t.Errorf("Unable to create client, %v", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "deployment_worker_workflow",
		TaskQueue: application.WorkflowTaskQueueName,
	}

	p := []model.RequestPayload{{PluginName: pluginName, Action: action, Data: dataJSON}}
	payload, _ := json.Marshal(p)
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.Workflow, payload)
	if err != nil {
		t.Errorf("Unable to execute workflow, %v", err)
	}

	logger.Infof("Started workflow, WorkflowID: %v RunID: %v", we.GetID(), we.GetRunID())

	// Synchronously wait for the workflow completion.
	var result model.ResponsePayload
	err = we.Get(context.Background(), &result)
	if err != nil {
		t.Errorf("Unable get workflow result, %v", err)
	}
	logger.Infof("Workflow result: %+v", result.ToString())
}
