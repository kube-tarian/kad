package workflows

import (
	"time"

	"github.com/kube-tarian/kad/deployment-worker/pkg/activities"
	"go.temporal.io/sdk/workflow"
)

// Workflow is a deployment workflow definition.
func Workflow(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: 60 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("Deployment workflow started", "name", name)

	var result string
	var a *activities.Activities
	err := workflow.ExecuteActivity(ctx, a.Activity, name).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info("Deployment workflow completed.", "result", result)

	return result, nil
}
