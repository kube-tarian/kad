package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	agentmodel "github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/model"
	"github.com/otiai10/copy"
)

type CrossPlaneApp struct {
	*ConfigureApp
	data crossplanePluginDS
}

func NewCrossPlaneApp() (*CrossPlaneApp, error) {
	baseConfigureApp, err := NewConfigureApp()
	if err != nil {
		return nil, err
	}

	cpluginInfo, err := ReadCrossPlanePluginConfig(baseConfigureApp.config.CrossPlanePluginConfig)
	if err != nil {
		return nil, err
	}
	return &CrossPlaneApp{data: cpluginInfo, ConfigureApp: baseConfigureApp}, err
}

func (pc *CrossPlaneApp) GetConfigMainApp() string {
	return pc.data[ConfigMainApp]
}

func (pc *CrossPlaneApp) GetGitRepo() string {
	return pc.data[GitRepo]
}

func (pc *CrossPlaneApp) GetGitConfigPath() string {
	return pc.data[GitConfigPath]
}

func (pc *CrossPlaneApp) GetMap() map[string]string {
	return map[string]string{}
}

func (cp *CrossPlaneApp) updateConfigs(params *model.CrossplaneUseCase, templateDir, reqRepo string) error {
	err := createProviderConfigs(filepath.Join(templateDir, cp.GetGitConfigPath()), params,
		cp.GetMap())
	if err != nil {
		return fmt.Errorf("failed to create provider config, %v", err)
	}

	err = copy.Copy(filepath.Join(templateDir, cp.GetGitConfigPath()), filepath.Join(reqRepo, cp.GetGitConfigPath()),
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

	templateRepo, customerRepo, err := cp.cloneRepos(ctx, cp.GetGitRepo(), req.RepoURL, accessToken)
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

	err = cp.deployMainApp(ctx, filepath.Join(customerRepo, cp.GetConfigMainApp()))
	if err != nil {
		respPayload := model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed), Message: json.RawMessage("{\"error\": \"failed to deploy main app\"}")}
		return respPayload, err
	}

	return model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusCompleted)}, nil
}

func createProviderConfigs(dir string, params *model.CrossplaneUseCase, pluginMap map[string]string) error {
	for _, provider := range params.CrossplaneProviders {
		cloudType := provider.CloudType
		providerFile := filepath.Join(dir, fmt.Sprintf("%s-provider.yaml", cloudType))
		dir := filepath.Dir(providerFile)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("err while creating directories: %v", dir)
		}

		file, err := os.Create(providerFile)
		if err != nil {
			return fmt.Errorf("err while creating file for provider: %v", err)
		}

		providerConfigString, err := createProviderCrdString(provider, params, pluginMap)
		if err != nil {
			return fmt.Errorf("createProviderConfigs: err createProviderCrdString: %v", err)
		}

		if _, err := file.WriteString(providerConfigString); err != nil {
			return fmt.Errorf("err while writing to controllerconfig: %v", err)
		}

		file.Close()
	}
	return nil
}

func createProviderCrdString(provider agentmodel.CrossplaneProvider, params *model.CrossplaneUseCase, pluginMap map[string]string) (string, error) {
	cloudType := provider.CloudType
	pkg, found := pluginMap[fmt.Sprintf("%s_package", cloudType)]
	if !found {
		return "", fmt.Errorf("plugin package not found for cloudType: %s", cloudType)
	}

	secretPath := fmt.Sprintf("generic/CloudProvider/%s", provider.CloudProviderId)
	providerConfigString := fmt.Sprintf(
		crossplaneProviderTemplate,
		cloudType, secretPath, secretPath,
		cloudType, pkg, cloudType,
	)
	return providerConfigString, nil
}
