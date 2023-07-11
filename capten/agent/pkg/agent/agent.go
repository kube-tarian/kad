package agent

import (
	"context"
	"fmt"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/temporalclient"
	"github.com/kube-tarian/kad/capten/agent/pkg/workers"
	"github.com/kube-tarian/kad/capten/common-pkg/db-create/cassandra"
	"go.temporal.io/sdk/client"
)

var _ agentpb.AgentServer = &Agent{}

type Agent struct {
	agentpb.UnimplementedAgentServer

	client *temporalclient.Client
	store  cassandra.Store
	log    logging.Logger
}

type AgentOption func(*Agent) error

func NewAgent(log logging.Logger, opts ...AgentOption) (*Agent, error) {

	agent := &Agent{
		log: log,
	}
	for _, opt := range opts {
		if err := opt(agent); err != nil {
			log.Errorf("Agent creation failed, %v", err)
			return nil, err
		}
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

func (a *Agent) StoreCred(ctx context.Context, request *agentpb.StoreCredentialRequest) (*agentpb.StoreCredentialResponse, error) {
	credPath := fmt.Sprintf("%s/%s/%s", request.CredentialType, request.CredEntityName, request.CredIdentifier)
	err := StoreCredential(ctx, request)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to store credentail for %s", credPath)
		a.log.Errorf("failed to store credentail for %s, %v", credPath, err)
		return &agentpb.StoreCredentialResponse{
			Status:        *agentpb.StatusCode_INTERNRAL_ERROR.Enum(),
			StatusMessage: err.Error(),
		}, nil
	}

	a.log.Audit("security", "storecred", "success", "system", "credentail stored for %s", credPath)
	a.log.Infof("stored credentail for entity %s", credPath)
	return &agentpb.StoreCredentialResponse{
		Status: *agentpb.StatusCode_OK.Enum(),
	}, nil
}

func (a *Agent) SyncApp(ctx context.Context, request *agentpb.SyncAppRequest) (*agentpb.SyncAppResponse, error) {
	err := a.syncApp(ctx, request)
	if err != nil {
		return &agentpb.SyncAppResponse{
			Status:        agentpb.StatusCode(1),
			StatusMessage: "FAILED",
		}, err
	}

	return &agentpb.SyncAppResponse{
		Status:        agentpb.StatusCode(0),
		StatusMessage: "SUCCESS",
	}, nil

}

func WithTemporal(log logging.Logger) AgentOption {

	return func(a *Agent) error {
		temporalClient, err := temporalclient.NewClient(log)
		if err != nil {
			return err
		}
		a.client = temporalClient
		return nil
	}

}

func WithCassandra(log logging.Logger) AgentOption {
	return func(a *Agent) error {

		store := cassandra.NewCassandraStore(log, nil)
		config := &cassandra.DBConfig{}
		err := envconfig.Process("", config)
		if err != nil {
			return err
		}

		if err := store.Connect(
			config.DbAddresses,
			config.DbAdminUsername,
			config.DbAdminPassword); err != nil {
			return err
		}

		a.store = store

		return nil
	}
}
