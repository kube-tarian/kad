package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/pkg/errors"

	"github.com/intelops/go-common/credentials"
)

func StoreCredential(ctx context.Context, request *agentpb.StoreCredentialRequest) error {
	credAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		return errors.WithMessage(err, "error in initializing vault credential client")
	}

	return credAdmin.PutCredential(ctx, request.CredentialType, request.CredEntityName,
		request.CredIdentifier, request.Credential)
}
