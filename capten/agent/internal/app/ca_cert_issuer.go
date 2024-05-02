package app

import (
	"context"
	"fmt"
	"os"
	"time"

	v1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	cmclient "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"github.com/kube-tarian/kad/capten/common-pkg/cert"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"github.com/pkg/errors"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

const (
	certFileName         = "server.cert"
	keyFileName          = "server.key"
	agentIssuerName      = "agent-ca-issuer"
	namespace            = "capten"
	serverCertSecretName = "agent-server-mtls"
)

func setupCACertIssuser() error {
	k8sclient, err := k8s.NewK8SClient(log)
	if err != nil {
		log.Errorf("failed to initalize k8s client, %v", err)
		return err
	}

	_, err = setupCertificateIssuer(k8sclient)
	if err != nil {
		log.Errorf("Setup Certificates Issuer failed, %v", err)
		return err
	}

	err = generateServerCertificates(k8sclient)
	if err != nil {
		log.Errorf("Server certificates generation failed, %v", err)
		return err
	}

	// r.RunTLS(fmt.Sprintf("%s:%d", cfg.Host, cfg.RestPort), certFileName, keyFileName)
	// r.Run(fmt.Sprintf("%s:%d", cfg.Host, cfg.RestPort))
	return nil
}

// Setup agent certificate issuer
func setupCertificateIssuer(k8sclient *k8s.K8SClient) (*cert.CertificatesData, error) {
	// TODO: Check certificates exist in Vault and control plan cluster
	// If exist skip
	// Else
	// 1. generate root certificates
	// 2. Create Certificate Issuer
	// 3. Store in Vault
	certsData, err := cert.GenerateRootCerts()
	if err != nil {
		return nil, err
	}

	err = k8s.CreateOrUpdateClusterCAIssuerSecret(k8sclient, certsData.RootCert.CertData, certsData.RootKey.KeyData, certsData.CaChainCertData)
	if err != nil {
		return nil, fmt.Errorf("failed to create/update CA Issuer Secret: %v", err)
	}

	err = k8s.CreateOrUpdateClusterIssuer(agentIssuerName)
	if err != nil {
		return nil, fmt.Errorf("failed to create/update CA Issuer %s in cert-manager: %v", agentIssuerName, err)
	}

	return certsData, nil
}

func generateServerCertificates(k8sClient *k8s.K8SClient) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return errors.WithMessage(err, "error while building kubeconfig")
	}
	cmClient, err := cmclient.NewForConfig(config)
	if err != nil {
		return err
	}

	err = k8sClient.CreateNamespace(context.Background(), namespace)
	if err != nil {
		return fmt.Errorf("failed to create namespace: %v", err)
	}

	err = generateCertManagerServerCertificate(cmClient, namespace, serverCertSecretName, agentIssuerName)
	if err != nil {
		return fmt.Errorf("failed to genereate server certificate: %v", err)
	}

	// TODO: it may take some time for certificate to get create
	// So have to keep wait and retry
	time.Sleep(10 * time.Second)
	secretData, err := k8sClient.GetSecretData(namespace, serverCertSecretName)
	if err != nil {
		return fmt.Errorf("failed to fetch certificates from secret, %v", err)
	}

	// Write certificates to files
	os.WriteFile(certFileName, []byte(secretData.Data["cert"]), cert.FilePermission)
	os.WriteFile(keyFileName, []byte(secretData.Data["key"]), cert.FilePermission)
	return nil
}

func generateCertManagerServerCertificate(cmClient *cmclient.Clientset, namespace string, certName string, issuerRefName string) error {
	usages := []v1.KeyUsage{v1.UsageDigitalSignature, v1.UsageKeyEncipherment, v1.UsageServerAuth}

	_, err := cmClient.CertmanagerV1().Certificates(namespace).Create(
		context.TODO(),
		&v1.Certificate{
			ObjectMeta: metav1.ObjectMeta{
				Name: certName,
			},
			Spec: v1.CertificateSpec{
				IssuerRef: cmmeta.ObjectReference{
					Name: issuerRefName, // "capten-ca-issuer"
					Kind: v1.ClusterIssuerKind,
				},
				SecretName: certName,
				CommonName: certName,
				Usages:     usages,
				PrivateKey: &v1.CertificatePrivateKey{
					Algorithm: v1.RSAKeyAlgorithm,
					Size:      2048,
					Encoding:  v1.PKCS1,
				},
			},
		},
		metav1.CreateOptions{},
	)
	if k8serror.IsAlreadyExists(err) {
		log.Infof("%v Certificate already exists", certName)
		return nil
	}
	if err != nil {
		log.Infof("%v Certificate generation failed, reason: %v", certName, err)
	} else {
		log.Infof("%v Certificate generation successful", certName)
	}

	return err
}
