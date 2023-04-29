package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/client"
	"github.com/kube-tarian/kad/server/pkg/db"
	"github.com/kube-tarian/kad/server/pkg/logging"
	"github.com/kube-tarian/kad/server/pkg/types"
)

type APIHandler struct {
	log    logging.Logger
	agents map[string]*client.Agent
}

var (
	agentMutex sync.RWMutex
)

func NewAPIHandler(log logging.Logger) (*APIHandler, error) {
	return &APIHandler{
		log:    log,
		agents: make(map[string]*client.Agent),
	}, nil
}

func (a *APIHandler) ConnectClient(customerId string) error {
	if _, ok := a.agents[customerId]; ok {
		return nil
	}

	agentCfg, err := fetchConfiguration(customerId)
	if err != nil {
		return err
	}

	agent, err := client.NewAgent(a.log, agentCfg)
	if err != nil {
		a.log.Errorf("failed to connect agent internal error", err)
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

func fetchConfiguration(customerID string) (*types.AgentConfiguration, error) {
	cfg := &types.AgentConfiguration{}
	session, err := db.New()
	if err != nil {
		logrus.Error("failed to get db session", err)
		return nil, err
	}

	agentInfo, err := session.GetAgentInfo(customerID)
	if err != nil {
		logrus.Error("failed to get db session", err)
		return nil, err
	}

	cfg.Address = agentInfo.Endpoint
	cfg.Port, err = strconv.Atoi(os.Getenv("AGENT_PORT"))
	if err != nil {
		return nil, fmt.Errorf("failed to convert agent port to int, port got: %v", os.Getenv("AGENT_PORT"))
	}

	vaultSession, err := client.NewVault()
	if err != nil {
		return nil, fmt.Errorf("failed to get vault session: %w", err)
	}

	vaultMap, err := vaultSession.GetCert("secret", customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get certs")
	}

	cfg.CaCert = vaultMap[types.ClientCertChainFileName]
	cfg.Cert = vaultMap[types.ClientCertFileName]
	cfg.Key = vaultMap[types.ClientKeyFileName]
	return cfg, err
}
