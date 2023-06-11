package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/agent/pkg/temporalclient"
	"go.temporal.io/sdk/client"
)

type Climon struct {
	client *temporalclient.Client
	log    logging.Logger
}

func NewClimon(client *temporalclient.Client, log logging.Logger) *Climon {
	return &Climon{
		client: client,
		log:    log,
	}
}

func (c *Climon) GetWorkflowName() string {
	return DeployWorkflowName
}

func (c *Climon) SendEvent(ctx context.Context, action string, payload *agentpb.ClimonInstallRequest) (client.WorkflowRun, error) {
	options := client.StartWorkflowOptions{
		ID:        "helm-deploy-workflow",
		TaskQueue: ClimonHelmTaskQueue,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	log.Printf("payload climon: %v", string(payloadJSON))
	we, err := c.client.TemporalClient.ExecuteWorkflow(context.Background(), options, DeployWorkflowName, action, json.RawMessage(payloadJSON))
	if err != nil {
		log.Println("error starting climon workflow", err)
		return nil, err
	}
	//printResults(deployInfo, we.GetID(), we.GetRunID())

	c.log.Infof("Started workflow, ID: %v, WorkflowName: %v RunID: %v", we.GetID(), DeploymentWorkerWorkflowName, we.GetRunID())

	// Wait for 5mins till workflow finishes
	// Timeout with 5mins
	var result model.ResponsePayload
	err = we.Get(ctx, &result)
	if err != nil {
		c.log.Errorf("Result for workflow ID: %v, workflowName: %v, runID: %v", we.GetID(), DeploymentWorkerWorkflowName, we.GetRunID())
		c.log.Errorf("Workflow result failed, %v", err)
		return we, err
	}
	c.log.Infof("workflow finished success, %+v", result.ToString())

	return we, nil
}

func (c *Climon) SendDeleteEvent(ctx context.Context, action string, payload *agentpb.ClimonDeleteRequest) (client.WorkflowRun, error) {
	options := client.StartWorkflowOptions{
		ID:        "helm-deploy-workflow",
		TaskQueue: ClimonHelmTaskQueue,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	we, err := c.client.TemporalClient.ExecuteWorkflow(context.Background(), options, DeployWorkflowName, action, payloadJSON)
	if err != nil {
		log.Println("error starting climon workflow", err)
		return nil, err
	}

	c.log.Infof("Started workflow, ID: %v, WorkflowName: %v RunID: %v", we.GetID(), DeploymentWorkerWorkflowName, we.GetRunID())

	// Wait for 5mins till workflow finishes
	// Timeout with 5mins
	var result model.ResponsePayload
	err = we.Get(ctx, &result)
	if err != nil {
		c.log.Errorf("Result for workflow ID: %v, workflowName: %v, runID: %v", we.GetID(), DeploymentWorkerWorkflowName, we.GetRunID())
		c.log.Errorf("Workflow result failed, %v", err)
		return we, err
	}
	c.log.Infof("workflow finished success, %+v", result.ToString())

	return we, nil
}

func (c *Climon) getWorkflowStatusByLatestWorkflow(run client.WorkflowRun) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			err := c.getWorkflowInformation(run)
			if err != nil {
				c.log.Errorf("get state of workflow failed: %v, retrying .....", err)
				continue
			}
			return nil
		case <-time.After(5 * time.Minute):
			c.log.Errorf("Timed out waiting for state of workflow")
			return fmt.Errorf("timedout waiting for the workflow to finish")
		}
	}
}

func (c *Climon) getWorkflowInformation(run client.WorkflowRun) error {
	latestRun := c.client.TemporalClient.GetWorkflow(context.Background(), run.GetID(), "")

	var result model.ResponsePayload
	if err := latestRun.Get(context.Background(), &result); err != nil {
		c.log.Errorf("Unable to decode query result", err)
		return err
	}
	c.log.Debugf("Result info: %+v", result)
	return nil
}
