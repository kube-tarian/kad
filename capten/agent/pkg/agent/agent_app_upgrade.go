package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/workers"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func (a *Agent) UpgradeApp(ctx context.Context, req *agentpb.UpgradeAppRequest) (*agentpb.UpgradeAppResponse, error) {

	if req.ReleaseName == "" {
		return &agentpb.UpgradeAppResponse{
			Status:        agentpb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "release name empty",
		}, nil
	}
	a.log.Infof("Received request for UpgradeApp, app %s", req.ReleaseName)

	appConfig, err := a.as.GetAppConfig(req.ReleaseName)
	if err != nil {
		a.log.Errorf("failed to GetAppConfig for release_name: %s err: %v", req.ReleaseName, err)
		return &agentpb.UpgradeAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err fetching appConfig").Error(),
		}, nil
	}

	// replace template with new override values from the request
	templateValuesMapping, err := deriveTemplateValuesMapping(req.OverrideValues, appConfig.Values.TemplateValues)
	if err != nil {
		a.log.Errorf("failed to derive template values, err: %v", err)
		return &agentpb.UpgradeAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "failed to dervice template values").Error(),
		}, nil
	}

	launchUiMapping := map[string]any{}
	newAppConfig := *appConfig

	if len(req.LaunchUIValues) > 0 {
		// replace launchUiMapping with the new launchUi values from the request
		launchUiMapping, err = deriveTemplateValuesMapping(req.LaunchUIValues, appConfig.Values.LaunchUIValues)
		if err != nil {
			a.log.Errorf("failed to replaceTemplateValues, err: %v", err)
			return &agentpb.UpgradeAppResponse{
				Status:        agentpb.StatusCode_INTERNRAL_ERROR,
				StatusMessage: errors.WithMessage(err, "err replacing launchUiMapping").Error(),
			}, nil
		}
	}

	finalOverrideValuesMapping := mergeRecursive(convertKey(templateValuesMapping), convertKey(launchUiMapping))
	marshaledOverrideValues, err := yaml.Marshal(finalOverrideValuesMapping)
	if err != nil {
		a.log.Errorf("failed to Marshal, err: %v", err)
		return &agentpb.UpgradeAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err marshalling overrideValues").Error(),
		}, nil
	}

	newAppConfig.Config.InstallStatus = "Updating"

	if err := a.as.UpsertAppConfig(&newAppConfig); err != nil {
		a.log.Errorf("failed to UpsertAppConfig, err: %v", err)
		return &agentpb.UpgradeAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err upserting new appConfig").Error(),
		}, nil
	}

	go func() {
		wd := workers.NewDeployment(a.tc, a.log)
		_, err := wd.SendEvent(context.TODO(), "update",
			toAppDeployRequestFromSyncApp(&newAppConfig, marshaledOverrideValues))
		if err != nil {
			newAppConfig.Config.InstallStatus = "Update Failed"
			if err := a.as.UpsertAppConfig(&newAppConfig); err != nil {
				a.log.Errorf("failed to UpsertAppConfig, err: %v", err)
				return
			}
			a.log.Errorf("failed to SendEvent, err: %v", err)
			return
		}

		newAppConfig.Config.InstallStatus = "Updated"
		// update the new override values in db from the request
		newAppConfig.Values.OverrideValues = req.OverrideValues

		if err := a.as.UpsertAppConfig(&newAppConfig); err != nil {
			a.log.Errorf("failed to UpsertAppConfig, err: %v", err)
			return
		}
	}()

	a.log.Infof("Triggerred app [%s] update", newAppConfig.Config.ReleaseName)
	return &agentpb.UpgradeAppResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "Triggerred app upgrade",
	}, nil

}
