package tekton

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	appconfig "github.com/kube-tarian/kad/capten/config-worker/internal/app_config"
	"github.com/kube-tarian/kad/capten/config-worker/internal/crossplane"
	"github.com/kube-tarian/kad/capten/model"
	agentmodel "github.com/kube-tarian/kad/capten/model"
	"github.com/otiai10/copy"
	"github.com/pkg/errors"
)

var (
	gitCred                 = "gitcred"
	dockerCred              = "docker-credentials"
	githubWebhook           = "github-webhook-secret"
	argoCred                = "argocd"
	crossplaneProjectConfig = "extraconfig"
	cosignDockerSecret      = "cosign-docker-secret"
	secrets                 = []string{gitCred, dockerCred, githubWebhook, argoCred, crossplaneProjectConfig, cosignDockerSecret}
	pipelineNamespace       = "tekton-pipelines"
	tektonChildTasks        = []string{"tekton-cluster-tasks", "tekton-pipelines"}
	addPipeline             = "add"
	deletePipeline          = "delete"
	cosignEntityName        = "cosign"
	cosignVaultId           = "signer"
	cosignSecName           = "cosign-keys"
)

type Config struct {
	PluginConfigFile     string `envconfig:"TEKTON_PLUGIN_CONFIG_FILE" default:"/tekton_plugin_config.json"`
	CrossplaneConfigFile string `envconfig:"CROSSPLANE_PLUGIN_CONFIG_FILE" default:"/crossplane_plugin_config.json"`
	DomainName           string `envconfig:"DOMAIN_NAME" default:"capten"`
}

type TektonApp struct {
	helper          *appconfig.AppGitConfigHelper
	pluginConfig    *tektonPluginConfig
	crossplanConfig *crossplane.CrossplanePluginConfig
	cfg             Config
}

func NewTektonApp() (*TektonApp, error) {
	cfg := Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	helper, err := appconfig.NewAppGitConfigHelper()
	if err != nil {
		return nil, err
	}

	pluginConfig, err := readTektonPluginConfig(cfg.PluginConfigFile)
	if err != nil {
		return nil, err
	}

	crossplaneConfig, err := crossplane.ReadCrossPlanePluginConfig(cfg.CrossplaneConfigFile)
	if err != nil {
		return nil, err
	}
	return &TektonApp{pluginConfig: pluginConfig, helper: helper, cfg: cfg, crossplanConfig: crossplaneConfig}, err
}

func readTektonPluginConfig(pluginFile string) (*tektonPluginConfig, error) {
	data, err := os.ReadFile(filepath.Clean(pluginFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read pluginConfig File: %s, err: %w", pluginFile, err)
	}

	var pluginData tektonPluginConfig
	err = json.Unmarshal(data, &pluginData)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &pluginData, nil
}

func (cp *TektonApp) configureProjectAndApps(ctx context.Context, req *model.TektonProjectSyncUsecase) (status string, err error) {
	logger.Infof("cloning default templates %s to project %s", cp.pluginConfig.TemplateGitRepo, req.RepoURL)
	templateRepo, err := cp.helper.CloneTemplateRepo(ctx, cp.pluginConfig.TemplateGitRepo, req.VaultCredIdentifier)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to clone repos")
	}
	defer os.RemoveAll(templateRepo)

	customerRepo, err := cp.helper.CloneUserRepo(ctx, req.RepoURL, req.VaultCredIdentifier)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to clone repos")
	}
	logger.Infof("cloned default templates to project %s", req.RepoURL)
	defer os.RemoveAll(customerRepo)

	err = cp.synchTektonConfig(req, templateRepo, customerRepo)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to update configs to repo")
	}
	logger.Infof("added provider config resources to cloned project %s", req.RepoURL)

	// update git project url
	if err := replaceCaptenUrls(customerRepo, cp.pluginConfig.TemplateGitRepo, req.RepoURL); err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to replace template url")
	}
	logger.Infof("updated resource configurations in cloned project %s", req.RepoURL)

	err = updateArgoCDTemplate(filepath.Join(customerRepo, cp.pluginConfig.PipelineSyncUpdate.MainAppValues))
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to updateArgoCDTemplate")
	}

	err = cp.helper.AddFilesToRepo([]string{"."})
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to add git repo")
	}

	err = cp.helper.CommitRepoChanges()
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to commit git repo")
	}
	logger.Infof("added cloned project %s changed to git", req.RepoURL)

	err = cp.deployMainApps(ctx, customerRepo)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to depoy argoCD apps")
	}

	return string(agentmodel.WorkFlowStatusCompleted), nil
}

func (cp *TektonApp) synchTektonConfig(req *model.TektonProjectSyncUsecase, templateDir, reqRepo string) error {
	for _, config := range []string{cp.pluginConfig.TektonProject, filepath.Join(cp.pluginConfig.PipelineClusterConfigSyncPath)} {
		err := copy.Copy(filepath.Join(templateDir, config), filepath.Join(reqRepo, config),
			copy.Options{
				OnDirExists: func(src, dest string) copy.DirExistsAction {
					return copy.Replace
				}})
		if err != nil {
			return fmt.Errorf("failed to copy dir from template to user repo, %v", err)
		}
	}

	// Copy pipeline template config
	err := copy.Copy(filepath.Join(templateDir, cp.pluginConfig.TektonPipelinePath),
		filepath.Join(reqRepo, cp.pluginConfig.TektonPipelinePath),
		copy.Options{
			OnDirExists: func(src, dest string) copy.DirExistsAction {
				return copy.Replace
			}})
	if err != nil {
		return fmt.Errorf("failed to copy dir from template to user repo, %v", err)
	}

	return nil
}

func (cp *TektonApp) deployMainApps(ctx context.Context, customerRepo string) (err error) {
	logger.Infof("%d main apps to deploy", len(cp.pluginConfig.ArgoCDApps))

	for _, tektonArgoApp := range cp.pluginConfig.ArgoCDApps {
		appPath := filepath.Join(customerRepo, tektonArgoApp.MainAppGitPath)
		err = cp.deployArgoCDApp(ctx, appPath, tektonChildTasks, tektonArgoApp.SynchApp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cp *TektonApp) deployArgoCDApp(ctx context.Context, appPath string, childApps []string, synchApp bool) (err error) {
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

func (cp *TektonApp) syncArgoCDChildApps(ctx context.Context, namespace string, apps []string) error {
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

func updateArgoCDTemplate(valuesFileName string) error {
	data, err := os.ReadFile(valuesFileName)
	if err != nil {
		return err
	}

	jsonData, err := k8s.ConvertYamlToJson(data)
	if err != nil {
		return err
	}

	var tektonConfig TektonConfigValues
	err = json.Unmarshal(jsonData, &tektonConfig)
	if err != nil {
		return err
	}

	jsonBytes, err := json.Marshal(tektonConfig)
	if err != nil {
		return err
	}

	yamlBytes, err := k8s.ConvertJsonToYaml(jsonBytes)
	if err != nil {
		return err
	}

	err = os.WriteFile(valuesFileName, yamlBytes, os.ModeAppend)
	return err
}
