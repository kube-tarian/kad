package workflows

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/cluster-plugins/clusterpluginspb"
	"github.com/kube-tarian/kad/capten/deployment-worker/internal/activities"
	"github.com/kube-tarian/kad/capten/model"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func PluginWorkflow(ctx workflow.Context, action string, payload json.RawMessage) (model.ResponsePayload, error) {
	var result model.ResponsePayload
	logger := logging.NewLogger()

	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: 600 * time.Second,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)
	ctx = workflow.WithRetryPolicy(ctx, temporal.RetryPolicy{MaximumAttempts: 1})

	execution := workflow.GetInfo(ctx).WorkflowExecution
	logger.Infof("execution: %+v\n", execution)

	// var a *activities.Activities
	var a *activities.PluginActivities
	var err error
	switch action {
	case string(model.AppInstallAction), string(model.AppUpdateAction), string(model.AppUpgradeAction):
		req := &clusterpluginspb.DeployClusterPluginRequest{}
		err = json.Unmarshal(payload, req)
		if err == nil {
			err = workflow.ExecuteActivity(ctx, a.PluginDeployActivity, payload).Get(ctx, &result)
		}
	case string(model.AppUnInstallAction):
		req := &model.DeployerDeleteRequest{}
		err = json.Unmarshal(payload, req)
		if err == nil {
			err = workflow.ExecuteActivity(ctx, a.PluginUndeployActivity, payload).Get(ctx, &result)
		}
	default:
		err = fmt.Errorf("unknown action %v", action)
	}
	if err != nil {
		logger.Errorf("Deployment workflow failed: %v, Error: %v", string(payload), err)
		return result, err
	}

	logger.Infof("Deployment workflow completed., result: %s", (&result).ToString())
	return result, nil
}
