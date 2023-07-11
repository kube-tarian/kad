package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/intelops/go-common/credentials"
	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/log"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
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

	serviceCred := credentials.ServiceCredential{
		UserName: *req.Username,
		Password: *req.Password,
	}
	serviceCredMap := credentials.PrepareServiceCredentialMap(serviceCred)
	response, err := agent.GetClient().StoreCredential(ctx,
		&agentpb.StoreCredentialRequest{
			CredentialType: credentials.ServiceUserCredentialType,
			CredEntityName: *req.Credname,
			CredIdentifier: *req.Username,
			Credential:     serviceCredMap,
		},
	)
	if err != nil {
		a.setFailedResponse(c, "failed to store credentials", err)
		return
	}

	if response.Status != agentpb.StatusCode_OK {
		a.setFailedResponse(c, "failed to store credentials", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: "stored credentials"})

	logger.Debug("credentials is stored successfully")
}
