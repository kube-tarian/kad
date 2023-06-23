package handler

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"

	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/client"
	"github.com/kube-tarian/kad/server/pkg/config"
	"github.com/kube-tarian/kad/server/pkg/db"
	"github.com/kube-tarian/kad/server/pkg/log"
	"github.com/kube-tarian/kad/server/pkg/types"
)

type APIHandler struct {
	agents map[string]*client.Agent
}

var (
	agentMutex sync.RWMutex
)

func NewAPIHandler() (*APIHandler, error) {
	return &APIHandler{
		agents: make(map[string]*client.Agent),
	}, nil
}

func (a *APIHandler) ConnectClient(customerId string) error {
	if _, ok := a.agents[customerId]; ok {
		return nil
	}

	agentCfg, err := getAgentConfig(customerId)
	if err != nil {
		return err
	}

	logger := log.GetLogger()
	defer logger.Sync()
	agent, err := client.NewAgent(agentCfg)
	if err != nil {
		logger.Error("failed to connect agent internal error", zap.Error(err))
		return err
	}

	agentMutex.Lock()
	a.agents[customerId] = agent
	agentMutex.Unlock()
	return err
}

func (a *APIHandler) GetClient(customerId string) *client.Agent {
	agentMutex.RLock()
	if agent, ok := a.agents[customerId]; ok && agent != nil {
		return agent
	}
	agentMutex.RUnlock()
	return nil
}

func (a *APIHandler) Close(customerId string) {
	agent := a.GetClient(customerId)
	if agent == nil {
		return
	}

	agentMutex.Lock()
	a.agents[customerId].Close()
	delete(a.agents, customerId)
	agentMutex.Unlock()
}

func (a *APIHandler) CloseAll() {
	for customerId, _ := range a.agents {
		a.Close(customerId)
	}
}

func (a *APIHandler) GetApiDocs(c *gin.Context) {
	swagger, err := api.GetSwagger()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.IndentedJSON(http.StatusOK, swagger)
}

func (a *APIHandler) GetStatus(c *gin.Context) {
	c.String(http.StatusOK, "")
}

func getAgentConfig(customerID string) (*types.AgentConfiguration, error) {
	agentCfg := &types.AgentConfiguration{}
	cfg := config.GetConfig()
	session, err := db.New(cfg.GetString("server.db"))
	if err != nil {
		logrus.Error("failed to get db session", err)
		return nil, err
	}

	agentInfo, err := session.GetAgentInfo(customerID)
	if err != nil {
		logrus.Error("failed to get db session", err)
		return nil, err
	}

	agentCfg.Address = agentInfo.Endpoint
	agentCfg.Port = cfg.GetInt("agent.Port")
	agentCfg.TlsEnabled = cfg.GetBool("agent.tlsEnabled")
	agentCfg.CaCert, agentCfg.Key, agentCfg.Cert, err = client.GetCaptenClusterCertificate(context.TODO(), customerID)
	return agentCfg, err
}
