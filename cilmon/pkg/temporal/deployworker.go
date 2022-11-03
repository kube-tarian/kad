package temporal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"intelops.io/climon/pkg/pb/climonpb"
	"intelops.io/climon/pkg/plugins"
)

const (
	ClimonDeployTaskQueue = "CLIMON_DEPLOY_TASK_QUEUE"
)

type deployWorker struct {
	temporalClient *Client
}

func NewDeployWorker(address string) (Worker, error) {
	temporalClient, err := NewClient(address)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create Temporal client")
	}

	deployWorkerObj := deployWorker{temporalClient: temporalClient}
	temporalClient.CreateWorker(ClimonDeployTaskQueue)
	temporalClient.RegisterWorkflow(deployWorkerObj.DeployWorkflow)
	temporalClient.RegisterActivity(deployWorkerObj.DeployActivity)
	return &deployWorkerObj, nil
}

func (d *deployWorker) Start() error {
	return d.temporalClient.StartWorker()
}

func (d *deployWorker) Stop() error {
	d.temporalClient.Close()
	return nil
}

func (d *deployWorker) DeployWorkflow(ctx workflow.Context, request string) error {
	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    500,
	}

	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failures by default, this is just an example.
		RetryPolicy: retryPolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)
	err := workflow.ExecuteActivity(ctx, d.DeployActivity, request).Get(ctx, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to deploy")
	}

	return nil
}

func (d *deployWorker) DeployActivity(ctx context.Context, request string) error {
	var deployRequest climonpb.DeployRequest
	fmt.Println(request)
	if err := json.Unmarshal([]byte(request), &deployRequest); err != nil {
		log.Println("failed to unmarshal the deploy request", err)
		return errors.Wrapf(err, "failed to umarashal request")
	}

	plugin, err := plugins.GetPlugin(deployRequest.Plugin)
	if err != nil {
		log.Println("failed to get plugin", err)
		return errors.Wrapf(err, "failed to get plugin")
	}

	return plugin.Run(ctx, &deployRequest)
}
