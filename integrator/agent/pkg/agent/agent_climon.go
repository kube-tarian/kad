package agent

import (
	"context"

	"github.com/kube-tarian/kad/integrator/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/integrator/agent/pkg/workers"
)

func (a *Agent) ClimonAppInstall(ctx context.Context, request *agentpb.ClimonInstallRequest) (*agentpb.JobResponse, error) {
	a.log.Infof("Recieved Climon Install event %+v", request)
	worker := workers.NewClimon(a.client, a.log)

	if request.ClusterName == "" {
		request.ClusterName = "inbuilt"
	}
	run, err := worker.SendEvent(ctx, "install", request)
	if err != nil {
		return &agentpb.JobResponse{}, err
	}

	return prepareJobResponse(run, worker.GetWorkflowName()), err
}

func (a *Agent) ClimonAppDelete(ctx context.Context, request *agentpb.ClimonDeleteRequest) (*agentpb.JobResponse, error) {
	a.log.Infof("Recieved Climon delete event %+v", request)
	worker := workers.NewClimon(a.client, a.log)

	if request.ClusterName == "" {
		request.ClusterName = "inbuilt"
	}
	run, err := worker.SendDeleteEvent(ctx, "delete", request)
	if err != nil {
		return &agentpb.JobResponse{}, err
	}

	return prepareJobResponse(run, worker.GetWorkflowName()), err
}
