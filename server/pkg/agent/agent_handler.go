package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/server/pkg/config"
	"github.com/kube-tarian/kad/server/pkg/credential"
	oryclient "github.com/kube-tarian/kad/server/pkg/ory-client"
	"github.com/kube-tarian/kad/server/pkg/store"
	"github.com/pkg/errors"
)

type AgentHandler struct {
	log         logging.Logger
	cfg         config.ServiceConfig
	agentMutex  sync.RWMutex
	agents      map[string]*Agent
	serverStore store.ServerStore
	oryClient   oryclient.OryClient
}

func NewAgentHandler(log logging.Logger, cfg config.ServiceConfig,
	serverStore store.ServerStore, oryClient oryclient.OryClient) *AgentHandler {
	return &AgentHandler{log: log, cfg: cfg, serverStore: serverStore, agents: map[string]*Agent{}, oryClient: oryClient}
}

func (s *AgentHandler) AddAgent(clusterID string, agentCfg *Config) error {
	if _, ok := s.agents[clusterID]; ok {
		return nil
	}

	agentCfg.ServicName = s.cfg.ServiceName
	agentCfg.AuthEnabled = s.cfg.AuthEnabled
	agent, err := NewAgent(s.log, agentCfg, s.oryClient)
	if err != nil {
		return err
	}

	s.agentMutex.Lock()
	defer s.agentMutex.Unlock()
	s.agents[clusterID] = agent
	return err
}

func (s *AgentHandler) UpdateAgent(clusterID string, agentCfg *Config) error {
	if _, ok := s.agents[clusterID]; !ok {
		return s.AddAgent(clusterID, agentCfg)
	}

	s.RemoveAgent(clusterID)
	return s.AddAgent(clusterID, agentCfg)
}

func (s *AgentHandler) GetAgent(orgId, clusterID string) (*Agent, error) {
	/*agent := s.getAgent(clusterID)
	if agent != nil {
		return agent, nil
	}*/

	cfg, err := s.getAgentConfig(orgId, clusterID)
	if err != nil {
		return nil, err
	}

	if err := s.AddAgent(clusterID, cfg); err != nil {
		return nil, err
	}

	agent := s.getAgent(clusterID)
	if agent != nil {
		return agent, nil
	}
	return nil, fmt.Errorf("failed to get agent")
}

func (s *AgentHandler) GetAgentClusterDetail(clusterID string) *Config {
	s.agentMutex.RLock()
	defer s.agentMutex.RUnlock()
	if agent, ok := s.agents[clusterID]; ok && agent != nil {
		return agent.cfg
	}

	return &Config{}
}

func (s *AgentHandler) getAgent(clusterID string) *Agent {
	s.agentMutex.RLock()
	defer s.agentMutex.RUnlock()
	if agent, ok := s.agents[clusterID]; ok && agent != nil {
		return agent
	}
	return nil
}

func (s *AgentHandler) RemoveAgent(clusterID string) {
	s.removeAgentEntry(clusterID)
}

func (s *AgentHandler) removeAgentEntry(clusterKey string) {
	s.agentMutex.Lock()
	defer s.agentMutex.Unlock()
	agent, ok := s.agents[clusterKey]
	if ok {
		if agent != nil {
			s.agents[clusterKey].Close()
		}
		delete(s.agents, clusterKey)
	}
}

func (s *AgentHandler) Close() {
	for clusterKey := range s.agents {
		s.removeAgentEntry(clusterKey)
	}
}

func (s *AgentHandler) getAgentConfig(orgId, clusterID string) (*Config, error) {
	agentCfg := &Config{}

	clusterDetail, err := s.serverStore.GetClusterDetails(orgId, clusterID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get cluster")
	}

	agentCfg.Address = clusterDetail.Endpoint
	agentCfg.ClusterName = clusterDetail.ClusterName

	certData, err := credential.GetClusterCerts(context.TODO(), clusterID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed get cert from vault")
	}

	agentCfg.CaCert = certData.CACert
	agentCfg.Key = certData.Key
	agentCfg.Cert = certData.Cert
	s.log.Infof("loaded agent certs for cluster %s", clusterID)
	return agentCfg, err
}
