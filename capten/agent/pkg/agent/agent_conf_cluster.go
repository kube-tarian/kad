package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/workers"
	"github.com/kube-tarian/kad/capten/model"
)

func (a *Agent) ClusterAdd(ctx context.Context, request *agentpb.ClusterRequest) (*agentpb.JobResponse, error) {
	a.log.Infof("Recieved Deployer Install event %+v", request)
	worker := workers.NewConfig(a.tc, a.log)

	run, err := worker.SendEvent(ctx, &model.ConfigureParameters{Resource: "cluster", Action: "add"}, request)
	if err != nil {
		return &agentpb.JobResponse{}, err
	}

	return prepareJobResponse(run, worker.GetWorkflowName()), err
}

func (a *Agent) ClusterDelete(ctx context.Context, request *agentpb.ClusterRequest) (*agentpb.JobResponse, error) {
	a.log.Infof("Recieved Deployer delete event %+v", request)
	worker := workers.NewConfig(a.tc, a.log)

	run, err := worker.SendEvent(ctx, &model.ConfigureParameters{Resource: "cluster", Action: "delete"}, request)
	if err != nil {
		return &agentpb.JobResponse{}, err
	}

	return prepareJobResponse(run, worker.GetWorkflowName()), err
}
