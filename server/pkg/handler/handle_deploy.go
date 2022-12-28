package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kube-tarian/kad/server/pkg/model"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"google.golang.org/protobuf/types/known/anypb"
)

func (s *APIHanlder) PostDeploy(c *gin.Context) {
	//TODO get address from database based on CustomerInfo
	s.log.Infof("deploy api invocation started")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req model.DeployPayload
	if err := c.BindJSON(&req); err != nil {
		sendResponse(c, "Failed to parse deploy payload", err)
		return
	}

	payload, err := json.Marshal(req.Payload)
	if err != nil {
		sendResponse(c, "Deploy request prepration failed", err)
		return
	}

	response, err := s.client.SubmitJob(
		ctx,
		&agentpb.JobRequest{
			Operation: req.Operation,
			Payload:   &anypb.Any{Value: payload},
		},
	)
	if err != nil {
		sendResponse(c, "failed to submit job", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &model.DeployResponse{
		Status:  "SUCCESS",
		Message: "submitted Job"})

	s.log.Infof("response received", response)
	s.log.Infof("deploy api invocation finished")
}

func sendResponse(c *gin.Context, msg string, err error) {
	log.Println("failed to submit job", err)
	c.IndentedJSON(http.StatusInternalServerError, &model.DeployResponse{
		Status:  "FAILED",
		Message: fmt.Sprintf("%s, %v", msg, err),
	})
}
