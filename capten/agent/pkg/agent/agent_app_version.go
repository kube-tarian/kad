package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/pkg/errors"
)

func (a *Agent) UpgradeAppWithVersion(ctx context.Context, req *agentpb.UpgradeAppWithVersionRequest) (*agentpb.UpgradeAppWithVersionResponse, error) {

	a.log.Infof("Received request for UpgradeApp, app %s", req.ReleaseName)

	appConfig, err := a.as.GetAppConfig(req.ReleaseName)
	if err != nil {
		a.log.Errorf("failed to GetAppConfig for release_name: %s err: %v", req.ReleaseName, err)
		return &agentpb.UpgradeAppWithVersionResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err fetching appConfig").Error(),
		}, nil
	}

	launchUiValues, err := GetSSOvalues(req.ReleaseName)
	if err != nil {
		a.log.Errorf("failed to getLanchUiValues for release: %s err: %v", req.ReleaseName, err)
		return &agentpb.UpgradeAppWithVersionResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err in getLanchUiValues").Error(),
		}, nil
	}

	newAppConfig, marshaledOverrideValues, err := PopulateTemplateValues(appConfig, nil, launchUiValues, a.log)
	if err != nil {
		return &agentpb.UpgradeAppWithVersionResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err populating template values").Error(),
		}, nil
	}

	installReq := toAppDeployRequestFromSyncApp(newAppConfig, marshaledOverrideValues)
	installReq.Version = req.GetVersion()
	go a.DeployApp(installReq, newAppConfig, []byte("update"))

	a.log.Infof("Triggerred app [%s] update", newAppConfig.Config.ReleaseName)
	return &agentpb.UpgradeAppWithVersionResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "Triggerred app upgrade",
	}, nil
}
