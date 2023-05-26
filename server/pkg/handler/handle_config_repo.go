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

func (a *APIHandler) PostAgentRepository(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.RepositoryPostRequest
	if err := c.BindJSON(&req); err != nil {
		a.sendResponse(c, "Failed to parse config payload", err)
		return
	}

	if err := a.ConnectClient("1"); err != nil {
		a.setFailedResponse(c, "agent connection failed", err)
		return
	}

	agent := a.GetClient("1")
	if agent == nil {
		a.setFailedResponse(c, fmt.Sprintf("unregistered customer %v", "1"), errors.New(""))
	}

	_, err := agent.GetClient().RepositoryAdd(ctx, &agentpb.RepositoryAddRequest{
		PluginName: req.PluginName,
		RepoName:   req.RepoName,
		RepoUrl:    req.RepoUrl,
	})
	if err != nil {
		a.sendResponse(c, "failed to submit job", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: "submitted Job"})
}
func (a *APIHandler) PutAgentRepository(c *gin.Context) {
	a.PostAgentRepository(c)
}

func (a *APIHandler) DeleteAgentRepository(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.RepositoryPostRequest
	if err := c.BindJSON(&req); err != nil {
		a.sendResponse(c, "Failed to parse config payload", err)
		return
	}

	if err := a.ConnectClient("1"); err != nil {
		a.setFailedResponse(c, "agent connection failed", err)
		return
	}

	agent := a.GetClient("1")
	if agent == nil {
		a.setFailedResponse(c, fmt.Sprintf("unregistered customer %v", "1"), errors.New(""))
	}

	_, err := agent.GetClient().RepositoryDelete(ctx, &agentpb.RepositoryDeleteRequest{
		PluginName: req.PluginName,
		RepoName:   req.RepoName,
	})

	if err != nil {
		a.sendResponse(c, "failed to submit job", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: "submitted Job"})
}
