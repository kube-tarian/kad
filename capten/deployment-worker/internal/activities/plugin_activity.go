package activities

import (
	"context"

	"github.com/kube-tarian/kad/capten/common-pkg/plugin_store/pluginstorepb"
)

type PluginActivities struct {
	state string
}

func NewPluginActivities() *PluginActivities {
	return &PluginActivities{
		state: "initialized",
	}
}

func (p *PluginActivities) PluginDeployActivity(ctx context.Context, req *pluginstorepb.DeployPluginRequest) (*pluginstorepb.DeployPluginResponse, error) {
	logger.Infof("state: %v", p.state)
	return nil, nil
}

func (p *PluginActivities) PluginUndeployActivity(ctx context.Context, req *pluginstorepb.DeployPluginRequest) (*pluginstorepb.DeployPluginResponse, error) {
	logger.Infof("state: %v", p.state)
	return nil, nil
}
