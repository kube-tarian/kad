package api

import (
	"context"
	"fmt"
	"log"

	//logger "github.com/kube-tarian/kad/integrator/common-pkg/logging"
	managedcluster "github.com/kube-tarian/kad/capten/common-pkg/managed-cluster"
	vaultcred "github.com/kube-tarian/kad/capten/common-pkg/vault-cred"
	v1 "k8s.io/api/core/v1"

	"github.com/kube-tarian/kad/capten/agent/internal/pb/agentpb"
	"github.com/kube-tarian/kad/capten/common-pkg/credential"

	"github.com/intelops/go-common/credentials"
)

var (
	clusterCredVaultPaths = map[string]string{"NATS": "generic/nats/auth-token"}
	//natsNameSpace         = "observability"
	//vaultAppRoleToken     = "vault-capten-token"
	vaultAddress      = "http://vault.%s"
	cluserAppRoleName = "approle-%s"
)

func (a *Agent) StoreCredential(ctx context.Context, request *agentpb.StoreCredentialRequest) (*agentpb.StoreCredentialResponse, error) {
	credPath := fmt.Sprintf("%s/%s/%s", request.CredentialType, request.CredEntityName, request.CredIdentifier)
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to intialize credentails client for %s", credPath)
		a.log.Errorf("failed to store credentail for %s, %v", credPath, err)
		return &agentpb.StoreCredentialResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: err.Error(),
		}, nil
	}

	err = credAdmin.PutCredential(ctx, request.CredentialType, request.CredEntityName,
		request.CredIdentifier, request.Credential)
	if err != nil {
		a.log.Audit("security", "storecred", "failed", "system", "failed to store credentail for %s", credPath)
		a.log.Errorf("failed to store credentail for %s, %v", credPath, err)
		return &agentpb.StoreCredentialResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: err.Error(),
		}, nil
	}

	a.log.Audit("security", "storecred", "success", "system", "credentail stored for %s", credPath)
	a.log.Infof("stored credentail for entity %s", credPath)
	return &agentpb.StoreCredentialResponse{
		Status: agentpb.StatusCode_OK,
	}, nil
}

func (a *Agent) GetClusterGlobalValues(ctx context.Context, _ *agentpb.GetClusterGlobalValuesRequest) (*agentpb.GetClusterGlobalValuesResponse, error) {
	values, err := credential.GetClusterGlobalValues(ctx)
	if err != nil {
		a.log.Errorf("%v", err)
		return &agentpb.GetClusterGlobalValuesResponse{
			Status:        agentpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: err.Error(),
		}, nil
	}

	a.log.Infof("fetched cluster global values")
	return &agentpb.GetClusterGlobalValuesResponse{
		Status:       agentpb.StatusCode_OK,
		GlobalValues: []byte(values),
	}, nil
}

func (a *Agent) ConfigureVaultSecret(ctx context.Context, request *agentpb.ConfigureVaultSecretRequest) (*agentpb.ConfigureVaultSecretResponse, error) {
	// return nil, fmt.Errorf("not supported")
	credentialPaths := []string{request.SecretPath}

	// logger.Logger.Infof("Cluster Id",request.ManagedClusterId)
	clusterName := request.ManagedClusterId
	log.Println("Cluster Id", request.ManagedClusterId)
	cluserAppRoleNameStr := fmt.Sprintf(cluserAppRoleName, clusterName)
	token, err := vaultcred.GetAppRoleToken(cluserAppRoleNameStr, credentialPaths)
	log.Println("App role token", token)
	if err != nil {
		fmt.Errorf("failed to  get app role token, %v", err)

		return &agentpb.ConfigureVaultSecretResponse{}, err
	}
	log.Printf("approle token created for cluster %s", request.ManagedClusterId)

	k8sclient, err := managedcluster.GetClusterK8SClient(ctx, request.ManagedClusterId)
	if err != nil {
		fmt.Errorf("failed to initalize k8s client, %v", err)
		return &agentpb.ConfigureVaultSecretResponse{}, err
	}

	cred := map[string][]byte{"token": []byte(token)}
	err = k8sclient.CreateOrUpdateSecret(ctx, request.Namespace, request.SecretName, v1.SecretTypeOpaque, cred, nil)
	if err != nil {
		fmt.Errorf("failed to create cluter vault token secret, %v", err)
		return &agentpb.ConfigureVaultSecretResponse{}, err
	}
	log.Printf("Secret %v created in the ns %v", request.SecretName, request.Namespace)
	//vaultAddressStr := fmt.Sprintf(vaultAddress, cp.cfg.DomainName)
	vaultAddressStr := "http://capten-dev-vault.platform.svc.cluster.local:8200"
	vaultStoreCRData := fmt.Sprintf(vaultStore, request.Namespace, vaultAddressStr, request.SecretName)
	ns, resource, err := k8sclient.DynamicClient.CreateResource(ctx, []byte(vaultStoreCRData))
	if err != nil {
		fmt.Errorf("failed to create cluter vault token secret, %v", err)
		return &agentpb.ConfigureVaultSecretResponse{}, err
	}
	log.Printf("create %s on cluster cluster %s/%s", ns, resource, clusterName)

	natsVaultExternalSecretData := fmt.Sprintf(natsVaultExternalSecret, request.Namespace, clusterCredVaultPaths["NATS"])
	ns, resource, err = k8sclient.DynamicClient.CreateResource(ctx, []byte(natsVaultExternalSecretData))
	if err != nil {
		fmt.Errorf("failed to create vault external secret, %v", err)
		return &agentpb.ConfigureVaultSecretResponse{}, err
	}
	log.Printf("create %s on cluster cluster %s/%s", ns, resource, clusterName)
	return &agentpb.ConfigureVaultSecretResponse{}, nil
}

func (a *Agent) CreateVaultRole(ctx context.Context, request *agentpb.CreateVaultRoleRequest) (*agentpb.CreateVaultRoleResponse, error) {
	return nil, fmt.Errorf("not supported")
}

func (a *Agent) UpdateVaultRole(ctx context.Context, request *agentpb.UpdateVaultRoleRequest) (*agentpb.UpdateVaultRoleResponse, error) {
	return nil, fmt.Errorf("not supported")
}

func (a *Agent) DeleteVaultRole(ctx context.Context, request *agentpb.DeleteVaultRoleRequest) (*agentpb.DeleteVaultRoleResponse, error) {
	return nil, fmt.Errorf("not supported")
}
