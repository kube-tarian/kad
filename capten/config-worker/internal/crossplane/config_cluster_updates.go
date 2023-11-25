package crossplane

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"github.com/kube-tarian/kad/capten/model"
	agentmodel "github.com/kube-tarian/kad/capten/model"
	"github.com/pkg/errors"
)

func getAppNameNamespace(ctx context.Context, fileName string) (string, string, error) {
	k8sclient, err := k8s.NewK8SClient(logging.NewLogger())
	if err != nil {
		return "", "", fmt.Errorf("failed to initalize k8s client: %v", err)
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", "", err
	}

	jsonData, err := k8s.ConvertYamlToJson(data)
	if err != nil {
		return "", "", err
	}

	// For the testing change the reqrepo to template one
	ns, resName, err := k8sclient.DynamicClient.GetNameNamespace(jsonData)
	if err != nil {
		return "", "", fmt.Errorf("failed to create the k8s custom resource: %v", err)
	}

	return ns, resName, nil

}

func (cp *CrossPlaneApp) configureClusterUpdate(ctx context.Context, req *model.CrossplaneClusterUpdate) (status string, err error) {
	logger.Infof("configuring the cluster endpoint for %s", req.RepoURL)
	endpoint, err := cp.helper.CreateCluster(ctx, req.ManagedClusterId, req.Name)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to CreateCluster in argocd app")
	}

	logger.Infof("CreateCluster argocd err: ", err)
	accessToken, err := cp.helper.GetAccessToken(ctx, req.GitProjectId)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to get token from vault")
	}

	logger.Infof("cloning default templates %s to project %s", cp.pluginConfig.TemplateGitRepo, req.RepoURL)
	templateRepo, customerRepo, err := cp.helper.CloneRepos(ctx, cp.pluginConfig.TemplateGitRepo, req.RepoURL, accessToken)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to clone repos")
	}
	logger.Infof("cloned default templates to project %s", req.RepoURL)

	defer os.RemoveAll(templateRepo)
	defer os.RemoveAll(customerRepo)

	fileName := filepath.Join(customerRepo, cp.pluginConfig.ClusterEndpointUpdates.File)
	// replace cluster endpoint
	err = updateClusterEndpointDetials(fileName, req.Name, endpoint, cp.cfg.ClusterDefaultAppsFile)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to replace the file")
	}

	err = cp.helper.AddToGit(ctx, model.CrossPlaneClusterUpdate, req.RepoURL, accessToken)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to add git repo")
	}

	logger.Infof("added cloned project %s changed to git", req.RepoURL)
	ns, resName, err := getAppNameNamespace(ctx, filepath.Join(customerRepo, cp.pluginConfig.ClusterEndpointUpdates.MainAppGitPath))
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to get name and namespace from")
	}

	err = cp.helper.SyncArgoCDApp(ctx, ns, resName)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to sync argocd app")
	}
	logger.Infof("synched provider config main-app %s", resName)

	err = cp.helper.WaitForArgoCDToSync(ctx, ns, resName)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to fetch argocd app")
	}

	return string(agentmodel.WorkFlowStatusCompleted), nil
}

func updateClusterEndpointDetials(filename, clusterName, clusterEndpoint, defaultAppFile string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	jsonData, err := k8s.ConvertYamlToJson(data)
	if err != nil {
		return err
	}

	var argoCDAppValue ArgoCDAppValue

	err = json.Unmarshal(jsonData, &argoCDAppValue)
	if err != nil {
		return err
	}

	clusters := *argoCDAppValue.Clusters
	for index := range clusters {
		cluster := &clusters[index]
		if cluster.Name == clusterName {
			defaultApps, err := readClusterDefaultApps(defaultAppFile)
			if err != nil {
				return err
			}

			for index := range defaultApps {
				localObj := &defaultApps[index]
				strings.ReplaceAll(localObj.ValuesPath, clusterNameSub, clusterName)
			}

			logger.Infof("udpated the req endpoint details to %s for name %s ", clusterEndpoint, clusterName)
			cluster.Server = clusterEndpoint

			break
		}
	}

	argoCDAppValue.Clusters = &clusters

	jsonBytes, err := json.Marshal(argoCDAppValue)
	if err != nil {
		return err
	}

	yamlBytes, err := k8s.ConvertJsonToYaml(jsonBytes)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, yamlBytes, os.ModeAppend)

	return err
}
