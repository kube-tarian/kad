package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/workers"
)

func (a *Agent) DeployerAppInstall(ctx context.Context, request *agentpb.ApplicationInstallRequest) (*agentpb.JobResponse, error) {
	a.log.Infof("Recieved Deployer Install event %+v", request)
	worker := workers.NewDeployment(a.client, a.log)

	if request.ClusterName == "" {
		request.ClusterName = "inbuilt"
	}
	run, err := worker.SendEvent(ctx, "install", request)
	if err != nil {
		return &agentpb.JobResponse{}, err
	}

	return prepareJobResponse(run, worker.GetWorkflowName()), err
}

func (a *Agent) DeployerAppDelete(ctx context.Context, request *agentpb.ApplicationDeleteRequest) (*agentpb.JobResponse, error) {
	a.log.Infof("Recieved Deployer delete event %+v", request)
	worker := workers.NewDeployment(a.client, a.log)

	if request.ClusterName == "" {
		request.ClusterName = "inbuilt"
	}
	run, err := worker.SendDeleteEvent(ctx, "delete", request)
	if err != nil {
		return &agentpb.JobResponse{}, err
	}

	return prepareJobResponse(run, worker.GetWorkflowName()), err
}
