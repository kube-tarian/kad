package helm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	helmclient "github.com/kube-tarian/kad/capten/common-pkg/plugins/helm/go-helm-client"
	"github.com/kube-tarian/kad/capten/model"
	"helm.sh/helm/v3/pkg/repo"
)

func (h *HelmCLient) Create(req *model.CreteRequestPayload) (json.RawMessage, error) {
	h.logger.Infof("Helm client Install invoke started")

	helmClient, err := h.getHelmClient(req.Namespace)
	if err != nil {
		h.logger.Errorf("helm client initialization failed, %v", err)
		return nil, err
	}

	err = h.addOrUpdate(helmClient, req)
	if err != nil {
		h.logger.Errorf("helm repo add failed, %v", err)
		return nil, err
	}

	chartSpec := helmclient.ChartSpec{
		ReleaseName:     req.ReleaseName,
		ChartName:       req.ChartName,
		Namespace:       req.Namespace,
		Version:         req.Version,
		Wait:            true,
		Timeout:         time.Duration(req.Timeout) * time.Minute,
		CreateNamespace: true,
	} // Use an unpacked chart directory.

	if req.ValuesYaml != "" {
		chartSpec.ValuesYaml = req.ValuesYaml
	}

	// Use the default rollback strategy offer by HelmClient (revert to the previous version).
	rel, err := helmClient.InstallOrUpgradeChart(
		context.Background(),
		&chartSpec,
		&helmclient.GenericHelmOptions{
			RollBack:              helmClient,
			InsecureSkipTLSverify: true,
		})
	if err != nil {
		h.logger.Errorf("helm install or update for request %+v failed, %v", req, err)
		return nil, err
	}

	h.logger.Infof("helm install of app %s successful in namespace: %v, status: %v", rel.Name, rel.Info.Status, rel.Namespace)
	h.logger.Infof("Helm client Install invoke finished")
	return json.RawMessage(fmt.Sprintf("{\"status\": \"Application %s install successful\"}", rel.Name)), nil
}

func (h *HelmCLient) getHelmClient(namespace string) (helmclient.Client, error) {
	opt := &helmclient.Options{
		Namespace:        namespace,
		RepositoryCache:  "/tmp/.helmcache",
		RepositoryConfig: "/tmp/.helmrepo",
		Debug:            true,
		Linting:          true,
		DebugLog:         h.logger.Debugf,
	}

	return helmclient.New(opt)
}

func (h *HelmCLient) addOrUpdate(client helmclient.Client, req *model.CreteRequestPayload) error {
	// Define a public chart repository.
	chartRepo := repo.Entry{
		Name:                  req.RepoName,
		URL:                   req.RepoURL,
		InsecureSkipTLSverify: true,
	}

	// Add a chart-repository to the client.
	if err := client.AddOrUpdateChartRepo(chartRepo); err != nil {
		h.logger.Errorf("helm repo add failed, %v", err)
		return err
	}
	return nil
}
