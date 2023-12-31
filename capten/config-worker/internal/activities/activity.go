package activities

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/config-worker/internal/crossplane"
	"github.com/kube-tarian/kad/capten/config-worker/internal/tekton"
	"github.com/kube-tarian/kad/capten/model"
	agentmodel "github.com/kube-tarian/kad/capten/model"
)

type Activities struct{}

var logger = logging.NewLogger()

func (c *Activities) ConfigurationActivity(ctx context.Context, params model.ConfigureParameters, payload json.RawMessage) (model.ResponsePayload, error) {
	logger.Infof("Activity: %s, %s", params.Resource, params.Action)

	switch params.Resource {
	case "crossplane":
		ca := crossplane.CrossPlaneActivities{}
		return ca.ConfigurationActivity(ctx, params, payload)
	case model.TektonPipelineConfigUseCase:
		tek := tekton.TektonpipelineActivty{}
		return tek.ConfigurationActivity(ctx, params, payload)
	default:
		logger.Errorf("unknown resource type: %s in configuration", params.Resource)
		return model.ResponsePayload{
			Status:  string(agentmodel.WorkFlowStatusFailed),
			Message: json.RawMessage("{\"error\": \"unknown resource type\"}"),
		}, fmt.Errorf("unknown resource type: %s in configuration", params.Resource)
	}
}
