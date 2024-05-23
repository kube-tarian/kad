package defaultplugindeployer

import (
	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/common-pkg/capten-store"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/pluginstorepb"
	pluginstore "github.com/kube-tarian/kad/capten/common-pkg/plugin-store"
)

type DefaultPluginsDeployer struct {
	pluginStore pluginstore.PluginStoreInterface
	log         logging.Logger
	frequency   string
}

func NewDefaultPluginsDeployer(
	log logging.Logger,
	frequency string,
	dbStore *captenstore.Store,
	handler pluginstore.PluginDeployHandler,
) (*DefaultPluginsDeployer, error) {
	pluginStore, err := pluginstore.NewPluginStore(log, dbStore, handler)
	if err != nil {
		return nil, err
	}
	return &DefaultPluginsDeployer{
		log:         log,
		frequency:   frequency,
		pluginStore: pluginStore,
	}, nil
}

func (p *DefaultPluginsDeployer) CronSpec() string {
	return p.frequency
}

func (p *DefaultPluginsDeployer) Run() {
	p.log.Debug("started default plugins deployer job")
	if err := p.pluginStore.SyncPlugins(pluginstorepb.StoreType_DEFAULT_STORE); err != nil {
		p.log.Errorf("failed to synch providers, %v", err)
	}

	p.deployPlugins()

	p.log.Debug("defualt plugins deployer job completed")
}

func (p *DefaultPluginsDeployer) deployPlugins() {
	plugins, err := p.pluginStore.GetPlugins(pluginstorepb.StoreType_DEFAULT_STORE)
	if err != nil {
		p.log.Errorf("failed to get plugins, %v", err)
	}

	for _, plugin := range plugins {
		if err := p.pluginStore.DeployPlugin(pluginstorepb.StoreType_DEFAULT_STORE, plugin.PluginName, plugin.Versions[0], []byte{}); err != nil {
			p.log.Errorf("failed to deploy plugin, %v", err)
		}
	}
}
