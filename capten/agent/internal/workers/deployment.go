package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/internal/temporalclient"
	"github.com/kube-tarian/kad/capten/model"
	"go.temporal.io/sdk/client"
)

const (
	DeploymentWorkerWorkflowName       = "Workflow"
	PluginDeploymentWorkerWorkflowName = "PluginWorkflow"
	DeploymentWorkerTaskQueue          = "Deployment"
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

func (d *Deployment) GetPluginWorkflowName() string {
	return PluginDeploymentWorkerWorkflowName
}

func (d *Deployment) SendEvent(
	ctx context.Context,
	workflowName string,
	action string,
	deployPayload *model.ApplicationInstallRequest,
) (client.WorkflowRun, error) {
	options := client.StartWorkflowOptions{
		ID:        uuid.NewString(),
		TaskQueue: DeploymentWorkerTaskQueue,
	}

	deployPayloadJSON, err := json.Marshal(deployPayload)
	if err != nil {
		return nil, err
	}

	d.log.Infof("Event sent to temporal: workflow: %s, action: %s, payload: %+v", action, deployPayload.String())
	var run client.WorkflowRun
	run, err = d.client.ExecuteWorkflow(ctx, options, workflowName, action, json.RawMessage(deployPayloadJSON))
	if err != nil {
		return nil, err
	}

	d.log.Infof("Started workflow, ID: %v, WorkflowName: %v RunID: %v", run.GetID(), DeploymentWorkerWorkflowName, run.GetRunID())

	// Wait for 5mins till workflow finishes
	// Timeout with 5mins
	// return run, d.getWorkflowStatusByLatestWorkflow(run)
	var result model.ResponsePayload
	err = run.Get(ctx, &result)
	if err != nil {
		d.log.Errorf("Result for workflow ID: %v, workflowName: %v, runID: %v", run.GetID(), DeploymentWorkerWorkflowName, run.GetRunID())
		d.log.Errorf("Workflow result failed, %v", err)
		return run, err
	}
	d.log.Infof("workflow finished success, %+v", result.ToString())
	return run, nil
}

func (d *Deployment) SendEventV2(
	ctx context.Context,
	workflowName string,
	action string,
	deployPayload model.DeployRequest,
) (client.WorkflowRun, error) {

	options := client.StartWorkflowOptions{
		ID:        uuid.NewString(),
		TaskQueue: DeploymentWorkerTaskQueue,
	}

	deployPayloadJSON, err := json.Marshal(deployPayload)
	if err != nil {
		return nil, err
	}

	d.log.Infof("Sending event to temporal: workflow: %s, action: %s, payload: %+v", workflowName, action, deployPayload.String())
	run, err := d.client.ExecuteWorkflow(ctx, options, workflowName, action, json.RawMessage(deployPayloadJSON))
	if err != nil {
		d.log.Errorf("failed to send event to workflow for plugin %s, %v", deployPayload.String(), err)
		return nil, err
	}

	d.log.Infof("Started workflow, ID: %v, WorkflowName: %v RunID: %v", run.GetID(), DeploymentWorkerWorkflowName, run.GetRunID())

	// Wait for 5mins till workflow finishes
	// Timeout with 5mins
	// return run, d.getWorkflowStatusByLatestWorkflow(run)
	var result model.ResponsePayload
	err = run.Get(ctx, &result)
	if err != nil {
		d.log.Errorf("Result for workflow ID: %v, workflowName: %v, runID: %v", run.GetID(), DeploymentWorkerWorkflowName, run.GetRunID())
		d.log.Errorf("Workflow result failed, %v", err)
		return run, err
	}
	d.log.Infof("workflow finished success, %+v", result.ToString())
	return run, nil
}

func (d *Deployment) SendDeleteEvent(
	ctx context.Context,
	workflowName string,
	action string,
	deployPayload any,
) (client.WorkflowRun, error) {
	options := client.StartWorkflowOptions{
		ID:        uuid.NewString(),
		TaskQueue: DeploymentWorkerTaskQueue,
	}

	payloadJSON, err := json.Marshal(deployPayload)
	if err != nil {
		return nil, err
	}

	d.log.Infof("Event sent to temporal: workflow: %s, action: %s, payload: %+v", workflowName, action, deployPayload)
	run, err := d.client.ExecuteWorkflow(ctx, options, workflowName, action, json.RawMessage(payloadJSON))
	if err != nil {
		return nil, err
	}

	d.log.Infof("Started workflow, ID: %v, WorkflowName: %v RunID: %v", run.GetID(), DeploymentWorkerWorkflowName, run.GetRunID())

	// Wait for 5mins till workflow finishes
	// Timeout with 5mins
	// return run, d.getWorkflowStatusByLatestWorkflow(run)
	var result model.ResponsePayload
	err = run.Get(ctx, &result)
	if err != nil {
		d.log.Errorf("Result for workflow ID: %v, workflowName: %v, runID: %v", run.GetID(), DeploymentWorkerWorkflowName, run.GetRunID())
		d.log.Errorf("Workflow result failed, %v", err)
		return run, err
	}
	d.log.Infof("workflow finished success, %+v", result.ToString())
	return run, nil
}

func (d *Deployment) getWorkflowStatusByLatestWorkflow(run client.WorkflowRun) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			err := d.getWorkflowInformation(run)
			if err != nil {
				d.log.Errorf("get state of workflow failed: %v, retrying .....", err)
				continue
			}
			return nil
		case <-time.After(5 * time.Minute):
			d.log.Errorf("Timed out waiting for state of workflow")
			return fmt.Errorf("timedout waiting for the workflow to finish")
		}
	}
}

func (d *Deployment) getWorkflowInformation(run client.WorkflowRun) error {
	latestRun := d.client.TemporalClient.GetWorkflow(context.Background(), run.GetID(), "")

	var result model.ResponsePayload
	if err := latestRun.Get(context.Background(), &result); err != nil {
		d.log.Errorf("Unable to decode query result", err)
		return err
	}
	d.log.Debugf("Result info: %+v", result)
	return nil
}
