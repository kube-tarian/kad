package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	//	"path"
	"time"

	//"github.com/argoproj/argo-cd/v2/util/grpc"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	//	"github.com/kube-tarian/kad/integrator/agent/pkg/vaultservpb"
	"github.com/kube-tarian/kad/server/api"
	client "github.com/kube-tarian/kad/server/pkg/client"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	//"k8s.io/kubernetes/pkg/kubelet/secret"
)

func (a *APIHandler) PostStoreAgentCred(c *gin.Context) {
	a.log.Debugf("Install Deploy applicaiton api invocation started")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.StoreAgentCredRequest
	if err := c.BindJSON(&req); err != nil {
		a.setFailedResponse(c, "Failed to parse deploy payload", err)
		return
	}

	if err := a.ConnectClient(*req.CustomerId); err != nil {
		a.setFailedResponse(c, "agent connection failed", err)
		return
	}

	agent := a.GetClient(*req.CustomerId)
	if agent == nil {
		a.setFailedResponse(c, fmt.Sprintf("unregistered customer %v", 1), errors.New(""))
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

	a.log.Infof("response received from agent", response)
	a.log.Debugf("credentials is stored successfully")
}
func (a *APIHandler) PostStoreCred(c *gin.Context) {
	a.log.Debugf("Storing credentials started")

	var req api.StoreCredRequest
	if err := c.BindJSON(&req); err != nil {
		a.setFailedResponse(c, "Failed to parse payload", err)
		return
	}

	secretName := *req.Credname
	username := *req.Username
	password := *req.Password
	pathname := *req.Pathname
	dbcredspath := secretName + "/" + pathname

	// Store the credentials in Vault

	client := client.Vault{}

	if err := client.PostCredentials(secretName, username, password, dbcredspath); err != nil {

		a.setFailedResponse(c, "Failed to store credentials in Vault", err)
		return
	}
	// Set a success response
	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: "Credentials stored successfully",
	})

	a.log.Debugf("Credentials stored successfully in Vault")
}
func (a *APIHandler) GetStoreCred(c *gin.Context) {
	a.log.Debugf("Retrieving credentials started")
	_, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	secretName := c.Query("secretName")
	pathName := c.Query("dbcreds")
	dbcredspath := secretName + pathName

	// Retrieve the credentials from Vault
	client := client.Vault{}

	username, password, err := client.GetCredentials(secretName, dbcredspath)
	if err != nil {
		a.setFailedResponse(c, "Failed to retrieve credentials from Vault", err)
		return
	}

	// Set a success response with the username and password
	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: "Credentials retrieved successfully",
	})
	a.log.Infof("Username:%v,Password:%v", username, password)

	a.log.Debugf("Credentials retrieved successfully from Vault")
}

func (a *APIHandler) PostStoreSecret(c *gin.Context) {
	// Parse the request body to get the name and kubeconfig
	var requestBody api.StoreSecretRequest

	client := agentpb.NewAgentClient(&grpc.ClientConn{})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp := agentpb.StoreSecretRequest{
		Name:       *requestBody.Name,
		Kubeconfig: *requestBody.Kubeconfig,
	}

	//response,err:=agentpb.AgentClient.StoreSecret(ctx,resp,)
	response, err := client.StoreSecret(ctx, &resp)
	if err != nil {
		a.setFailedResponse(c, "failed to store credentials", err)
		return
	}
	fmt.Println("Response is ", response)
	if response.Status != "SUCCESS" {
		a.setFailedResponse(c, "failed to store credentials", err)
		return
	}

	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: "stored credentials"})

	a.log.Infof("response received from agent", response)
	a.log.Debugf("credentials is stored successfully")

}

func (a *APIHandler) GetStoreSecret(c *gin.Context) {

}
