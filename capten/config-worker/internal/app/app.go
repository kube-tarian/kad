package app

import (
	"github.com/intelops/go-common/logging"
	workerframework "github.com/kube-tarian/kad/capten/common-pkg/worker-framework"
	"github.com/kube-tarian/kad/capten/config-worker/internal/activities"
	"github.com/kube-tarian/kad/capten/config-worker/internal/workflows"
)

const (
	WorkflowTaskQueueName = "Configure"
)

func Start() {
	logger := logging.NewLogger()
	logger.Infof("Starting config worker..\n")

	worker, err := workerframework.NewWorkerV2(WorkflowTaskQueueName, logger)
	if err != nil {
		logger.Fatalf("Worker initialization failed, Reason: %v", err)
	}

	worker.RegisterWorkflows([]interface{}{workflows.Workflow}...)

	worker.RegisterActivities([]interface{}{&activities.Activities{}}...)

	logger.Infof("Running config worker..\n")
	if err := worker.Run(); err != nil {
		logger.Fatalf("failed to start the config-worker, err: %v", err)
	}

	logger.Infof("Exiting config worker\n")
}
