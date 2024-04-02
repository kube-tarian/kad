package crossplane

import (
	"context"
	"fmt"

	managedcluster "github.com/kube-tarian/kad/capten/common-pkg/managed-cluster"
	vaultcred "github.com/kube-tarian/kad/capten/common-pkg/vault-cred"
	v1 "k8s.io/api/core/v1"
)

var (
	vaultAppRoleTokenSecret = "approle-vault-token"
	vaultAddress            = "http://vault.%s"
	cluserAppRoleName       = "capten-approle-%s"
	secretStoreName         = "approle-vault-store"
)

func (cp *CrossPlaneApp) configureExternalSecretsOnCluster(ctx context.Context,
	clusterName, clusterID string, extSecrets []clusterExternalSecret) error {
	logger.Infof("configure external secrets for cluster %s/%s", clusterName, clusterID)

	credentialPaths, namespaces := getUniqueSecretPathsAndNamespaces(extSecrets)
	if len(namespaces) == 0 {
		logger.Infof("no external secrets defined for cluster %s/%s", clusterName, clusterID)
		return nil
	}

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

	vaultAddressStr := fmt.Sprintf(vaultAddress, cp.cfg.DomainName)
	logger.Infof("Vault Address String", vaultAddressStr)
	for _, namespace := range namespaces {
		cred := map[string][]byte{"token": []byte(token)}
		err = k8sclient.CreateOrUpdateSecret(ctx, namespace, vaultAppRoleTokenSecret, v1.SecretTypeOpaque, cred, nil)
		if err != nil {
			logger.Infof("failed to create cluter vault token secret %s/%s, %v", namespace, vaultAppRoleTokenSecret, err)
			continue
		}

		err := k8sclient.CreateOrUpdateSecretStore(ctx, secretStoreName, namespace,
			vaultAddressStr, vaultAppRoleTokenSecret, "token")
		if err != nil {
			return fmt.Errorf("failed to create cluter vault token secret, %v", err)
		}
		logger.Infof("created %s/%s on cluster cluster %s", namespace, secretStoreName, clusterName)
	}

	for _, extSecret := range extSecrets {
		externalSecretName := "external-" + extSecret.SecretName
		vaultSecretData := map[string]string{}
		for _, secretData := range extSecret.VaultSecrets {
			vaultSecretData[secretData.SecretKey] = secretData.SecretPath
		}
		err := k8sclient.CreateOrUpdateExternalSecret(ctx, externalSecretName, extSecret.Namespace,
			secretStoreName, extSecret.SecretName, extSecret.SecretType, vaultSecretData)
		if err != nil {
			logger.Infof("failed to create vault external secret, %v", err)
			continue
		}
		logger.Infof("created %s/%s on cluster cluster %s", extSecret.Namespace, externalSecretName, clusterName)
	}
	return nil
}

func getUniqueSecretPathsAndNamespaces(extSecrets []clusterExternalSecret) ([]string, []string) {
	credentialPaths := map[string]bool{}
	namspaces := map[string]bool{}
	for _, extSecret := range extSecrets {
		for _, secretData := range extSecret.VaultSecrets {
			credentialPaths[secretData.SecretPath] = true
		}
		namspaces[extSecret.Namespace] = true
	}
	return getKeysFromBoolMap(credentialPaths), getKeysFromBoolMap(namspaces)
}

func getKeysFromBoolMap(inputMap map[string]bool) []string {
	var keys []string

	for key := range inputMap {
		keys = append(keys, key)
	}

	return keys
}
