package workers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/kube-tarian/kad/integrator/agent/pkg/temporalclient"
	"go.temporal.io/sdk/client"
)

type climon struct {
	client *temporalclient.Client
}

func NewClimon(client *temporalclient.Client) *climon {
	return &climon{
		client: client,
	}
}

func (c *climon) GetWorkflowName() string {
	return DeployWorkflowName
}

func (c *climon) SendEvent(ctx context.Context, deployPayload json.RawMessage) (client.WorkflowRun, error) {
	options := client.StartWorkflowOptions{
		ID:        "helm-deploy-workflow",
		TaskQueue: ClimonHelmTaskQueue,
	}

	/*
		deployInfo := helm.DeployInfo{
			Version:     "1.0",
			RepoUrl:     "https://charts.bitnami.com/bitnami",
			RepoName:    "bitnami",
			Namespace:   "web",
			ChartName:   "bitnami/wordpress",
			ReleaseName: "intelops",
			ReferenceID: uuid.New().String(),
		}
	*/

	we, err := c.client.ExecuteWorkflow(context.Background(), options, DeployWorkflowName, deployPayload)
	if err != nil {
		log.Println("error starting TransferMoney workflow", err)
		return nil, err
	}
	//printResults(deployInfo, we.GetID(), we.GetRunID())

	return we, nil
}
