package helmplugin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kube-tarian/kad/integrator/model"
	helmclient "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/repo"
)

func (h *HelmCLient) Create(payload model.RequestPayload) (json.RawMessage, error) {
	h.logger.Infof("Helm client Install invoke started")

	req := &model.Request{}
	err := json.Unmarshal(payload.Data, req)
	if err != nil {
		h.logger.Errorf("payload unmarshal failed, %v", err)
		return nil, err
	}

	opt := &helmclient.Options{
		Namespace:        req.Namespace, // Change this to the namespace you wish the client to operate in.
		RepositoryCache:  "/tmp/.helmcache",
		RepositoryConfig: "/tmp/.helmrepo",
		Debug:            true,
		Linting:          true,
		DebugLog:         h.logger.Debugf,
	}

	helmClient, err := helmclient.New(opt)
	if err != nil {
		h.logger.Errorf("helm client initialization failed, %v", err)
		return nil, err
	}

	err = h.addOrUpdate(helmClient, req)
	if err != nil {
		h.logger.Errorf("helm repo add failed, %v", err)
		return nil, err
	}

	// Use an unpacked chart directory.
	chartSpec := helmclient.ChartSpec{
		ReleaseName: req.ReleaseName,
		ChartName:   fmt.Sprintf("%s/%s", req.RepoName, req.ChartName),
		Namespace:   req.Namespace,
		Wait:        true,
		Timeout:     time.Duration(req.Timeout) * time.Minute,
	}

	// Use the default rollback strategy offer by HelmClient (revert to the previous version).
	rel, err := helmClient.InstallOrUpgradeChart(
		context.Background(),
		&chartSpec,
		&helmclient.GenericHelmOptions{
			RollBack: helmClient,
		})
	if err != nil {
		h.logger.Errorf("helm install or update for request %+v failed, %v", req, err)
		return nil, err
	}

	h.logger.Infof("helm install of app %s successful in namespace: %v, status: %v", rel.Name, rel.Info.Status, rel.Namespace)
	h.logger.Infof("Helm client Install invoke finished")
	return json.RawMessage(fmt.Sprintf("{\"status\": \"Application %s install successful\"}", rel.Name)), nil
}

func (h *HelmCLient) addOrUpdate(client helmclient.Client, req *model.Request) error {
	// Define a public chart repository.
	chartRepo := repo.Entry{
		Name: req.RepoName,
		URL:  req.RepoURL,
	}

	// Add a chart-repository to the client.
	if err := client.AddOrUpdateChartRepo(chartRepo); err != nil {
		h.logger.Errorf("helm repo add failed, %v", err)
		return err
	}
	return nil
}
