package helmplugin

import (
	"encoding/json"
	"fmt"

	"github.com/kube-tarian/kad/integrator/deployment-worker/pkg/model"
	helmclient "github.com/mittwald/go-helm-client"
)

func (h *HelmCLient) Delete(payload model.RequestPayload) error {
	h.logger.Infof("Helm client Install invoke started")

	req := &model.Request{}
	err := json.Unmarshal(payload.Data, req)
	if err != nil {
		h.logger.Errorf("payload unmarshal failed, %v", err)
		return err
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
		return err
	}

	err = h.addOrUpdate(helmClient, req)
	if err != nil {
		h.logger.Errorf("helm repo add failed, %v", err)
		return err
	}

	// Define the released chart to be uninstalled.
	chartSpec := helmclient.ChartSpec{
		ReleaseName: req.ReleaseName,
		ChartName:   fmt.Sprintf("%s/%s", req.RepoName, req.ChartName),
		Namespace:   req.Namespace,
		Wait:        true,
	}

	// Uninstall the chart release.
	// Note that helmclient.Options.Namespace should ideally match the namespace in chartSpec.Namespace.
	err = helmClient.UninstallRelease(&chartSpec)
	if err != nil {
		h.logger.Errorf("helm uninitialization for request %+v failed, %v", req, err)
		return err
	}

	h.logger.Infof("helm uninstall of app %s successful in namespace: %v", req.ReleaseName, req.Namespace)
	h.logger.Infof("Helm client Install invoke finished")
	return nil
}
