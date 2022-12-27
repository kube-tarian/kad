package workers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/kube-tarian/kad/agent/pkg/logging"
	"github.com/kube-tarian/kad/agent/pkg/model"
	"github.com/kube-tarian/kad/agent/pkg/temporalclient"
	"go.temporal.io/sdk/client"
)

const (
	DeploymentWorkerWorkflowName = "Workflow"
	DeploymentWorkerTaskQueue    = "Deployment"
)

type Deployment struct {
	client *temporalclient.Client
	log    logging.Logger
}

func NewDeployment(client *temporalclient.Client, log logging.Logger) *Deployment {
	return &Deployment{
		client: client,
		log:    log,
	}
}

func (d *Deployment) GetWorkflowName() string {
	return DeploymentWorkerWorkflowName
}

func (d *Deployment) SendEvent(ctx context.Context, deployPayload json.RawMessage) (client.WorkflowRun, error) {
	options := client.StartWorkflowOptions{
		ID:        uuid.NewString(),
		TaskQueue: DeploymentWorkerTaskQueue,
	}

	log.Printf("Event sent to temporal: %v", string(deployPayload))
	run, err := d.client.ExecuteWorkflow(ctx, options, DeploymentWorkerWorkflowName, deployPayload)
	if err != nil {
		return nil, err
	}

	d.log.Infof("Started workflow, ID: %v, WorkflowName: %v RunID: %v", run.GetID(), DeploymentWorkerWorkflowName, run.GetRunID())

	// Asynchronously wait for the workflow completion.
	// TODO: To be fixed context deadline/failed
	// go func() {
	var result model.ResponsePayload
	err = run.Get(ctx, &result)
	if err != nil {
		d.log.Errorf("Result for workflow ID: %v, workflowName: %v, runID: %v", run.GetID(), DeploymentWorkerWorkflowName, run.GetRunID())
		d.log.Errorf("Workflow result failed, %v", err)
		return run, err
	}
	d.log.Infof("Result for workflow ID: %v, workflowName: %v, runID: %v", run.GetID(), DeploymentWorkerWorkflowName, run.GetRunID())
	// }()

	return run, nil
}
