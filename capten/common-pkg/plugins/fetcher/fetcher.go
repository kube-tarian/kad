package fetcher

// This package is to fetch the credentials for plugins used in capten
// Fetches the plugin details from Cassandra
// It uses kubernetes clientset to fetch plugin credential details from underlying kubernetes using details fetched from cassandra

import (
	"fmt"

	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"github.com/kube-tarian/kad/capten/common-pkg/logging"
)

type CredentialFetcher struct {
	k8sClient *k8s.K8SClient
	log       logging.Logger
}

func NewCredentialFetcher(log logging.Logger) (*CredentialFetcher, error) {
	// Initialize kubernetes client
	k8sClient, err := k8s.NewK8SClient(log)
	if err != nil {
		log.Errorf("K8S client initialization failed: %v", err)
		return nil, fmt.Errorf("k8 client initialization failed, %v", err)
	}

	// TODO: Initialze Cassandra client

	return &CredentialFetcher{
		k8sClient: k8sClient,
		log:       log,
	}, nil
}

func (c *CredentialFetcher) FetchPluginDetails(req *PluginRequest) (*PluginResponse, error) {
	// Fetch the plugin details from Cassandra
	pluginDetails, err := FetchPluginDetails(c.log, req.PluginName)
	if err != nil {
		c.log.Errorf("Failed to fetch plugin details from store, %v", err)
		return nil, err
	}

	// Fetch the plugin credentials from Kubernetes
	switch req.PluginName {
	case "argocd":
		return c.FetchArgoCDDetails(
			pluginDetails.Namespace,
			pluginDetails.ReleaseName+"-server",
		)
	}

	return nil, fmt.Errorf("unsupported plugin: %s", req.PluginName)
}

func (c *CredentialFetcher) FetchArgoCDDetails(namespace, releaseName string) (*PluginResponse, error) {
	service, err := c.k8sClient.FetchServiceDetails(&k8s.ServiceDetailsRequest{
		Namespace:   namespace,
		ServiceName: releaseName,
	})
	if err != nil {
		c.log.Errorf("Fetching plugin credentials failed: %v", err)
		return nil, err
	}
	// Depending on the service port details isSSLEnabled can be prepared. For now it is set to false default scenario
	isSSLEnabled := false

	credentialDetails, err := c.k8sClient.FetchSecretDetails(&k8s.SecretDetailsRequest{
		Namespace:  namespace,
		SecretName: "argocd-initial-admin-secret",
	})
	if err != nil {
		c.log.Errorf("Fetching plugin credentials failed: %v", err)
		return nil, err
	}
	return &PluginResponse{
		ServiceURL:   fmt.Sprintf("%s.%s.svc.cluster.local", service.Name, service.Namespace),
		IsSSLEnabled: isSSLEnabled,
		Username:     "admin", // admin user is not available in secret
		Password:     credentialDetails.Data["password"],
	}, nil
}
