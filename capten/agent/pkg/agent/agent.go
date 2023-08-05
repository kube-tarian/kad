package agent

import (
	"context"
	"fmt"
	"os"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	captenstore "github.com/kube-tarian/kad/capten/agent/pkg/capten-store"
	"github.com/kube-tarian/kad/capten/agent/pkg/temporalclient"
	"github.com/kube-tarian/kad/capten/agent/pkg/workers"

	"go.temporal.io/sdk/client"
)

var _ agentpb.AgentServer = &Agent{}

type Agent struct {
	agentpb.UnimplementedAgentServer
	tc  *temporalclient.Client
	as  *captenstore.Store
	log logging.Logger
}

func NewAgent(log logging.Logger) (*Agent, error) {
	var tc *temporalclient.Client
	var err error

	if os.Getenv("ENV") != "LOCAL" {
		tc, err = temporalclient.NewClient(log)
		if err != nil {
			return nil, err
		}
	}
	// Note how lack of dependecy injection leads to codesmell

	as, err := captenstore.NewStore(log)
	if err != nil {
		return nil, err
	}

	agent := &Agent{
		tc:  tc,
		as:  as,
		log: log,
	}
	return agent, nil
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

	return prepareJobResponse(run, worker.GetWorkflowName()), err
}

func (a *Agent) getWorker(operatoin string) (workers.Worker, error) {
	switch operatoin {
	default:
		return nil, fmt.Errorf("unsupported operation %s", operatoin)
	}
}

func prepareJobResponse(run client.WorkflowRun, name string) *agentpb.JobResponse {
	if run != nil {
		return &agentpb.JobResponse{Id: run.GetID(), RunID: run.GetRunID(), WorkflowName: name}
	}
	return &agentpb.JobResponse{}
}
