package server

import (
	"context"
	"fmt"

	"github.com/kube-tarian/kad/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/agent/pkg/logging"
	"github.com/kube-tarian/kad/agent/pkg/temporalclient"
	"github.com/kube-tarian/kad/agent/pkg/workers"
)

type Agent struct {
	agentpb.UnimplementedAgentServer
	client *temporalclient.Client
	log    logging.Logger
}

func NewAgent(log logging.Logger) (*Agent, error) {
	clnt, err := temporalclient.NewClient(log)
	if err != nil {
		log.Errorf("Agent creation failed, %v", err)
		return nil, err
	}

	return &Agent{
		client: clnt,
		log:    log,
	}, nil
}

func (a *Agent) SubmitJob(ctx context.Context, request *agentpb.JobRequest) (*agentpb.JobResponse, error) {
	a.log.Infof("Recieved event %+v", request)
	worker, err := a.getWorker(request.Operation)
	if err != nil {
		return &agentpb.JobResponse{}, err
	}

	run, err := worker.SendEvent(ctx, request.Payload.GetValue())
	if err != nil {
		return &agentpb.JobResponse{}, err
	}

	return &agentpb.JobResponse{Id: run.GetID(), RunID: run.GetRunID(), WorkflowName: worker.GetWorkflowName()}, err
}

func (a *Agent) getWorker(operatoin string) (workers.Worker, error) {
	switch operatoin {
	case "climon":
		return workers.NewClimon(a.client), nil
	case "deployment":
		return workers.NewDeployment(a.client, a.log), nil
	default:
		return nil, fmt.Errorf("unsupported operation %s", operatoin)
	}
}
