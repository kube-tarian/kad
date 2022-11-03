package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"

	"go.temporal.io/sdk/client"

	"github.com/kube-tarian/kad/climon/pkg/pb/climonpb"
	"github.com/kube-tarian/kad/climon/pkg/temporal"
)

func main() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	options := client.StartWorkflowOptions{
		ID:        "helm-deploy-workflow",
		TaskQueue: temporal.ClimonDeployTaskQueue,
	}

	deployInfo := climonpb.DeployRequest{
		Version:     "1.0",
		RepoUrl:     "https://charts.bitnami.com/bitnami",
		RepoName:    "bitnami",
		Namespace:   "web",
		ChartName:   "bitnami/wordpress",
		ReleaseName: "intelops",
		ReferenceID: uuid.New().String(),
		Plugin:      "helm",
	}

	deployInfoBytes, err := json.Marshal(deployInfo)
	if err != nil {
		fmt.Println("failed to marshal", err)
		return
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, "DeployWorkflow", string(deployInfoBytes))
	if err != nil {
		log.Fatalln("error starting TransferMoney workflow", err)
	}
	printResults(&deployInfo, we.GetID(), we.GetRunID())
}

func printResults(deployInfo *climonpb.DeployRequest, workflowID, runID string) {
	log.Printf("deploy information %+v\n", deployInfo)
	log.Printf(
		"\nWorkflowID: %s RunID: %s\n",
		workflowID,
		runID,
	)
}
