package client

import (
	"context"

	"github.com/intelops/go-common/credentials"
	"github.com/pkg/errors"
)

const (
	captenClusterCertType = "capten-cert"
)

func PutCaptenClusterCertificate(ctx context.Context, certIdentifier, caCertData, keyData, certData string) error {
	certAdmin, err := credentials.NewCredentialAdmin(ctx)
	if err != nil {
		return errors.WithMessage(err, "error in initializing vault credential client")
	}

	err = certAdmin.PutCertificateData(ctx, captenClusterCertType, certIdentifier,
		credentials.CertificateData{
			CACert: caCertData,
			Key:    keyData,
			Cert:   certData,
		})
	if err != nil {
		err = errors.WithMessage(err, "error in storing ceritification")
	}
	return err
}

func GetCaptenClusterCertificate(ctx context.Context, certIdentifier string) (caCertData, keyData, certData string, err error) {
	certReader, err := credentials.NewCredentialReader(ctx)
	if err != nil {
		err = errors.WithMessage(err, "error in initializing vault credential client")
		return
	}

	resCertData, err := certReader.GetCertificateData(ctx, captenClusterCertType, certIdentifier)
	if err != nil {
		err = errors.WithMessage(err, "error in reading certificate")
		return
	}
	caCertData = resCertData.CACert
	keyData = resCertData.Key
	certData = resCertData.Cert
	return
}
