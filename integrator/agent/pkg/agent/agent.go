package agent

import (
	"context"
	"fmt"
	"log"

	"github.com/kube-tarian/kad/integrator/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/integrator/agent/pkg/temporalclient"
	"github.com/kube-tarian/kad/integrator/agent/pkg/vaultservpb"
	"github.com/kube-tarian/kad/integrator/agent/pkg/workers"
	"github.com/kube-tarian/kad/integrator/climon/pkg/db/cassandra"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/kube-tarian/kad/integrator/model"
	"go.temporal.io/sdk/client"
)

type Agent struct {
	agentpb.UnimplementedAgentServer
	client *temporalclient.Client
	log    logging.Logger
}

func NewAgent(log logging.Logger) (*Agent, error) {
	temporalClient, err := temporalclient.NewClient(log)
	if err != nil {
		log.Errorf("Agent creation failed, %v", err)
		return nil, err
	}

	return &Agent{
		client: temporalClient,
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

func (a *Agent) StoreCred(ctx context.Context, request *agentpb.StoreCredRequest) (*agentpb.StoreCredResponse, error) {
	vaultServ, err := GetVaultServClient()
	if err != nil {
		log.Println("failed to connect vaultserv", err)
		return &agentpb.StoreCredResponse{
			Status: "FAILED",
		}, err
	}

	response, err := vaultServ.StoreCred(ctx, &vaultservpb.StoreCredRequest{
		Username: request.Username,
		Password: request.Password,
		Credname: request.Credname,
	})

	if err != nil {
		log.Println("failed to store creds", err)
		return nil, err
	}

	return &agentpb.StoreCredResponse{
		Status: response.Status,
	}, nil
}

func (a *Agent) GetAppInfo(ctx context.Context, request *agentpb.AppInfoRequest) ([]*agentpb.AppInfoResponse, error) {
	dbConf, err := cassandra.GetDbConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to store data in database: %w", err)
	}
	var logger logging.Logger
	db, err := cassandra.NewCassandraStore(logger, dbConf.DbAddresses, dbConf.DbAdminUsername, dbConf.DbAdminPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to store data in database %w", err)
	}

	appsInfo, err := db.GetAppInfo(ctx, &model.GetAppInfoRequest{Type: request.Type})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all apps from database %w", err)
	}

	var apps []*agentpb.AppInfoResponse
	for _, app := range appsInfo {
		apps = append(apps, &agentpb.AppInfoResponse{
			ReleaseName: app.ReleaseName,
			RepoName:    app.RepoName,
			RepoUrl:     app.RepoUrl,
			ChartName:   app.ChartName,
			Namespace:   app.Namespace,
			Version:     app.Version,
		})
	}
	return apps, nil
}
