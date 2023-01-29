package handler

import (
	"github.com/kube-tarian/kad/server/pkg/db"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/model"
)

func (s *APIHanlder) PostRegisterAgent(c *gin.Context) {
	s.log.Infof("Register agent api invocation started")

	var req api.AgentRequest
	if err := c.BindJSON(&req); err != nil {
		s.sendResponse(c, "Failed to parse deploy payload", err)
		return
	}

	//TODO Save in DB and internal cache
	session, err := db.New()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, &model.DeployResponse{
			Status:  "FAILED",
			Message: "failed to get db session"})
		logrus.Error("failed to get db session", err)
		return
	}

	err = session.RegisterEndpoint(req.CustomerId, req.Endpoint)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, &model.DeployResponse{
			Status:  "FAILED",
			Message: "failed to store data"})
		logrus.Error("failed to get db session", err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	s.log.Infof("Register agent api invocation finished")
}

func (s *APIHanlder) GetRegisterAgent(c *gin.Context) {
	s.log.Infof("Get all registered agents api invocation started")

	//TODO Get all agents from DB

	c.IndentedJSON(http.StatusOK, &model.AgentsResponse{})

	s.log.Infof("Get all registered agents api invocation finished")
}

func (s *APIHanlder) PutRegisterAgent(c *gin.Context) {
	s.log.Infof("Update register agent api invocation started")

	var req api.AgentRequest
	if err := c.BindJSON(&req); err != nil {
		s.sendResponse(c, "Failed to parse deploy payload", err)
		return
	}

	//TODO Update in DB and internal cache

	c.Writer.WriteHeader(http.StatusOK)
	s.log.Infof("Update register agent api invocation finished")
}
