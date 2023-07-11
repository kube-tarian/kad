package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/config"
	"github.com/kube-tarian/kad/server/pkg/db"
	"github.com/kube-tarian/kad/server/pkg/log"
	"github.com/kube-tarian/kad/server/pkg/model"
	"github.com/kube-tarian/kad/server/pkg/types"
)

func (a *APIHandler) PostAgentEndpoint(c *gin.Context) {
	logger := log.GetLogger()
	defer logger.Sync()

	logger.Info("Register agent api invocation started")

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

	cfg := config.GetConfig()
	session, err := db.New(cfg.GetString("server.db"))
	if err != nil {
		a.setFailedResponse(c, "failed to get db session", nil)
		logger.Error("failed to get db session", zap.Error(err))
		return
	}

	err = session.RegisterCluster(customerId, customerId, endpoint)
	if err != nil {
		a.setFailedResponse(c, "failed to store data", nil)
		logger.Error("failed to get db session", zap.Error(err))
		return
	}

	err = a.vault.PutCert(c, customerId, customerId,
		fileContentsMap[types.ClientCertChainFileName],
		fileContentsMap[types.ClientKeyFileName],
		fileContentsMap[types.ClientCertFileName],
	)

	if err != nil {
		a.setFailedResponse(c, "failed to register", nil)
		logger.Error("failed to store cert in vault", zap.Error(err))
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	logger.Info("registered new agent endpoint")
}

func (a *APIHandler) GetAgentEndpoint(c *gin.Context) {
	logger := log.GetLogger()
	defer logger.Sync()

	logger.Debug("get all registered agents api")
	c.IndentedJSON(http.StatusOK, &model.AgentsResponse{})
	logger.Debug("get all registered agents api")
}

func (a *APIHandler) PutAgentEndpoint(c *gin.Context) {
	logger := log.GetLogger()
	defer logger.Sync()

	logger.Debug("update register agent api invocation started")
	var req api.AgentRequest
	if err := c.BindJSON(&req); err != nil {
		a.setFailedResponse(c, "Failed to parse deploy payload", err)
		return
	}

	//TODO Update in DB and internal cache
	c.Writer.WriteHeader(http.StatusOK)
	logger.Debug("update register agent api invocation finished")
}

func (a *APIHandler) PostAgentApps(c *gin.Context) {
	//var req api.AgentAppsRequest
	logger := log.GetLogger()
	defer logger.Sync()

	//if err := c.BindJSON(&req); err != nil {
	//	a.setFailedResponse(c, "failed to parse apps payload", err)
	//	return
	//}

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		a.setFailedResponse(c, "failed to read payload", err)
		return
	}

	fmt.Println("body is", string(jsonData))
	syncData := &agentpb.SyncRequest{
		Type: "app-data",
		Data: string(jsonData),
	}

	customerId := c.GetHeader("customer_id")
	if customerId == "" {
		a.setFailedResponse(c, "missing customer id header", errors.New(""))
		return
	}

	if err := a.ConnectClient(customerId); err != nil {
		a.setFailedResponse(c, "agent connection failed", err)
		return
	}

	agent := a.GetClient(customerId)
	if agent == nil {
		a.setFailedResponse(c, fmt.Sprintf("unregistered customer %v", "1"), errors.New(""))
	}

	response, err := agent.GetClient().Sync(context.Background(), syncData)

	logger.Debug("response ")
	fmt.Printf("response %+v, err: %v\n", response, err)
	//fmt.Printf("%+v\n", req.Tools)
}
