package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/log"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
)

func (a *APIHandler) PostAgentDeploy(c *gin.Context) {
	logger := log.GetLogger()
	defer logger.Sync()

	logger.Debug("deploying application")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.DeployerPostRequest
	if err := c.BindJSON(&req); err != nil {
		a.setFailedResponse(c, "failed to parse deploy payload", err)
		return
	}

	if err := a.ConnectClient("1"); err != nil {
		a.setFailedResponse(c, "agent connection failed", err)
		return
	}

	agent := a.GetClient("1")
	if agent == nil {
		a.setFailedResponse(c, fmt.Sprintf("unregistered customer %v", 1), errors.New(""))
	}

	response, err := agent.GetClient().DeployerAppInstall(
		ctx,
		&agentpb.ApplicationInstallRequest{
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

	logger.Debug("deployed application successfully")
}

func (a *APIHandler) PutAgentDeploy(c *gin.Context) {
	a.PostAgentDeploy(c)
}

func (a *APIHandler) DeleteAgentDeploy(c *gin.Context) {
	logger := log.GetLogger()
	logger.Debug("deleting application")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.DeployerDeleteRequest
	if err := c.BindJSON(&req); err != nil {
		a.setFailedResponse(c, "failed to parse deploy payload", err)
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

	response, err := agent.GetClient().DeployerAppDelete(
		ctx,
		&agentpb.ApplicationDeleteRequest{
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

	logger.Debug("application deleted successfully")
}

func (a *APIHandler) sendResponse(c *gin.Context, msg string, err error) {
	c.IndentedJSON(http.StatusInternalServerError, &api.Response{
		Status:  "FAILED",
		Message: fmt.Sprintf("%s, %v", msg, err),
	})
}

func toString(resp *agentpb.JobResponse) string {
	return fmt.Sprintf("Workflow details, ID: %v, RUN-ID: %v, NAME: %v", resp.Id, resp.RunID, resp.WorkflowName)
}
