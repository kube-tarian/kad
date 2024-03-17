package workflows

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"time"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/cluster-plugins/clusterpluginspb"
	"github.com/kube-tarian/kad/capten/deployment-worker/internal/activities"
	"github.com/kube-tarian/kad/capten/model"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"gopkg.in/yaml.v2"
)

func PluginWorkflow(ctx workflow.Context, action string, payload json.RawMessage) (model.ResponsePayload, error) {
	result := &model.ResponsePayload{}
	logger := logging.NewLogger()

	// var a *activities.Activities
	var a *activities.PluginActivities
	var err error
	switch action {
	case string(model.AppInstallAction), string(model.AppUpdateAction), string(model.AppUpgradeAction):
		result, err = hanldeDeployWorkflow(ctx, payload, logger)
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

func hanldeDeployWorkflow(ctx workflow.Context, payload json.RawMessage, log logging.Logger) (*model.ResponsePayload, error) {
	var a *activities.PluginActivities
	result := &model.ResponsePayload{}
	ctx = setContext(ctx, 600, log)

	req := &clusterpluginspb.Plugin{}
	err := json.Unmarshal(payload, req)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"wrong content: %s\"}", err.Error())),
		}, err
	}

	for _, capability := range req.Capabilities {
		switch capability {
		case "capten-sdk":
			err = workflow.ExecuteActivity(ctx, a.PluginDeployPreActionMTLSActivity, req).Get(ctx, result)
			if err != nil {
				return result, err
			}

		case "vault-store":
			err = workflow.ExecuteActivity(ctx, a.PluginDeployPreActionVaultStoreActivity, req).Get(ctx, result)
			if err != nil {
				return result, err
			}

		case "postgres-store":
			err = workflow.ExecuteActivity(ctx, a.PluginDeployPreActionPostgresStoreActivity, req).Get(ctx, result)
			if err != nil {
				return result, err
			}

		default:
			log.Infof("Unsupported capability %s", capability)
		}
	}

	result, err = executeAppDeployment(ctx, req, a)
	if err != nil {
		return result, err
	}
	err = workflow.ExecuteActivity(ctx, a.PluginDeployPostActionActivity, req).Get(ctx, result)
	if err != nil {
		return result, err
	}

	return result, err
}

func executeAppDeployment(ctx workflow.Context, req *clusterpluginspb.Plugin, a *activities.PluginActivities) (result *model.ResponsePayload, err error) {
	err = workflow.ExecuteActivity(ctx, a.PluginDeployUpdateStatusActivity, req.PluginName, "plugin-app-installing").Get(ctx, result)
	if err != nil {
		return result, err
	}

	cwo := workflow.ChildWorkflowOptions{
		WorkflowID: "APP-DEPLOY-CHILD-WORKFLOW-ID",
	}
	ctx = workflow.WithChildOptions(ctx, cwo)

	templateValues, err := deriveTemplateValues(req.OverrideValues, req.Values)
	if err != nil {
		err = fmt.Errorf("failed to derive template values for app %s, %v", req.PluginName, err)
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \" failed to prepare app values, %s\"}", err.Error())),
		}, nil
	}

	appDeployReq := prepareAppDeployRequestFromPlugin(req, templateValues)
	result = &model.ResponsePayload{}
	err = workflow.ExecuteChildWorkflow(ctx, Workflow, appDeployReq).Get(ctx, &result)
	if err != nil {
		err = workflow.ExecuteActivity(ctx, a.PluginDeployUpdateStatusActivity, req.PluginName, "plugin-app-installfailed").Get(ctx, result)
		if err != nil {
			return result, err
		}
		return result, err
	}

	err = workflow.ExecuteActivity(ctx, a.PluginDeployUpdateStatusActivity, req.PluginName, "plugin-app-installed").Get(ctx, result)
	if err != nil {
		return result, err
	}

	return
}

func deriveTemplateValues(overrideValues, templateValues []byte) ([]byte, error) {
	overrideValuesMapping := map[string]any{}
	if err := yaml.Unmarshal(overrideValues, &overrideValuesMapping); err != nil {
		return nil, errors.WithMessagef(err, "failed to Unmarshal override values")
	}

	templateValues, err := executeTemplateValuesTemplate(templateValues, overrideValuesMapping)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to exeute template values update")
	}
	return templateValues, nil
}

func executeTemplateValuesTemplate(data []byte, values map[string]any) (transformedData []byte, err error) {
	if len(data) == 0 {
		return
	}

	tmpl, err := template.New("templateVal").Parse(string(data))
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, values)
	if err != nil {
		return
	}

	transformedData = buf.Bytes()
	return
}

func prepareAppDeployRequestFromPlugin(data *clusterpluginspb.Plugin, values []byte) *model.ApplicationInstallRequest {
	return &model.ApplicationInstallRequest{
		PluginName:     "helm",
		RepoName:       data.ChartName,
		RepoURL:        data.ChartRepo,
		ChartName:      data.ChartName,
		Namespace:      data.DefaultNamespace,
		ReleaseName:    data.ChartName,
		Version:        data.Version,
		ClusterName:    "capten",
		OverrideValues: string(values),
		Timeout:        10,
	}
}
