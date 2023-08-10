package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/model"
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

func (a *Agent) InstallApp(ctx context.Context, request *agentpb.InstallAppRequest) (*agentpb.JobResponse, error) {
	a.log.Infof("Recieved App Install event %+v", request)
	worker := workers.NewClimon(a.tc, a.log)

	config := &model.AppConfig{
		ReleaseName:         request.AppConfig.ReleaseName,
		AppName:             request.AppConfig.AppName,
		Version:             request.AppConfig.Version,
		Category:            request.AppConfig.Category,
		Description:         request.AppConfig.Description,
		ChartName:           request.AppConfig.ChartName,
		RepoName:            request.AppConfig.RepoName,
		RepoURL:             request.AppConfig.RepoURL,
		Namespace:           request.AppConfig.Namespace,
		CreateNamespace:     request.AppConfig.CreateNamespace,
		PrivilegedNamespace: request.AppConfig.PrivilegedNamespace,
		Icon:                string(request.AppConfig.Icon),
		LaunchURL:           request.AppConfig.LaunchURL,
		LaunchRedirectURL:   request.AppConfig.LaunchRedirectURL,
		ClusterName:         "inbuilt",
	}
	run, err := worker.SendInstallAppEvent(ctx, "install", config)
	if err != nil {
		return &agentpb.JobResponse{}, err
	}

	return prepareJobResponse(run, worker.GetWorkflowName()), err
}
