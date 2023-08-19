package workflows

import (
	"time"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/config-worker/pkg/activities"
	"github.com/kube-tarian/kad/capten/model"
	"go.temporal.io/sdk/internal"
	"go.temporal.io/sdk/workflow"
)

// Workflow is a config workflow definition.
func Workflow(ctx workflow.Context, params model.ConfigureParameters, string, payload interface{}) (model.ResponsePayload, error) {
	var result model.ResponsePayload
	logger := logging.NewLogger()

	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: 60 * time.Second,
		RetryPolicy:            &internal.RetryPolicy{MaximumAttempts: 1},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	execution := workflow.GetInfo(ctx).WorkflowExecution
	logger.Infof("execution: %+v\n", execution)

	var a *activities.Activities
	err := workflow.ExecuteActivity(ctx, a.ConfigurationActivity, params, payload).Get(ctx, &result)
	if err != nil {
		logger.Errorf("Activity failed, Error: %v", err)
		return result, err
	}

	logger.Infof("Configuration workflow completed., result: %s", (&result).ToString())

	return result, nil
}
