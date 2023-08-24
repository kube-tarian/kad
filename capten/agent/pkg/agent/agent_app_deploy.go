package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/workers"
)

func (a *Agent) ClimonAppInstall(ctx context.Context, request *agentpb.ClimonInstallRequest) (*agentpb.JobResponse, error) {
	a.log.Infof("Recieved Climon Install event %+v", request)
	worker := workers.NewClimon(a.tc, a.log)

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
	worker := workers.NewClimon(a.tc, a.log)

	if request.ClusterName == "" {
		request.ClusterName = "inbuilt"
	}
	run, err := worker.SendDeleteEvent(ctx, "delete", request)
	if err != nil {
		return &agentpb.JobResponse{}, err
	}

	return prepareJobResponse(run, worker.GetWorkflowName()), err
}

func (a *Agent) InstallApp(ctx context.Context, request *agentpb.InstallAppRequest) (*agentpb.InstallAppResponse, error) {
	a.log.Infof("Recieved App Install request %+v", request)
	worker := workers.NewDeployment(a.tc, a.log)

	config := &agentpb.ApplicationInstallRequest{
		PluginName:  "helm",
		RepoName:    request.AppConfig.RepoName,
		RepoUrl:     request.AppConfig.RepoURL,
		ChartName:   request.AppConfig.ChartName,
		Namespace:   request.AppConfig.Namespace,
		ReleaseName: request.AppConfig.ReleaseName,
		Version:     request.AppConfig.Version,
		ClusterName: "capten",
		ValuesYaml:  string(request.OverrideValues),
	}

	run, err := worker.SendEvent(ctx, "install", config)
	if err != nil {
		return &agentpb.InstallAppResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: "Internall error in create install app job",
		}, err
	}

	return &agentpb.InstallAppResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "success",
		JobResponse:   &agentpb.JobResponse{Id: run.GetID(), RunID: run.GetRunID(), WorkflowName: worker.GetWorkflowName()},
	}, nil
}
