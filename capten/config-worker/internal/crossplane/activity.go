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
	config, err := NewCrossPlaneApp()
	if err != nil {
		return model.ResponsePayload{
			Status:  string(model.WorkFlowStatusFailed),
			Message: json.RawMessage("{\"error\": \"failed to initialize crossplane plugin\"}"),
		}, err
	}

	status := ""

	switch params.Action {
	case model.CrossPlaneClusterUpdate:
		reqLocal := &model.CrossplaneClusterUpdate{}
		if err = json.Unmarshal(payload, reqLocal); err != nil {
			logger.Errorf("failed to unmarshall the crossplane req for %s, %v", model.CrossPlaneClusterUpdate, err)
			err = fmt.Errorf("failed to unmarshall the crossplane req for %s", model.CrossPlaneClusterUpdate)
		}
		status, err = config.configureClusterUpdate(ctx, reqLocal)
		if err != nil {
			logger.Errorf("failed to configure crossplane project for %s, %v", model.CrossPlaneClusterUpdate, err)
			err = fmt.Errorf("failed to configure crossplane project for %s", model.CrossPlaneClusterUpdate)
		}
	default:
		reqLocal := &model.CrossplaneUseCase{}
		if err = json.Unmarshal(payload, reqLocal); err != nil {
			logger.Errorf("failed to unmarshall the crossplane req, %v", err)
			err = fmt.Errorf("failed to unmarshall the crossplane req")
		}
		status, err = config.configureProjectAndApps(ctx, reqLocal)
		if err != nil {
			logger.Errorf("failed to configure crossplane project, %v", err)
			err = fmt.Errorf("failed to configure crossplane project")
		}
	}
	logger.Infof("crossplane plugin configured")
	return model.ResponsePayload{Status: status}, err
}
