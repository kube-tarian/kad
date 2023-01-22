package handler

import (
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
