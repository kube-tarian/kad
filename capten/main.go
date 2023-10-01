package main

import (
	"context"
	"fmt"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/temporalclient"
	"github.com/kube-tarian/kad/capten/agent/pkg/workers"
	"github.com/kube-tarian/kad/capten/model"
	"go.temporal.io/sdk/client"
)

var (
	log = logging.NewLogger()
)

func prepareJobResponse(run client.WorkflowRun, name string) *agentpb.JobResponse {
	if run != nil {
		return &agentpb.JobResponse{Id: run.GetID(), RunID: run.GetRunID(), WorkflowName: name}
	}
	return &agentpb.JobResponse{}
}

func Config(ctx context.Context) error {

	tc, err := temporalclient.NewClient(log)
	if err != nil {
		log.Fatal("faield to get temporal CLient: %v", err)
	}
	log.Infof("Recieved Deployer Install event")
	worker := workers.NewConfig(tc, log)

	ci := model.UseCase{Type: "git",
		RepoURL: "https://github.com/indresh-28/test-tekton.git"}
	_, err = worker.SendEvent(ctx,
		&model.ConfigureParameters{Resource: "configure-ci-cd", Action: "tekton"}, ci)
	if err != nil {
		return err
	}

	return nil
}

// func DeployerAppInstall(ctx context.Context, request *agentpb.ApplicationInstallRequest) (*agentpb.JobResponse, error) {

// 	tc, err := temporalclient.NewClient(log)
// 	if err != nil {
// 		log.Fatal("faield to get temporal CLient: %v", err)
// 	}
// 	log.Infof("Recieved Deployer Install event %+v", request)
// 	worker := workers.NewDeployment(tc, log)

// 	if request.ClusterName == "" {
// 		request.ClusterName = "inbuilt"
// 	}
// 	run, err := worker.SendEvent(ctx, "install", request)
// 	if err != nil {
// 		return &agentpb.JobResponse{}, err
// 	}

// 	return prepareJobResponse(run, worker.GetWorkflowName()), err
// }

// func DeploymentInstallActivity(ctx context.Context, req *model.DeployerPostRequest) error {
// 	log.Infof("Activity, name: %+v", req)
// 	// e := activity.GetInfo(ctx)
// 	// logger.Infof("activity info: %+v", e)

// 	plugin, err := plugins.GetPlugin(req.PluginName, log)
// 	if err != nil {
// 		log.Errorf("Get plugin  failed: %v", err)
// 		return err
// 	}

// 	deployerPlugin, ok := plugin.(workerframework.DeploymentWorker)
// 	if !ok {
// 		return fmt.Errorf("plugin not supports deployment activities")
// 	}

// 	emptyVersion := ""
// 	if req.Version == nil {
// 		req.Version = &emptyVersion
// 	}
// 	if req.ValuesYaml == nil {
// 		req.ValuesYaml = &emptyVersion
// 	}
// 	msg, err := deployerPlugin.Create(&model.CreteRequestPayload{
// 		RepoName:    req.RepoName,
// 		RepoURL:     req.RepoUrl,
// 		ChartName:   req.ChartName,
// 		Namespace:   req.Namespace,
// 		ReleaseName: req.ReleaseName,
// 		Timeout:     req.Timeout,
// 		Version:     *req.Version,
// 		ValuesYaml:  *req.ValuesYaml,
// 	})
// 	if err != nil {
// 		fmt.Println("ERR Details: ", err)
// 		return err
// 	}

// 	fmt.Println("MSG:", msg)
// 	return nil
// }

func main() {
	fmt.Println("statrted: ")
	// resp, err := DeployerAppInstall(context.TODO(), &agentpb.ApplicationInstallRequest{
	// 	PluginName:  "helm",
	// 	RepoName:    "tools",
	// 	RepoUrl:     "https://kube-tarian.github.io/helmrepo-supporting-tools",
	// 	Namespace:   "external-secrets",
	// 	ClusterName: "NewCluster",
	// 	ReleaseName: "external-secrets",
	// 	ChartName:   "external-secrets",
	// 	Timeout:     5, // in minutes
	// 	Version:     "1.0.0",
	// })
	// clusterName := "NewCluster"
	// ver := "12.8.4"
	// err := DeploymentInstallActivity(context.TODO(), &model.DeployerPostRequest{
	// 	PluginName:  "helm",
	// 	RepoName:    "bitnami",
	// 	RepoUrl:     "https://charts.bitnami.com/bitnami",
	// 	Namespace:   "default",
	// 	ClusterName: &clusterName,
	// 	ReleaseName: "postgresql",
	// 	ChartName:   "postgresql",
	// 	Timeout:     2, // in minutes
	// 	Version:     &ver,
	// })
	err := Config(context.TODO())
	if err != nil {
		fmt.Println("faield to get temporal CLient: ", err)
		return
	}
	fmt.Println("done: ")
}
