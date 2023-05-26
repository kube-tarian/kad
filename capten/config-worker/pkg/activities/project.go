package activities

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kube-tarian/kad/capten/model"
)

func handleProject(ctx context.Context, params model.ConfigureParameters, payload interface{}) (model.ResponsePayload, error) {
	var err error
	switch params.Action {
	case "add":
		if req, ok := payload.(model.ProjectPostRequest); ok {
			return addProject(req)
		}
		err = fmt.Errorf("wrong payload")
	case "delete":
		if req, ok := payload.(model.ProjectDeleteRequest); ok {
			return deleteProject(req)
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

func addProject(req model.ProjectPostRequest) (model.ResponsePayload, error) {
	configPlugin, err := getConfigPlugin(req.PluginName, logger)
	if err != nil {
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", err)),
		}, err
	}

	msg, err := configPlugin.ProjectAdd(req)
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

func deleteProject(req model.ProjectDeleteRequest) (model.ResponsePayload, error) {
	configPlugin, err := getConfigPlugin(req.PluginName, logger)
	if err != nil {
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", err)),
		}, err
	}

	msg, err := configPlugin.ProjectDelete(req)
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
