package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/client"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
)

func (s *APIHanlder) PostConfigatorCluster(c *gin.Context) {
	s.log.Debugf("Add cluster api invocation started")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.ClusterRequest
	if err := c.BindJSON(&req); err != nil {
		s.sendResponse(c, "Failed to parse config payload", err)
		return
	}

	agentClient, err := client.NewAgent(s.log)
	if err != nil {
		s.log.Errorf("failed to connect agent internal error", err)
		s.sendResponse(c, "agent connection failed", err)
		return
	}
	defer agentClient.Close()

	response, err := agentClient.GetClient().ClusterAdd(ctx, &agentpb.ClusterRequest{
		PluginName:  req.PluginName,
		ClusterName: req.ClusterName,
	})
	if err != nil {
		s.sendResponse(c, "failed to submit job", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: "submitted Job"})

	s.log.Infof("response received", response)
	s.log.Debugf("Add cluster api invocation finished")
}

func (s *APIHanlder) DeleteConfigatorCluster(c *gin.Context) {
	s.log.Debugf("Delete cluster from plugin api invocation started")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.ClusterRequest
	if err := c.BindJSON(&req); err != nil {
		s.sendResponse(c, "Failed to parse config payload", err)
		return
	}

	agentClient, err := client.NewAgent(s.log)
	if err != nil {
		s.log.Errorf("failed to connect agent internal error", err)
		s.sendResponse(c, "agent connection failed", err)
		return
	}
	defer agentClient.Close()

	response, err := agentClient.GetClient().ClusterDelete(ctx, &agentpb.ClusterRequest{
		PluginName:  req.PluginName,
		ClusterName: req.ClusterName,
	})
	if err != nil {
		s.sendResponse(c, "failed to submit job", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: "submitted Job"})

	s.log.Infof("response received", response)
	s.log.Debugf("Delete cluster from plugin api invocation finished")
}
