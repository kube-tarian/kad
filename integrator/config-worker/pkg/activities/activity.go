package activities

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/kube-tarian/kad/integrator/model"
)

type Activities struct {
}

var logger = logging.NewLogger()

func (a *Activities) ConfigurationActivity(ctx context.Context, params model.ConfigureParameters, payload interface{}) (model.ResponsePayload, error) {
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
	default:
		logger.Errorf("unknown resource type in configuration")
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage("{\"error\": \"unknown resource type in configuration\"}"),
		}, fmt.Errorf("unknown resource type in configuration")
	}
}
