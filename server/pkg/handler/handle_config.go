package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"google.golang.org/protobuf/types/known/anypb"
)

func (a *APIHandler) PostConfig(c *gin.Context) {
	//TODO get address from database based on CustomerInfo
	a.log.Infof("config api invocation started")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.ConfigRequestPayload
	if err := c.BindJSON(&req); err != nil {
		a.setFailedResponse(c, "failed to parse config payload", err)
		return
	}

	payload, err := json.Marshal(req.Payload)
	if err != nil {
		a.setFailedResponse(c, "Config request preparation failed", err)
		return
	}

	agent := a.GetClient(req.CustomerId)
	if agent == nil {
		a.setFailedResponse(c, fmt.Sprintf("unregistered customer %v", req.CustomerId), errors.New(""))
	}

	// TODO: currently climon payload is submitted to temporal via agent.
	// This flow has to modified as per the understanding.
	response, err := agent.SubmitJob(
		ctx,
		&agentpb.JobRequest{
			Operation: req.Operation,
			Payload:   &anypb.Any{Value: payload},
		},
	)
	if err != nil {
		a.setFailedResponse(c, "failed to submit job", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: "submitted Job"})

	a.log.Infof("response received", response)
	a.log.Infof("config api invocation finished")
}
