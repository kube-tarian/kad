package crossplane

import (
	"context"
	"encoding/json"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/model"
)

type CrossPlaneActivities struct{}

var logger = logging.NewLogger()

func (c *CrossPlaneActivities) ConfigurationActivity(ctx context.Context, params model.ConfigureParameters, payload json.RawMessage) (model.ResponsePayload, error) {
	logger.Infof("Activity: %s, %s", params.Resource, params.Action)

	req := &model.CrossplaneUseCase{}
	if err := json.Unmarshal(payload, req); err != nil {
		return model.ResponsePayload{
			Status:  string(model.WorkFlowStatusFailed),
			Message: json.RawMessage("{\"error\": \"failed to read payload\"}"),
		}, err
	}

	config, err := NewCrossPlaneApp()
	if err != nil {
		return model.ResponsePayload{
			Status:  string(model.WorkFlowStatusFailed),
			Message: json.RawMessage("{\"error\": \"failed to initialize crossplane plugin\"}"),
		}, err
	}

	status, err := config.Configure(ctx, req)
	if err != nil {
		logger.Errorf("crossplane plugin configure failed, %v", err)
		return model.ResponsePayload{
			Status:  status,
			Message: json.RawMessage("{\"error\": \"failed to configure crossplane plugin\"}"),
		}, err
	}
	logger.Infof("crossplane plugin configured")
	return model.ResponsePayload{
		Status: status,
	}, err
}
