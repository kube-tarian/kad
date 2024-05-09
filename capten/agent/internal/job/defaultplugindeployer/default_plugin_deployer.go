package defaultplugindeployer

import (
	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/agent/internal/capten-store"
	"github.com/kube-tarian/kad/capten/agent/internal/temporalclient"
	pluginconfigtore "github.com/kube-tarian/kad/capten/common-pkg/pluginconfig-store"
)

type DefaultPluginsDeployer struct {
	pluginStoreHandler *PluginStore
	log                logging.Logger
	frequency          string
}

func NewDefaultPluginsDeployer(
	log logging.Logger,
	frequency string,
	dbStore *captenstore.Store,
	pas *pluginconfigtore.Store,
	tc *temporalclient.Client,
) (*DefaultPluginsDeployer, error) {
	pluginStoreHandler, err := NewPluginStore(log, dbStore, pas, tc)
	if err != nil {
		return nil, err
	}
	return &DefaultPluginsDeployer{
		log:                log,
		frequency:          frequency,
		pluginStoreHandler: pluginStoreHandler,
	}, nil
}

func (p *DefaultPluginsDeployer) CronSpec() string {
	return p.frequency
}

func (p *DefaultPluginsDeployer) Run() {
	p.log.Debug("started default plugins deployer job")
	if err := p.pluginStoreHandler.SyncPlugins(); err != nil {
		p.log.Errorf("failed to synch providers, %v", err)
	}

	p.pluginStoreHandler.DeployPlugins()

	p.log.Debug("defualt plugins deployer job completed")
}
