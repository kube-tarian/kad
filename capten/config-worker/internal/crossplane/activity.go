package crossplane

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/model"
)

type CrossPlaneActivities struct{}

var logger = logging.NewLogger()

func (c *CrossPlaneActivities) ConfigurationActivity(ctx context.Context, params model.ConfigureParameters, payload json.RawMessage) (model.ResponsePayload, error) {
	logger.Infof("Activity: %s, %s", params.Resource, params.Action)
	status, err := processConfigurationActivity(ctx, params, payload)
	if err != nil {
		return model.ResponsePayload{
			Status: status,
			Message: json.RawMessage(
				fmt.Sprintf("{\"error\": \"%s\"}", err.Error())),
		}, err
	}

	logger.Infof("crossplane plugin action %s configured", params.Action)
	return model.ResponsePayload{Status: status}, err
}

func processConfigurationActivity(ctx context.Context, params model.ConfigureParameters, payload json.RawMessage) (string, error) {
	cp, err := NewCrossPlaneApp()
	if err != nil {
		return string(model.WorkFlowStatusFailed), fmt.Errorf("failed to initialize crossplane plugin")
	}

	switch params.Action {
	case model.CrossPlaneClusterUpdate:
		reqLocal := &model.CrossplaneClusterUpdate{}
		if err := json.Unmarshal(payload, reqLocal); err != nil {
			logger.Errorf("failed to unmarshall the crossplane req for %s, %v", model.CrossPlaneClusterUpdate, err)
			return string(model.WorkFlowStatusFailed), fmt.Errorf("failed to unmarshall the crossplane req for %s", model.CrossPlaneClusterUpdate)
		}
		status, err := cp.configureClusterUpdate(ctx, reqLocal)
		if err != nil {
			logger.Errorf("failed to configure crossplane project for %s, %v", model.CrossPlaneClusterUpdate, err)
			return status, fmt.Errorf("failed to configure crossplane project for %s", model.CrossPlaneClusterUpdate)
		}
		return status, nil
	case model.CrossPlaneProjectSync:
		reqLocal := &model.CrossplaneUseCase{}
		if err := json.Unmarshal(payload, reqLocal); err != nil {
			logger.Errorf("failed to unmarshall the crossplane req, %v", err)
			return string(model.WorkFlowStatusFailed), fmt.Errorf("failed to unmarshall the crossplane req")
		}
		status, err := cp.configureProjectAndApps(ctx, reqLocal)
		if err != nil {
			logger.Errorf("failed to configure crossplane project, %v", err)
			return string(model.WorkFlowStatusFailed), fmt.Errorf("failed to configure crossplane project")
		}
		return status, nil
	case model.CrossPlaneProjectDelete:
		reqLocal := &model.CrossplaneClusterUpdate{}
		if err := json.Unmarshal(payload, reqLocal); err != nil {
			logger.Errorf("failed to unmarshall the crossplane req for %s, %v", model.CrossPlaneClusterUpdate, err)
			return string(model.WorkFlowStatusFailed), fmt.Errorf("failed to unmarshall the crossplane req for %s", model.CrossPlaneClusterUpdate)
		}
		status, err := cp.configureClusterDelete(ctx, reqLocal)
		if err != nil {
			logger.Errorf("failed to configure crossplane project for %s, %v", model.CrossPlaneClusterUpdate, err)
			return status, fmt.Errorf("failed to configure crossplane project for %s", model.CrossPlaneClusterUpdate)
		}
		return status, nil
	default:
		return string(model.WorkFlowStatusFailed), fmt.Errorf("invalid crossplane action")
	}
}
