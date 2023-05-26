package workers

import (
	"context"
	"github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/agent/pkg/temporalclient"
	"go.temporal.io/sdk/client"
	"log"
)

const (
	SyncTaskQueue    = "SYNC_TASK_QUEUE"
	SyncWorkFlowName = "SyncWorkflow"
)

type Sync struct {
	client *temporalclient.Client
}

type App struct {
	Name            string                 `json:"name"`
	ChartName       string                 `json:"chartName"`
	RepoName        string                 `json:"repoName"`
	RepoURL         string                 `json:"repoURL"`
	Namespace       string                 `json:"namespace"`
	ReleaseName     string                 `json:"releaseName"`
	Version         string                 `json:"version"`
	CreateNamespace bool                   `json:"createNamespace"`
	Override        map[string]interface{} `json:"override"`
}

type SyncDataRequest struct {
	Type string `json:"type"`
	Apps []App  `json:"apps"`
}

func NewSync(client *temporalclient.Client) *Sync {
	return &Sync{
		client: client,
	}
}

func (s *Sync) GetWorkflowName() string {
	return DeployWorkflowName
}

func (s *Sync) SendEvent(ctx context.Context, request SyncDataRequest) (client.WorkflowRun, error) {
	options := client.StartWorkflowOptions{
		ID:        "sync-workflow",
		TaskQueue: SyncTaskQueue,
	}

	/*
		var appDatas []AppData
		if err := json.Unmarshal([]byte(syncRequest.Data), &appDatas); err != nil {
			return nil, err
		}

		syncDataRequest := SyncDataRequest{
			Type: syncRequest.Type,
			AppData: appDatas,
		}*/

	we, err := s.client.TemporalClient.ExecuteWorkflow(context.Background(), options, SyncWorkFlowName, request)
	if err != nil {
		log.Println("error starting climon workflow", err)
		return nil, err
	}
	//printResults(deployInfo, we.GetID(), we.GetRunID())

	log.Printf("Started workflow, ID: %v, WorkflowName: %v RunID: %v", we.GetID(), SyncWorkFlowName, we.GetRunID())

	// Wait for 5mins till workflow finishes
	// Timeout with 5mins
	var result model.ResponsePayload
	err = we.Get(ctx, &result)
	if err != nil {
		log.Printf("Result for workflow ID: %v, workflowName: %v, runID: %v", we.GetID(), DeploymentWorkerWorkflowName, we.GetRunID())
		log.Printf("Workflow result failed, %v", err)
		return we, err
	}

	log.Printf("workflow finished success, %+v", result.ToString())
	return we, nil
}

/*
func (s *Sync) getWorkflowStatusByLatestWorkflow(run client.WorkflowRun) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			err := c.getWorkflowInformation(run)
			if err != nil {
				c.log.Errorf("get state of workflow failed: %v, retrying .....", err)
				continue
			}
			return nil
		case <-time.After(5 * time.Minute):
			c.log.Errorf("Timed out waiting for state of workflow")
			return fmt.Errorf("timedout waiting for the workflow to finish")
		}
	}
}

func (c *Climon) getWorkflowInformation(run client.WorkflowRun) error {
	latestRun := c.client.TemporalClient.GetWorkflow(context.Background(), run.GetID(), "")

	var result model.ResponsePayload
	if err := latestRun.Get(context.Background(), &result); err != nil {
		c.log.Errorf("Unable to decode query result", err)
		return err
	}
	c.log.Debugf("Result info: %+v", result)
	return nil
}
*/
