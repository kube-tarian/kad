package app

import (
	"github.com/intelops/go-common/logging"
	workerframework "github.com/kube-tarian/kad/capten/common-pkg/worker-framework"
	"github.com/kube-tarian/kad/capten/deployment-worker/internal/activities"
	"github.com/kube-tarian/kad/capten/deployment-worker/internal/workflows"
)

const (
	WorkflowTaskQueueName = "Deployment"
)

func Start() {
	logger := logging.NewLogger()
	logger.Infof("Started deployment worker\n")

	worker, err := workerframework.NewWorker(WorkflowTaskQueueName, workflows.Workflow, &activities.Activities{}, logger)
	if err != nil {
		logger.Fatalf("Worker initialization failed, Reason: %v\n", err)
	}

	logger.Infof("Running deployment worker..\n")
	if err := worker.Run(); err != nil {
		logger.Fatalf("failed to start the deployment-worker, err: %v", err)
	}

	logger.Infof("Exiting deployment worker\n")
}
