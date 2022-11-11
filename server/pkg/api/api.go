package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kube-tarian/kad/server/pkg/client"
	"github.com/kube-tarian/kad/server/pkg/model"
	"github.com/kube-tarian/kad/server/pkg/pb/agentpb"
	"github.com/kube-tarian/kad/server/pkg/pb/climonpb"
	"github.com/kube-tarian/kad/server/pkg/server"
	"google.golang.org/protobuf/types/known/anypb"
)

type DeployPayload struct {
	ChartName string `json:"chart_name"`
}

type DeployResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Setup(ginEngine *gin.Engine, s *server.Server) {
	ginEngine.POST("/deploy", s.Deploy)
}

func deploy(c *gin.Context) {
	//TODO get address from database based on CustomerInfo
	log.Println("deploy api invocation started")

	agent, err := client.NewAgent()
	if err != nil {
		log.Println("failed to connect agent internal error", err)
		c.IndentedJSON(http.StatusInternalServerError, &model.DeployResponse{
			Status:  "FAILED",
			Message: "Failed to connect to agent"})
	}
	defer agent.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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

	reqJSON := climonRequest.String()
	jobRequest := &agentpb.JobRequest{
		Operation: "DEPLOY",
		Payload:   &anypb.Any{Value: []byte(reqJSON)},
	}

	response, err := agent.SubmitJob(ctx, jobRequest)
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
	log.Println("deploy api invocation finished")
}
