package api

import (
	"context"
	"fmt"

	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/agentpb"
	vaultcred "github.com/kube-tarian/kad/capten/common-pkg/vault-cred"
	
	v1 "k8s.io/api/core/v1"

	"github.com/kube-tarian/kad/capten/common-pkg/credential"

	"github.com/intelops/go-common/credentials"
)

var (
	vaultAddress     = "http://vault.%s"
	kadAppRolePrefix = "kad-approle-"
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
	a.log.Infof("Configure Vault Secret Request recieved for secret ", request.SecretName)

	secretPathsData := map[string]string{}
	secretPaths := []string{}
	for _, secretPathData := range request.SecretPathData {
		secretPathsData[secretPathData.SecretKey] = secretPathData.SecretPath
		secretPaths = append(secretPaths, secretPathData.SecretPath)
	}

	appRoleName := kadAppRolePrefix + request.SecretName
	token, err := vaultcred.GetAppRoleToken(appRoleName, secretPaths)
	if err != nil {
		a.log.Errorf("failed to create app role token for %s, %v", appRoleName, err)
		return &agentpb.ConfigureVaultSecretResponse{Status: agentpb.StatusCode_INTERNRAL_ERROR}, err
	}

	k8sclient, err := k8s.NewK8SClient(a.log)
	if err != nil {
		a.log.Errorf("failed to initalize k8s client, %v", err)
		return &agentpb.ConfigureVaultSecretResponse{Status: agentpb.StatusCode_INTERNRAL_ERROR}, err
	}

	cred := map[string][]byte{"token": []byte(token)}
	vaultTokenSecretName := "vault-token-" + request.SecretName
	err = k8sclient.CreateOrUpdateSecret(ctx, request.Namespace, vaultTokenSecretName, v1.SecretTypeOpaque, cred, nil)
	if err != nil {
		a.log.Errorf("failed to create cluter vault token secret, %v", err)
		return &agentpb.ConfigureVaultSecretResponse{Status: agentpb.StatusCode_INTERNRAL_ERROR}, err
	}

	vaultAddressStr := fmt.Sprintf(vaultAddress, a.cfg.DomainName)
	secretStoreName := "ext-store-" + request.SecretName
	err = k8sclient.CreateOrUpdateSecretStore(ctx, secretStoreName, request.Namespace, vaultAddressStr, vaultTokenSecretName, "token")
	if err != nil {
		a.log.Errorf("failed to create cluter vault token secret, %v", err)
		return &agentpb.ConfigureVaultSecretResponse{Status: agentpb.StatusCode_INTERNRAL_ERROR}, err
	}
	a.log.Infof("created secret store %s/%s", request.Namespace, secretStoreName)

	externalSecretName := "ext-secret-" + request.SecretName
	err = k8sclient.CreateOrUpdateExternalSecret(ctx, externalSecretName, request.Namespace, secretStoreName,
		request.SecretName, "", secretPathsData)
	if err != nil {
		a.log.Errorf("failed to create vault external secret, %v", err)
		return &agentpb.ConfigureVaultSecretResponse{Status: agentpb.StatusCode_INTERNRAL_ERROR}, err
	}
	a.log.Infof("created external secret %s/%s", request.Namespace, externalSecretName)
	return &agentpb.ConfigureVaultSecretResponse{Status: agentpb.StatusCode_OK}, nil
}

func (a *Agent) CreateVaultRole(ctx context.Context, request *agentpb.CreateVaultRoleRequest) (*agentpb.CreateVaultRoleResponse, error) {
	a.log.Infof("vault role %s creat request for cluster %s", request.RoleName, request.ManagedClusterName)
	secretPolicies := map[string]int{}
	for _, secretPolicy := range request.SecretPolicy {
		secretPolicies[secretPolicy.SecretPath] = int(secretPolicy.Access)
	}

	err := vaultcred.CreateVaultRole(request.ManagedClusterName, request.RoleName,
		request.Namespaces, request.ServiceAccounts, secretPolicies)
	if err != nil {
		a.log.Errorf("failed to create vault role, %v", err)
		return &agentpb.CreateVaultRoleResponse{Status: agentpb.StatusCode_INTERNRAL_ERROR}, err
	}
	a.log.Infof("created vault role %s for cluster %s", request.RoleName, request.ManagedClusterName)
	return &agentpb.CreateVaultRoleResponse{Status: agentpb.StatusCode_OK}, nil
}

func (a *Agent) UpdateVaultRole(ctx context.Context, request *agentpb.UpdateVaultRoleRequest) (*agentpb.UpdateVaultRoleResponse, error) {
	return nil, fmt.Errorf("not supported")
}

func (a *Agent) DeleteVaultRole(ctx context.Context, request *agentpb.DeleteVaultRoleRequest) (*agentpb.DeleteVaultRoleResponse, error) {
	return nil, fmt.Errorf("not supported")
}
