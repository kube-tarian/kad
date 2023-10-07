package agent

import (
	"context"
	"fmt"

	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/credential"

	"github.com/intelops/go-common/credentials"
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
