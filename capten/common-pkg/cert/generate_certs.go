package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"

	"github.com/pkg/errors"
)

const (
	FilePermission          os.FileMode = 0644
	caBitSize                           = 4096
	OrgName                             = "Intelops"
	RootCACommonName                    = "Capten Agent Root CA"
	ClusterCACertSecretName             = "agent-ca-cert"
	CertManagerNamespace                = "cert-manager"
)

type Key struct {
	Key     *rsa.PrivateKey
	KeyData []byte
}

type Cert struct {
	Cert     *x509.Certificate
	CertData []byte
}

type CertificatesData struct {
	RootKey         *Key
	RootCert        *Cert
	CaChainCertData []byte
}

func GenerateRootCerts() (*CertificatesData, error) {
	rootKey, rootCertTemplate, err := generateCACert()
	if err != nil {
		return nil, err
	}

	return &CertificatesData{
		RootKey:         rootKey,
		RootCert:        rootCertTemplate,
		CaChainCertData: rootCertTemplate.CertData,
	}, nil
}

func generateCACert() (*Key, *Cert, error) {
	rootKey, err := rsa.GenerateKey(rand.Reader, caBitSize)
	if err != nil {
		err = errors.WithMessage(err, "failed to generate RSA key for root certificate")
		return nil, nil, err
	}

	rootCertTemplate := &x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{OrgName},
			CommonName:   RootCACommonName,
		},
		SerialNumber:          big.NewInt(1),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(5, 0, 0),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	rootCert, err := x509.CreateCertificate(rand.Reader, rootCertTemplate, rootCertTemplate, &rootKey.PublicKey, rootKey)
	if err != nil {
		err = errors.WithMessage(err, "failed to create root CA certificate")
		return nil, nil, err
	}

	rootCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rootCert})
	rootKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(rootKey)})

	return &Key{
			Key:     rootKey,
			KeyData: rootKeyPEM,
		},
		&Cert{
			Cert:     rootCertTemplate,
			CertData: rootCertPEM,
		}, nil
}
