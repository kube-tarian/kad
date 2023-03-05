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

func (a *APIHandler) PostClimon(c *gin.Context) {
	a.log.Debugf("Install climon application api invocation started")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.ClimonPostRequest
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

	a.log.Infof("response received", response)
	a.log.Debugf("Install climon application api invocation finished")
}

func (a *APIHandler) PutClimon(c *gin.Context) {
	a.log.Debugf("Update climon application api invocation started")

	a.PostClimon(c)
	a.log.Debugf("Update climon application api invocation finished")
}

func (a *APIHandler) DeleteClimon(c *gin.Context) {
	a.log.Debugf("Delete climon application api invocation started")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.ClimonDeleteRequest
	if err := c.BindJSON(&req); err != nil {
		a.setFailedResponse(c, "Failed to parse deploy payload", err)
		return
	}

	agent := a.GetClient("1")
	if agent == nil {
		a.setFailedResponse(c, fmt.Sprintf("unregistered customer %v", req.CustomerId), errors.New(""))
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

	a.log.Infof("response received", response)
	a.log.Debugf("Delete climon application api invocation finished")
}
