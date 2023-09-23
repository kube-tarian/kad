package activities

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/model"
)

type Activities struct {
}

var logger = logging.NewLogger()

func (a *Activities) ConfigurationActivity(ctx context.Context, params model.ConfigureParameters, payload json.RawMessage) (model.ResponsePayload, error) {
	logger.Infof("Activity, name: %+v", payload)
	// e := activity.GetInfo(ctx)
	// logger.Infof("activity info: %+v", e)

	switch params.Resource {
	case "cluster":
		return handleCluster(ctx, params, payload)
	case "repository":
		return handleRepository(ctx, params, payload)
	case "project":
		return handleProject(ctx, params, payload)
	case "git":
		return handleGit(ctx, params, payload)
	default:
		logger.Errorf("unknown resource type: %s in configuration", params.Resource)
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage("{\"error\": \"unknown resource type in configuration\"}"),
		}, fmt.Errorf("unknown resource type in configuration")
	}
}
