package workflows

import (
	"log"
	"time"

	"github.com/kube-tarian/kad/integrator/deployment-worker/pkg/activities"
	"github.com/kube-tarian/kad/integrator/deployment-worker/pkg/model"
	"go.temporal.io/sdk/workflow"
)

// Workflow is a deployment workflow definition.
func Workflow(ctx workflow.Context, req model.RequestPayload) (model.ResponsePayload, error) {
	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: 60 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	execution := workflow.GetInfo(ctx).WorkflowExecution
	logger.Info("Deployment workflow started", "name", req)
	log.Printf("execution: %+v\n", execution)

	var result model.ResponsePayload
	var a *activities.Activities
	err := workflow.ExecuteActivity(ctx, a.DeploymentActivity, req).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return result, err
	}

	logger.Info("Deployment workflow completed.", "result", result)

	return result, nil
}
