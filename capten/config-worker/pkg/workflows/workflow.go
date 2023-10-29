package workflows

import (
	"encoding/json"
	"time"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/config-worker/pkg/activities"
	"github.com/kube-tarian/kad/capten/model"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// Workflow is a config workflow definition.
func Workflow(ctx workflow.Context, params model.ConfigureParameters, payload json.RawMessage) (model.ResponsePayload, error) {
	var result model.ResponsePayload
	logger := logging.NewLogger()

	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: 300 * time.Second,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)
	ctx = workflow.WithRetryPolicy(ctx, temporal.RetryPolicy{MaximumAttempts: 1})

	execution := workflow.GetInfo(ctx).WorkflowExecution
	logger.Infof("workflow execution information: %+v\n", execution)

	var a *activities.Activities
	err := workflow.ExecuteActivity(ctx, a.ConfigurationActivity, params, payload).Get(ctx, &result)
	if err != nil {
		logger.Errorf("Activity execution failed, Error: %v", err)
		return result, err
	}

	logger.Infof("successfully executed workflow., result: %s", (&result).ToString())

	return result, nil
}
