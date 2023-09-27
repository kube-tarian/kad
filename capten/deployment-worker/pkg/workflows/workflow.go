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

	fmt.Println("Action => ", action)
	fmt.Println("Incomming Request payload => ", string(payload))
	reqX := &model.ApplicationDeployRequest{}
	_ = json.Unmarshal(payload, reqX)
	x, _ := json.Marshal(reqX)
	fmt.Println("Outgoing Request payload => ", string(x))

	var a *activities.Activities
	var err error
	switch action {
	case "install", "update":
		req := &model.ApplicationDeployRequest{}
		err = json.Unmarshal(payload, req)
		if err == nil {
			err = workflow.ExecuteActivity(ctx, a.DeploymentInstallActivity, payload).Get(ctx, &result)
		}
	case "delete":
		req := &model.DeployerDeleteRequest{}
		err = json.Unmarshal(payload, req)
		if err == nil {
			fmt.Println("Incomming Request payload => ", string(payload))
			x, _ := json.Marshal(req)
			fmt.Println("Outgoing Request payload => ", string(x))
			err = workflow.ExecuteActivity(ctx, a.DeploymentDeleteActivity, payload).Get(ctx, &result)
		}
		fmt.Println("Error unmarshalling => ", err.Error())
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
