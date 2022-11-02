package temporal

import (
	"context"
	"log"

	tc "go.temporal.io/sdk/client"
)

type climon struct {
	client
}

func NewClimon(temporalClient client) *climon {
	return &climon{
		client: temporalClient,
	}
}

func (c *climon) Deploy(deployPayload string) (tc.WorkflowRun, error) {
	options := tc.StartWorkflowOptions{
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

	we, err := c.ExecuteWorkflow(context.Background(), options, DeployWorkflowName, deployPayload)
	if err != nil {
		log.Println("error starting TransferMoney workflow", err)
		return nil, err
	}
	//printResults(deployInfo, we.GetID(), we.GetRunID())

	return we, nil
}
