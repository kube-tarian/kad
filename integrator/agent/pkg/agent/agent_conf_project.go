package agent

import (
	"context"

	"github.com/kube-tarian/kad/integrator/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/integrator/agent/pkg/workers"
	"github.com/kube-tarian/kad/integrator/model"
)

func (a *Agent) ProjectAdd(ctx context.Context, request *agentpb.ProjectAddRequest) (*agentpb.JobResponse, error) {
	a.log.Infof("Recieved Deployer Install event %+v", request)
	worker := workers.NewConfig(a.client, a.log)

	run, err := worker.SendEvent(ctx, &model.ConfigureParameters{Resource: "project", Action: "add"}, request)
	if err != nil {
		return &agentpb.JobResponse{}, err
	}

	return prepareJobResponse(run, worker.GetWorkflowName()), err
}

func (a *Agent) ProjectDelete(ctx context.Context, request *agentpb.ProjectDeleteRequest) (*agentpb.JobResponse, error) {
	a.log.Infof("Recieved Deployer delete event %+v", request)
	worker := workers.NewConfig(a.client, a.log)

	run, err := worker.SendEvent(ctx, &model.ConfigureParameters{Resource: "project", Action: "delete"}, request)
	if err != nil {
		return &agentpb.JobResponse{}, err
	}

	return prepareJobResponse(run, worker.GetWorkflowName()), err
}
