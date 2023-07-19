package api

import (
	"context"
	"fmt"
	"sync"

	"github.com/kube-tarian/kad/server/pkg/client"
	"github.com/kube-tarian/kad/server/pkg/config"
	"github.com/kube-tarian/kad/server/pkg/db"
	"github.com/kube-tarian/kad/server/pkg/log"
	"github.com/kube-tarian/kad/server/pkg/types"

	"go.uber.org/zap"
)

var (
	agentMutex sync.RWMutex
)

func (a *Api) ConnectClient(orgId, clusterName string, agentCfg *types.AgentConfiguration) error {
	clusterKey := getClusterKey(orgId, clusterName)
	if _, ok := a.agents[clusterKey]; ok {
		return nil
	}

	logger := log.GetLogger()
	defer logger.Sync()
	agent, err := client.NewAgent(agentCfg)
	if err != nil {
		logger.Error("failed to connect agent internal error", zap.Error(err))
		return err
	}

	agentMutex.Lock()
	a.agents[clusterKey] = agent
	agentMutex.Unlock()
	return err
}

func (a *Api) ReConnect(orgId, clusterName string, agentCfg *types.AgentConfiguration) error {
	clusterKey := getClusterKey(orgId, clusterName)
	if _, ok := a.agents[clusterKey]; ok {
		return a.ConnectClient(orgId, clusterName, agentCfg)
	}

	a.Close(clusterKey)
	return a.ConnectClient(orgId, clusterName, agentCfg)
}

func (a *Api) GetClient(orgId, clusterName string) *client.Agent {
	agentMutex.RLock()
	clusterKey := getClusterKey(orgId, clusterName)
	if agent, ok := a.agents[clusterKey]; ok && agent != nil {
		return agent
	}
	agentMutex.RUnlock()
	return nil
}

func (a *Api) Close(clusterKey string) {
	agentMutex.Lock()
	agent, ok := a.agents[clusterKey]
	if ok {
		if agent != nil {
			a.agents[clusterKey].Close()
		}
		delete(a.agents, clusterKey)
	}
	agentMutex.Unlock()
}

func (a *Api) CloseAll() {
	for clusterKey, _ := range a.agents {
		a.Close(clusterKey)
	}
}

func (a *Api) getAgentConfig(orgId, clusterName string) (*types.AgentConfiguration, error) {
	logger := log.GetLogger()
	defer logger.Sync()

	agentCfg := &types.AgentConfiguration{}
	cfg := config.GetConfig()
	session, err := db.New(cfg.GetString(types.ServerDbCfgKey))
	if err != nil {
		logger.Error("failed to get db session", zap.Error(err))
		return nil, err
	}

	endpoint, err := session.GetClusterEndpoint(orgId, clusterName)
	if err != nil {
		logger.Error("failed to get db session", zap.Error(err))
		return nil, err
	}

	agentCfg.Address = endpoint
	agentCfg.Port = cfg.GetInt(types.AgentPortCfgKey)
	agentCfg.TlsEnabled = cfg.GetBool(types.AgentTlsEnabledCfgKey)
	certDataMap, err := a.vault.GetCert(context.TODO(), orgId, clusterName)
	if err != nil {
		logger.Error("failed get cert from vault", zap.Error(err))
		return nil, err
	}

	agentCfg.CaCert = certDataMap[types.ClientCertChainFileName]
	agentCfg.Key = certDataMap[types.ClientKeyFileName]
	agentCfg.Cert = certDataMap[types.ClientCertFileName]
	return agentCfg, err
}

func getClusterKey(orgId, clusterName string) string {
	return fmt.Sprintf("%s-%s", orgId, clusterName)
}
