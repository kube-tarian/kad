package temporal

import (
	"context"
	"fmt"
	"time"

	"github.com/kube-tarian/kad/capten/climon/pkg/db/cassandra"
	"github.com/kube-tarian/kad/capten/climon/pkg/types"

	"github.com/pkg/errors"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	SyncTaskQueue = "SYNC_TASK_QUEUE"
)

type syncWorker struct {
	temporalClient *Client
	db             cassandra.Store
}

type SyncDataRequest struct {
	Type string      `json:"type"`
	Apps []types.App `json:"apps"`
}

func NewSyncWorker(address string) (Worker, error) {
	temporalClient, err := NewClient(address)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create Temporal client")
	}

	db := cassandra.GetStore()
	syncWorkerObj := syncWorker{temporalClient: temporalClient, db: db}
	temporalClient.CreateWorker(SyncTaskQueue)
	temporalClient.RegisterWorkflow(syncWorkerObj.SyncWorkflow)
	temporalClient.RegisterActivity(syncWorkerObj.SyncActivity)
	return &syncWorkerObj, nil
}

func (s *syncWorker) Start() error {
	return s.temporalClient.StartWorker()
}

func (s *syncWorker) Stop() error {
	s.temporalClient.Close()
	return nil
}

func (s *syncWorker) SyncWorkflow(ctx workflow.Context, request SyncDataRequest) error {
	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    5,
	}

	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failures by default, this is just an example.
		RetryPolicy: retryPolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)
	err := workflow.ExecuteActivity(ctx, s.SyncActivity, request).Get(ctx, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to deploy")
	}

	return nil
}

func (s *syncWorker) SyncActivity(ctx context.Context, request SyncDataRequest) error {
	fmt.Printf("%+v", request)
	return s.db.InsertApps(request.Apps)
}
