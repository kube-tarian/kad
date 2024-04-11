package vaultcred

import (
	"context"
	"fmt"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"github.com/kelseyhightower/envconfig"
	managedcluster "github.com/kube-tarian/kad/capten/common-pkg/managed-cluster"
	"github.com/kube-tarian/kad/capten/common-pkg/vault-cred/vaultcredpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	v1 "k8s.io/api/core/v1"
)

type config struct {
	VaultCredAddress string `envconfig:"VAULT_CRED_ADDR" default:"vault-cred:8080"`
}

func GetAppRoleToken(appRoleName string, credentialPaths []string) (string, error) {
	conf := &config{}
	if err := envconfig.Process("", conf); err != nil {
		return "", fmt.Errorf("vault cred config read failed, %v", err)
	}

	vc, err := grpc.Dial(conf.VaultCredAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(timeout.UnaryClientInterceptor(60*time.Second)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    30, // seconds
			Timeout: 10, // seconds
		}))
	if err != nil {
		return "", fmt.Errorf("failed to connect vauld-cred server, %v", err)
	}
	vcClient := vaultcredpb.NewVaultCredClient(vc)

	tokenData, err := vcClient.CreateAppRoleToken(context.Background(), &vaultcredpb.CreateAppRoleTokenRequest{
		AppRoleName: appRoleName,
		SecretPaths: credentialPaths,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate app role token for %s, %v", appRoleName, err)
	}
	return tokenData.Token, nil
}

func DeleteAppRole(appRoleName string) error {
	conf := &config{}
	if err := envconfig.Process("", conf); err != nil {
		return fmt.Errorf("vault cred config read failed, %v", err)
	}

	vc, err := grpc.Dial(conf.VaultCredAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(timeout.UnaryClientInterceptor(60*time.Second)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    30, // seconds
			Timeout: 10, // seconds
		}))
	if err != nil {
		return fmt.Errorf("failed to connect vauld-cred server, %v", err)
	}
	vcClient := vaultcredpb.NewVaultCredClient(vc)

	resp, err := vcClient.DeleteAppRole(context.Background(), &vaultcredpb.DeleteAppRoleRequest{
		RoleName: appRoleName,
	})
	if err != nil {
		return fmt.Errorf("failed to delete app role %s, reason %v", appRoleName, err)
	} else if resp.Status != vaultcredpb.StatusCode_OK {
		return fmt.Errorf("failed to delete app role %s, stauts %v, message: %v", appRoleName, resp.Status, resp.StatusMessage)
	}
	return nil
}

func RegisterClusterVaultAuth(clusterID, clusterName string) error {
	conf := &config{}
	if err := envconfig.Process("", conf); err != nil {
		return fmt.Errorf("vault cred config read failed, %v", err)
	}

	vc, err := grpc.Dial(conf.VaultCredAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect vauld-cred server, %v", err)
	}
	vcClient := vaultcredpb.NewVaultCredClient(vc)

	k8sClient, err := managedcluster.GetClusterK8SClient(context.TODO(), clusterID)
	if err != nil {
		return fmt.Errorf("failed to connect cluster %s k8s server, %v", clusterName, err)
	}

	err = k8sClient.CreateOrUpdateServiceAccount(context.TODO(), "default", "vault-auth")
	if err != nil {
		return fmt.Errorf("failed to create service account on cluster %s, %v", clusterName, err)
	}

	err = k8sClient.CreateOrUpdateSecret(context.TODO(), "default", "vault-auth",
		v1.SecretTypeServiceAccountToken, nil, map[string]string{"kubernetes.io/service-account.name": "vault-auth"})
	if err != nil {
		return fmt.Errorf("failed to create secret on cluster %s, %v", clusterName, err)
	}

	err = k8sClient.CreateOrUpdateClusterRoleBinding(context.TODO(), map[string]string{"vault-auth": "default"}, "system:auth-delegator")
	if err != nil {
		return fmt.Errorf("failed to create cluster role binding on cluster %s, %v", clusterName, err)
	}

	secretData, err := k8sClient.GetSecretData("default", "vault-auth")
	if err != nil {
		return fmt.Errorf("failed to read secret from cluster %s, %v", clusterName, err)
	}

	_, err = vcClient.AddClusterK8SAuth(context.Background(), &vaultcredpb.AddClusterK8SAuthRequest{
		ClusterName: clusterName,
		Host:        k8sClient.Config.Host,
		CaCert:      secretData.Data["ca.crt"],
		JwtToken:    secretData.Data["token"],
	})
	if err != nil {
		return fmt.Errorf("failed to add k8s auth for cluster %s, %v", clusterName, err)
	}
	return nil
}

func CreateVaultRole(clusterName, roleName string, namespaces, serviceAccounts []string, securityPolicies map[string]int) error {
	conf := &config{}
	if err := envconfig.Process("", conf); err != nil {
		return fmt.Errorf("vault cred config read failed, %v", err)
	}

	vc, err := grpc.Dial(conf.VaultCredAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect vauld-cred server, %v", err)
	}

	vaultSecurityPolicies := []*vaultcredpb.SecretPolicy{}
	for secretPath, access := range securityPolicies {
		vaultSecurityPolicy := &vaultcredpb.SecretPolicy{
			SecretPath: secretPath,
			Access:     vaultcredpb.SecretAccess(access),
		}
		vaultSecurityPolicies = append(vaultSecurityPolicies, vaultSecurityPolicy)
	}

	vcClient := vaultcredpb.NewVaultCredClient(vc)
	_, err = vcClient.CreateK8SAuthRole(context.TODO(), &vaultcredpb.CreateK8SAuthRoleRequest{
		RoleName:        roleName,
		ClusterName:     clusterName,
		Namespaces:      namespaces,
		ServiceAccounts: serviceAccounts,
		SecretPolicy:    vaultSecurityPolicies,
	})
	if err != nil {
		return fmt.Errorf("failed to create vault role %s for cluster %s, %v", roleName, clusterName, err)
	}
	return nil
}
