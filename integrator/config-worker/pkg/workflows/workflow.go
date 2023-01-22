package workflows

import (
	"encoding/json"
	"time"

	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/kube-tarian/kad/integrator/config-worker/pkg/activities"
	"github.com/kube-tarian/kad/integrator/model"
	"go.temporal.io/sdk/workflow"
)

// Workflow is a config workflow definition.
func Workflow(ctx workflow.Context, payload json.RawMessage) (model.ResponsePayload, error) {
	var result model.ResponsePayload
	logger := logging.NewLogger()

	logger.Infof("Configuration workflow started, req: %+v", string(payload))
	req := []model.ConfigPayload{}
	err := json.Unmarshal(payload, &req)
	if err != nil {
		logger.Errorf("Deployer worker payload unmarshall failed, Error: %v", err)
		return result, err
	}

	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: 60 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	execution := workflow.GetInfo(ctx).WorkflowExecution
	logger.Infof("execution: %+v\n", execution)

	var a *activities.Activities
	err = workflow.ExecuteActivity(ctx, a.ConfigurationActivity, req[0]).Get(ctx, &result)
	if err != nil {
		logger.Errorf("Activity failed, Error: %v", err)
		return result, err
	}

	logger.Infof("Configuration workflow completed., result: %s", (&result).ToString())

	return result, nil
}
