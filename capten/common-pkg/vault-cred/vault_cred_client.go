package vaultcred

import (
	"context"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/vault-cred/vaultcredpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type config struct {
	VaultCredAddress string `envconfig:"VAULT_CRED_ADDR" default:"vault-cred:8080"`
}

func GetAppRoleToken(appRoleName string, credentialPaths []string) (string, error) {
	conf := &config{}
	if err := envconfig.Process("", conf); err != nil {
		return "", fmt.Errorf("vault cred config read failed, %v", err)
	}

	vc, err := grpc.Dial(conf.VaultCredAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
