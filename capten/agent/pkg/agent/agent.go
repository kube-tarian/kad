package agent

import (
	"context"
	"fmt"

	"github.com/intelops/go-common/logging"
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

func NewAgent(log logging.Logger, temporalClient *temporalclient.Client, store cassandra.Store) (*Agent, error) {
	return &Agent{
		client: temporalClient,
		store:  store,
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
			Status: "FAILED",
		}, err
	}

	return &agentpb.SyncAppResponse{
		Status: "SUCCESS",
	}, nil

}
