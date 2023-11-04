package activities

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/intelops/go-common/logging"
	agentmodel "github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/model"
)

type Activities struct{}

var logger = logging.NewLogger()

func (a *Activities) ConfigurationActivity(ctx context.Context, params model.ConfigureParameters, payload json.RawMessage) (model.ResponsePayload, error) {
	logger.Infof("Activity: %s, %s", params.Resource, params.Action)

	switch params.Resource {
	case "cluster":
		return handleCluster(ctx, params, payload)
	case "repository":
		return handleRepository(ctx, params, payload)
	case "project":
		return handleProject(ctx, params, payload)
	case crossPlaneUsecase:
		config, err := NewCrossPlaneApp()
		if err != nil {
			return model.ResponsePayload{
				Status:  string(agentmodel.WorkFlowStatusFailed),
				Message: json.RawMessage("{\"error\": \"failed to initialize crossplane plugin\"}"),
			}, err
		}
		return config.ExecuteSteps(ctx, params, payload)
	default:
		logger.Errorf("unknown resource type: %s in configuration", params.Resource)
		return model.ResponsePayload{
			Status:  string(agentmodel.WorkFlowStatusFailed),
			Message: json.RawMessage("{\"error\": \"unknown resource type in configuration\"}"),
		}, fmt.Errorf("unknown resource type in configuration")
	}
}
