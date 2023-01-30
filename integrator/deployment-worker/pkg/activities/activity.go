package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/kube-tarian/kad/integrator/common-pkg/plugins"
	workerframework "github.com/kube-tarian/kad/integrator/common-pkg/worker-framework"
	"github.com/kube-tarian/kad/integrator/model"
)

type Activities struct {
}

func (a *Activities) DeploymentActivity(ctx context.Context, req model.RequestPayload) (model.ResponsePayload, error) {
	logger := logging.NewLogger()
	logger.Infof("Activity, name: %+v", req.ToString())
	// e := activity.GetInfo(ctx)
	// logger.Infof("activity info: %+v", e)

	plugin, err := plugins.GetPlugin(req.PluginName, logger)
	if err != nil {
		logger.Errorf("Get plugin  failed: %v", err)
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", strings.ReplaceAll(err.Error(), "\"", "\\\""))),
		}, err
	}
	deployerPlugin, ok := plugin.(workerframework.DeploymentWorker)
	if !ok {
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", strings.ReplaceAll(err.Error(), "\"", "\\\""))),
		}, fmt.Errorf("plugin not supports deployment activities")
	}
	msg, err := deployerPlugin.DeployActivities(req)
	if err != nil {
		logger.Errorf("Deploy activities failed %s: %v", req.Action, err)
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", strings.ReplaceAll(err.Error(), "\"", "\\\""))),
		}, err
	}

	return model.ResponsePayload{
		Status:  "Success",
		Message: msg,
	}, nil
}
