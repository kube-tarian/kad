package pluginstore

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/server/pkg/agent"
	iamclient "github.com/kube-tarian/kad/server/pkg/iam-client"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/pb/captenpluginspb"
	"github.com/kube-tarian/kad/server/pkg/pb/clusterpluginspb"
	"github.com/kube-tarian/kad/server/pkg/pb/pluginstorepb"
	"github.com/kube-tarian/kad/server/pkg/store"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type PluginStore struct {
	log          logging.Logger
	cfg          *Config
	dbStore      store.ServerStore
	agentHandler *agent.AgentHandler
	iam          iamclient.IAMRegister
}

func NewPluginStore(log logging.Logger, dbStore store.ServerStore,
	agentHandler *agent.AgentHandler, iam iamclient.IAMRegister) (*PluginStore, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}

	return &PluginStore{
		log:          log,
		cfg:          cfg,
		dbStore:      dbStore,
		agentHandler: agentHandler,
		iam:          iam,
	}, nil
}

func (p *PluginStore) ConfigureStore(clusterId string, config *pluginstorepb.PluginStoreConfig) error {
	return p.dbStore.WritePluginStoreConfig(clusterId, config)
}

func (p *PluginStore) GetStoreConfig(clusterId string, storeType pluginstorepb.StoreType) (*pluginstorepb.PluginStoreConfig, error) {
	if storeType == pluginstorepb.StoreType_LOCAL_STORE {
		config, err := p.dbStore.ReadPluginStoreConfig(clusterId, storeType)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return &pluginstorepb.PluginStoreConfig{
					StoreType: pluginstorepb.StoreType_LOCAL_STORE,
				}, nil
			}
			return nil, err
		}
		return config, nil
	} else if storeType == pluginstorepb.StoreType_CENTRAL_STORE {
		return &pluginstorepb.PluginStoreConfig{
			StoreType:     pluginstorepb.StoreType_CENTRAL_STORE,
			GitProjectId:  p.cfg.PluginStoreProjectID,
			GitProjectURL: p.cfg.PluginStoreProjectURL,
		}, nil
	} else {
		return nil, fmt.Errorf("not supported store type")
	}
}

func (p *PluginStore) SyncPlugins(orgId, clusterId string, storeType pluginstorepb.StoreType) error {
	config, err := p.GetStoreConfig(clusterId, storeType)
	if err != nil {
		return err
	}

	pluginStoreDir, err := p.clonePluginStoreProject(orgId, clusterId, config.GitProjectURL, config.GitProjectId, storeType)
	if err != nil {
		return err
	}
	defer os.RemoveAll(pluginStoreDir)

	pluginListFilePath := p.getPluginListFilePath(pluginStoreDir)
	p.log.Infof("Loading plugin data from %s", pluginListFilePath)
	pluginListData, err := os.ReadFile(pluginListFilePath)
	if err != nil {
		return errors.WithMessage(err, "failed to read store config file")
	}

	var plugins PluginListData
	if err := yaml.Unmarshal(pluginListData, &plugins); err != nil {
		return errors.WithMessage(err, "failed to unmarshall store config file")
	}

	addedPlugins := map[string]bool{}
	for _, pluginName := range plugins.Plugins {
		err := p.addPluginApp(config.GitProjectId, pluginStoreDir, pluginName)
		if err != nil {
			p.log.Errorf("%v", err)
			continue
		}
		addedPlugins[pluginName] = true
		p.log.Infof("stored plugin data for plugin %s for cluster %s", pluginName, clusterId)
	}

	dbPlugins, err := p.dbStore.ReadPlugins(config.GitProjectId)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
	}

	for _, dbPlugin := range dbPlugins {
		if _, ok := addedPlugins[dbPlugin.PluginName]; !ok {
			if err = p.dbStore.DeletePlugin(config.GitProjectId, dbPlugin.PluginName); err != nil {
				p.log.Infof("failed to deleted plugin data for plugin %s for cluster %s",
					dbPlugin.PluginName, clusterId)
			}
			p.log.Infof("deleted plugin data for plugin %s for cluster %s", dbPlugin.PluginName, clusterId)
		}
	}

	return nil
}

func (p *PluginStore) clonePluginStoreProject(orgId, clusterId, projectURL, projectId string,
	storeType pluginstorepb.StoreType) (pluginStoreDir string, err error) {
	pluginStoreDir, err = os.MkdirTemp(p.cfg.PluginsStoreProjectMount, tmpGitProjectCloneStr)
	if err != nil {
		err = fmt.Errorf("failed to create plugin store tmp dir, err: %v", err)
		return
	}

	accessToken, err := p.getGitProjectAccessToken(orgId, clusterId, projectId, storeType)
	if err != nil {
		err = fmt.Errorf("failed to get git project credentias, %v", err)
		return
	}

	p.log.Infof("cloning plugin store project %s to %s", projectURL, pluginStoreDir)
	gitClient := NewGitClient()
	if err = gitClient.Clone(pluginStoreDir, projectURL, accessToken); err != nil {
		os.RemoveAll(pluginStoreDir)
		err = fmt.Errorf("failed to Clone plugin store project, err: %v", err)
		return
	}
	return
}

func (p *PluginStore) addPluginApp(gitProjectId, pluginStoreDir, pluginName string) error {
	appData, err := os.ReadFile(p.getPluginFilePath(pluginStoreDir, pluginName))
	if err != nil {
		return errors.WithMessagef(err, "failed to read store plugin %s", pluginName)
	}

	var pluginData Plugin
	if err := yaml.Unmarshal(appData, &pluginData); err != nil {
		return errors.WithMessagef(err, "failed to unmarshall store plugin %s", pluginName)
	}

	var iconData []byte
	if len(pluginData.Icon) != 0 {
		iconData, err = os.ReadFile(p.getPluginIconFilePath(pluginStoreDir, pluginName, pluginData.Icon))
		if err != nil {
			return errors.WithMessagef(err, "failed to read icon %s for plugin %s", pluginData.Icon, pluginName)
		}
	}

	if pluginData.PluginName == "" || len(pluginData.Versions) == 0 {
		return fmt.Errorf("app name/version is missing for %s", pluginName)
	}

	plugin := &pluginstorepb.PluginData{
		PluginName:  pluginData.PluginName,
		Description: pluginData.Description,
		Category:    pluginData.Category,
		Versions:    pluginData.Versions,
		Icon:        iconData,
	}

	if err := p.dbStore.WritePluginData(gitProjectId, plugin); err != nil {
		return errors.WithMessagef(err, "failed to store plugin %s", pluginName)
	}
	return nil
}

func (p *PluginStore) GetPlugins(clusterId string, storeType pluginstorepb.StoreType) ([]*pluginstorepb.Plugin, error) {
	config, err := p.GetStoreConfig(clusterId, storeType)
	if err != nil {
		return nil, err
	}

	plugins, err := p.dbStore.ReadPlugins(config.GitProjectId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return []*pluginstorepb.Plugin{}, nil
		}
		return nil, err
	}
	return plugins, nil
}

func (p *PluginStore) GetPluginData(clusterId string, storeType pluginstorepb.StoreType, pluginName string) (*pluginstorepb.PluginData, error) {
	config, err := p.GetStoreConfig(clusterId, storeType)
	if err != nil {
		return nil, err
	}

	return p.dbStore.ReadPluginData(config.GitProjectId, pluginName)
}

func (p *PluginStore) GetPluginValues(orgId, clusterId string, storeType pluginstorepb.StoreType,
	pluginName, version string) ([]byte, error) {
	config, err := p.GetStoreConfig(clusterId, storeType)
	if err != nil {
		return nil, err
	}

	pluginStoreDir, err := p.clonePluginStoreProject(orgId, clusterId, config.GitProjectURL, config.GitProjectId, storeType)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(pluginStoreDir)

	pluginConfig, err := p.getPluginConfig(pluginStoreDir, pluginName, version)
	if err != nil {
		return nil, err
	}

	return p.getPluginValues(pluginConfig, pluginStoreDir, pluginName, version)
}

func (p *PluginStore) getPluginValues(pluginConfig *PluginConfig, pluginStoreDir, pluginName, version string) ([]byte, error) {
	pluginValuesPath := p.getPluginDeployValuesFilePath(pluginStoreDir, pluginName, version,
		pluginConfig.Deployment.ControlplaneCluster.ValuesFile)
	p.log.Infof("Loading %s plugin values from %s", pluginName, pluginValuesPath)
	pluginValuesData, err := os.ReadFile(pluginValuesPath)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read plugins values file")
	}
	return pluginValuesData, nil
}

func (p *PluginStore) getPluginConfig(pluginStoreDir, pluginName, version string) (*PluginConfig, error) {
	pluginConfigPath := p.getPluginConfigFilePath(pluginStoreDir, pluginName, version)
	pluginConfigData, err := os.ReadFile(pluginConfigPath)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read store config file")
	}

	pluginConfig := &PluginConfig{}
	if err := yaml.Unmarshal(pluginConfigData, pluginConfig); err != nil {
		return nil, errors.WithMessage(err, "failed to unmarshall store config file")
	}
	if pluginConfig.Deployment.ControlplaneCluster == nil {
		return nil, errors.WithMessage(err, "no deployment found")
	}
	return pluginConfig, nil
}

func (p *PluginStore) DeployPlugin(orgId, clusterId string, storeType pluginstorepb.StoreType,
	pluginName, version string, values []byte) error {
	config, err := p.GetStoreConfig(clusterId, storeType)
	if err != nil {
		return fmt.Errorf("faild to fetch store config, %v", err)
	}

	pluginData, err := p.dbStore.ReadPluginData(config.GitProjectId, pluginName)
	if err != nil {
		return fmt.Errorf("faild to read plugin %s from project %s, %v", pluginName, config.GitProjectId, err)
	}

	if !stringContains(pluginData.Versions, version) {
		return fmt.Errorf("version %s not supported", version)
	}

	pluginStoreDir, err := p.clonePluginStoreProject(orgId, clusterId, config.GitProjectURL, config.GitProjectId, storeType)
	if err != nil {
		return err
	}
	defer os.RemoveAll(pluginStoreDir)

	pluginConfig, err := p.getPluginConfig(pluginStoreDir, pluginName, version)
	if err != nil {
		return err
	}

	validCapabilities, invalidCapabilities := filterSupporttedCapabilties(pluginConfig.Capabilities)
	if len(invalidCapabilities) > 0 {
		p.log.Infof("skipped plugin %s invalid capabilities %v", pluginName, invalidCapabilities)
	}

	if len(values) == 0 {
		values, err = p.getPluginValues(pluginConfig, pluginStoreDir, pluginName, version)
		if err != nil {
			p.log.Infof("no values defined for plugin %s", pluginName)
		}
	}

	overrideValuesMapping, err := p.getOverrideTemplateValues(orgId, clusterId)
	if err != nil {
		return err
	}

	apiEndpoint, uiEndpoint, err := p.getPluginDataAPIValues(pluginConfig, overrideValuesMapping)
	if err != nil {
		return err
	}

	if isUISSOCapabilitySupported(validCapabilities) {

		clientId, clientSecret, err := p.registerPluginSSO(orgId, clusterId, pluginName, uiEndpoint)
		if err != nil {
			return err
		}

		overrideValuesMapping[oAuthBaseURLName] = p.cfg.CaptenOAuthURL
		overrideValuesMapping[oAuthClientIdName] = clientId
		overrideValuesMapping[oAuthClientSecretName] = clientSecret
	}

	overrideValues, err := yaml.Marshal(overrideValuesMapping)
	if err != nil {
		return err
	}

	plugin := &clusterpluginspb.Plugin{
		StoreType:           clusterpluginspb.StoreType(pluginData.StoreType),
		PluginName:          pluginData.PluginName,
		Description:         pluginData.Description,
		Category:            pluginData.Category,
		Icon:                pluginData.Icon,
		Version:             version,
		ChartName:           pluginConfig.Deployment.ControlplaneCluster.ChartName,
		ChartRepo:           pluginConfig.Deployment.ControlplaneCluster.ChartRepo,
		DefaultNamespace:    pluginConfig.Deployment.ControlplaneCluster.DefaultNamespace,
		PrivilegedNamespace: pluginConfig.Deployment.ControlplaneCluster.PrivilegedNamespace,
		ApiEndpoint:         apiEndpoint,
		UiEndpoint:          uiEndpoint,
		Capabilities:        validCapabilities,
		Values:              values,
		OverrideValues:      overrideValues,
	}

	agent, err := p.agentHandler.GetAgent(orgId, clusterId)
	if err != nil {
		return err
	}

	p.log.Infof("Sending plugin %s deploy request to cluster %s", pluginName, clusterId)
	client := agent.GetClusterPluginsClient()
	_, err = client.DeployClusterPlugin(context.Background(), &clusterpluginspb.DeployClusterPluginRequest{Plugin: plugin})
	if err != nil {
		return err
	}
	return nil
}

func (p *PluginStore) UnDeployPlugin(orgId, clusterId string, storeType pluginstorepb.StoreType, pluginName string) error {
	agent, err := p.agentHandler.GetAgent(orgId, clusterId)
	if err != nil {
		return err
	}

	p.log.Infof("Sending plugin %s undeploy request to cluster %s", pluginName, clusterId)
	client := agent.GetClusterPluginsClient()
	_, err = client.UnDeployClusterPlugin(context.Background(),
		&clusterpluginspb.UnDeployClusterPluginRequest{StoreType: clusterpluginspb.StoreType(storeType), PluginName: pluginName})
	if err != nil {
		return err
	}
	return nil
}

func stringContains(arr []string, target string) bool {
	for _, str := range arr {
		if str == target {
			return true
		}
	}
	return false
}

func (p *PluginStore) getGitProjectAccessToken(orgId, clusterId, projectId string,
	storeType pluginstorepb.StoreType) (string, error) {
	if storeType == pluginstorepb.StoreType_CENTRAL_STORE {
		return p.cfg.PluginStoreProjectAccess, nil
	}

	agent, err := p.agentHandler.GetAgent(orgId, clusterId)
	if err != nil {
		return "", err
	}

	projectsResp, err := agent.GetCaptenPluginsClient().GetGitProjects(context.Background(), &captenpluginspb.GetGitProjectsRequest{})
	if err != nil {
		return "", err
	}

	var accessToken string
	for _, project := range projectsResp.Projects {
		if project.Id == projectId {
			accessToken = project.AccessToken
			break
		}
	}

	if len(accessToken) == 0 {
		return "", fmt.Errorf("project not found")
	}

	return accessToken, nil
}

func filterSupporttedCapabilties(pluginCapabilties []string) (validCapabilties, invalidCapabilities []string) {
	validCapabilties = []string{}
	invalidCapabilities = []string{}
	for _, pluginCapability := range pluginCapabilties {
		_, ok := supporttedCapabilities[pluginCapability]
		if ok {
			validCapabilties = append(validCapabilties, pluginCapability)
		} else {
			invalidCapabilities = append(invalidCapabilities, pluginCapability)
		}
	}
	return
}

func isUISSOCapabilitySupported(pluginCapabilties []string) bool {
	for _, pluginCapability := range pluginCapabilties {
		if pluginCapability == uiSSOCapabilityName {
			return true
		}
	}
	return false
}

func (p *PluginStore) registerPluginSSO(orgId, clusterId, pluginName, uiSSOURL string) (clientID, clientSecret string, err error) {
	pluginClientName := fmt.Sprintf("%s-%s", pluginName, clusterId)
	p.log.Infof("Register plugin %s as app-client %s with IAM, clusterId: %s, [org: %s]",
		pluginName, pluginClientName, clusterId, orgId)
	clientID, clientSecret, err = p.iam.RegisterAppClientSecrets(context.Background(),
		pluginClientName, uiSSOURL, orgId)
	if err != nil {
		err = errors.WithMessagef(err, "failed to register plugin %s on cluster %s with IAM", pluginName, clusterId)
	}
	return
}

func (p *PluginStore) getPluginDataAPIValues(pluginConfig *PluginConfig, overrideValues map[string]string) (string, string, error) {
	apiEndpoint, err := replaceTemplateValuesInString(pluginConfig.ApiEndpoint, overrideValues)
	if err != nil {
		return "", "", fmt.Errorf("failed to update template values in plguin data, %v", err)
	}

	uiEndpoint, err := replaceTemplateValuesInString(pluginConfig.UIEndpoint, overrideValues)
	if err != nil {
		return "", "", fmt.Errorf("failed to update template values in plguin data, %v", err)
	}
	return apiEndpoint, uiEndpoint, nil
}

func (p *PluginStore) getOverrideTemplateValues(orgId, clusterID string) (map[string]string, error) {
	clusterGlobalValues, err := p.getClusterGlobalValues(orgId, clusterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster global values, %v", err)
	}

	overrideValues := map[string]string{}
	for key, val := range clusterGlobalValues {
		overrideValues[key] = fmt.Sprintf("%v", val)

	}

	return overrideValues, nil
}

func (p *PluginStore) getClusterGlobalValues(orgId, clusterID string) (map[string]interface{}, error) {
	agent, err := p.agentHandler.GetAgent(orgId, clusterID)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to initialize agent for cluster %s", clusterID)
	}
	resp, err := agent.GetClient().GetClusterGlobalValues(context.TODO(), &agentpb.GetClusterGlobalValuesRequest{})
	if err != nil {
		return nil, err
	}
	if resp.Status != agentpb.StatusCode_OK {
		return nil, fmt.Errorf("failed to get global values for cluster %s", clusterID)
	}

	var globalValues map[string]interface{}
	err = yaml.Unmarshal(resp.GlobalValues, &globalValues)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to unmarshal cluster values")
	}
	p.log.Debugf("cluster %s globalValues: %+v", clusterID, globalValues)
	return globalValues, nil
}

func replaceTemplateValuesInByteData(data []byte,
	values map[string]interface{}) (transformedData []byte, err error) {
	tmpl, err := template.New("templateVal").Parse(string(data))
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, values)
	if err != nil {
		return
	}

	transformedData = buf.Bytes()
	return
}

func replaceTemplateValuesInString(data string, values map[string]string) (transformedData string, err error) {
	if len(data) == 0 {
		return
	}

	tmpl, err := template.New("templateVal").Parse(data)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, values)
	if err != nil {
		return
	}

	transformedData = string(buf.Bytes())

	return
}

func prepareFilePath(parts ...string) string {
	return filepath.Join(parts...)
}

func (p *PluginStore) getPluginListFilePath(parentFolder string) string {
	return prepareFilePath(parentFolder, p.cfg.PluginsStorePath, p.cfg.PluginsFileName)
}

func (p *PluginStore) getPluginFilePath(parentFolder, pluginName string) string {
	return prepareFilePath(parentFolder, p.cfg.PluginsStorePath, pluginName, p.cfg.PluginFileName)
}

func (p *PluginStore) getPluginIconFilePath(parentFolder, pluginName, iconFileName string) string {
	return prepareFilePath(parentFolder, p.cfg.PluginsStorePath, pluginName, iconFileName)
}

func (p *PluginStore) getPluginConfigFilePath(parentFolder, pluginName, version string) string {
	return prepareFilePath(parentFolder, p.cfg.PluginsStorePath, pluginName, version, p.cfg.PluginConfigFileName)
}

func (p *PluginStore) getPluginDeployValuesFilePath(parentFolder, pluginName, version, valuesFile string) string {
	return prepareFilePath(parentFolder, p.cfg.PluginsStorePath, pluginName, version, valuesFile)
}
