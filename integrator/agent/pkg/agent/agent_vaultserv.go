package agent

import (
	"context"
	"fmt"

	"github.com/kube-tarian/kad/integrator/agent/pkg/vaultservpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Vault struct {
	client vaultservpb.VaultClient
}

func GetVaultServClient() (*Vault, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "127.0.0.1", 9098),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("failed to connect: ", err)
		return nil, err
	}

	fmt.Printf("gRPC connection started to %s:%d", "127.0.0.1", 9098)
	vaultClient := vaultservpb.NewVaultClient(conn)
	return &Vault{
		client: vaultClient,
	}, nil
}

func (v *Vault) StoreCred(ctx context.Context, request *vaultservpb.StoreCredRequest) (*vaultservpb.StoreCredResponse, error) {
	return v.client.StoreCred(ctx, request)
}
func (v *Vault) StoreSecret(ctx context.Context, request *vaultservpb.StoreSecretRequest) (*vaultservpb.StoreSecretResponse, error) {
	return v.client.StoreSecret(ctx, request)
}
