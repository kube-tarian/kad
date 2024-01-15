package crossplane

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	fileutil "github.com/kube-tarian/kad/capten/config-worker/internal/file_util"
	"github.com/kube-tarian/kad/capten/model"
	agentmodel "github.com/kube-tarian/kad/capten/model"
	"github.com/pkg/errors"
)

func (cp *CrossPlaneApp) configureConfigProviderUpdate(ctx context.Context, req *model.CrossplaneClusterUpdate) (status string, err error) {
	logger.Infof("configuring config provider %s update", req.ManagedClusterName)

	x, _ := json.Marshal(req)

	fmt.Println("configureConfigProviderUpdate request")

	fmt.Println(string(x))

	return "", nil

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

	err = cp.configureExternalSecretsOnCluster(ctx, req.ManagedClusterId, req.ManagedClusterName)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to create cluster secrets")
	}
	return string(agentmodel.WorkFlowStatusCompleted), nil
}
