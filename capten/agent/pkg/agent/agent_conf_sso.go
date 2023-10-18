package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/credential"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/agentpb"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func (a *Agent) ConfigureAppSSO(ctx context.Context,
	req *agentpb.ConfigureAppSSORequest) (*agentpb.ConfigureAppSSOResponse, error) {
	a.log.Infof("Received ConfigureAppSSO request, %+v", req)

	appConfig, err := a.as.GetAppConfig(req.ReleaseName)
	if err != nil {
		a.log.Errorf("failed to read app %s config data, %v", req.ReleaseName, err)
		return &agentpb.ConfigureAppSSOResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "failed to read app config data").Error(),
		}, nil
	}

	if err := credential.StoreAppOauthCredential(ctx, req.ReleaseName, req.ClientId, req.ClientSecret); err != nil {
		a.log.Errorf("failed to store oauth credential for app %s, %v", req.ReleaseName, err)
		return &agentpb.ConfigureAppSSOResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "failed to store SSO credential").Error(),
		}, nil
	}

	ssoOverwriteMapping := map[string]any{
		"ClientId":     req.ClientId,
		"ClientSecret": req.ClientSecret,
		"OAuthBaseURL": req.OAuthBaseURL,
	}

	// save OAuthBaseURL in the db as part of the override values
	overrideValuesMapping := map[string]any{}
	if err := yaml.Unmarshal(appConfig.Values.OverrideValues, &overrideValuesMapping); err != nil {
		a.log.Errorf("failed to ummrashal override values for app %s, %v", req.ReleaseName, err)
		return &agentpb.ConfigureAppSSOResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "failed to prepare app values").Error(),
		}, nil
	}
	overrideValuesMapping["OAuthBaseURL"] = req.OAuthBaseURL

	for key, val := range overrideValuesMapping {
		ssoOverwriteMapping[key] = val
	}

	ssoOverwriteBytes, _ := yaml.Marshal(ssoOverwriteMapping)
	overrideValuesBytes, _ := yaml.Marshal(overrideValuesMapping)
	updateAppConfig, marshaledOverrideValues, err := populateTemplateValues(appConfig, overrideValuesBytes, ssoOverwriteBytes, a.log)
	if err != nil {
		a.log.Errorf("failed to populate template values for app %s, %v", req.ReleaseName, err)
		return &agentpb.ConfigureAppSSOResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "failed to prepare app values").Error(),
		}, nil
	}

	updateAppConfig.Config.InstallStatus = string(appUpgradingStatus)
	if err := a.as.UpsertAppConfig(updateAppConfig); err != nil {
		a.log.Errorf("failed to update app config data for app %s, %v", req.ReleaseName, err)
		return &agentpb.ConfigureAppSSOResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "failed to update app config data").Error(),
		}, nil
	}

	deployReq := prepareAppDeployRequestFromSyncApp(updateAppConfig, marshaledOverrideValues)
	go a.upgradeAppWithWorkflow(deployReq, updateAppConfig)

	a.log.Infof("Triggerred app [%s] upgrade with SSO configure", updateAppConfig.Config.ReleaseName)
	return &agentpb.ConfigureAppSSOResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "Triggerred app upgrade",
	}, nil
}
