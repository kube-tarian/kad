package activities

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/kube-tarian/kad/integrator/common-pkg/plugins"
	workerframework "github.com/kube-tarian/kad/integrator/common-pkg/worker-framework"
	"github.com/kube-tarian/kad/integrator/model"
)

type Activities struct {
}

func (a *Activities) ConfigurationActivity(ctx context.Context, req model.ConfigPayload) (model.ResponsePayload, error) {
	logger := logging.NewLogger()
	logger.Infof("Activity, name: %+v", req)
	// e := activity.GetInfo(ctx)
	// logger.Infof("activity info: %+v", e)

	plugin, err := plugins.GetPlugin(req.PluginName, logger)
	if err != nil {
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", err)),
		}, err
	}

	configPlugin, ok := plugin.(workerframework.ConfigurationWorker)
	if !ok {
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", err)),
		}, fmt.Errorf("plugin not supports Configuration activities")
	}

	msg, err := configPlugin.ConfigurationActivities(req)
	if err != nil {
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", err)),
		}, err
	}

	return model.ResponsePayload{
		Status:  "Success",
		Message: msg,
	}, nil
}
