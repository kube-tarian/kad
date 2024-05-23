package api

import (
	"context"
	"fmt"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/internal/config"
	captenstore "github.com/kube-tarian/kad/capten/common-pkg/capten-store"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/agentpb"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/clusterpluginspb"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/pluginstorepb"
	pluginstore "github.com/kube-tarian/kad/capten/common-pkg/plugin-store"
	"github.com/kube-tarian/kad/capten/common-pkg/temporalclient"
)

var _ agentpb.AgentServer = &Agent{}

type Agent struct {
	agentpb.UnimplementedAgentServer
	captenpluginspb.UnimplementedCaptenPluginsServer
	clusterpluginspb.UnimplementedClusterPluginsServer
	pluginstorepb.UnimplementedPluginStoreServer
	tc       *temporalclient.Client
	as       *captenstore.Store
	log      logging.Logger
	cfg      *config.SericeConfig
	plugin   pluginstore.PluginStoreInterface
	createPr bool
}

func NewAgent(log logging.Logger, cfg *config.SericeConfig, as *captenstore.Store) (agent *Agent, err error) {
	tc, err := temporalclient.NewClient(log)
	if err != nil {
		log.Errorf("failed to initialize temporal client, %v", err)
		return
	}

	agent = &Agent{
		tc:  tc,
		as:  as,
		cfg: cfg,
		log: log,
	}

	agent.plugin, err = pluginstore.NewPluginStore(log, as, tc)
	if err != nil {
		return nil, err
	}
	return agent, nil
}

func (a *Agent) Ping(ctx context.Context, request *agentpb.PingRequest) (*agentpb.PingResponse, error) {
	a.log.Infof("Ping request received")
	return &agentpb.PingResponse{Status: agentpb.StatusCode_OK}, nil
}

func validateArgs(args ...any) error {
	for index, arg := range args {
		switch item := arg.(type) {
		case string:
			if len(item) == 0 {
				return fmt.Errorf("empty string not allowed for arg index: %v", index)
			}
		case map[string]string:
			for k, v := range item {
				if len(v) == 0 {
					return fmt.Errorf("map value empty for key: %v", k)
				}
			}
		case []string:
			if len(item) == 0 {
				return fmt.Errorf("empty []string not allowed for arg index: %v", index)
			}
		default:
			return fmt.Errorf("validation not implemented for this type")
		}

	}
	return nil
}
