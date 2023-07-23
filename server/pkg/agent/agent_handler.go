package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/kube-tarian/kad/server/pkg/credential"
	"github.com/kube-tarian/kad/server/pkg/store"
	"github.com/pkg/errors"
)

type AgentHandler struct {
	agentMutex  sync.RWMutex
	agents      map[string]*Agent
	serverStore store.ServerStore
}

func NewAgentHandler(serverStore store.ServerStore) *AgentHandler {
	return &AgentHandler{serverStore: serverStore}
}

func (s *AgentHandler) AddAgent(orgId, clusterName string, agentCfg *Config) error {
	clusterKey := getClusterAgentKey(orgId, clusterName)
	if _, ok := s.agents[clusterKey]; ok {
		return nil
	}

	agent, err := NewAgent(agentCfg)
	if err != nil {
		return err
	}

	s.agentMutex.Lock()
	defer s.agentMutex.Unlock()
	s.agents[clusterKey] = agent
	return err
}

func (s *AgentHandler) UpdateAgent(orgId, clusterName string, agentCfg *Config) error {
	clusterKey := getClusterAgentKey(orgId, clusterName)
	if _, ok := s.agents[clusterKey]; ok {
		return s.AddAgent(orgId, clusterName, agentCfg)
	}

	s.RemoveAgent(orgId, clusterName)
	return s.AddAgent(orgId, clusterName, agentCfg)
}

func (s *AgentHandler) GetAgent(orgId, clusterName string) (*Agent, error) {
	agent := s.getAgent(orgId, clusterName)
	if agent != nil {
		return agent, nil
	}

	cfg, err := s.getAgentConfig(orgId, clusterName)
	if err != nil {
		return nil, err
	}

	if err := s.AddAgent(orgId, clusterName, cfg); err != nil {
		return nil, err
	}

	agent = s.getAgent(orgId, clusterName)
	if agent != nil {
		return agent, nil
	}
	return nil, fmt.Errorf("failed to get agent")
}

func (s *AgentHandler) getAgent(orgId, clusterName string) *Agent {
	s.agentMutex.RLock()
	defer s.agentMutex.RUnlock()
	clusterKey := getClusterAgentKey(orgId, clusterName)
	if agent, ok := s.agents[clusterKey]; ok && agent != nil {
		return agent
	}
	return nil
}

func (s *AgentHandler) RemoveAgent(orgId, clusterName string) {
	clusterKey := getClusterAgentKey(orgId, clusterName)
	s.removeAgentEntry(clusterKey)
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

func (s *AgentHandler) getAgentConfig(orgId, clusterName string) (*Config, error) {
	agentCfg := &Config{}
	endpoint, err := s.serverStore.GetClusterEndpoint(orgId, clusterName)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get cluster")
	}

	agentCfg.Address = endpoint
	certData, err := credential.GetClusterCerts(context.TODO(), orgId, clusterName)
	if err != nil {
		return nil, errors.WithMessage(err, "failed get cert from vault")
	}

	agentCfg.CaCert = certData.CACert
	agentCfg.Key = certData.Key
	agentCfg.Cert = certData.Cert
	return agentCfg, err
}

func getClusterAgentKey(orgId, clusterName string) string {
	return fmt.Sprintf("%s-%s", orgId, clusterName)
}
