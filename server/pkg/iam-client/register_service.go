package iamclient

import (
	"github.com/intelops/go-common/logging"
	oryclient "github.com/kube-tarian/kad/server/pkg/ory-client"
	"github.com/pkg/errors"
)

func RegisterService(log logging.Logger) error {
	oryclient, err := oryclient.NewOryClient(log)
	if err != nil {
		return errors.WithMessage(err, "OryClient initialization failed")
	}

	IC, err := NewClient(oryclient, log)
	if err != nil {
		return errors.WithMessage(err, "Error occured while created IAM client")
	}

	err = IC.RegisterWithIam()
	if err != nil {
		return errors.WithMessage(err, "Registering capten server as oauth client failed")
	}

	err = IC.RegisterRolesActions()
	if err != nil {
		return errors.WithMessage(err, "Registering Roles and Actions in IAM failed")
	}
	return nil
}
