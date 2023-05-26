package workers

import (
	"context"
	"encoding/json"

	"go.temporal.io/sdk/client"
)

const (
	ClimonHelmTaskQueue = "CLIMON_HELM_TASK_QUEUE"
	DeployWorkflowName  = "Workflow"
)

type Worker interface {
	SendEvent(ctx context.Context, payload json.RawMessage) (client.WorkflowRun, error)
	GetWorkflowName() string
}
