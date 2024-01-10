package crossplane

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	fileutil "github.com/kube-tarian/kad/capten/config-worker/internal/file_util"
	"github.com/kube-tarian/kad/capten/model"
	agentmodel "github.com/kube-tarian/kad/capten/model"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
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
	logger.Infof("configuring crossplane project for cluster %s update", req.ManagedClusterName)
	endpoint, err := cp.helper.CreateCluster(ctx, req.ManagedClusterId, req.ManagedClusterName)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to CreateCluster in argocd app")
	}

	logger.Infof("cloning default templates %s to project %s", cp.pluginConfig.TemplateGitRepo, req.RepoURL)
	templateRepo, err := cp.helper.CloneTemplateRepo(ctx, cp.pluginConfig.TemplateGitRepo, req.GitProjectId)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to clone repos")
	}
	defer os.RemoveAll(templateRepo)

	customerRepo, err := cp.helper.CloneUserRepo(ctx, req.RepoURL, req.GitProjectId)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to clone repos")
	}
	logger.Infof("cloned default templates to project %s", req.RepoURL)

	defer os.RemoveAll(customerRepo)

	clusterValuesFile := filepath.Join(customerRepo, cp.pluginConfig.ClusterEndpointUpdates.ClusterValuesFile)
	defaultAppListFile := filepath.Join(templateRepo, cp.pluginConfig.ClusterEndpointUpdates.DefaultAppListFile)
	err = updateClusterEndpointDetials(clusterValuesFile, req.ManagedClusterName, endpoint, defaultAppListFile)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to replace the file")
	}

	defaultAppValPath := filepath.Join(templateRepo, cp.pluginConfig.ClusterEndpointUpdates.DefaultAppValuesPath)
	clusterDefaultAppValPath := filepath.Join(customerRepo,
		cp.pluginConfig.ClusterEndpointUpdates.ClusterDefaultAppValuesPath, req.ManagedClusterName)
	err = cp.syncDefaultAppVaules(req.ManagedClusterName, defaultAppValPath, clusterDefaultAppValPath)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to sync default app value files")
	}

	templateValues := cp.prepareTemplateVaules(req.ManagedClusterName)
	if err := fileutil.UpdateFilesInFolderWithTempaltes(clusterDefaultAppValPath, templateValues); err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to update default app template values")
	}

	if err := fileutil.UpdateFileWithTempaltes(clusterValuesFile, templateValues); err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to update cluster config template values")
	}

	logger.Infof("default app vaules synched for cluster %s", req.ManagedClusterName)

	err = cp.helper.AddFilesToRepo([]string{"."})
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to add git repo")
	}

	err = cp.helper.CommitRepoChanges()
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to commit git repo")
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

func updateClusterEndpointDetials(valuesFileName, clusterName, clusterEndpoint, defaultAppFile string) error {
	data, err := os.ReadFile(valuesFileName)
	if err != nil {
		return err
	}

	jsonData, err := k8s.ConvertYamlToJson(data)
	if err != nil {
		return err
	}

	var clusterConfig ClusterConfigValues
	err = json.Unmarshal(jsonData, &clusterConfig)
	if err != nil {
		return err
	}

	defaultApps, err := readClusterDefaultApps(defaultAppFile)
	if err != nil {
		return err
	}

	clusters := []Cluster{}
	if clusterConfig.Clusters != nil {
		clusters = *clusterConfig.Clusters
	}

	var clusterFound bool
	for index := range clusters {
		if clusters[index].Name == clusterName {
			clusters[index] = prepareClusterData(clusterName, clusterEndpoint, defaultApps)
			clusterFound = true
			break
		}
	}

	if !clusterFound {
		clusters = append(clusters, prepareClusterData(clusterName, clusterEndpoint, defaultApps))
	}

	clusterConfig.Clusters = &clusters
	jsonBytes, err := json.Marshal(clusterConfig)
	if err != nil {
		return err
	}

	yamlBytes, err := k8s.ConvertJsonToYaml(jsonBytes)
	if err != nil {
		return err
	}

	err = os.WriteFile(valuesFileName, yamlBytes, os.ModeAppend)
	return err
}

func (cp *CrossPlaneApp) syncDefaultAppVaules(clusterName, defaultAppVaulesPath, clusterDefaultAppVaulesPath string) error {
	if err := fileutil.CreateFolderIfNotExist(clusterDefaultAppVaulesPath); err != nil {
		return err
	}

	if err := fileutil.SyncFiles(defaultAppVaulesPath, clusterDefaultAppVaulesPath); err != nil {
		return err
	}

	return nil
}

func (cp *CrossPlaneApp) configureClusterDelete(ctx context.Context, req *model.CrossplaneClusterUpdate) (status string, err error) {
	logger.Infof("configuring crossplane project for cluster %s delete", req.ManagedClusterName)

	customerRepo, err := cp.helper.CloneUserRepo(ctx, req.RepoURL, req.GitProjectId)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to clone repos")
	}
	logger.Infof("cloned default templates to project %s", req.RepoURL)

	defer os.RemoveAll(customerRepo)

	clusterValuesFile := filepath.Join(customerRepo, cp.pluginConfig.ClusterEndpointUpdates.ClusterValuesFile)
	err = removeClusterValues(clusterValuesFile, req.ManagedClusterName)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to replace the file")
	}

	dirToDeleteInRepoPath := filepath.Join(".", cp.pluginConfig.ClusterEndpointUpdates.ClusterDefaultAppValuesPath, req.ManagedClusterName)
	logger.Infof("removing the cluster '%s' from git repo path %s", req.ManagedClusterName, dirToDeleteInRepoPath)

	err = cp.helper.RemoveFilesFromRepo([]string{dirToDeleteInRepoPath})
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to remove from git repo")
	}

	dirToDelete := filepath.Join(customerRepo, cp.pluginConfig.ClusterEndpointUpdates.ClusterDefaultAppValuesPath, req.ManagedClusterName)
	if err := os.RemoveAll(dirToDelete); err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to remove cluster folder")
	}

	if err := cp.helper.AddFilesToRepo([]string{"."}); err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to add git repo")
	}

	err = cp.helper.CommitRepoChanges()
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to commit git repo")
	}

	logger.Infof("commited project %s changed to git", req.RepoURL)
	return string(agentmodel.WorkFlowStatusCompleted), nil
}

func removeClusterValues(valuesFileName, clusterName string) error {
	logger.Infof("for the culster %s, removing the cluster values from %s file", clusterName, valuesFileName)
	data, err := os.ReadFile(valuesFileName)
	if err != nil {
		return err
	}

	jsonData, err := k8s.ConvertYamlToJson(data)
	if err != nil {
		return err
	}

	var clusterConfig ClusterConfigValues
	err = json.Unmarshal(jsonData, &clusterConfig)
	if err != nil {
		return err
	}

	clusters := []Cluster{}
	if clusterConfig.Clusters != nil {
		clusters = *clusterConfig.Clusters
	}

	newclusters := []Cluster{}
	for _, cluster := range clusters {
		if cluster.Name != clusterName {
			newclusters = append(newclusters, cluster)
		}
	}

	clusterConfig.Clusters = &newclusters
	jsonBytes, err := json.Marshal(clusterConfig)
	if err != nil {
		return err
	}

	yamlBytes, err := k8s.ConvertJsonToYaml(jsonBytes)
	if err != nil {
		return err
	}

	err = os.WriteFile(valuesFileName, yamlBytes, os.ModeAppend)
	return err
}

func (cp *CrossPlaneApp) prepareTemplateVaules(clusterName string) map[string]string {
	val := map[string]string{
		"DomainName":  cp.cfg.DomainName,
		"ClusterName": clusterName,
	}
	return val
}

func prepareClusterData(clusterName, endpoint string, defaultApps []DefaultApps) Cluster {
	return Cluster{
		Name:    clusterName,
		Server:  endpoint,
		DefApps: defaultApps,
	}
}

func readClusterDefaultApps(clusterDefaultAppsFile string) ([]DefaultApps, error) {
	data, err := os.ReadFile(filepath.Clean(clusterDefaultAppsFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read default applist File: %s, err: %w", clusterDefaultAppsFile, err)
	}

	var appList DefaultAppList
	err = yaml.Unmarshal(data, &appList)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall default apps file: %s, err: %w", clusterDefaultAppsFile, err)
	}

	return appList.DefaultApps, nil
}
