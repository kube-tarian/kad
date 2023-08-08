package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/plugins"
	workerframework "github.com/kube-tarian/kad/capten/common-pkg/worker-framework"
	"github.com/kube-tarian/kad/capten/model"
)

type Activities struct {
}

var logger = logging.NewLogger()

func (a *Activities) ClimonInstallActivity(ctx context.Context, req *model.ClimonPostRequest) (model.ResponsePayload, error) {
	logger.Infof("Activity, name: %+v", req)
	plugin, err := plugins.GetPlugin(req.PluginName, logger)
	if err != nil {
		logger.Errorf("Get plugin  failed: %v", err)
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", strings.ReplaceAll(err.Error(), "\"", "\\\""))),
		}, err
	}

	climonPlugin, ok := plugin.(workerframework.ClimonWorker)
	if !ok {
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage("{\"error\": \"not impmented climon worker plugin\"}"),
		}, fmt.Errorf("plugin not supports deployment activities")
	}

	emptyVersion := ""
	if req.Version == nil {
		req.Version = &emptyVersion
	}
	msg, err := climonPlugin.Create(&model.CreteRequestPayload{
		RepoName:    req.RepoName,
		RepoURL:     req.RepoUrl,
		ChartName:   req.ChartName,
		Namespace:   req.Namespace,
		ReleaseName: req.ReleaseName,
		Timeout:     req.Timeout,
		Version:     *req.Version,
	})
	if err != nil {
		logger.Errorf("Deploy activities failed %v", err)
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

func (a *Activities) ClimonDeleteActivity(ctx context.Context, req *model.ClimonDeleteRequest) (model.ResponsePayload, error) {
	logger.Infof("Activity, name: %+v", req)
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

	msg, err := deployerPlugin.Delete(&model.DeleteRequestPayload{
		Namespace:   req.Namespace,
		ReleaseName: req.ReleaseName,
		Timeout:     req.Timeout,
		ClusterName: *req.ClusterName,
	})

	if err != nil {
		logger.Errorf("Deploy activities failed %v", err)
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
