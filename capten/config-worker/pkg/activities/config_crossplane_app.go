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

	crossPlaneGitRepoAttribute                 = "git_repo"
	crossPlaneSyncjConfigPathAttribute         = "synch_config_path"
	crossPlaneProviderConfigPathAttribute      = "provider_config_path"
	crossPlaneProviderConfigMainAppAttribute   = "provider_config_main_app"
	crossPlaneProviderConfigChildAppsAttribute = "provider_config_child_apps"
	tmpCloneStr                                = "clone*"
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
	logger.Infof("cloned default templates to project %s", req.RepoURL)

	defer os.RemoveAll(templateRepo)
	defer os.RemoveAll(customerRepo)

	err = cp.updateConfigs(req, templateRepo, customerRepo)
	if err != nil {
		respPayload := model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed), Message: json.RawMessage("{\"error\": \"failed to update configs to repo\"}")}
		return respPayload, err
	}
	logger.Infof("added provider config resources to cloned project %s", req.RepoURL)

	// update git project url
	if err := replaceCaptenUrls(customerRepo, req.RepoURL); err != nil {
		respPayload := model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed), Message: json.RawMessage("{\"error\": \"failed to replace template url\"}")}
		return respPayload, err
	}
	logger.Infof("updated resource configurations in cloned project %s", req.RepoURL)

	err = cp.addToGit(ctx, req.Type, req.RepoURL, accessToken, req.PushToDefaultBranch)
	if err != nil {
		respPayload := model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed), Message: json.RawMessage("{\"error\": \"failed to add git repo\"}")}
		return respPayload, err
	}
	logger.Infof("added cloned project %s changed to git", req.RepoURL)

	if !req.PushToDefaultBranch {
		logger.Info("requested to create PR.. skipping the further steps")
		return model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusCompleted)}, nil
	}

	appPath := filepath.Join(customerRepo, cp.pluginConfig[crossPlaneProviderConfigMainAppAttribute])
	childApps := strings.Split(cp.pluginConfig[crossPlaneProviderConfigChildAppsAttribute], ",")
	respPayload, err := cp.deployArgoCDAppAndSyncWithChilds(ctx, appPath, childApps)
	if err != nil {
		return respPayload, err
	}
	return model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusCompleted)}, nil
}

func (cp *CrossPlaneApp) deployArgoCDAppAndSyncWithChilds(ctx context.Context, appPath string, childApps []string) (model.ResponsePayload, error) {
	ns, resName, err := cp.deployMainApp(ctx, appPath)
	if err != nil {
		respPayload := model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed), Message: json.RawMessage("{\"error\": \"failed to deploy main app\"}")}
		return respPayload, err
	}
	logger.Infof("deployed provider config main-app %s", resName)

	err = cp.syncArgoCDApp(ctx, ns, resName)
	if err != nil {
		respPayload := model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed), Message: json.RawMessage("{\"error\": \"failed to sync argocd app\"}")}
		return respPayload, err
	}
	logger.Infof("synched provider config main-app %s", resName)

	err = cp.waitForArgoCDToSync(ctx, ns, resName)
	if err != nil {
		respPayload := model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed), Message: json.RawMessage("{\"error\": \"failed to fetch argocd app\"}")}
		return respPayload, err
	}

	err = cp.syncArgoCDChildApps(ctx, ns, childApps)
	if err != nil {
		respPayload := model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed), Message: json.RawMessage("{\"error\": \"failed to synch argocd child app\"}")}
		return respPayload, err
	}
	logger.Infof("synched provider config child apps")
	return model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusCompleted)}, nil
}

func (cp *CrossPlaneApp) syncArgoCDChildApps(ctx context.Context, namespace string, apps []string) error {
	for _, appName := range apps {
		err := cp.syncArgoCDApp(ctx, namespace, appName)
		if err != nil {
			return fmt.Errorf("failed to sync app %s, %v", appName, err)
		}
		logger.Infof("synched provider config child-app %s", appName)

		err = cp.waitForArgoCDToSync(ctx, namespace, appName)
		if err != nil {
			return fmt.Errorf("failed to get sync status of app %s, %v", appName, err)
		}
	}
	return nil
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

	secretPath := fmt.Sprintf("generic/cloud-provider/%s", provider.CloudProviderId)
	providerConfigString := fmt.Sprintf(
		crossplaneProviderTemplate,
		providerName, secretPath, secretPath,
		providerName, pkg, providerName,
	)
	return providerConfigString, nil
}

func replaceCaptenUrls(dir string, replacement string) error {

	target := "https://github.com/intelops/capten-templates.git"
	if strings.HasSuffix(replacement, ".git") {
		replacement += ".git"
	}

	// List all files in the directory
	fileList := []string{}

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".yaml") {
			fileList = append(fileList, path)
		}
		return nil
	}); err != nil {
		return err
	}

	// Replace the string in each file
	for _, filePath := range fileList {
		err := replaceInFile(filePath, target, replacement)
		if err != nil {
			fmt.Printf("Error replacing in %s: %v\n", filePath, err)
		}
	}

	return nil
}

func replaceInFile(filePath, target, replacement string) error {
	// Read the file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Perform the string replacement
	newData := strings.Replace(string(data), target, replacement, -1)

	// Write the modified content back to the file
	err = os.WriteFile(filePath, []byte(newData), 0644)
	if err != nil {
		return err
	}

	return nil
}
