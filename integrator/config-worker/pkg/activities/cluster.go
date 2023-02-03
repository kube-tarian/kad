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

func handleCluster(ctx context.Context, params model.ConfigureParameters, payload interface{}) (model.ResponsePayload, error) {
	var err error
	switch params.Action {
	case "add":
		if req, ok := payload.(model.ClusterRequest); ok {
			return addCluster(req)
		}
		err = fmt.Errorf("wrong payload")
	case "delete":
		if req, ok := payload.(model.ClusterRequest); ok {
			return deleteCluster(req)
		}
		err = fmt.Errorf("wrong payload")
	default:
		err = fmt.Errorf("unknown action %s for resouce %s", params.Action, params.Resource)
	}
	return model.ResponsePayload{
		Status:  "Failed",
		Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", err.Error())),
	}, err
}

func getConfigPlugin(pluginName string, log logging.Logger) (workerframework.ConfigurationWorker, error) {
	plugin, err := plugins.GetPlugin(pluginName, logger)
	if err != nil {
		return nil, err
	}

	configPlugin, ok := plugin.(workerframework.ConfigurationWorker)
	if !ok {
		return nil, fmt.Errorf("plugin not supports Configuration activities")
	}

	return configPlugin, nil
}

func addCluster(req model.ClusterRequest) (model.ResponsePayload, error) {
	configPlugin, err := getConfigPlugin(req.PluginName, logger)
	if err != nil {
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", err)),
		}, err
	}

	msg, err := configPlugin.ClusterAdd(req)
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

func deleteCluster(req model.ClusterRequest) (model.ResponsePayload, error) {
	configPlugin, err := getConfigPlugin(req.PluginName, logger)
	if err != nil {
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", err)),
		}, err
	}

	msg, err := configPlugin.ClusterDelete(req)
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
