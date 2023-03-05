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

func (a *APIHandler) PostConfigatorRepository(c *gin.Context) {
	a.log.Debugf("Add repository api invocation started")

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

	response, err := agent.GetClient().RepositoryAdd(ctx, &agentpb.RepositoryAddRequest{
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

	a.log.Infof("response received", response)
	a.log.Debugf("Add repository api invocation finished")
}
func (a *APIHandler) PutConfigatorRepository(c *gin.Context) {
	a.log.Debugf("Update repositoy from plugin api invocation started")

	a.PostConfigatorRepository(c)
	a.log.Debugf("Delete repositoy from plugin api invocation finished")
}

func (a *APIHandler) DeleteConfigatorRepository(c *gin.Context) {
	a.log.Debugf("Delete repository from plugin api invocation started")

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

	response, err := agent.GetClient().RepositoryDelete(ctx, &agentpb.RepositoryDeleteRequest{
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

	a.log.Infof("response received", response)
	a.log.Debugf("Delete repository from plugin api invocation finished")
}
