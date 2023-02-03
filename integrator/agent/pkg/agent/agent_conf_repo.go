package agent

import (
	"context"

	"github.com/kube-tarian/kad/integrator/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/integrator/agent/pkg/workers"
	"github.com/kube-tarian/kad/integrator/model"
)

func (a *Agent) RepositoryAdd(ctx context.Context, request *agentpb.RepositoryAddRequest) (*agentpb.JobResponse, error) {
	a.log.Infof("Recieved Deployer Install event %+v", request)
	worker := workers.NewConfig(a.client, a.log)

	run, err := worker.SendEvent(ctx, &model.ConfigureParameters{Resource: "repo", Action: "add"}, request)
	if err != nil {
		return &agentpb.JobResponse{}, err
	}

	return prepareJobResponse(run, worker.GetWorkflowName()), err
}

func (a *Agent) RepositoryDelete(ctx context.Context, request *agentpb.RepositoryDeleteRequest) (*agentpb.JobResponse, error) {
	a.log.Infof("Recieved Deployer delete event %+v", request)
	worker := workers.NewConfig(a.client, a.log)

	run, err := worker.SendEvent(ctx, &model.ConfigureParameters{Resource: "repo", Action: "delete"}, request)
	if err != nil {
		return &agentpb.JobResponse{}, err
	}

	return prepareJobResponse(run, worker.GetWorkflowName()), err
}
