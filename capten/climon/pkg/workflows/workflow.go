package workflows

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kube-tarian/kad/capten/climon/pkg/activities"
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/model"

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
	execution := workflow.GetInfo(ctx).WorkflowExecution
	logger.Infof("execution: %+v\n", execution)

	var a *activities.Activities
	var err error
	switch action {
	case "install", "update":
		req := &model.ClimonPostRequest{}
		err = json.Unmarshal(payload, req)
		if err == nil {
			err = workflow.ExecuteActivity(ctx, a.ClimonInstallActivity, payload).Get(ctx, &result)
		}
	case "delete":
		req := &model.ClimonDeleteRequest{}
		err = json.Unmarshal(payload, req)
		if err == nil {
			err = workflow.ExecuteActivity(ctx, a.ClimonDeleteActivity, payload).Get(ctx, &result)
		}
	default:
		err = fmt.Errorf("unknown action %v", action)
	}
	if err != nil {
		logger.Errorf("Activity failed, Error: %v", err)
		return result, err
	}

	logger.Infof("Deployment workflow completed., result: %s", (&result).ToString())
	return result, nil
}
