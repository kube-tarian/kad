package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/credential"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func (a *Agent) ConfigureAppSSO(
	ctx context.Context, req *agentpb.ConfigureAppSSORequest) (*agentpb.ConfigureAppSSOResponse, error) {

	a.log.Infof("Received request for ConfigureAppSSO, app %s", req.ReleaseName)

	appConfig, err := a.as.GetAppConfig(req.ReleaseName)
	if err != nil {
		a.log.Errorf("failed to GetAppConfig for release_name: %s err: %v", req.ReleaseName, err)
		return &agentpb.ConfigureAppSSOResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err fetching appConfig").Error(),
		}, nil
	}

	if err := credential.StoreAppOauthCredential(ctx, req.ReleaseName, req.ClientId, req.ClientSecret); err != nil {
		a.log.Errorf("failed to store credential for ClientId: %s, %v", req.ClientId, err)
		return &agentpb.ConfigureAppSSOResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err saving SSO credentials in vault").Error(),
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
		return nil, errors.WithMessagef(err, "failed to Unmarshal override values")
	}
	overrideValuesMapping["OAuthBaseURL"] = req.OAuthBaseURL

	for key, val := range overrideValuesMapping {
		ssoOverwriteMapping[key] = val
	}

	ssoOverwriteBytes, _ := yaml.Marshal(ssoOverwriteMapping)
	overrideValuesBytes, _ := yaml.Marshal(overrideValuesMapping)
	newAppConfig, marshaledOverrideValues, err := PopulateTemplateValues(appConfig, overrideValuesBytes, ssoOverwriteBytes, a.log)
	if err != nil {
		return nil, errors.WithMessage(err, "err PopulateTemplateValues")
	}

	newAppConfig.Config.InstallStatus = "Updating"

	if err := a.as.UpsertAppConfig(newAppConfig); err != nil {
		a.log.Errorf("failed to UpsertAppConfig, err: %v", err)
		return &agentpb.ConfigureAppSSOResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: errors.WithMessage(err, "err upserting new appConfig").Error(),
		}, nil
	}

	installReq := toAppDeployRequestFromSyncApp(newAppConfig, marshaledOverrideValues)
	go a.DeployApp(installReq, newAppConfig, []byte("update"))

	a.log.Infof("Triggerred app [%s] update", newAppConfig.Config.ReleaseName)
	return &agentpb.ConfigureAppSSOResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "Triggerred app upgrade",
	}, nil
}
