package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/kube-tarian/kad/agent/pkg/agentpb"

	//"github.com/kube-tarian/kad/agent/pkg/agentpb"
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
	// clnt, err := temporalclient.NewClient(log)
	// if err != nil {
	// 	log.Errorf("Agent creation failed, %v", err)
	// 	return nil, err
	// }

	return &Agent{
		//	client: clnt,
		log: log,
	}, nil
}

func (a *Agent) SubmitJob(ctx context.Context, request *agentpb.JobRequest) (*agentpb.JobResponse, error) {
	a.log.Infof("Recieved event %+v", request)
	worker, err := a.getWorker(request.WorkerType)
	if err != nil {
		return &agentpb.JobResponse{}, err
	}
	var payload string
	switch request.Operation {

	case "Install":
		helm := &agentpb.HelmAppInstallRequest{}
		Bytes, err := json.Marshal(helm)
		if err != nil {
			log.Printf("Error while marshelling helminstall app")
		}
		payload = string(Bytes)
	case "Update":
		helm := &agentpb.HelmAppUpdateRequest{}
		Bytes, err := json.Marshal(helm)
		if err != nil {
			log.Printf("Error while marshelling helminstall app")
		}
		payload = string(Bytes)
	}
	fmt.Println(payload)
	// helm := &agentpb.HelmAppInstallRequest{}
	// //	bytes, err := json.Marshal(HelmAppInstallRequest)
	// bytes, err := json.Marshal(helm)

	// bytes, err := json.Marshal(HelmAppInstallRequest)
	// // jobRequest := JobRequest {
	// // 	//operation: "Install",
	// // 	//WorkerType: "Integrator",
	// // 	payload: string(bytes),
	// // }
	// jobreq:=request.Payload
	//j //obreq := &agentpb.JobRequest{}
	// jobreq:=request.Payload

	//run, err := worker.SendEvent(ctx, request.Payload.GetValue())
	//	run, err := worker.SendEvent(ctx, json.RawMessage(jobreq.GetPayload()))
	//	run, err := worker.SendEvent(ctx, json.RawMessage(jobreq.GetPayload()))
	run, err := worker.SendEvent(ctx, json.RawMessage(request.GetPayload()))
	if err != nil {
		return &agentpb.JobResponse{}, err
	}

	// run, err := worker.SendEvent(ctx, request.Payload.GetValue())
	// if err != nil {
	// 	return &agentpb.JobResponse{}, err
	// }

	return &agentpb.JobResponse{Id: run.GetID(), RunID: run.GetRunID(), WorkflowName: worker.GetWorkflowName()}, err
}

func (a *Agent) getWorker(WorkerType string) (workers.Worker, error) {
	switch WorkerType {
	case "climon":
		return workers.NewClimon(a.client), nil
	case "deployment":
		return workers.NewDeployment(a.client, a.log), nil
	default:
		return nil, fmt.Errorf("unsupported operation %s", WorkerType)
	}
}
