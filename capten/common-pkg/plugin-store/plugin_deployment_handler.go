package pluginstore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"text/template"

	"github.com/kube-tarian/kad/capten/common-pkg/pb/clusterpluginspb"
	"github.com/kube-tarian/kad/capten/common-pkg/workers"
	"github.com/kube-tarian/kad/capten/model"
	"gopkg.in/yaml.v2"
)

func (p *PluginStore) DeployClusterPlugin(ctx context.Context, pluginData *clusterpluginspb.Plugin) error {
	p.log.Infof("Recieved Plugin Deploy request for plugin %s, version %+v", pluginData.PluginName, pluginData.Version)

	values, err := replaceTemplateValueBytesInByteData(pluginData.Values, pluginData.OverrideValues)
	if err != nil {
		return fmt.Errorf("failed to derive template values for plugin %s, %v", pluginData.PluginName, err)
	}

	pluginData.Values = values
	pluginData.InstallStatus = string(model.AppIntallingStatus)
	if err := p.dbStore.UpsertClusterPluginConfig(pluginData); err != nil {
		return fmt.Errorf("failed to update plugin config data for plugin %s, %v", pluginData.PluginName, err)
	}

	go p.deployPluginWithWorkflow(pluginData)
	p.log.Infof("Triggerred plugin [%s] install", pluginData.PluginName)
	return nil
}

func (p *PluginStore) UnDeployClusterPlugin(ctx context.Context, request *clusterpluginspb.UnDeployClusterPluginRequest) error {
	pluginConfigdata, err := p.dbStore.GetClusterPluginConfig(request.PluginName)
	if err != nil {
		return fmt.Errorf("failed to fetch plugin config record %s, %v", request.PluginName, err)
	}

	pluginConfigdata.InstallStatus = string(model.AppUnInstallingStatus)
	if err := p.dbStore.UpsertClusterPluginConfig(pluginConfigdata); err != nil {
		return fmt.Errorf("failed to update plugin config status with UnInstalling for plugin %s, %v", request.PluginName, err)
	}

	go p.unInstallPluginWithWorkflow(request, pluginConfigdata)

	p.log.Infof("Triggerred plugin [%s] un install", request.PluginName)
	return nil
}

// func (p *PluginStore) deployPluginWithWorkflow(pluginData *clusterpluginspb.Plugin) {
// 	wd := workers.NewDeployment(p.tc, p.log)
// 	_, err := wd.SendEventV2(context.TODO(), wd.GetPluginWorkflowName(), string(model.AppInstallAction), pluginData)
// 	if err != nil {
// 		// pluginConfig.InstallStatus = string(model.AppIntallFailedStatus)
// 		// if err := a.pas.UpsertPluginConfig(pluginConfig); err != nil {
// 		// 	a.log.Errorf("failed to update plugin config for plugin %s, %v", pluginConfig.PluginName, err)
// 		// 	return
// 		// }
// 		p.log.Errorf("sendEventV2 failed, plugin: %s, reason: %v", pluginData.PluginName, err)
// 		return
// 	}
// 	// TODO: workflow will update the final status
// 	// Write a periodic scheduler which will go through all apps not in installed status and check the status till either success or failed.
// 	// Make SendEventV2 asynchrounous so that periodic scheduler will take care of monitoring.
// }

func (p *PluginStore) deployPluginWithWorkflow(pluginData *clusterpluginspb.Plugin) {
	wd := workers.NewDeployment(p.tc, p.log)

	// Convert pluginData to a JSON string
	pluginDataJSON, err := json.Marshal(pluginData)
	if err != nil {
		p.log.Errorf("failed to marshal pluginData: %s, reason: %v", pluginData.PluginName, err)
		return
	}

	// Create a PluginDeployRequest instance
	pluginDeployRequest := &model.PluginDeployRequest{
		Data: string(pluginDataJSON),
	}
	log.Println("Plugin Deploy Req", pluginDeployRequest)

	// Ensure the payload is a model.DeployRequest
	_, err = wd.SendEventV2(context.TODO(), wd.GetPluginWorkflowName(), string(model.AppInstallAction), pluginDeployRequest)
	if err != nil {
		// Uncomment and update the plugin configuration if needed
		// pluginConfig.InstallStatus = string(model.AppIntallFailedStatus)
		// if err := a.pas.UpsertPluginConfig(pluginConfig); err != nil {
		// 	a.log.Errorf("failed to update plugin config for plugin %s, %v", pluginConfig.PluginName, err)
		// 	return
		// }
		p.log.Errorf("sendEventV2 failed, plugin: %s, reason: %v", pluginData.PluginName, err)
		return
	}
}

func (p *PluginStore) unInstallPluginWithWorkflow(request *clusterpluginspb.UnDeployClusterPluginRequest, plugin *clusterpluginspb.Plugin) {
	wd := workers.NewDeployment(p.tc, p.log)
	_, err := wd.SendDeleteEvent(context.TODO(), wd.GetPluginWorkflowName(), string(model.AppUnInstallAction), request)
	if err != nil {
		p.log.Errorf("failed to send delete event to workflow for plugin %s, %v", request.PluginName, err)

		plugin.InstallStatus = string(model.AppUnUninstallFailedStatus)
		if err := p.dbStore.UpsertClusterPluginConfig(plugin); err != nil {
			p.log.Errorf("failed to update plugin config status with Installed for plugin %s, %v", request.PluginName, err)
		}
	}
}

func replaceTemplateValueBytesInByteData(data []byte,
	values []byte) (transformedData []byte, err error) {
	tmpl, err := template.New("templateVal").Parse(string(data))
	if err != nil {
		return
	}

	mapValues := map[string]any{}
	if err = yaml.Unmarshal(values, &mapValues); err != nil {
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, mapValues)
	if err != nil {
		return
	}

	transformedData = buf.Bytes()
	return
}
