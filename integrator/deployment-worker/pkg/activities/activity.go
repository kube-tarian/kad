package activities

import (
	"context"
	"log"

	"github.com/kube-tarian/kad/integrator/deployment-worker/pkg/model"
	"go.temporal.io/sdk/activity"
)

type Activities struct {
}

func (a *Activities) DeploymentActivity(ctx context.Context, req model.RequestPayload) (model.ResponsePayload, error) {
	logger := activity.GetLogger(ctx)
	e := activity.GetInfo(ctx)
	logger.Info("Activity", "name", req)
	log.Printf("activity info: %+v\n", e)
	return model.ResponsePayload{
		Status: "Success",
	}, nil
}
