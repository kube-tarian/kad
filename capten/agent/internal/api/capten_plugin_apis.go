package api

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/internal/pb/captenpluginspb"
)

func (a *Agent) GetCaptenPlugins(ctx context.Context, request *captenpluginspb.GetCaptenPluginsRequest) (
	*captenpluginspb.GetCaptenPluginsResponse, error) {
	a.log.Infof("Get Capten Plugins request recieved")

	res, err := a.as.GetAllApps()
	if err != nil {
		a.log.Errorf("failed to fetch plugins, %v", err)
		return &captenpluginspb.GetCaptenPluginsResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "failed to fetch plugins",
		}, nil
	}

	plugins := make([]*captenpluginspb.CaptenPlugin, 0)
	for _, r := range res {
		appConfig := r.GetConfig()
		if len(appConfig.PluginName) == 0 {
			continue
		}

		plugins = append(plugins, &captenpluginspb.CaptenPlugin{
			PluginName:        r.GetConfig().GetPluginName(),
			PluginDescription: r.GetConfig().GetPluginDescription(),
			Icon:              r.GetConfig().GetIcon(),
			UiEndpoint:        r.GetConfig().GetUiEndpoint(),
			UiModuleEndpoint:  r.GetConfig().GetUiModuleEndpoint(),
			ApiEndpoint:       r.GetConfig().GetApiEndpoint(),
			InstallStatus:     r.GetConfig().GetInstallStatus(),
			RuntimeStatus:     r.GetConfig().GetRuntimeStatus(),
		})
	}

	a.log.Infof("Fetched %d capten plugins", len(plugins))
	return &captenpluginspb.GetCaptenPluginsResponse{Status: captenpluginspb.StatusCode_OK,
		StatusMessage: "successfully fetched plugin",
		Plugins:       plugins}, nil
}
