package defaultplugindeployer

import (
	"context"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/agentpb"
)

type defaultPluginsDeployer interface {
	DeployDefaultApps(ctx context.Context, request *agentpb.DeployDefaultAppsRequest) (*agentpb.DeployDefaultAppsResponse, error)
}

type DefaultPluginsDeployer struct {
	agent     defaultPluginsDeployer
	log       logging.Logger
	frequency string
}

func NewDefaultPluginsDeployer(
	log logging.Logger,
	frequency string,
	agent defaultPluginsDeployer,
) (*DefaultPluginsDeployer, error) {
	return &DefaultPluginsDeployer{
		log:       log,
		frequency: frequency,
		agent:     agent,
	}, nil
}

func (p *DefaultPluginsDeployer) CronSpec() string {
	return p.frequency
}

func (p *DefaultPluginsDeployer) Run() {
	p.log.Debug("started default plugins deployer job")
	resp, _ := p.agent.DeployDefaultApps(context.TODO(), &agentpb.DeployDefaultAppsRequest{
		Upgrade: false,
	})
	if resp.Status != agentpb.StatusCode_OK {
		p.log.Errorf("failed to deploy default apps, %s", resp.StatusMessage)
		return
	}
	p.log.Debug("defualt plugins deployer job completed")
}
