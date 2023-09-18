package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/pkg/errors"
)

func (a *Agent) UpgradeAppWithValues(ctx context.Context, req *agentpb.UpgradeAppWithValuesRequest) (*agentpb.UpgradeAppWithValuesResponse, error) {

	a.log.Infof("Received request for UpgradeApp, app %s", req.ReleaseName)

	// Get the config templates for release name
	appConfig, err := a.as.GetAppConfig(req.ReleaseName)
	if err != nil {
		a.log.Errorf("failed to GetAppConfig for release_name: %s err: %v", req.ReleaseName, err)
		return &agentpb.UpgradeAppWithValuesResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err fetching appConfig").Error(),
		}, nil
	}

	launchUiValues, err := GetSSOvalues(req.ReleaseName)
	if err != nil {
		a.log.Errorf("failed to getLanchUiValues for release: %s err: %v", req.ReleaseName, err)
		return &agentpb.UpgradeAppWithValuesResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err in getLanchUiValues").Error(),
		}, nil
	}

	// populate template values, overriding with launchUiValues if needed
	newAppConfig, marshaledOverrideValues, err := PopulateTemplateValues(appConfig, req.OverrideValues, launchUiValues, a.log)
	if err != nil {
		return &agentpb.UpgradeAppWithValuesResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err populating template values").Error(),
		}, nil
	}
	appConfig.Config.InstallStatus = "updating"

	// Upsert the config with status as updating
	if err := a.as.UpsertAppConfig(appConfig); err != nil {
		a.log.Errorf("failed to UpsertAppConfig, err: %v", err)
		return &agentpb.UpgradeAppWithValuesResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err upserting new appConfig").Error(),
		}, nil
	}

	installReq := toAppDeployRequestFromSyncApp(newAppConfig, marshaledOverrideValues)
	go a.DeployApp(installReq, newAppConfig, []byte("update"))

	a.log.Infof("Triggerred app [%s] update", newAppConfig.Config.ReleaseName)
	return &agentpb.UpgradeAppWithValuesResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "Triggerred app upgrade",
	}, nil

}
