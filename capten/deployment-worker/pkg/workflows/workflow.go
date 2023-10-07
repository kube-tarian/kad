package workflows

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/deployment-worker/pkg/activities"
	"github.com/kube-tarian/kad/capten/model"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// Workflow is a deployment workflow definition.
func Workflow(ctx workflow.Context, action string, payload json.RawMessage) (model.ResponsePayload, error) {
	var result model.ResponsePayload
	logger := logging.NewLogger()

	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: 600 * time.Second,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)
	ctx = workflow.WithRetryPolicy(ctx, temporal.RetryPolicy{MaximumAttempts: 1})

	execution := workflow.GetInfo(ctx).WorkflowExecution
	logger.Infof("execution: %+v\n", execution)

	var a *activities.Activities
	var err error
	switch action {
	case "install", "update", "upgrade":
		req := &model.ApplicationDeployRequest{}
		err = json.Unmarshal(payload, req)
		if err == nil {
			err = workflow.ExecuteActivity(ctx, a.DeploymentInstallActivity, payload).Get(ctx, &result)
		}
	case "delete":
		req := &model.DeployerDeleteRequest{}
		err = json.Unmarshal(payload, req)
		if err == nil {
			err = workflow.ExecuteActivity(ctx, a.DeploymentDeleteActivity, payload).Get(ctx, &result)
		}
	default:
		err = fmt.Errorf("unknown action %v", action)
	}
	if err != nil {
		logger.Errorf("Activity failed: %v, Error: %v", string(payload), err)
		return result, err
	}

	logger.Infof("Deployment workflow completed., result: %s", (&result).ToString())
	return result, nil
}
