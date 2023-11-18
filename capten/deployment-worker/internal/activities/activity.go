package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/plugins/helm"
	"github.com/kube-tarian/kad/capten/model"
)

type Activities struct {
}

var logger = logging.NewLogger()

func (a *Activities) DeploymentInstallActivity(ctx context.Context, req *model.ApplicationDeployRequest) (model.ResponsePayload, error) {
	logger.Infof("Activity, name: %+v", req)

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

func (a *Activities) DeploymentDeleteActivity(ctx context.Context, req *model.DeployerDeleteRequest) (model.ResponsePayload, error) {
	logger.Infof("Activity, name: %+v", req)

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
