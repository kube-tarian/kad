package workflows

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"time"

	"github.com/google/uuid"
	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/common-pkg/capten-store"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/agentpb"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/clusterpluginspb"
	"github.com/kube-tarian/kad/capten/deployment-worker/internal/activities"
	"github.com/kube-tarian/kad/capten/model"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"gopkg.in/yaml.v2"
)

type WorkflowHeader struct {
	Action string `json:"action"`
}

func PluginWorkflow(ctx workflow.Context, payload json.RawMessage) (model.ResponsePayload, error) {
	logger := logging.NewLogger()
	logger.Infof("plugin deployment workflow started. header payload: %s", string(payload))

	workflowData := []interface{}{}
	err := json.Unmarshal(payload, &workflowData)
	if err != nil {
		return model.ResponsePayload{
			Status:  "error",
			Message: []byte(err.Error()),
		}, err
	}
	logger.Infof("plugin deployment workflow started. header payload: %+v", workflowData)

	req := &WorkflowHeader{}
	err = json.Unmarshal(workflowData[0].(json.RawMessage), req)
	if err != nil {
		return model.ResponsePayload{
			Status:  "error",
			Message: []byte(err.Error()),
		}, err
	}

	action := req.Action
	result := &model.ResponsePayload{}

	as, err := captenstore.NewStore(logger)
	if err != nil {
		logger.Errorf("failed to initialize plugin app store, %v", err)
	}

	switch action {
	case string(model.AppInstallAction), string(model.AppUpdateAction), string(model.AppUpgradeAction):
		result, err = hanldeDeployWorkflow(ctx, payload, logger, as)
	case string(model.AppUnInstallAction):
		result, err = hanldeUndeployWorkflow(ctx, payload, logger, as)
	default:
		err = fmt.Errorf("unknown action %v", action)
	}
	if err != nil {
		logger.Errorf("plugin deployment action %s workflow failed, Error: %v", action, err)
		return *result, err
	}

	logger.Infof("plugin deployment action %s workflow completed., result: %s", action, (result).ToString())
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

func hanldeDeployWorkflow(ctx workflow.Context, payload json.RawMessage, log logging.Logger, pas *captenstore.Store) (*model.ResponsePayload, error) {
	log.Info("Starting deploy plugin workflow")
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

	log.Infof("Started plugin workflow for %s with capabilities: %v", req.PluginName, req.Capabilities)
	for _, capability := range req.Capabilities {
		switch capability {
		case "capten-sdk":
			err = workflow.ExecuteActivity(ctx, a.PluginDeployPreActionMTLSActivity, req).Get(ctx, result)
			if err != nil {
				log.Errorf("pre-installation capten-sdk failed, %s, reason: %v", req.PluginName, err)
				return result, err
			}

		case "vault-store":
			err = workflow.ExecuteActivity(ctx, a.PluginDeployPreActionVaultStoreActivity, req).Get(ctx, result)
			if err != nil {
				log.Errorf("pre-installation vault-store failed, %s, reason: %v", req.PluginName, err)
				return result, err
			}

		case "postgres-store":
			err = workflow.ExecuteActivity(ctx, a.PluginDeployPreActionPostgresStoreActivity, req).Get(ctx, result)
			if err != nil {
				log.Errorf("pre-installation postgres failed, %s, reason: %v", req.PluginName, err)
				return result, err
			}

		default:
			log.Infof("Unsupported capability %s for plugin %s", capability, req.PluginName)
		}
	}

	result, err = executeAppDeployment(ctx, req, a, log, pas)
	if err != nil {
		log.Errorf("App installation failed, %s, reason: %v", req.PluginName, err)
		return result, err
	}
	err = workflow.ExecuteActivity(ctx, a.PluginDeployPostActionActivity, req).Get(ctx, result)
	if err != nil {
		log.Errorf("post-action failed, %s, reason: %v", req.PluginName, err)
		return result, err
	}
	log.Infof("Finsihed plugin workflow for %s", req.PluginName)

	return result, err
}

func hanldeUndeployWorkflow(ctx workflow.Context, payload json.RawMessage, log logging.Logger, pas *captenstore.Store) (*model.ResponsePayload, error) {
	var a *activities.PluginActivities
	result := &model.ResponsePayload{}
	ctx = setContext(ctx, 600, log)

	log.Infof("Payload: %+v", payload)
	req := &clusterpluginspb.UnDeployClusterPluginRequest{}
	err := json.Unmarshal(payload, req)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"wrong content: %s\"}", err.Error())),
		}, err
	}
	log.Infof("undeployclusterplugin request: %+v", req)

	pluginConfig, err := pas.GetClusterPluginConfig(req.PluginName)
	if err != nil {
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \"Failed to fetch plugin configuration: %s\"}", err.Error())),
		}, err
	}

	result, err = executeAppUndeployment(ctx, pluginConfig, a, log, pas)
	if err != nil {
		return result, err
	}

	for _, capability := range pluginConfig.Capabilities {
		switch capability {
		case "capten-sdk":
			err = workflow.ExecuteActivity(ctx, a.PluginUndeployPreActionMTLSActivity, pluginConfig).Get(ctx, result)
			if err != nil {
				return result, err
			}

		case "vault-store":
			err = workflow.ExecuteActivity(ctx, a.PluginUndeployPreActionVaultStoreActivity, pluginConfig).Get(ctx, result)
			if err != nil {
				return result, err
			}

		case "postgres-store":
			err = workflow.ExecuteActivity(ctx, a.PluginUndeployPreActionPostgresStoreActivity, pluginConfig).Get(ctx, result)
			if err != nil {
				return result, err
			}

		default:
			log.Infof("Unsupported capability %s", capability)
		}
	}

	err = workflow.ExecuteActivity(ctx, a.PluginUndeployPostActionActivity, pluginConfig).Get(ctx, result)
	if err != nil {
		return result, err
	}

	return result, err
}

func executeAppDeployment(
	ctx workflow.Context,
	req *clusterpluginspb.Plugin,
	a *activities.PluginActivities,
	log logging.Logger,
	pas *captenstore.Store,
) (result *model.ResponsePayload, err error) {
	result = &model.ResponsePayload{}
	err = workflow.ExecuteActivity(ctx, a.PluginDeployUpdateStatusActivity, req.PluginName, "plugin-app-installing").Get(ctx, result)
	if err != nil {
		return result, err
	}

	cwo := workflow.ChildWorkflowOptions{
		WorkflowID: uuid.NewString(),
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

	uiEndpoint, err := executeStringTemplateValues(req.UiEndpoint, req.OverrideValues)
	if err != nil {
		log.Errorf("failed to derive template launch URL for app %s, %v", req.PluginName, err)
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \" failed to prepare app values, %s\"}", err.Error())),
		}, nil
	}

	apiEndpoint, err := executeStringTemplateValues(req.ApiEndpoint, req.OverrideValues)
	if err != nil {
		log.Errorf("failed to derive template launch URL for app %s, %v", req.PluginName, err)
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \" failed to prepare app values, %s\"}", err.Error())),
		}, nil
	}

	syncConfig := &agentpb.SyncAppData{
		Config: &agentpb.AppConfig{
			ReleaseName:         req.ChartName,
			AppName:             req.ChartName,
			Version:             req.Version,
			Category:            req.Category,
			Description:         req.Description,
			ChartName:           req.ChartName + "/" + req.ChartName,
			RepoName:            req.ChartName,
			RepoURL:             req.ChartRepo,
			Namespace:           req.DefaultNamespace,
			CreateNamespace:     true,
			PrivilegedNamespace: req.PrivilegedNamespace,
			Icon:                req.Icon,
			UiEndpoint:          uiEndpoint,
			UiModuleEndpoint:    req.UiModuleEndpoint,
			InstallStatus:       string(model.AppIntallingStatus),
			DefualtApp:          false,
			PluginName:          req.PluginName,
			PluginDescription:   req.Description,
			ApiEndpoint:         apiEndpoint,
			PluginStoreType:     agentpb.PluginStoreType(req.StoreType),
		},
		Values: &agentpb.AppValues{
			OverrideValues: req.OverrideValues,
			TemplateValues: req.Values,
		},
	}

	if err := pas.UpsertAppConfig(syncConfig); err != nil {
		log.Errorf("failed to update app config data for app %s, %v", req.PluginName, err)
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \" failed to update app config data, %s\"}", err.Error())),
		}, nil
	}

	log.Infof("Triggering child work flow for plugin %v, app: %v", req.PluginName, req.ChartName)
	appDeployReq := prepareAppDeployRequestFromPlugin(req, templateValues)
	result = &model.ResponsePayload{}
	err = workflow.ExecuteChildWorkflow(ctx, Workflow, string(model.AppInstallAction), appDeployReq).Get(ctx, result)
	if err != nil {
		syncConfig.Config.InstallStatus = string(model.AppIntallFailedStatus)
		if err1 := pas.UpsertAppConfig(syncConfig); err1 != nil {
			log.Errorf("failed to update app config data for app %s, %v", req.PluginName, err1)
		}

		result1 := &model.ResponsePayload{}
		err1 := workflow.ExecuteActivity(ctx, a.PluginDeployUpdateStatusActivity, req.PluginName, "plugin-app-installfailed").Get(ctx, result1)
		if err1 != nil {
			log.Errorf("failed to update app config data for app %s, %v, result: %v", req.PluginName, err1, result1)
		}
		return result, err
	}

	syncConfig.Config.InstallStatus = string(model.AppIntalledStatus)
	if err := pas.UpsertAppConfig(syncConfig); err != nil {
		log.Errorf("failed to update app config data for app %s, %v", req.PluginName, err)
		// return &model.ResponsePayload{
		// 	Status:  "FAILED",
		// 	Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \" failed to update app config data, %s\"}", err.Error())),
		// }, nil
	}

	err = workflow.ExecuteActivity(ctx, a.PluginDeployUpdateStatusActivity, req.PluginName, "plugin-app-installed").Get(ctx, result)
	if err != nil {
		log.Errorf("failed to update app config data for app %s, %v", req.PluginName, err)
		// return result, err
	}

	return &model.ResponsePayload{
		Status: "SUCCESS",
	}, nil
}

func executeAppUndeployment(
	ctx workflow.Context,
	req *clusterpluginspb.Plugin,
	a *activities.PluginActivities,
	log logging.Logger,
	pas *captenstore.Store,
) (result *model.ResponsePayload, err error) {
	result = &model.ResponsePayload{}
	err = workflow.ExecuteActivity(ctx, a.PluginDeployUpdateStatusActivity, req.PluginName, "plugin-app-undeploying").Get(ctx, result)
	if err != nil {
		return result, err
	}

	cwo := workflow.ChildWorkflowOptions{
		WorkflowID: "APP-UNDEPLOY-CHILD-WORKFLOW-ID",
	}
	ctx = workflow.WithChildOptions(ctx, cwo)

	syncConfig, err := pas.GetAppConfig(req.PluginName)
	if err != nil {
		log.Errorf("appconfig fetch failed, %v", err)
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \" failed to update app config data, %s\"}", err.Error())),
		}, nil
	}

	syncConfig.Config.InstallStatus = string(model.AppUnInstallingStatus)
	if err := pas.UpsertAppConfig(syncConfig); err != nil {
		log.Errorf("failed to update app config data for app %s, %v", req.PluginName, err)
		return &model.ResponsePayload{
			Status:  "FAILED",
			Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \" failed to update app config data, %s\"}", err.Error())),
		}, nil
	}

	log.Info("invoking child workflow")
	appDeployReq := prepareAppUndeployRequestFromPlugin(syncConfig)
	result = &model.ResponsePayload{}
	err = workflow.ExecuteChildWorkflow(ctx, Workflow, string(model.AppUnInstallAction), appDeployReq).Get(ctx, result)
	if err != nil {
		syncConfig.Config.InstallStatus = string(model.AppUnUninstallFailedStatus)
		if err1 := pas.UpsertAppConfig(syncConfig); err1 != nil {
			log.Errorf("failed to update app config data for app %s, %v", req.PluginName, err1)
		}

		result1 := &model.ResponsePayload{}
		err1 := workflow.ExecuteActivity(ctx, a.PluginDeployUpdateStatusActivity, req.PluginName, "plugin-app-uninstallfailed").Get(ctx, result1)
		if err1 != nil {
			log.Errorf("failed to update app config data for app %s, error: %v, response: %+v", req.PluginName, err1, result1)
		}
		return result, err
	}

	log.Info("finished child workflow")
	syncConfig.Config.InstallStatus = string(model.AppUnInstalledStatus)
	if err := pas.DeleteAppConfigByReleaseName(req.PluginName); err != nil {
		log.Errorf("failed to update app config data for app %s, %v", req.PluginName, err)
		// return &model.ResponsePayload{
		// 	Status:  "FAILED",
		// 	Message: json.RawMessage(fmt.Sprintf("{ \"reason\": \" failed to update app config data, %s\"}", err.Error())),
		// }, nil
	}

	err = workflow.ExecuteActivity(ctx, a.PluginDeployUpdateStatusActivity, req.PluginName, "plugin-app-uninstalled").Get(ctx, result)
	if err != nil {
		log.Errorf("failed to update app config data for app %s, %v", req.PluginName, err)
		// return result, err
	}

	return &model.ResponsePayload{
		Status: "SUCCESS",
	}, nil
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
		ChartName:      data.ChartName + "/" + data.ChartName,
		Namespace:      data.DefaultNamespace,
		ReleaseName:    data.ChartName,
		Version:        data.Version,
		ClusterName:    "capten",
		OverrideValues: string(values),
		Timeout:        10,
	}
}

func prepareAppUndeployRequestFromPlugin(data *agentpb.SyncAppData) *model.ApplicationDeleteRequest {
	return &model.ApplicationDeleteRequest{
		PluginName:  "helm",
		Namespace:   data.Config.Namespace,
		ReleaseName: data.Config.ReleaseName,
		ClusterName: "capten",
		Timeout:     10,
	}
}

func executeStringTemplateValues(data string, values []byte) (transformedData string, err error) {
	if len(data) == 0 {
		return
	}

	tmpl, err := template.New("templateVal").Parse(data)
	if err != nil {
		return
	}

	mapValues := map[string]any{}
	if err = yaml.Unmarshal(values, &mapValues); err != nil {
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, mapValues)
	if err != nil {
		return
	}

	transformedData = string(buf.Bytes())
	return
}
