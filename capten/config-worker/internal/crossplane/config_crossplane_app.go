package crossplane

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/intelops/go-common/credentials"
	"github.com/kelseyhightower/envconfig"
	appconfig "github.com/kube-tarian/kad/capten/config-worker/internal/app_config"
	"github.com/kube-tarian/kad/capten/model"
	agentmodel "github.com/kube-tarian/kad/capten/model"
	"github.com/otiai10/copy"
	"github.com/pkg/errors"
)

type Config struct {
	PluginConfigFile        string `envconfig:"CROSSPLANE_PLUGIN_CONFIG_FILE" default:"/crossplane_plugin_config.json"`
	CloudProviderEntityName string `envconfig:"CLOUD_PROVIDER_ENTITY_NAME" default:"cloud-provider"`
	DomainName              string `envconfig:"DOMAIN_NAME" default:"capten"`
}

type CrossPlaneApp struct {
	helper       *appconfig.AppGitConfigHelper
	pluginConfig *CrossplanePluginConfig
	cfg          Config
}

func NewCrossPlaneApp() (*CrossPlaneApp, error) {
	cfg := Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	helper, err := appconfig.NewAppGitConfigHelper()
	if err != nil {
		return nil, err
	}

	pluginConfig, err := ReadCrossPlanePluginConfig(cfg.PluginConfigFile)
	if err != nil {
		return nil, err
	}
	return &CrossPlaneApp{pluginConfig: pluginConfig, helper: helper, cfg: cfg}, err
}

func ReadCrossPlanePluginConfig(pluginFile string) (*CrossplanePluginConfig, error) {
	data, err := os.ReadFile(filepath.Clean(pluginFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read pluginConfig File: %s, err: %w", pluginFile, err)
	}

	var pluginData CrossplanePluginConfig
	err = json.Unmarshal(data, &pluginData)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &pluginData, nil
}

func (cp *CrossPlaneApp) configureProjectAndApps(ctx context.Context, req *model.CrossplaneUseCase) (status string, err error) {
	logger.Infof("cloning default templates %s to project %s", cp.pluginConfig.TemplateGitRepo, req.RepoURL)

	customerRepo, err := cp.helper.CloneUserRepo(ctx, req.RepoURL, req.VaultCredIdentifier)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessagef(err, "failed to clone repo %s", req.RepoURL)
	}
	defer os.RemoveAll(customerRepo)

	logger.Infof("cloned project %s", req.RepoURL)
	templateRepo, err := cp.helper.CloneTemplateRepo(ctx, cp.pluginConfig.TemplateGitRepo, req.VaultCredIdentifier)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessagef(err, "failed to clone repo %s", cp.pluginConfig.TemplateGitRepo)
	}
	defer os.RemoveAll(templateRepo)

	err = cp.synchProviders(req, templateRepo, customerRepo)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to update configs to repo")
	}
	logger.Infof("added provider config resources to cloned project %s", req.RepoURL)

	// update git project url
	if err := replaceCaptenUrls(customerRepo, cp.pluginConfig.TemplateGitRepo, req.RepoURL); err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to replace template url")
	}
	logger.Infof("updated resource configurations in cloned project %s", req.RepoURL)

	err = cp.helper.AddFilesToRepo([]string{"."})
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to add git repo")
	}

	err = cp.helper.CommitRepoChanges()
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to commit git repo")
	}
	logger.Infof("added cloned project %s changed to git", req.RepoURL)

	err = cp.deployArgoCDApps(ctx, customerRepo)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to depoy argoCD apps")
	}

	return string(agentmodel.WorkFlowStatusCompleted), nil
}

func (cp *CrossPlaneApp) synchProviders(req *model.CrossplaneUseCase, templateDir, reqRepo string) error {

	if err := cp.deleteProviderConfigs(reqRepo, req); err != nil {
		return fmt.Errorf("failed to delete provider configs, %v", err)
	}

	err := cp.createProviderConfigs(filepath.Join(templateDir, cp.pluginConfig.ProviderConfigSyncPath), req)
	if err != nil {
		return fmt.Errorf("failed to create provider config, %v", err)
	}

	err = copy.Copy(filepath.Join(templateDir, cp.pluginConfig.CrossplaneConfigSyncPath),
		filepath.Join(reqRepo, cp.pluginConfig.CrossplaneConfigSyncPath),
		copy.Options{
			OnDirExists: func(src, dest string) copy.DirExistsAction {
				return copy.Replace
			}})
	if err != nil {
		return fmt.Errorf("failed to copy dir from template to user repo, %v", err)
	}

	return nil
}

func (cp *CrossPlaneApp) deployArgoCDApps(ctx context.Context, customerRepo string) (err error) {
	logger.Infof("%d main apps to deploy", len(cp.pluginConfig.ArgoCDApps))

	for _, argoApp := range cp.pluginConfig.ArgoCDApps {
		appPath := filepath.Join(customerRepo, argoApp.MainAppGitPath)
		err = cp.deployArgoCDApp(ctx, appPath, argoApp.ChildAppNames, argoApp.SynchApp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cp *CrossPlaneApp) deployArgoCDApp(ctx context.Context, appPath string, childApps []string, synchApp bool) (err error) {
	ns, resName, err := cp.helper.DeployMainApp(ctx, appPath)
	if err != nil {
		return errors.WithMessage(err, "failed to deploy main app")
	}
	logger.Infof("deployed provider config main-app %s", resName)

	if synchApp {
		err = cp.helper.SyncArgoCDApp(ctx, ns, resName)
		if err != nil {
			return errors.WithMessage(err, "failed to sync argocd app")
		}
		logger.Infof("synched provider config main-app %s", resName)

		err = cp.helper.WaitForArgoCDToSync(ctx, ns, resName)
		if err != nil {
			return errors.WithMessage(err, "failed to fetch argocd app")
		}

		err = cp.syncArgoCDChildApps(ctx, ns, childApps)
		if err != nil {
			return errors.WithMessage(err, "failed to synch argocd child app")
		}
		logger.Infof("synched provider config child apps")
	}
	return nil
}

func (cp *CrossPlaneApp) syncArgoCDChildApps(ctx context.Context, namespace string, apps []string) error {
	for _, appName := range apps {
		err := cp.helper.SyncArgoCDApp(ctx, namespace, appName)
		if err != nil {
			return fmt.Errorf("failed to sync app %s, %v", appName, err)
		}
		logger.Infof("synched provider config child-app %s", appName)

		err = cp.helper.WaitForArgoCDToSync(ctx, namespace, appName)
		if err != nil {
			return fmt.Errorf("failed to get sync status of app %s, %v", appName, err)
		}
	}
	return nil
}

func (cp *CrossPlaneApp) createProviderConfigs(dir string, req *model.CrossplaneUseCase) error {
	logger.Infof("processing %d crossplane providers to generate provider config", len(req.CrossplaneProviders))
	for _, provider := range req.CrossplaneProviders {
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

		providerConfigString, err := cp.createProviderConfigResource(provider, req)
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

func (cp *CrossPlaneApp) deleteProviderConfigs(reqRepoDir string, req *model.CrossplaneUseCase) error {
	logger.Infof("processing %d crossplane providers to sync & delete provider configs", len(req.CrossplaneProviders))

	providerMap := map[string]agentmodel.CrossplaneProvider{}
	for _, provider := range req.CrossplaneProviders {
		providerMap[strings.ToLower(provider.ProviderName)] = provider
	}

	fi, err := os.ReadDir(filepath.Join(reqRepoDir, cp.pluginConfig.ProviderConfigSyncPath))
	if err != nil {
		return err
	}

	for _, file := range fi {
		provider := strings.TrimSuffix(file.Name(), "-provider.yaml")

		_, ok := providerMap[provider]
		if !ok {
			providerRepoFilePath := filepath.Join(reqRepoDir, cp.pluginConfig.ProviderConfigSyncPath, fmt.Sprintf("%s-provider.yaml", provider))
			logger.Infof("removing the provider '%s' from git repo path %s", provider, providerRepoFilePath)

			if err := os.Remove(providerRepoFilePath); err != nil {
				logger.Errorf("failed to remove from the file", err)
				continue
			}

			fileToDeleteInRepoPath := filepath.Join(".", cp.pluginConfig.ProviderConfigSyncPath, fmt.Sprintf("%s-provider.yaml", provider))
			if err = cp.helper.RemoveFilesFromRepo([]string{fileToDeleteInRepoPath}); err != nil {
				logger.Errorf("failed to remove from git repo", err)
				continue
			}
		}
	}

	return nil
}

func (cp *CrossPlaneApp) createProviderConfigResource(provider agentmodel.CrossplaneProvider, req *model.CrossplaneUseCase) (string, error) {
	cloudType := strings.ToLower(provider.CloudType)
	pkg, found := cp.pluginConfig.ProviderPackages[cloudType]
	if !found {
		return "", fmt.Errorf("plugin package not found")
	}

	secretPath := fmt.Sprintf("%s/%s/%s", credentials.GenericCredentialType, cp.cfg.CloudProviderEntityName, provider.CloudProviderId)

	switch provider.CloudType {
	case "AWS":
		providerConfigString := fmt.Sprintf(
			crossplaneAWSProviderTemplate,
			cloudType, secretPath, secretPath,
			cloudType, pkg, cloudType,
		)
		return providerConfigString, nil
	case "GCP":
		providerConfigString := fmt.Sprintf(
			crossplaneGCPProviderTemplate,
			cloudType, secretPath, secretPath,
			cloudType, pkg, cloudType,
		)
		return providerConfigString, nil
	case "AZURE":
		providerConfigString := fmt.Sprintf(
			crossplaneAzureProviderTemplate,
			cloudType, secretPath, secretPath,
			cloudType, pkg, cloudType,
		)
		return providerConfigString, nil
	default:
		return "", fmt.Errorf("cloud type %s not supported", provider.CloudType)
	}

}

func replaceCaptenUrls(dir string, src, target string) error {
	if !strings.HasSuffix(src, ".git") {
		src += ".git"
	}

	if !strings.HasSuffix(target, ".git") {
		target += ".git"
	}

	fileList := []string{}
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".yaml") {
			fileList = append(fileList, path)
		}
		return nil
	}); err != nil {
		return err
	}

	for _, filePath := range fileList {
		err := replaceInFile(filePath, src, target)
		if err != nil {
			logger.Errorf("Error replacing in %s: %v\n", filePath, err)
		}
	}
	return nil
}

func replaceInFile(filePath, target, replacement string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	newData := strings.Replace(string(data), target, replacement, -1)
	err = os.WriteFile(filePath, []byte(newData), 0644)
	if err != nil {
		return err
	}
	return nil
}
