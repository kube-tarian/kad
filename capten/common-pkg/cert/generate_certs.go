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

	"github.com/intelops/go-common/logging"
	"github.com/pkg/errors"
)

const (
	FilePermission           os.FileMode = 0644
	caBitSize                            = 4096
	certBitSize                          = 2048
	rootCAKeyFileName                    = "root.key"
	rootCACertFileName                   = "root.crt"
	interCAKeyFileName                   = "inter-ca.key"
	interCACertFileName                  = "inter-ca.crt"
	CAFileName                           = "ca.crt"
	OrgName                              = "Intelops"
	RootCACommonName                     = "Capten Agent Root CA"
	IntermediateCACommonName             = "Capten Agent Cluster CA"
	ClusterCACertSecretName              = "agent-ca-cert"
	CertManagerNamespace                 = "cert-manager"
)

var (
	log = logging.NewLogger()
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
	InterKey        *Key
	InterCert       *Cert
	CaChainCertData []byte
}

func GenerateRootCerts() (*CertificatesData, error) {
	rootKey, rootCertTemplate, err := generateCACert()
	if err != nil {
		return nil, err
	}

	interKey, interCACertTemplate, err := generateIntermediateCACert(rootKey.Key, rootCertTemplate.Cert)
	if err != nil {
		return nil, err
	}

	caCertChain, err := generateCACertChain(rootCertTemplate.CertData, interCACertTemplate.CertData)
	if err != nil {
		return nil, err
	}
	log.Infof("%v\n%v\n", interKey, interCACertTemplate, caCertChain)
	return &CertificatesData{
		RootKey:         rootKey,
		RootCert:        rootCertTemplate,
		InterKey:        interKey,
		InterCert:       interCACertTemplate,
		CaChainCertData: caCertChain,
	}, nil
}

func generateCACert() (*Key, *Cert, error) { //(rootKey *rsa.PrivateKey, rootCertTemplate *x509.Certificate, err error) {
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

func generateIntermediateCACert(rootKey *rsa.PrivateKey, rootCertTemplate *x509.Certificate) (*Key, *Cert, error) {
	interKey, err := rsa.GenerateKey(rand.Reader, caBitSize)
	if err != nil {
		err = errors.WithMessage(err, "failed to generate RSA key for intermediate certificate")
		return nil, nil, err
	}

	interCACertTemplate := &x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{OrgName},
			CommonName:   IntermediateCACommonName,
			Locality:     []string{"agent"},
		},
		SerialNumber:          big.NewInt(1),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(2, 0, 0),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	interCert, err := x509.CreateCertificate(rand.Reader, interCACertTemplate, rootCertTemplate, &interKey.PublicKey, rootKey)
	if err != nil {
		err = errors.WithMessage(err, "failed to create intermediate CA certificate")
		return nil, nil, err
	}

	interCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: interCert})
	interKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(interKey)})

	return &Key{
			Key:     interKey,
			KeyData: interKeyPEM,
		},
		&Cert{
			Cert:     interCACertTemplate,
			CertData: interCertPEM,
		}, nil
}

func generateCACertChain(caCertPEMFromFile, interCACertPEMFromFile []byte) ([]byte, error) {
	return append(caCertPEMFromFile, interCACertPEMFromFile...), nil
}
