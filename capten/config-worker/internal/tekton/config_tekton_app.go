package tekton

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"github.com/kube-tarian/kad/capten/common-pkg/plugins/argocd"
	appconfig "github.com/kube-tarian/kad/capten/config-worker/internal/app_config"
	"github.com/kube-tarian/kad/capten/config-worker/internal/crossplane"
	"github.com/kube-tarian/kad/capten/model"
	agentmodel "github.com/kube-tarian/kad/capten/model"
	"github.com/otiai10/copy"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	gitCred           = "gitcred"
	dockerCred        = "docker-credentials"
	githubWebhook     = "github-webhook-secret"
	argoCred          = "argocd"
	extraConfig       = "extraconfig"
	secrets           = []string{gitCred, dockerCred, githubWebhook, argoCred, extraConfig}
	pipelineNamespace = "tekton-pipelines"
	tektonChildTasks  = []string{"tekton-cluster-tasks"}
	addPipeline       = "add"
	deletePipeline    = "delete"
	mainAppName       = "tekton-apps"
	cosignEntityName  = "cosign"
	cosignVaultId     = "signer"
	cosignSecName     = "cosign-keys"
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

func (cp *TektonApp) configureProjectAndApps(ctx context.Context, req *model.TektonPipelineUseCase) (status string, err error) {
	logger.Infof("cloning default templates %s to project %s", cp.pluginConfig.TemplateGitRepo, req.RepoURL)
	templateRepo, err := cp.helper.CloneTemplateRepo(ctx, cp.pluginConfig.TemplateGitRepo, req.CredentialIdentifiers[agentmodel.Git].Id)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to clone repos")
	}
	defer os.RemoveAll(templateRepo)

	customerRepo, err := cp.helper.CloneUserRepo(ctx, req.RepoURL, req.CredentialIdentifiers[agentmodel.Git].Id)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to clone repos")
	}
	logger.Infof("cloned default templates to project %s", req.RepoURL)

	defer os.RemoveAll(customerRepo)

	err = cp.synchPipelineConfig(req, templateRepo, customerRepo)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to update configs to repo")
	}
	logger.Infof("added provider config resources to cloned project %s", req.RepoURL)

	// update git project url
	if err := replaceCaptenUrls(customerRepo, cp.pluginConfig.TemplateGitRepo, req.RepoURL); err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to replace template url")
	}
	logger.Infof("updated resource configurations in cloned project %s", req.RepoURL)

	err = updateArgoCDTemplate(filepath.Join(customerRepo, cp.pluginConfig.PipelineSyncUpdate.MainAppValues), req.PipelineName, addPipeline)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to updateArgoCDTemplate")
	}

	err = updatePipelineTemplate(filepath.Join(customerRepo,
		strings.ReplaceAll(cp.pluginConfig.PipelineSyncUpdate.PipelineValues, "<NAME>", req.PipelineName)), req.PipelineName, cp.cfg.DomainName)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to updatePipelineTemplate")
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

	err = cp.createOrUpdateSecrets(ctx, req)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to create k8s secrets")
	}

	err = cp.deployArgoCDApps(ctx, customerRepo, req.PipelineName)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to depoy argoCD apps")
	}

	return string(agentmodel.WorkFlowStatusCompleted), nil
}

func (cp *TektonApp) deleteProjectAndApps(ctx context.Context, req *model.TektonPipelineUseCase) (status string, err error) {
	logger.Infof("cloning user repo %s", req.RepoURL)
	customerRepo, err := cp.helper.CloneUserRepo(ctx, req.RepoURL, req.CredentialIdentifiers[agentmodel.Git].Id)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to clone repos")
	}

	defer os.RemoveAll(customerRepo)

	logger.Infof("removing pipeline directory from %s", req.RepoURL)
	err = cp.helper.RemoveFilesFromRepo([]string{filepath.Join(cp.pluginConfig.TektonPipelinePath, req.PipelineName)})
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to remove pipeline from repo")
	}
	logger.Infof("removed pipeline resources from project %s", req.RepoURL)

	logger.Infof("update main resource values.yaml in cloned project %s", req.RepoURL)

	err = updateArgoCDTemplate(filepath.Join(customerRepo, cp.pluginConfig.PipelineSyncUpdate.MainAppValues), req.PipelineName, deletePipeline)
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

	err = cp.deleteSecrets(ctx, req)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to create k8s secrets")
	}

	// err = cp.helper.DeleteArgoCDApp(ctx, pipelineNamespace, req.PipelineName, mainAppName)
	// if err != nil {
	// 	return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to delete argoCD apps")
	// }

	return string(agentmodel.WorkFlowStatusCompleted), nil
}

func (cp *TektonApp) synchPipelineConfig(req *model.TektonPipelineUseCase, templateDir, reqRepo string) error {
	if _, err := os.Stat(filepath.Join(reqRepo, cp.pluginConfig.TektonProject)); err != nil {
		for _, config := range []string{cp.pluginConfig.TektonProject, filepath.Join(cp.pluginConfig.TektonPipelinePath, cp.pluginConfig.PipelineClusterConfigSyncPath)} {
			err := copy.Copy(filepath.Join(templateDir, config), filepath.Join(reqRepo, config),
				copy.Options{
					OnDirExists: func(src, dest string) copy.DirExistsAction {
						return copy.Replace
					}})
			if err != nil {
				return fmt.Errorf("failed to copy dir from template to user repo, %v", err)
			}
		}
	}

	// Copy pipeline specific config
	err := copy.Copy(filepath.Join(templateDir, cp.pluginConfig.TektonPipelinePath, cp.pluginConfig.PipelineConfigSyncPath),
		filepath.Join(reqRepo, cp.pluginConfig.TektonPipelinePath, req.PipelineName),
		copy.Options{
			OnDirExists: func(src, dest string) copy.DirExistsAction {
				return copy.Replace
			}})
	if err != nil {
		return fmt.Errorf("failed to copy dir from template to user repo, %v", err)
	}

	return nil
}

func (cp *TektonApp) deployArgoCDApps(ctx context.Context, customerRepo, pipelineName string) (err error) {
	logger.Infof("%d main apps to deploy", len(cp.pluginConfig.ArgoCDApps))

	for _, argoApp := range cp.pluginConfig.ArgoCDApps {
		appPath := filepath.Join(customerRepo, argoApp.MainAppGitPath)
		err = cp.deployArgoCDApp(ctx, appPath, append(tektonChildTasks, pipelineName), argoApp.SynchApp)
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

func (cp *TektonApp) createOrUpdateSecrets(ctx context.Context, req *model.TektonPipelineUseCase) error {
	log := logging.NewLogger()
	k8sclient, err := k8s.NewK8SClient(log)
	if err != nil {
		return fmt.Errorf("failed to initalize k8s client, %v", err)
	}

	k8sclient.Clientset.CoreV1().Namespaces().Create(ctx, &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: pipelineNamespace}}, metav1.CreateOptions{})

	// One time activity
	key, pub, err := cp.helper.GetCosingKeys(ctx, cosignEntityName, cosignVaultId)
	if err != nil {
		return fmt.Errorf("failed to get cosign keys from vault, %v", err)
	}

	if err := k8sclient.CreateOrUpdateSecret(ctx, pipelineNamespace, cosignSecName,
		v1.SecretTypeOpaque, map[string][]byte{appconfig.CosignKey: []byte(key),
			appconfig.CosignPub: []byte(pub)}, map[string]string{}); err != nil {
		return fmt.Errorf("failed to create/update cosign-keys k8s secret, %v", err)
	}

	for _, secret := range secrets {
		strdata := make(map[string][]byte)
		secName := secret + "-" + req.PipelineName
		switch secret {
		case dockerCred:
			username, password, err := cp.helper.GetContainerRegCreds(ctx,
				req.CredentialIdentifiers[agentmodel.Container].Identifier, req.CredentialIdentifiers[agentmodel.Container].Id)
			if err != nil {
				return fmt.Errorf("failed to get docker cfg secret, %v", err)
			}

			data, err := handleDockerCfgJSONContent(username, password, req.CredentialIdentifiers[agentmodel.Container].Url)
			if err != nil {
				return fmt.Errorf("failed to get docker cfg secret, %v", err)
			}
			strdata[".dockerconfigjson"] = data
			strdata["config.json"] = data
			if err := k8sclient.CreateOrUpdateSecret(ctx, pipelineNamespace, secName,
				v1.SecretTypeDockerConfigJson, strdata, map[string]string{}); err != nil {
				return fmt.Errorf("failed to create/update k8s secret, %v", err)
			}

		case gitCred, githubWebhook:
			username, token, err := cp.helper.GetGitCreds(ctx, req.CredentialIdentifiers[agentmodel.Git].Id)
			if err != nil {
				return fmt.Errorf("failed to get git secret, %v", err)
			}
			strdata["username"] = []byte(username)
			strdata["password"] = []byte(token)
			if err := k8sclient.CreateOrUpdateSecret(ctx, pipelineNamespace, secName,
				v1.SecretTypeBasicAuth, strdata, nil); err != nil {
				return fmt.Errorf("failed to create/update k8s secret, %v", err)
			}
		case argoCred:
			cfg, err := argocd.GetConfig(log)
			if err != nil {
				return fmt.Errorf("failed to get argo-cd secret, %v", err)
			}
			strdata["SERVER_URL"] = []byte(cfg.ServiceURL)
			strdata["USERNAME"] = []byte(cfg.Username)
			strdata["PASSWORD"] = []byte(cfg.Password)
			if err := k8sclient.CreateOrUpdateSecret(ctx, pipelineNamespace, secName,
				v1.SecretTypeOpaque, strdata, map[string]string{}); err != nil {
				return fmt.Errorf("failed to create/update k8s secret, %v", err)
			}
		case extraConfig:
			username, token, err := cp.helper.GetGitCreds(ctx, req.CredentialIdentifiers[agentmodel.ExtraGitProject].Id)
			if err != nil {
				return fmt.Errorf("failed to get git secret, %v", err)
			}

			kubeConfig, kubeCa, kubeEndpoint, err := cp.helper.GetClusterCreds(ctx, req.CredentialIdentifiers[agentmodel.ManagedCluster].Identifier, req.CredentialIdentifiers[agentmodel.ManagedCluster].Id)
			if err != nil {
				return fmt.Errorf("failed to get GetClusterCreds, %v", err)
			}
			strdata["GIT_USER_NAME"] = []byte(username)
			strdata["GIT_TOKEN"] = []byte(token)
			strdata["GIT_PROJECT_URL"] = []byte(req.CredentialIdentifiers[agentmodel.ExtraGitProject].Url)
			strdata["APP_CONFIG_PATH"] = []byte(filepath.Join(cp.crossplanConfig.ClusterEndpointUpdates.DefaultAppValuesPath, req.CredentialIdentifiers[agentmodel.ManagedCluster].Url))
			strdata["CLUSTER_CA"] = []byte(kubeCa)
			strdata["CLUSTER_ENDPOINT"] = []byte(kubeEndpoint)
			strdata["CLUSTER_CONFIG"] = []byte(kubeConfig)

			if err := k8sclient.CreateOrUpdateSecret(ctx, pipelineNamespace, secName,
				v1.SecretTypeOpaque, strdata, nil); err != nil {
				return fmt.Errorf("failed to create/update k8s secret, %v", err)
			}

		default:
			return fmt.Errorf("secret step: %s type not found", secret)
		}
	}

	return nil
}

func (cp *TektonApp) deleteSecrets(ctx context.Context, req *model.TektonPipelineUseCase) error {
	k8sclient, err := k8s.NewK8SClient(logging.NewLogger())
	if err != nil {
		return fmt.Errorf("failed to initalize k8s client, %v", err)
	}

	for _, secret := range secrets {
		if err := k8sclient.DeleteSecret(ctx, pipelineNamespace, secret+"-"+req.PipelineName); err != nil {
			return err
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

func updateArgoCDTemplate(valuesFileName, pipelineName, action string) error {
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

	if tektonConfig.TektonPipelines == nil {
		tektonConfig.TektonPipelines = &[]TektonPipeline{}
	}

	switch action {
	case addPipeline:
		tektonPipelines := []TektonPipeline{{Name: pipelineName}}

		for _, pipeline := range *tektonConfig.TektonPipelines {
			tektonPipelines = append(tektonPipelines, TektonPipeline{Name: pipeline.Name})
		}

		tektonConfig.TektonPipelines = &tektonPipelines
	case deletePipeline:
		tektonPipelines := []TektonPipeline{}
		for _, pipeline := range *tektonConfig.TektonPipelines {
			if pipeline.Name == pipelineName {
				continue
			}
			tektonPipelines = append(tektonPipelines, TektonPipeline{Name: pipeline.Name})
		}

		tektonConfig.TektonPipelines = &tektonPipelines
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

func updatePipelineTemplate(valuesFileName, pipelineName, domainName string) error {
	data, err := os.ReadFile(valuesFileName)
	if err != nil {
		return err
	}

	jsonData, err := k8s.ConvertYamlToJson(data)
	if err != nil {
		return err
	}

	var tektonPipelineConfig TektonPieplineConfigValues
	err = json.Unmarshal(jsonData, &tektonPipelineConfig)
	if err != nil {
		return err
	}

	// GET dashboard and ingress domain suffix.
	tektonPipelineConfig.IngressDomainName = model.TektonHostName + "." + domainName
	tektonPipelineConfig.PipelineName = pipelineName
	tektonPipelineConfig.TektonDashboard = "http://" + tektonPipelineConfig.IngressDomainName
	secretName := []SecretNames{}

	for _, secret := range secrets {
		secretName = append(secretName, SecretNames{Name: secret + "-" + pipelineName})
	}

	tektonPipelineConfig.SecretName = &secretName

	jsonBytes, err := json.Marshal(tektonPipelineConfig)
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

// handleDockerCfgJSONContent serializes a ~/.docker/config.json file
func handleDockerCfgJSONContent(username, password, server string) ([]byte, error) {
	dockerConfigAuth := DockerConfigEntry{
		Username: username,
		Password: password,
		Auth:     encodeDockerConfigFieldAuth(username, password),
	}
	dockerConfigJSON := DockerConfigJSON{
		Auths: map[string]DockerConfigEntry{server: dockerConfigAuth},
	}

	return json.Marshal(dockerConfigJSON)
}

// encodeDockerConfigFieldAuth returns base64 encoding of the username and password string
func encodeDockerConfigFieldAuth(username, password string) string {
	fieldValue := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(fieldValue))
}
