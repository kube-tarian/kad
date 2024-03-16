package pluginstore

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/server/pkg/agent"
	"github.com/kube-tarian/kad/server/pkg/credential"
	iamclient "github.com/kube-tarian/kad/server/pkg/iam-client"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
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
		config, err := p.dbStore.ReadPluginStoreConfig(clusterId)
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

func (p *PluginStore) SyncPlugins(clusterId string, storeType pluginstorepb.StoreType) error {
	config, err := p.GetStoreConfig(clusterId, storeType)
	if err != nil {
		return err
	}

	pluginStoreDir, err := p.clonePluginStoreProject(config.GitProjectURL, config.GitProjectId, storeType)
	if err != nil {
		return err
	}
	defer os.RemoveAll(pluginStoreDir)

	pluginListFilePath := pluginStoreDir + "/" + p.cfg.PluginsStorePath + "/" + p.cfg.PluginsFileName
	p.log.Infof("Loading plugin data from %s", pluginListFilePath)
	pluginListData, err := os.ReadFile(pluginListFilePath)
	if err != nil {
		return errors.WithMessage(err, "failed to read store config file")
	}

	var plugins PluginListData
	if err := yaml.Unmarshal(pluginListData, &plugins); err != nil {
		return errors.WithMessage(err, "failed to unmarshall store config file")
	}

	for _, pluginName := range plugins.Plugins {
		err := p.addPluginApp(config.GitProjectId, pluginStoreDir, pluginName)
		if err != nil {
			p.log.Errorf("%v", err)
			continue
		}
		p.log.Infof("stored plugin data for plugin %s for cluster %s", pluginName, clusterId)
	}
	return nil
}

func (p *PluginStore) clonePluginStoreProject(projectURL, projectId string,
	storeType pluginstorepb.StoreType) (pluginStoreDir string, err error) {
	pluginStoreDir, err = os.MkdirTemp(p.cfg.PluginsStoreProjectMount, tmpGitProjectCloneStr)
	if err != nil {
		err = fmt.Errorf("failed to create plugin store tmp dir, err: %v", err)
		return
	}

	accessToken, err := p.getGitProjectAccessToken(projectId, storeType)
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
	appData, err := os.ReadFile(pluginStoreDir + "/" + p.cfg.PluginsStorePath + "/" + pluginName + "/plugin.yaml")
	if err != nil {
		return errors.WithMessagef(err, "failed to read store plugin %s", pluginName)
	}

	var pluginData Plugin
	if err := yaml.Unmarshal(appData, &pluginData); err != nil {
		return errors.WithMessagef(err, "failed to unmarshall store plugin %s", pluginName)
	}

	if pluginData.PluginName == "" || len(pluginData.DeploymentConfig.Versions) == 0 {
		return fmt.Errorf("app name/version is missing for %s", pluginName)
	}

	plugin := &pluginstorepb.PluginData{
		PluginName:          pluginData.PluginName,
		Description:         pluginData.Description,
		Category:            pluginData.Category,
		ChartName:           pluginData.DeploymentConfig.ChartName,
		ChartRepo:           pluginData.DeploymentConfig.ChartRepo,
		Versions:            pluginData.DeploymentConfig.Versions,
		DefaultNamespace:    pluginData.DeploymentConfig.DefaultNamespace,
		PrivilegedNamespace: pluginData.DeploymentConfig.PrivilegedNamespace,
		ApiEndpoint:         pluginData.PluginConfig.ApiEndpoint,
		UiEndpoint:          pluginData.PluginConfig.UiEndpoint,
		Capabilities:        pluginData.PluginConfig.Capabilities,
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

func (p *PluginStore) GetPluginValues(clusterId string, storeType pluginstorepb.StoreType,
	pluginName, version string) ([]byte, error) {
	config, err := p.GetStoreConfig(clusterId, storeType)
	if err != nil {
		return nil, err
	}

	pluginStoreDir, err := p.clonePluginStoreProject(config.GitProjectURL, config.GitProjectId, storeType)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(pluginStoreDir)

	pluginValuesPath := pluginStoreDir + "/" + p.cfg.PluginsStorePath + "/" + pluginName + "/" + version + "/" + "values.yaml"
	p.log.Infof("Loading %s plugin values from %s", pluginName, pluginValuesPath)
	pluginListData, err := os.ReadFile(pluginValuesPath)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to read plugins values file")
	}
	return pluginListData, nil
}

func (p *PluginStore) DeployPlugin(orgId, clusterId string, storeType pluginstorepb.StoreType,
	pluginName, version string, values []byte) error {
	config, err := p.GetStoreConfig(clusterId, storeType)
	if err != nil {
		return err
	}

	pluginData, err := p.dbStore.ReadPluginData(config.GitProjectId, pluginName)
	if err != nil {
		return err
	}

	if !stringContains(pluginData.Versions, version) {
		return fmt.Errorf("version %s not supported", version)
	}

	validCapabilities, invalidCapabilities := filterSupporttedCapabilties(pluginData.Capabilities)
	if len(invalidCapabilities) > 0 {
		p.log.Infof("skipped plugin %s invalid capabilities %v", pluginName, invalidCapabilities)
	}

	pluginData, updatedValues, err := p.updatePluginDataTemplateValues(orgId, clusterId, pluginData, values)
	if err != nil {
		return err
	}

	overrideValues := map[string]string{}
	if isUISSOCapabilitySupported(validCapabilities) {
		clientId, clientSecret, err := p.registerPluginSSO(orgId, clusterId, pluginName, pluginData.UiEndpoint)
		if err != nil {
			return err
		}
		overrideValues[oAuthBaseURLName] = p.cfg.CaptenOAuthURL
		overrideValues[oAuthClientIdName] = clientId
		overrideValues[oAuthClientSecretName] = clientSecret
	}

	plugin := &clusterpluginspb.Plugin{
		StoreType:           clusterpluginspb.StoreType(pluginData.StoreType),
		PluginName:          pluginData.PluginName,
		Description:         pluginData.Description,
		Category:            pluginData.Category,
		Icon:                pluginData.Icon,
		Version:             version,
		ChartName:           pluginData.ChartName,
		ChartRepo:           pluginData.ChartRepo,
		DefaultNamespace:    pluginData.DefaultNamespace,
		PrivilegedNamespace: pluginData.PrivilegedNamespace,
		ApiEndpoint:         pluginData.ApiEndpoint,
		UiEndpoint:          pluginData.UiEndpoint,
		Capabilities:        validCapabilities,
		Values:              updatedValues,
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

func (p *PluginStore) getGitProjectAccessToken(projectId string,
	storeType pluginstorepb.StoreType) (string, error) {
	if storeType == pluginstorepb.StoreType_CENTRAL_STORE {
		return "", nil
	}
	cred, err := credential.GetGenericCredential(context.Background(), p.cfg.GitVaultEntityName, projectId)
	if err != nil {
		err = errors.WithMessagef(err, "error while reading credential %s/%s from the vault",
			p.cfg.GitVaultEntityName, projectId)
		return "", err
	}

	return cred[gitProjectAccessTokenAttribute], nil
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

func (p *PluginStore) updatePluginDataTemplateValues(orgId, clusterID string,
	pluginData *pluginstorepb.PluginData, values []byte) (*pluginstorepb.PluginData, []byte, error) {
	clusterGlobalValues, err := p.getClusterGlobalValues(orgId, clusterID)
	if err != nil {
		return pluginData, values, fmt.Errorf("failed to get cluster global values, %v", err)
	}

	apiEndpoint, err := replaceTemplateValuesInString(pluginData.ApiEndpoint, clusterGlobalValues)
	if err != nil {
		return pluginData, values, fmt.Errorf("failed to update template values in plguin data, %v", err)
	}

	uiEndpoint, err := replaceTemplateValuesInString(pluginData.UiEndpoint, clusterGlobalValues)
	if err != nil {
		return pluginData, values, fmt.Errorf("failed to update template values in plguin data, %v", err)
	}

	updatedValues, err := replaceTemplateValuesInByteData(values, clusterGlobalValues)
	if err != nil {
		return pluginData, values, fmt.Errorf("failed to update template values in plguin values, %v", err)
	}

	pluginData.ApiEndpoint = apiEndpoint
	pluginData.UiEndpoint = uiEndpoint
	return pluginData, updatedValues, nil
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

func replaceTemplateValuesInString(data string, values map[string]interface{}) (transformedData string, err error) {
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
