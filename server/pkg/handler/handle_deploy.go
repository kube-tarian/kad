package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/client"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
)

func (s *APIHanlder) PostDeployer(c *gin.Context) {
	s.log.Debugf("Install Deploy applicaiton api invocation started")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.DeployerPostRequest
	if err := c.BindJSON(&req); err != nil {
		s.sendResponse(c, "Failed to parse deploy payload", err)
		return
	}

	agentClient, err := client.NewAgent(s.log)
	if err != nil {
		s.log.Errorf("failed to connect agent internal error", err)
		s.sendResponse(c, "agent connection failed", err)
		return
	}
	defer agentClient.Close()

	response, err := agentClient.GetClient().DeployerAppInstall(
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
		s.sendResponse(c, "failed to submit job", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: toString(response)})

	s.log.Infof("response received", response)
	s.log.Debugf("Install Deploy application api invocation finished")
}

func (s *APIHanlder) PutDeployer(c *gin.Context) {
	s.log.Debugf("Update Deploy application api invocation started")
	s.PostDeployer(c)
	s.log.Debugf("Update Deploy application api invocation finished")
}

func (s *APIHanlder) DeleteDeployer(c *gin.Context) {
	s.log.Debugf("Delete climon application api invocation started")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.DeployerDeleteRequest
	if err := c.BindJSON(&req); err != nil {
		s.sendResponse(c, "Failed to parse deploy payload", err)
		return
	}

	agentClient, err := client.NewAgent(s.log)
	if err != nil {
		s.log.Errorf("failed to connect agent internal error", err)
		s.sendResponse(c, "agent connection failed", err)
		return
	}
	defer agentClient.Close()

	response, err := agentClient.GetClient().DeployerAppDelete(
		ctx,
		&agentpb.ApplicationDeleteRequest{
			PluginName:  req.PluginName,
			Namespace:   req.Namespace,
			ReleaseName: req.ReleaseName,
			Timeout:     uint32(req.Timeout),
		},
	)
	if err != nil {
		s.sendResponse(c, "failed to submit job", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: toString(response)})

	s.log.Infof("response received", response)
	s.log.Debugf("Delete climon application api invocation finished")
}

func (s *APIHanlder) sendResponse(c *gin.Context, msg string, err error) {
	s.log.Errorf("failed to submit job", err)
	c.IndentedJSON(http.StatusInternalServerError, &api.Response{
		Status:  "FAILED",
		Message: fmt.Sprintf("%s, %v", msg, err),
	})
}

func toString(resp *agentpb.JobResponse) string {
	return fmt.Sprintf("Workflow details, ID: %v, RUN-ID: %v, NAME: %v", resp.Id, resp.RunID, resp.WorkflowName)
}
