package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kube-tarian/kad/server/api"
	"github.com/kube-tarian/kad/server/pkg/client"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
)

func (a *APIHandler) PostStoreAgentCred(c *gin.Context) {
	a.log.Debugf("Install Deploy applicaiton api invocation started")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	//var req api.StoreCredRequest
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
	_, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req api.StoreCredRequest
	if err := c.BindJSON(&req); err != nil {
		a.setFailedResponse(c, "Failed to parse payload", err)
		return
	}

	secretName := *req.Credname
	username := *req.Username
	password := *req.Password
	astradbcreds := "secretname/astrabdcreds"

	// Store the credentials in Vault
	client := client.Vault{}
	// client.GetClient().StoreCred(ctx,&agentpb.StoreCredRequest{
	// 	username: *req.Username,

	// })
	if err := client.PostCredentials(secretName, username, password, astradbcreds); err != nil {
		a.setFailedResponse(c, "Failed to store credentials in Vault", err)
		return
	}
	// if err := client.PostCredentials(secretName, username, password, astradbcreds); err != nil {
	//     a.setFailedResponse(c, "Failed to store credentials in Vault", err)
	//     return
	// }
	//   //client.NewVault()
	// Set a success response
	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: "Credentials stored successfully",
	})

	a.log.Debugf("Credentials stored successfully in Vault")
}

func (a *APIHandler) GetCredentials(c gin.Context) {
	a.log.Debugf("Retrieving credentials started")
	_, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	var req api.StoreCredRequest
	if err := c.BindJSON(&req); err != nil {
		a.setFailedResponse(&c, "Failed to parse payload", err)
		return
	}

	secretName := *req.Credname
	astradbcreds := "secretname/astrabdcreds"

	// Retrieve the credentials from Vault
	client := client.Vault{}
	username, password, err := client.GetCredentials(secretName, astradbcreds)
	if err != nil {
		a.setFailedResponse(&c, "Failed to retrieve credentials from Vault", err)
		return
	}
	c.IndentedJSON(http.StatusOK, &api.Response{
		Status:  "SUCCESS",
		Message: "Credentials retrieved successfully",
	})
	// Return the credentials
	// c.IndentedJSON(http.StatusOK, &api.CredentialsResponse{
	// 	Status:  "SUCCESS",
	// 	Message: "Credentials retrieved successfully",
	// 	Data: &api.CredentialsData{
	// 		Username: &username,
	// 		Password: &password,
	// 	},
	// })
	a.log.Infof("Username:%v,Password:%v", username, password)
	a.log.Debugf("Credentials retrieved successfully from Vault")
}

// func (a *APIHandler) PostStoreCred(c *gin.Context) {
// 	a.log.Debugf("Store credentials API invocation started")
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
// 	defer cancel()

// 	var req api.StoreCredRequest
// 	if err := c.BindJSON(&req); err != nil {
// 		a.setFailedResponse(c, "Failed to parse credentials payload", err)
// 		return
// 	}

// 	// Call the StoreCred function of the first agent in the agents map.
// 	var response *agentpb.StoreCredResponse
// 	var err error
// 	for _, agent := range a.agents {

// 		response, err = agent.GetClient().StoreCred(ctx, &agentpb.StoreCredRequest{
// 			Credname: *req.Credname,
// 			Username: *req.Username,
// 			Password: *req.Password,
// 		})
// 		if err == nil && response.Status == "SUCCESS" {
// 			break
// 		}
// 	}

// 	if response == nil || response.Status != "SUCCESS" {
// 		a.setFailedResponse(c, "failed to store credentials", err)
// 		return
// 	}

// 	c.IndentedJSON(http.StatusOK, &api.Response{
// 		Status:  "SUCCESS",
// 		Message: "stored credentials"})

// 	a.log.Infof("response received from agent", response)
// 	a.log.Debugf("credentials is stored successfully")
// }
