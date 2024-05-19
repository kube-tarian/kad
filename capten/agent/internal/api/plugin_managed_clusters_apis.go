package api

import (
	"context"
	"fmt"

	"github.com/intelops/go-common/credentials"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/captenpluginspb"
)

const ManagedClusterEntityName = "managedcluster"

func (a *Agent) GetManagedClusters(ctx context.Context, request *captenpluginspb.GetManagedClustersRequest) (
	*captenpluginspb.GetManagedClustersResponse, error) {
	a.log.Infof("Get Managed Clusters request recieved")

	managedClusters, err := a.as.GetManagedClusters()
	if err != nil {
		a.log.Errorf("failed to get managedClusters from db, %v", err)
		return &captenpluginspb.GetManagedClustersResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "couldn't fetch managed clusters",
		}, nil
	}

	for _, r := range managedClusters {
		_, secretPath, secretKeys, err := a.getManagedClusterCredential(ctx, r.Id)
		if err != nil {
			a.log.Errorf("failed to get credential, %v", err)
			return &captenpluginspb.GetManagedClustersResponse{
				Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
				StatusMessage: "failed to fetch managed clusters",
			}, err
		}
		r.SecretePath = secretPath
		r.SecreteKeys = secretKeys
	}

	a.log.Infof("Fetched %d Managed Clusters", len(managedClusters))
	return &captenpluginspb.GetManagedClustersResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successfully fetched the Crossplane projects",
		Clusters:      managedClusters,
	}, nil
}

func (a *Agent) GetManagedClusterKubeconfig(ctx context.Context, request *captenpluginspb.GetManagedClusterKubeconfigRequest) (
	*captenpluginspb.GetManagedClusterKubeconfigResponse, error) {
	if err := validateArgs(request.Id); err != nil {
		a.log.Infof("request validation failed", err)
		return &captenpluginspb.GetManagedClusterKubeconfigResponse{
			Status:        captenpluginspb.StatusCode_INVALID_ARGUMENT,
			StatusMessage: "request validation failed",
		}, nil
	}

	a.log.Infof("Get Managed Cluster %s kubeconfig request recieved", request.Id)

	creds, _, _, err := a.getManagedClusterCredential(ctx, request.GetId())
	if err != nil {
		a.log.Errorf("failed to get managedClusters kubeconfig from vault, %v", err)
		return &captenpluginspb.GetManagedClusterKubeconfigResponse{
			Status:        captenpluginspb.StatusCode_INTERNAL_ERROR,
			StatusMessage: "couldn't fetch managed clusters",
		}, nil
	}

	a.log.Infof("Fetched %d Managed Clusters", 0)
	return &captenpluginspb.GetManagedClusterKubeconfigResponse{
		Status:        captenpluginspb.StatusCode_OK,
		StatusMessage: "successfully fetched the Crossplane projects",
		Kubeconfig:    creds["kubeconfig"],
	}, nil
}

func (a *Agent) getManagedClusterCredential(ctx context.Context, id string) (map[string]string, string, []string, error) {
	credPath := fmt.Sprintf("%s/%s/%s", credentials.GenericCredentialType, ManagedClusterEntityName, id)
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to intialize credentials client for %s", credPath)
		a.log.Errorf("failed to get crendential for %s, %v", credPath, err)
		return nil, "", nil, err
	}

	cred, err := credAdmin.GetCredential(ctx, credentials.GenericCredentialType, ManagedClusterEntityName, id)
	if err != nil {
		a.log.Errorf("failed to get credential for %s, %v", credPath, err)
		return nil, "", nil, err
	}

	secretKeys := []string{}
	for key := range cred {
		secretKeys = append(secretKeys, key)
	}
	return cred, credPath, secretKeys, nil
}

// store managed cluster kubeconfig and endpoint in vault
func (a *Agent) StoreManagedClusterCredential(ctx context.Context, id string, kubeconfig string, endpoint string) error {
	credPath := fmt.Sprintf("%s/%s/%s", credentials.GenericCredentialType, ManagedClusterEntityName, id)
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to intialize credentials client for %s", credPath)
		a.log.Errorf("failed to store credential for %s, %v", credPath, err)
		return err
	}

	credentialMap := map[string]string{
		"kubeconfig": kubeconfig,
		"endpoint":   endpoint,
	}
	err = credAdmin.PutCredential(ctx, credentials.GenericCredentialType, ManagedClusterEntityName,
		id, credentialMap)

	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to store crendential for %s", credPath)
		a.log.Errorf("failed to store credential for %s, %v", credPath, err)
		return err
	}
	a.log.Audit("security", "storecred", "success", "system", "credential stored for %s", credPath)
	a.log.Infof("stored credential for entity %s", credPath)
	return nil
}
