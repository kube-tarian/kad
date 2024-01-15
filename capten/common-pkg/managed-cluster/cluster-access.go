package managedcluster

import (
	"context"
	"fmt"

	"github.com/intelops/go-common/credentials"
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/credential"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
)

const (
	kubeConfig               = "kubeconfig"
	k8sEndpoint              = "endpoint"
	k8sClusterCA             = "clusterCA"
	managedClusterEntityName = "managedcluster"
	clusterSecretName        = "%s-cluster"
)

var logger = logging.NewLogger()

func StoreClusterAccessData(ctx context.Context, namespace, clusterID string) error {
	k8sclient, err := k8s.NewK8SClient(logger)
	if err != nil {
		return fmt.Errorf("failed to get k8s client, %v", err)
	}

	secretName := fmt.Sprintf(clusterSecretName, clusterID)
	resp, err := k8sclient.GetSecretData(namespace, secretName)
	if err != nil {
		return fmt.Errorf("failed to get secret %s/%s, %v", namespace, secretName, err)
	}

	cred := map[string]string{}
	cred[kubeConfig] = resp.Data[kubeConfig]
	cred[k8sClusterCA] = resp.Data[k8sClusterCA]
	cred[k8sEndpoint] = resp.Data[k8sEndpoint]

	err = credential.PutGenericCredential(context.TODO(), managedClusterEntityName, clusterID, cred)
	if err != nil {
		return fmt.Errorf("failed to store cluster access data for cluster %s, %v", clusterID, err)
	}
	logger.Infof("stored cluster access data for cluster %s", clusterID)
	return nil
}

func DeleteClusterAccessData(ctx context.Context, clusterID string) error {
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		logger.Errorf("failed to delete credential for cluster %s, %v", clusterID, err)
		return err
	}

	err = credAdmin.DeleteCredential(ctx, credentials.GenericCredentialType, managedClusterEntityName, clusterID)
	if err != nil {
		logger.Errorf("failed to delete credential for cluster %s, %v", clusterID, err)
		return err
	}
	logger.Infof("deleted cluster access data for cluster %s", clusterID)
	return nil
}

func GetClusterAccessData(ctx context.Context, clusterID string) (string, string, string, error) {
	cred, err := credential.GetGenericCredential(ctx, managedClusterEntityName, clusterID)
	if err != nil {
		logger.Errorf("failed to delete credential for cluster %s, %v", clusterID, err)
		return "", "", "", err
	}

	logger.Infof("fetch cluster access data for cluster %s", clusterID)
	return cred[kubeConfig], cred[k8sClusterCA], cred[k8sEndpoint], nil
}

func GetClusterK8SClient(ctx context.Context, clusterID string) (*k8s.K8SClient, error) {
	kubeConfig, clusterCA, k8sEndpoint, err := GetClusterAccessData(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	k8sclient, err := k8s.NewK8SClientForCluster(logging.NewLogger(), kubeConfig, clusterCA, k8sEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to initalize k8s client: %v", err)
	}
	return k8sclient, nil
}
