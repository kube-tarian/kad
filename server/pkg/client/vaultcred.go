package client

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	vaultcredclient "github.com/intelops/go-common/vault-cred-client"
)

func PutCaptenClusterCertificate(ctx context.Context, certIdentifier, caCertData, keyData, certData string) error {
	certAdmin, err := vaultcredclient.NewCertificateAdmin()
	if err != nil {
		return errors.WithMessage(err, "error in initializing vault credential client")
	}

	err = certAdmin.StoreCertificate(ctx, vaultcredclient.CaptenClusterCert, certIdentifier,
		vaultcredclient.CertificateData{
			CACert: caCertData,
			Key:    keyData,
			Cert:   certData,
		})
	if err != nil {
		err = errors.WithMessage(err, "error in reading ceritification from vault")
	}
	return err
}

func GetCaptenClusterCertificate(ctx context.Context, certIdentifier string) (caCertData, keyData, certData string, err error) {
	certReader, err := vaultcredclient.NewCertificateReader()
	if err != nil {
		err = errors.WithMessage(err, "error in initializing vault credential client")
		return
	}

	resCertData, err := certReader.GetCertificate(ctx, vaultcredclient.CaptenClusterCert, certIdentifier)
	if err != nil {
		fmt.Printf("read failed %v\n", err)
		return
	}
	caCertData = resCertData.CACert
	keyData = resCertData.Key
	certData = resCertData.Cert
	return
}
