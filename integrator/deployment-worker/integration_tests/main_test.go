package integrationtests

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/kube-tarian/kad/integrator/deployment-worker/pkg/model"
	"github.com/kube-tarian/kad/integrator/deployment-worker/pkg/workflows"
	"go.temporal.io/sdk/client"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestIntegrationDeploymentEvent(t *testing.T) {
	data := setup()

	stop := startMain()

	sendDeploymentEvent(t)
	log.Println("Sleeping now for 5 seconds")
	time.Sleep(5 * time.Second)

	log.Println("Starting teardown")
	tearDown(data)
	stop <- true
}

func sendDeploymentEvent(t *testing.T) {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		t.Errorf("Unable to create client, %v", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "deployment_workflowID",
		TaskQueue: "Deployment",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.Workflow, model.RequestPayload{SubAction: "Temporal"})
	if err != nil {
		t.Errorf("Unable to execute workflow, %v", err)
	}

	t.Logf("Started workflow, WorkflowID: %v RunID: %v", we.GetID(), we.GetRunID())

	// Synchronously wait for the workflow completion.
	var result model.ResponsePayload
	err = we.Get(context.Background(), &result)
	if err != nil {
		t.Errorf("Unable get workflow result, %v", err)
	}
	t.Logf("Workflow result: %+v\n", result)
}
