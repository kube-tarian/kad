package client

import (
	"context"
	"fmt"
	"github.com/kube-tarian/kad/server/pkg/config"
	"github.com/kube-tarian/kad/server/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	vaultpb "github.com/intelops/vault-cred/proto/pb/vaultcredpb"
)

const (
	vaultAddressCfgKey = "vault.address"
	vaultPortCfgKey    = "vault.port"
)

type Vault struct {
	client vaultpb.VaultCredClient
}

func NewVault() (*Vault, error) {
	cfg := config.GetConfig()
	target := fmt.Sprintf("%s:%d", cfg.GetString(vaultAddressCfgKey), cfg.GetInt(vaultPortCfgKey))
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := vaultpb.NewVaultCredClient(conn)
	return &Vault{
		client: client,
	}, nil
}

func (v *Vault) PutCert(ctx context.Context, orgId, clusterName, caCertData, keyData, certData string) error {
	certReqMap := map[string]string{
		types.ClientCertChainFileName: caCertData,
		types.ClientKeyFileName:       keyData,
		types.ClientCertFileName:      certData,
	}

	putCredReq := vaultpb.PutCredRequest{
		CredentialType: "cert",
		CredEntityName: orgId,
		CredIdentifier: clusterName,
		Credential:     certReqMap,
	}

	_, err := v.client.PutCred(ctx, &putCredReq)
	return err
}

func (v *Vault) GetCert(ctx context.Context, orgId, clusterName string) (map[string]string, error) {
	getCredReq := vaultpb.GetCredRequest{
		CredentialType: "cert",
		CredEntityName: orgId,
		CredIdentifier: clusterName,
	}

	response, err := v.client.GetCred(ctx, &getCredReq)
	return response.Credential, err
}
