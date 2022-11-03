package helmplugin

import (
	helmclient "github.com/mittwald/go-helm-client"
)

func (h *HelmCLient) List() error {
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
		return err
	}

	// List all deployed releases.
	results, err := helmClient.ListDeployedReleases()
	if err != nil {
		h.logger.Errorf("Fetching deployed applications failed, %v", err)
		return err
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
	h.logger.Infof("Helm client List invoke finished")
	return nil
}
