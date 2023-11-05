package agent

import (
	"context"
	"fmt"

	"github.com/intelops/go-common/logging"
	captenstore "github.com/kube-tarian/kad/capten/agent/pkg/capten-store"
	"github.com/kube-tarian/kad/capten/agent/pkg/config"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/agent/pkg/temporalclient"
)

var _ agentpb.AgentServer = &Agent{}

type Agent struct {
	agentpb.UnimplementedAgentServer
	captenpluginspb.UnimplementedCaptenPluginsServer
	tc       *temporalclient.Client
	as       *captenstore.Store
	log      logging.Logger
	createPr bool
}

func NewAgent(log logging.Logger, cfg *config.SericeConfig) (*Agent, error) {
	var tc *temporalclient.Client
	var err error

	tc, err = temporalclient.NewClient(log)
	if err != nil {
		return nil, err
	}

	as, err := captenstore.NewStore(log)
	if err != nil {
		// ignoring store failure until DB user creation working
		// return nil, err
		log.Errorf("failed to initialize store, %v", err)
	}

	agent := &Agent{
		tc:  tc,
		as:  as,
		log: log,
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
		default:
			return fmt.Errorf("validation not implemented for this type")
		}

	}
	return nil
}
