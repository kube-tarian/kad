package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
)

func (a *APIHandler) PostAgentProject(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.ProjectPostRequest
	if err := c.BindJSON(&req); err != nil {
		a.sendResponse(c, "Failed to parse config payload", err)
		return
	}

	agent, err := a.agentHandler.GetAgent("", "")
	if err != nil {
		a.setFailedResponse(c, fmt.Sprintf("unregistered customer %v", "1"), errors.New(""))
		return
	}

	_, err = agent.GetClient().ProjectAdd(ctx, &agentpb.ProjectAddRequest{
		PluginName:  req.PluginName,
		ProjectName: req.ProjectName,
	})
	if err != nil {
		a.sendResponse(c, "failed to submit job", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: "submitted Job"})

}
func (a *APIHandler) PutAgentProject(c *gin.Context) {
	a.PostAgentProject(c)
}

func (a *APIHandler) DeleteAgentProject(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.ProjectDeleteRequest
	if err := c.BindJSON(&req); err != nil {
		a.sendResponse(c, "Failed to parse config payload", err)
		return
	}

	agent, err := a.agentHandler.GetAgent("", "")
	if err != nil {
		a.setFailedResponse(c, fmt.Sprintf("unregistered customer %v", "1"), errors.New(""))
		return
	}

	_, err = agent.GetClient().ProjectDelete(ctx, &agentpb.ProjectDeleteRequest{
		PluginName:  req.PluginName,
		ProjectName: req.ProjectName,
	})
	if err != nil {
		a.sendResponse(c, "failed to submit job", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: "submitted Job"})
}
