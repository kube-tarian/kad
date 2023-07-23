package handler

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/agent"
	"github.com/kube-tarian/kad/server/pkg/store"
)

type APIHandler struct {
	agentHandler *agent.AgentHandler
	serverStore  store.ServerStore
	log          logging.Logger
}

var (
	agentMutex sync.RWMutex
)

func NewAPIHandler(log logging.Logger, serverStore store.ServerStore) (*APIHandler, error) {
	return &APIHandler{
		log:          log,
		agentHandler: agent.NewAgentHandler(serverStore),
	}, nil
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

func (a *APIHandler) Close() {
	a.agentHandler.Close()
}
