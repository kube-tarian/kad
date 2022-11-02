package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"intelops.io/server/pkg/client"
	"intelops.io/server/pkg/pb/agentpb"
	"intelops.io/server/pkg/pb/climonpb"
)

type DeployPayload struct {
	ChartName string `json:"chart_name"`
}

type DeployResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Setup(ginEngine *gin.Engine) {
	ginEngine.POST("/deploy", deploy)
}

func deploy(c *gin.Context) {
	//TODO get address from database based on CustomerInfo
	agent, err := client.NewAgent("localhost:50012")
	if err != nil {
		log.Println("failed to connect agent internal error", err)
		c.IndentedJSON(http.StatusInternalServerError, &DeployResponse{
			Status:  "FAILED",
			Message: "Failed to connect to agent"})
	}

	defer agent.Close()
	agentClient := agent.GetClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var deployPayload DeployPayload
	if err := c.BindJSON(&deployPayload); err != nil {
		log.Println("failed to parse payload", err)
		return
	}

	climonRequest := climonpb.DeployRequest{
		Version:     "",
		RepoUrl:     "https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml",
		RepoName:    "",
		Namespace:   "",
		ChartName:   deployPayload.ChartName,
		ReleaseName: "",
		ReferenceID: "",
	}

	jobRequest := &agentpb.JobRequest{
		Operation: "DEPLOY",
		Payload:   climonRequest.String(),
	}

	response, err := agentClient.SubmitJob(ctx, jobRequest)
	if err != nil {
		log.Println("failed to submit job", err)
		c.IndentedJSON(http.StatusInternalServerError, &DeployResponse{
			Status:  "FAILED",
			Message: "Failed to connect to agent"})
		return
	}

	c.IndentedJSON(http.StatusOK, &DeployResponse{
		Status:  "SUCCESS",
		Message: "submitted Job"})

	log.Println("response received", response)
	log.Println("called deploy")
}
