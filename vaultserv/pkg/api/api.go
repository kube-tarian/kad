package api

import (
	"context"

	"github.com/kube-tarian/kad/vaultserv/pkg/client"
	"github.com/kube-tarian/kad/vaultserv/pkg/pb/vaultservpb"
	"go.uber.org/zap"
)

type VaultServ struct {
	vaultservpb.UnimplementedVaultServer
	client *client.Vault
}

func NewVaultServ() (*VaultServ, error) {
	vaultClient, err := client.NewVault()
	if err != nil {
		return nil, err
	}
	return &VaultServ{
		client: vaultClient,
	}, nil
}

func (v *VaultServ) StoreCred(ctx context.Context, request *vaultservpb.StoreCredRequest) (*vaultservpb.StoreCredResponse, error) {
	log, _ := zap.NewProduction()
	defer log.Sync()

	err := v.client.PutCredential("secret", request)
	if err != nil {
		log.Error("failed to store cred", zap.Error(err))
		return &vaultservpb.StoreCredResponse{
			Status: "FAILED",
		}, err
	}

	return &vaultservpb.StoreCredResponse{
		Status: "SUCCESS",
	}, nil
}
func (v *VaultServ) StoreSecret(ctx context.Context, request *vaultservpb.StoreSecretRequest) (*vaultservpb.StoreSecretResponse, error) {
	log, _ := zap.NewProduction()
	defer log.Sync()

	err := v.client.PutSecret("secret", request)
	if err != nil {
		log.Error("failed to store cred", zap.Error(err))
		return &vaultservpb.StoreSecretResponse{
			Status: "FAILED",
		}, err
	}

	return &vaultservpb.StoreSecretResponse{
		Status: "SUCCESS",
	}, nil
}
