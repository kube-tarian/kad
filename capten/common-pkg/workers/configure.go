package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/temporalclient"
	"github.com/kube-tarian/kad/capten/model"
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

func (d *Config) SendEvent(ctx context.Context, confParams *model.ConfigureParameters, deployPayload interface{}) (client.WorkflowRun, error) {
	options := client.StartWorkflowOptions{
		ID:        uuid.NewString(),
		TaskQueue: ConfigWorkerTaskQueue,
	}

	deployPayloadJson, err := json.Marshal(deployPayload)
	if err != nil {
		return nil, err
	}

	d.log.Debugf("Event sent to temporal: %+v", deployPayload)
	run, err := d.client.TemporalClient.ExecuteWorkflow(ctx, options, ConfigWorkerWorkflowName, confParams, json.RawMessage(deployPayloadJson))
	if err != nil {
		return nil, err
	}

	d.log.Infof("Started workflow, ID: %v, WorkflowName: %v RunID: %v", run.GetID(), ConfigWorkerWorkflowName, run.GetRunID())

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

func (d *Config) SendAsyncEvent(ctx context.Context, confParams *model.ConfigureParameters, deployPayload interface{}) (string, error) {
	options := client.StartWorkflowOptions{
		ID:        uuid.NewString(),
		TaskQueue: ConfigWorkerTaskQueue,
	}

	deployPayloadJson, err := json.Marshal(deployPayload)
	if err != nil {
		return "", err
	}

	d.log.Debugf("Event sent to temporal: %+v", deployPayload)
	run, err := d.client.TemporalClient.ExecuteWorkflow(ctx, options, ConfigWorkerWorkflowName, confParams, json.RawMessage(deployPayloadJson))
	if err != nil {
		return "", err
	}

	d.log.Infof("Started Async workflow, ID: %v, WorkflowName: %v RunID: %v", run.GetID(), ConfigWorkerWorkflowName, run.GetRunID())

	return run.GetID(), nil
}

func (d *Config) getWorkflowStatusByLatestWorkflow(ctx context.Context, run client.WorkflowRun) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			_, err := d.GetWorkflowInformation(ctx, run.GetID())
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

func (d *Config) GetWorkflowInformation(ctx context.Context, workFlowId string) (model.ResponsePayload, error) {
	d.log.Debugf("Fetching workflow Id: %s, status...", workFlowId)

	latestRun := d.client.TemporalClient.GetWorkflow(ctx, workFlowId, "")

	var result model.ResponsePayload
	if err := latestRun.Get(ctx, &result); err != nil {
		d.log.Errorf("failed to get the workflow Id: %s, status: ", workFlowId, err)
		return result, err
	}

	d.log.Debugf("Result workflow Id: %s, status: %v", workFlowId, result)

	return result, nil
}
