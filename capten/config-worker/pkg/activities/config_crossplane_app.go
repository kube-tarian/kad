package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	agentmodel "github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/model"
	"github.com/otiai10/copy"
)

const (
	tektonUsecase     = "tekton"
	crossPlaneUsecase = "crossplane"

	crossPlaneGitRepoAttribute               = "git_repo"
	crossPlaneSyncjConfigPathAttribute       = "synch_config_path"
	crossPlaneProviderConfigPathAttribute    = "provider_config_path"
	crossPlaneProviderConfigMainAppAttribute = "provider_config_main_app"
	tmpCloneStr                              = "clone*"
)

type CrossPlaneApp struct {
	*ConfigureApp
	pluginConfig crossplanePluginConfig
}

func NewCrossPlaneApp() (*CrossPlaneApp, error) {
	baseConfigureApp, err := NewConfigureApp()
	if err != nil {
		return nil, err
	}

	pluginConfig, err := readCrossPlanePluginConfig(baseConfigureApp.config.CrossPlanePluginConfig)
	if err != nil {
		return nil, err
	}
	return &CrossPlaneApp{pluginConfig: pluginConfig, ConfigureApp: baseConfigureApp}, err
}

func (cp *CrossPlaneApp) updateConfigs(params *model.CrossplaneUseCase, templateDir, reqRepo string) error {
	providerConfigPath := cp.pluginConfig[crossPlaneProviderConfigPathAttribute]
	err := cp.createProviderConfigs(filepath.Join(templateDir, providerConfigPath), params)
	if err != nil {
		return fmt.Errorf("failed to create provider config, %v", err)
	}

	synchConfigPath := cp.pluginConfig[crossPlaneSyncjConfigPathAttribute]
	err = copy.Copy(filepath.Join(templateDir, synchConfigPath), filepath.Join(reqRepo, synchConfigPath),
		copy.Options{
			OnDirExists: func(src, dest string) copy.DirExistsAction {
				return copy.Replace
			}})
	if err != nil {
		return fmt.Errorf("failed to copy dir from template to user repo, %v", err)
	}
	return nil
}

func (cp *CrossPlaneApp) ExecuteSteps(ctx context.Context, params model.ConfigureParameters, payload json.RawMessage) (model.ResponsePayload, error) {
	req := &model.CrossplaneUseCase{}
	err := json.Unmarshal(payload, req)
	if err != nil {
		respPayload := model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed), Message: json.RawMessage("{\"error\": \"requested payload is wrong\"}")}
		return respPayload, err
	}

	accessToken, err := cp.getAccessToken(ctx, req.VaultCredIdentifier)
	if err != nil {
		respPayload := model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed), Message: json.RawMessage("{\"error\": \"failed to get token from vault\"}")}
		return respPayload, nil
	}

	templateRepo, customerRepo, err := cp.cloneRepos(ctx, cp.pluginConfig[crossPlaneGitRepoAttribute], req.RepoURL, accessToken)
	if err != nil {
		respPayload := model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed), Message: json.RawMessage("{\"error\": \"failed to clone repos\"}")}
		return respPayload, err
	}

	defer os.RemoveAll(templateRepo)
	defer os.RemoveAll(customerRepo)

	err = cp.updateConfigs(req, templateRepo, customerRepo)
	if err != nil {
		respPayload := model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed), Message: json.RawMessage("{\"error\": \"failed to update configs to repo\"}")}
		return respPayload, err
	}

	err = cp.addToGit(ctx, req.Type, req.RepoURL, accessToken, req.PushToDefaultBranch)
	if err != nil {
		respPayload := model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed), Message: json.RawMessage("{\"error\": \"failed to add git repo\"}")}
		return respPayload, err
	}

	if !req.PushToDefaultBranch {
		logger.Info("requested to create PR.. skipping the further steps")
		return model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusCompleted)}, nil
	}

	ns, resName, err := cp.deployMainApp(ctx, filepath.Join(customerRepo, cp.pluginConfig[crossPlaneProviderConfigMainAppAttribute]))
	if err != nil {
		respPayload := model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed), Message: json.RawMessage("{\"error\": \"failed to deploy main app\"}")}
		return respPayload, err
	}

	// force Sync the app and monitor the status of the app.
	// then wait for cluster-claims.

	err = cp.syncArgoCDApp(ctx, ns, resName)
	if err != nil {
		respPayload := model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed), Message: json.RawMessage("{\"error\": \"failed to sync argocd app\"}")}
		return respPayload, err
	}

	err = cp.waitForArgoCDToSync(ctx, ns, resName)
	if err != nil {
		respPayload := model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed), Message: json.RawMessage("{\"error\": \"failed to fetch argocd app\"}")}
		return respPayload, err
	}

	return model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusCompleted)}, nil
}

func (cp *CrossPlaneApp) createProviderConfigs(dir string, params *model.CrossplaneUseCase) error {
	logger.Infof("processing %d crossplane providers to generate provider config", len(params.CrossplaneProviders))
	for _, provider := range params.CrossplaneProviders {
		providerName := strings.ToLower(provider.ProviderName)
		providerFile := filepath.Join(dir, fmt.Sprintf("%s-provider.yaml", providerName))
		dir := filepath.Dir(providerFile)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create dir %s, %v", dir, err)
		}

		file, err := os.Create(providerFile)
		if err != nil {
			return fmt.Errorf("failed to create file %s, %v", providerFile, err)
		}
		defer file.Close()

		providerConfigString, err := cp.createProviderConfigResource(provider, params)
		if err != nil {
			return fmt.Errorf("failed prepare provider %s config: %v", providerName, err)
		}

		if _, err := file.WriteString(providerConfigString); err != nil {
			return fmt.Errorf("failed to write provider %s config to %s, %v", providerName, providerFile, err)
		}
		logger.Infof("crossplane provider %s config written to %s", providerName, providerFile)
	}
	return nil
}

func (cp *CrossPlaneApp) createProviderConfigResource(provider agentmodel.CrossplaneProvider, params *model.CrossplaneUseCase) (string, error) {
	cloudType := strings.ToLower(provider.CloudType)
	providerName := strings.ToLower(provider.ProviderName)
	packageAttribute := fmt.Sprintf("%s_package", cloudType)
	pkg, found := cp.pluginConfig[packageAttribute]
	if !found {
		return "", fmt.Errorf("plugin package attribute %s not found", packageAttribute)
	}

	secretPath := fmt.Sprintf("generic/CloudProvider/%s", provider.CloudProviderId)
	providerConfigString := fmt.Sprintf(
		crossplaneProviderTemplate,
		providerName, secretPath, secretPath,
		providerName, pkg, providerName,
	)
	return providerConfigString, nil
}
