package activities

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kube-tarian/kad/capten/common-pkg/plugins/helm"
	"github.com/kube-tarian/kad/capten/model"
)

func installApplication(req *model.ApplicationDeployRequest) (model.ResponsePayload, error) {
	helmClient, err := helm.NewClient(logger)
	if err != nil {
		logger.Errorf("Get plugin  failed: %v", err)
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", strings.ReplaceAll(err.Error(), "\"", "\\\""))),
		}, err
	}

	msg, err := helmClient.Create(&model.CreteRequestPayload{
		RepoName:    req.RepoName,
		RepoURL:     req.RepoURL,
		ChartName:   req.ChartName,
		Namespace:   req.Namespace,
		ReleaseName: req.ReleaseName,
		Timeout:     int(req.Timeout),
		Version:     req.Version,
		ValuesYaml:  req.OverrideValues,
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

func uninstallApplication(req *model.DeployerDeleteRequest) (model.ResponsePayload, error) {
	helmClient, err := helm.NewClient(logger)
	if err != nil {
		logger.Errorf("Get helm client  failed: %v", err)
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", strings.ReplaceAll(err.Error(), "\"", "\\\""))),
		}, err
	}

	msg, err := helmClient.Delete(&model.DeleteRequestPayload{
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
