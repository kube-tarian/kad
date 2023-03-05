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

func (a *APIHandler) PostConfigatorProject(c *gin.Context) {
	a.log.Debugf("Add project api invocation started")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.ProjectPostRequest
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

	response, err := agent.GetClient().ProjectAdd(ctx, &agentpb.ProjectAddRequest{
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

	a.log.Infof("response received", response)
	a.log.Debugf("Add project api invocation finished")
}
func (a *APIHandler) PutConfigatorProject(c *gin.Context) {
	a.log.Debugf("Update project from plugin api invocation started")

	a.PostConfigatorProject(c)
	a.log.Debugf("Delete project from plugin api invocation finished")
}

func (a *APIHandler) DeleteConfigatorProject(c *gin.Context) {
	a.log.Debugf("Delete project from plugin api invocation started")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.ProjectDeleteRequest
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

	response, err := agent.GetClient().ProjectDelete(ctx, &agentpb.ProjectDeleteRequest{
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

	a.log.Infof("response received", response)
	a.log.Debugf("Delete project from plugin api invocation finished")
}
