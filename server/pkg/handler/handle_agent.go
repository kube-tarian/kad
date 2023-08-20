package handler

import (
	"errors"
	"net/http"

	"github.com/kube-tarian/kad/server/pkg/credential"

	"github.com/gin-gonic/gin"

	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/model"
	"github.com/kube-tarian/kad/server/pkg/types"
)

func (a *APIHandler) PostAgentEndpoint(c *gin.Context) {
	a.log.Info("Register agent api invocation started")

	//var req api.AgentRequest
	customerId := c.GetHeader("customer_id")
	if customerId == "" {
		a.setFailedResponse(c, "missing customer id header", errors.New(""))
		return
	}

	endpoint := c.GetHeader("endpoint")
	if endpoint == "" {
		a.setFailedResponse(c, "missing endpoint in header", errors.New(""))
		return
	}

	fileContentsMap, err := a.getFileContent(c, map[string]string{
		"ca_crt":     types.ClientCertChainFileName,
		"client_crt": types.ClientCertFileName,
		"client_key": types.ClientKeyFileName})
	if err != nil {
		a.setFailedResponse(c, "failed to register agent", err)
		return
	}

	err = a.serverStore.AddCluster(customerId, customerId, customerId, endpoint)
	if err != nil {
		a.setFailedResponse(c, "failed to store data", nil)
		a.log.Error("failed to get db session", err)
		return
	}

	err = credential.PutClusterCerts(c, customerId,
		fileContentsMap[types.ClientCertChainFileName],
		fileContentsMap[types.ClientKeyFileName],
		fileContentsMap[types.ClientCertFileName],
	)

	if err != nil {
		a.setFailedResponse(c, "failed to register", nil)
		a.log.Error("failed to store cert in vault", err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	a.log.Info("registered new agent endpoint")
}

func (a *APIHandler) GetAgentEndpoint(c *gin.Context) {
	a.log.Debug("get all registered agents api")
	c.IndentedJSON(http.StatusOK, &model.AgentsResponse{})
	a.log.Debug("get all registered agents api")
}

func (a *APIHandler) PutAgentEndpoint(c *gin.Context) {
	a.log.Debug("update register agent api invocation started")
	var req api.AgentRequest
	if err := c.BindJSON(&req); err != nil {
		a.setFailedResponse(c, "Failed to parse deploy payload", err)
		return
	}

	//TODO Update in DB and internal cache
	c.Writer.WriteHeader(http.StatusOK)
	a.log.Debug("update register agent api invocation finished")
}

func (a *APIHandler) PostAgentApps(c *gin.Context) {
	a.setFailedResponse(c, "not implemented", errors.New(""))
}
