package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/kube-tarian/kad/integrator/agent/pkg/model"
	"github.com/kube-tarian/kad/integrator/agent/pkg/temporalclient"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
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

func (c *Climon) SendEvent(ctx context.Context, deployPayload json.RawMessage) (client.WorkflowRun, error) {
	options := client.StartWorkflowOptions{
		ID:        "helm-deploy-workflow",
		TaskQueue: ClimonHelmTaskQueue,
	}

	/*
		deployInfo := helm.DeployInfo{
			Version:     "1.0",
			RepoUrl:     "https://charts.bitnami.com/bitnami",
			RepoName:    "bitnami",
			Namespace:   "web",
			ChartName:   "bitnami/wordpress",
			ReleaseName: "intelops",
			ReferenceID: uuid.New().String(),
		}
	*/

	we, err := c.client.ExecuteWorkflow(context.Background(), options, DeployWorkflowName, deployPayload)
	if err != nil {
		log.Println("error starting TransferMoney workflow", err)
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

func (d *Climon) getWorkflowStatusByLatestWorkflow(run client.WorkflowRun) error {
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

func (d *Climon) getWorkflowInformation(run client.WorkflowRun) error {
	latestRun := d.client.TemporalClient.GetWorkflow(context.Background(), run.GetID(), "")

	var result model.ResponsePayload
	if err := latestRun.Get(context.Background(), &result); err != nil {
		d.log.Errorf("Unable to decode query result", err)
		return err
	}
	d.log.Debugf("Result info: %+v", result)
	return nil
}
