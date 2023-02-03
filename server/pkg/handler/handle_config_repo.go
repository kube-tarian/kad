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

func (s *APIHanlder) PostConfigatorRepository(c *gin.Context) {
	s.log.Debugf("Add repository api invocation started")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.RepositoryPostRequest
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

	response, err := agentClient.GetClient().RepositoryAdd(ctx, &agentpb.RepositoryAddRequest{
		PluginName: req.PluginName,
		RepoName:   req.RepoName,
		RepoUrl:    req.RepoUrl,
	})
	if err != nil {
		s.sendResponse(c, "failed to submit job", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: "submitted Job"})

	s.log.Infof("response received", response)
	s.log.Debugf("Add repository api invocation finished")
}
func (s *APIHanlder) PutConfigatorRepository(c *gin.Context) {
	s.log.Debugf("Update repositoy from plugin api invocation started")

	s.PostConfigatorRepository(c)
	s.log.Debugf("Delete repositoy from plugin api invocation finished")
}

func (s *APIHanlder) DeleteConfigatorRepository(c *gin.Context) {
	s.log.Debugf("Delete repository from plugin api invocation started")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.RepositoryPostRequest
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

	response, err := agentClient.GetClient().RepositoryDelete(ctx, &agentpb.RepositoryDeleteRequest{
		PluginName: req.PluginName,
		RepoName:   req.RepoName,
	})
	if err != nil {
		s.sendResponse(c, "failed to submit job", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: "submitted Job"})

	s.log.Infof("response received", response)
	s.log.Debugf("Delete repository from plugin api invocation finished")
}
