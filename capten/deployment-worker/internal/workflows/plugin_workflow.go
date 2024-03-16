package workflows

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/deployment-worker/internal/activities"
	"github.com/kube-tarian/kad/capten/model"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func PluginWorkflow(ctx workflow.Context, action string, payload json.RawMessage, capabilities []string) (model.ResponsePayload, error) {
	result := &model.ResponsePayload{}
	logger := logging.NewLogger()

	// var a *activities.Activities
	var a *activities.PluginActivities
	var err error
	switch action {
	case string(model.AppInstallAction), string(model.AppUpdateAction), string(model.AppUpgradeAction):
		result, err = hanldeDeployWorkflow(ctx, payload, capabilities, logger)
	case string(model.AppUnInstallAction):
		ctx = setContext(ctx, 600, logger)
		req := &model.DeployerDeleteRequest{}
		err = json.Unmarshal(payload, req)
		if err == nil {
			err = workflow.ExecuteActivity(ctx, a.PluginUndeployActivity, payload).Get(ctx, result)
		}
	default:
		err = fmt.Errorf("unknown action %v", action)
	}
	if err != nil {
		logger.Errorf("Deployment workflow failed: %v, Error: %v", string(payload), err)
		return *result, err
	}

	logger.Infof("Deployment workflow completed., result: %s", (result).ToString())
	return *result, nil
}

func setContext(ctx workflow.Context, timeInSeconds int, log logging.Logger) workflow.Context {
	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Duration(timeInSeconds) * time.Second,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)
	ctx = workflow.WithRetryPolicy(ctx, temporal.RetryPolicy{MaximumAttempts: 1})

	execution := workflow.GetInfo(ctx).WorkflowExecution
	log.Infof("execution: %+v\n", execution)
	return ctx
}

func hanldeDeployWorkflow(ctx workflow.Context, payload json.RawMessage, capabilities []string, log logging.Logger) (*model.ResponsePayload, error) {
	var a *activities.PluginActivities
	result := &model.ResponsePayload{}
	ctx = setContext(ctx, 600, log)

	req := &model.ApplicationDeployRequest{}
	err := json.Unmarshal(payload, req)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"wrong content: %s\"}", err.Error())),
		}, err
	}

	for _, capability := range capabilities {
		switch capability {
		case "capten-sdk":
			err = workflow.ExecuteActivity(ctx, a.PluginDeployPreActionMTLSActivity, payload).Get(ctx, result)
			if err != nil {
				return result, err
			}

		case "vault-store":
			err = workflow.ExecuteActivity(ctx, a.PluginDeployPreActionVaultStoreActivity, payload).Get(ctx, result)
			if err != nil {
				return result, err
			}

		case "postgres-store":
			err = workflow.ExecuteActivity(ctx, a.PluginDeployPreActionPostgresStoreActivity, payload).Get(ctx, result)
			if err != nil {
				return result, err
			}

		default:
			log.Infof("Unsupported capability %s", capability)
		}
	}

	err = workflow.ExecuteActivity(ctx, a.PluginDeployActivity, payload).Get(ctx, result)
	if err != nil {
		return result, err
	}

	err = workflow.ExecuteActivity(ctx, a.PluginDeployPostActionActivity, payload).Get(ctx, result)
	if err != nil {
		return result, err
	}

	return result, err
}
