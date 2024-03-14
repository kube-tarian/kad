package activities

import (
	"context"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/model"
)

type Activities struct {
}

var logger = logging.NewLogger()

func (a *Activities) DeploymentInstallActivity(ctx context.Context, req *model.ApplicationDeployRequest) (model.ResponsePayload, error) {
	logger.Infof("Activity, name: %+v", req)
	return installApplication(req)
}

func (a *Activities) DeploymentDeleteActivity(ctx context.Context, req *model.DeployerDeleteRequest) (model.ResponsePayload, error) {
	logger.Infof("Activity, name: %+v", req)

	return uninstallApplication(req)
}
