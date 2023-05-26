package activities

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kube-tarian/kad/capten/model"
)

func handleRepository(ctx context.Context, params model.ConfigureParameters, payload interface{}) (model.ResponsePayload, error) {
	var err error
	switch params.Action {
	case "add":
		if req, ok := payload.(model.RepositoryPostRequest); ok {
			return addRepository(req)
		}
		err = fmt.Errorf("wrong payload")
	case "delete":
		if req, ok := payload.(model.RepositoryDeleteRequest); ok {
			return deleteRepository(req)
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

func addRepository(req model.RepositoryPostRequest) (model.ResponsePayload, error) {
	configPlugin, err := getConfigPlugin(req.PluginName, logger)
	if err != nil {
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", err)),
		}, err
	}

	msg, err := configPlugin.RepositoryAdd(req)
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

func deleteRepository(req model.RepositoryDeleteRequest) (model.ResponsePayload, error) {
	configPlugin, err := getConfigPlugin(req.PluginName, logger)
	if err != nil {
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", err)),
		}, err
	}

	msg, err := configPlugin.RepositoryDelete(req)
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
