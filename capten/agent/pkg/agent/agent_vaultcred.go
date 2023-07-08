package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/pkg/errors"

	vaultcredclient "github.com/intelops/go-common/vault-cred-client"
)

func StoreCredential(ctx context.Context, request *agentpb.StoreCredRequest) error {
	credAdmin, err := vaultcredclient.NewServiceCredentailAdmin()
	if err != nil {
		return errors.WithMessage(err, "error in initializing vault credential client")
	}

	return credAdmin.PutServiceCredential(ctx, request.Credname, request.Username,
		vaultcredclient.ServiceCredentail{UserName: request.Username, Password: request.Password})
}

func StoreGerericCredential(ctx context.Context, request *agentpb.StoreGerericCredRequest) error {
	credAdmin, err := vaultcredclient.NewGerericCredentailAdmin()
	if err != nil {
		return errors.WithMessage(err, "error in initializing vault credential client")
	}
	return credAdmin.PutGenericCredential(ctx, request.CredentialType, request.CredEntityName, request.CredIdentifier,
		vaultcredclient.GerericCredentail{Credential: request.Credential})
}
