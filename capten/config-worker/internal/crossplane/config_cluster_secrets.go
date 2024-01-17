package crossplane

import (
	"context"
	"fmt"

	managedcluster "github.com/kube-tarian/kad/capten/common-pkg/managed-cluster"
	vaultcred "github.com/kube-tarian/kad/capten/common-pkg/vault-cred"
	v1 "k8s.io/api/core/v1"
)

var (
	clusterCredVaultPaths = map[string]string{"NATS": "generic/nats/auth-token"}
	natsNameSpace         = "observability"
	vaultAppRoleToken     = "vault-cluster-token"
	vaultAddress          = "http://vault.%s"
	cluserAppRoleName     = "approle-%s"
)

func (cp *CrossPlaneApp) configureExternalSecretsOnCluster(ctx context.Context, clusterName, clusterID string) error {
	credentialPaths := getCrdentialPaths()
	cluserAppRoleNameStr := fmt.Sprintf(cluserAppRoleName, clusterName)
	token, err := vaultcred.GetAppRoleToken(cluserAppRoleNameStr, credentialPaths)
	if err != nil {
		return err
	}
	logger.Infof("approle token created for cluster %s/%s", clusterName, clusterID)

	k8sclient, err := managedcluster.GetClusterK8SClient(ctx, clusterID)
	if err != nil {
		return fmt.Errorf("failed to initalize k8s client, %v", err)
	}

	cred := map[string][]byte{"token": []byte(token)}
	err = k8sclient.CreateOrUpdateSecret(ctx, natsNameSpace, vaultAppRoleToken, v1.SecretTypeOpaque, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create cluter vault token secret, %v", err)
	}

	vaultAddress = fmt.Sprintf(vaultAddress, cp.cfg.DomainName)
	vaultStoreCRData := fmt.Sprintf(vaultStore, natsNameSpace, vaultAddress)
	ns, resource, err := k8sclient.DynamicClient.CreateResource(ctx, []byte(vaultStoreCRData))
	if err != nil {
		return fmt.Errorf("failed to create cluter vault token secret, %v", err)
	}
	logger.Infof("create %s/%s on cluster cluster %s/%s", ns, resource, clusterName)

	natsVaultExternalSecretData := fmt.Sprintf(natsVaultExternalSecret, natsNameSpace, clusterCredVaultPaths["NATS"])
	ns, resource, err = k8sclient.DynamicClient.CreateResource(ctx, []byte(natsVaultExternalSecretData))
	if err != nil {
		return fmt.Errorf("failed to create vault external secret, %v", err)
	}
	logger.Infof("create %s/%s on cluster cluster %s/%s", ns, resource, clusterName)
	return nil
}

func getCrdentialPaths() []string {
	credentialPaths := []string{}
	for _, credentialPath := range clusterCredVaultPaths {
		credentialPaths = append(credentialPaths, credentialPath)
	}
	return credentialPaths
}
