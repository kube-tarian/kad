package crossplane

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kube-tarian/kad/capten/model"
	agentmodel "github.com/kube-tarian/kad/capten/model"
	"github.com/pkg/errors"
)

func (cp *CrossPlaneApp) configureConfigProviderUpdate(ctx context.Context, req *model.CrossplaneProviderUpdate) (status string, err error) {
	logger.Infof("configuring config provider %s update", req.ProviderName)

	customerRepo, err := cp.helper.CloneUserRepo(ctx, req.RepoURL, req.GitProjectId)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to clone repo")
	}
	logger.Infof("cloned customer to project %s", req.RepoURL)

	defer os.RemoveAll(customerRepo)

	cloudType := strings.ToLower(req.CloudType)
	syncPath := fmt.Sprintf("%s/%s-packages/%s-k8s-package.yaml", cp.pluginConfig.ProviderEndpointUpdates.SyncAppPath, cloudType, cloudType)

	if _, err := os.Stat(fmt.Sprintf("%s/%s-packages", cp.pluginConfig.ProviderEndpointUpdates.SyncAppPath, cloudType)); err != nil && os.IsNotExist(err) {
		logger.Errorf("provider package directory is not available, path - %s, provider - %s", syncPath, cloudType)
		return string(agentmodel.WorkFlowStatusCompleted), nil
	} else if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to get provider config directory")
	}

	ns, resName, err := getAppNameNamespace(ctx, filepath.Join(customerRepo, syncPath))
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to get name and namespace from")
	}

	err = cp.helper.SyncArgoCDApp(ctx, ns, resName)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to sync argocd app")
	}
	logger.Infof("synched provider config main-app %s, namespace %s", resName, ns)

	err = cp.helper.WaitForArgoCDToSync(ctx, ns, resName)
	if err != nil {
		return string(agentmodel.WorkFlowStatusFailed), errors.WithMessage(err, "failed to fetch argocd app")
	}

	return string(agentmodel.WorkFlowStatusCompleted), nil
}
