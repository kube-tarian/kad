package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/client"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"google.golang.org/protobuf/types/known/anypb"
)

func (s *APIHanlder) PostConfig(c *gin.Context) {
	//TODO get address from database based on CustomerInfo
	s.log.Infof("config api invocation started")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.ConfigRequestPayload
	if err := c.BindJSON(&req); err != nil {
		s.sendResponse(c, "Failed to parse config payload", err)
		return
	}

	payload, err := json.Marshal(req.Payload)
	if err != nil {
		s.sendResponse(c, "Config request prepration failed", err)
		return
	}

	// TODO: currently climon payload is submitted to temporal via agent.
	// This flow has to modified as per the understanding.
	agentClient, err := client.NewAgent(s.log)
	if err != nil {
		s.log.Errorf("failed to connect agent internal error", err)
		s.sendResponse(c, "agent connection failed", err)
		return
	}
	defer agentClient.Close()

	response, err := agentClient.SubmitJob(
		ctx,
		&agentpb.JobRequest{
			Operation: req.Operation,
			Payload:   &anypb.Any{Value: payload},
		},
	)
	if err != nil {
		s.sendResponse(c, "failed to submit job", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: "submitted Job"})

	s.log.Infof("response received", response)
	s.log.Infof("config api invocation finished")
}
