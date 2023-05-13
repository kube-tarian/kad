package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/client"
	"github.com/kube-tarian/kad/server/pkg/db"
	"github.com/kube-tarian/kad/server/pkg/model"
	"github.com/kube-tarian/kad/server/pkg/types"
)

func (a *APIHandler) PostRegisterAgent(c *gin.Context) {
	a.log.Infof("Register agent api invocation started")

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

	session, err := caasandra.New()
	if err != nil {
		a.setFailedResponse(c, "failed to get db session", nil)
		a.log.Error("failed to get db session", err)
		return
	}

	err = session.RegisterEndpoint(customerId, endpoint, map[string]string{})
	if err != nil {
		a.setFailedResponse(c, "failed to store data", nil)
		a.log.Error("failed to get db session", err)
		return
	}

	vaultSession, err := client.NewVault()
	if err != nil {
		a.setFailedResponse(c, "failed to register", nil)
		a.log.Error("failed to create vault session", err)
		return
	}

	err = vaultSession.PutCert("secret",
		fileContentsMap[types.ClientCertChainFileName],
		fileContentsMap[types.ClientCertFileName],
		fileContentsMap[types.ClientKeyFileName],
		customerId)

	if err != nil {
		a.setFailedResponse(c, "failed to register", nil)
		a.log.Error("failed to store cert in vault", err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	a.log.Infof("Register agent api invocation finished")
}

func (a *APIHandler) GetRegisterAgent(c *gin.Context) {
	a.log.Infof("Get all registered agents api invocation started")

	//TODO Get all agents from DB

	c.IndentedJSON(http.StatusOK, &model.AgentsResponse{})

	a.log.Infof("Get all registered agents api invocation finished")
}

func (a *APIHandler) PutRegisterAgent(c *gin.Context) {
	a.log.Infof("Update register agent api invocation started")

	var req api.AgentRequest
	if err := c.BindJSON(&req); err != nil {
		a.setFailedResponse(c, "Failed to parse deploy payload", err)
		return
	}

	//TODO Update in DB and internal cache
	c.Writer.WriteHeader(http.StatusOK)
	a.log.Infof("Update register agent api invocation finished")
}
