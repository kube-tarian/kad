package activities

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kube-tarian/kad/integrator/deployment-worker/pkg/model"
	"github.com/kube-tarian/kad/integrator/deployment-worker/pkg/plugins"
	"github.com/kube-tarian/kad/integrator/pkg/logging"
)

type Activities struct {
}

func (a *Activities) DeploymentActivity(ctx context.Context, req model.RequestPayload) (model.ResponsePayload, error) {
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
	msg, err := plugin.Exec(req)
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
