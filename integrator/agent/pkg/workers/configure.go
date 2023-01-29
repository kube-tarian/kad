package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kube-tarian/kad/integrator/agent/pkg/model"
	"github.com/kube-tarian/kad/integrator/agent/pkg/temporalclient"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"go.temporal.io/sdk/client"
)

const (
	ConfigWorkerWorkflowName = "Workflow"
	ConfigWorkerTaskQueue    = "Configure"
)

type Config struct {
	client *temporalclient.Client
	log    logging.Logger
}

func NewConfig(client *temporalclient.Client, log logging.Logger) *Config {
	return &Config{
		client: client,
		log:    log,
	}
}

func (d *Config) GetWorkflowName() string {
	return ConfigWorkerWorkflowName
}

func (d *Config) SendEvent(ctx context.Context, deployPayload json.RawMessage) (client.WorkflowRun, error) {
	options := client.StartWorkflowOptions{
		ID:        uuid.NewString(),
		TaskQueue: ConfigWorkerTaskQueue,
	}

	log.Printf("Event sent to temporal: %v", string(deployPayload))
	run, err := d.client.ExecuteWorkflow(ctx, options, ConfigWorkerWorkflowName, deployPayload)
	if err != nil {
		return nil, err
	}

	d.log.Infof("Started workflow, ID: %v, WorkflowName: %v RunID: %v", run.GetID(), ConfigWorkerWorkflowName, run.GetRunID())

	// Wait for 5mins till workflow finishes
	// Timeout with 5mins
	// return run, d.getWorkflowStatusByLatestWorkflow(run)
	var result model.ResponsePayload
	err = run.Get(ctx, &result)
	if err != nil {
		d.log.Errorf("Result for workflow ID: %v, workflowName: %v, runID: %v", run.GetID(), ConfigWorkerWorkflowName, run.GetRunID())
		d.log.Errorf("Workflow result failed, %v", err)
		return run, err
	}
	d.log.Infof("workflow finished success, %+v", result.ToString())
	return run, nil
}

func (d *Config) getWorkflowStatusByLatestWorkflow(run client.WorkflowRun) error {
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

func (d *Config) getWorkflowInformation(run client.WorkflowRun) error {
	latestRun := d.client.TemporalClient.GetWorkflow(context.Background(), run.GetID(), "")

	var result model.ResponsePayload
	if err := latestRun.Get(context.Background(), &result); err != nil {
		d.log.Errorf("Unable to decode query result", err)
		return err
	}
	d.log.Debugf("Result info: %+v", result)
	return nil
}
