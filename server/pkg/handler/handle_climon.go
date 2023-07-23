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

func (a *APIHandler) PostAgentClimondeploy(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.ClimonPostRequest
	if err := c.BindJSON(&req); err != nil {
		a.setFailedResponse(c, "failed to parse deploy payload", err)
		return
	}

	agent, err := a.agentHandler.GetAgent("", "")
	if err != nil {
		a.setFailedResponse(c, fmt.Sprintf("unregistered customer %v", "1"), errors.New(""))
		return
	}

	response, err := agent.GetClient().ClimonAppInstall(
		ctx,
		&agentpb.ClimonInstallRequest{
			PluginName:  req.PluginName,
			RepoName:    req.RepoName,
			RepoUrl:     req.RepoUrl,
			ChartName:   req.ChartName,
			Namespace:   req.Namespace,
			ReleaseName: req.ReleaseName,
			Timeout:     uint32(req.Timeout),
		},
	)
	if err != nil {
		a.setFailedResponse(c, "failed to submit job", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: toString(response)})

}

func (a *APIHandler) PutAgentClimondeploy(c *gin.Context) {
	a.PostAgentClimondeploy(c)
}

func (a *APIHandler) DeleteAgentClimondeploy(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.ClimonDeleteRequest
	if err := c.BindJSON(&req); err != nil {
		a.setFailedResponse(c, "Failed to parse deploy payload", err)
		return
	}

	agent, err := a.agentHandler.GetAgent("", "")
	if err != nil {
		a.setFailedResponse(c, fmt.Sprintf("unregistered customer %v", "1"), errors.New(""))
		return
	}

	response, err := agent.GetClient().ClimonAppDelete(
		ctx,
		&agentpb.ClimonDeleteRequest{
			PluginName:  req.PluginName,
			Namespace:   req.Namespace,
			ReleaseName: req.ReleaseName,
			Timeout:     uint32(req.Timeout),
		},
	)
	if err != nil {
		a.setFailedResponse(c, "failed to submit job", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: toString(response)})

}
