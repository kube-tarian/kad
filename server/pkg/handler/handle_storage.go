package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/log"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"net/http"
	"time"
)

func (a *APIHandler) PostAgentSecret(c *gin.Context) {
	logger := log.GetLogger()
	defer logger.Sync()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.StoreCredRequest
	if err := c.BindJSON(&req); err != nil {
		a.setFailedResponse(c, "failed to parse store-cred payload", err)
		return
	}

	if err := a.ConnectClient(*req.CustomerId); err != nil {
		a.setFailedResponse(c, "failed to connect agent", err)
		return
	}

	agent := a.GetClient(*req.CustomerId)
	if agent == nil {
		a.setFailedResponse(c, fmt.Sprintf("unregistered customer %v", *req.CustomerId), errors.New(""))
	}

	response, err := agent.GetClient().StoreCred(
		ctx,
		&agentpb.StoreCredRequest{
			Credname: *req.Credname,
			Username: *req.Username,
			Password: *req.Password,
		},
	)
	if err != nil {
		a.setFailedResponse(c, "failed to store credentials", err)
		return
	}

	if response.Status != "SUCCESS" {
		a.setFailedResponse(c, "failed to store credentials", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: "stored credentials"})

	logger.Debug("credentials is stored successfully")
}
