package helmplugin

import (
	"encoding/json"

	"github.com/kube-tarian/kad/integrator/model"
	helmclient "github.com/mittwald/go-helm-client"
)

func (h *HelmCLient) List(payload model.RequestPayload) (json.RawMessage, error) {
	h.logger.Infof("Helm client List invoke started")

	opt := &helmclient.Options{
		Namespace:        "default", // Change this to the namespace you wish the client to operate in.
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

	// List all deployed releases.
	results, err := helmClient.ListDeployedReleases()
	if err != nil {
		h.logger.Errorf("Fetching deployed applications failed, %v", err)
		return nil, err
	}

	for _, rel := range results {
		h.logger.Infof("Name: %v, Namespace: %v, Revision: %v, Updated at: %v, Status: %v, Chart: %v, AppVersion: %v",
			rel.Name,
			rel.Namespace,
			rel.Version,
			rel.Info.FirstDeployed.UTC(),
			rel.Info.Status,
			rel.Chart.Name(),
			rel.Chart.AppVersion(),
		)
	}
	respMsg, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}
	h.logger.Infof("Helm client List invoke finished")
	return respMsg, nil
}
