package clusterissuer

import (
	"context"
	"fmt"
	"os"
	"time"

	v1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	cmclient "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/cert"
	"github.com/kube-tarian/kad/capten/common-pkg/credential"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"github.com/pkg/errors"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

const (
	CertFileName         = "server.cert"
	KeyFileName          = "server.key"
	namespace            = "capten"
	serverCertSecretName = "agent-server-mtls"
)

func SetupCACertIssuser(clusterIssuerName string, log logging.Logger) error {
	k8sclient, err := k8s.NewK8SClient(log)
	if err != nil {
		log.Errorf("failed to initalize k8s client, %v", err)
		return err
	}

	err = setupCertificateIssuer(k8sclient, clusterIssuerName, log)
	if err != nil {
		log.Errorf("Setup Certificates Issuer failed, %v", err)
		return err
	}
	return nil
}

// Setup agent certificate issuer
func setupCertificateIssuer(k8sClient *k8s.K8SClient, clusterIssuerName string, log logging.Logger) error {
	// Create Agent Cluster Issuer
	certsData, err := k8s.CreateOrUpdateClusterIssuer(clusterIssuerName, k8sClient, false)
	if err != nil {
		return fmt.Errorf("failed to create/update CA Issuer %s in cert-manager: %v", clusterIssuerName, err)
	}

	// Update Vault
	err = credential.PutClusterCerts(context.TODO(), "kad-agent", "kad-agent", string(certsData.CaChainCertData), string(certsData.RootKey.KeyData), string(certsData.RootCert.CertData))
	if err != nil {
		log.Errorf("Failed to write to vault, %v", err)
		log.Infof("Continued to start the agent as these certs from vault are not used...")
	}
	return nil
}

func GenerateServerCertificates(clusterIssuerName string, log logging.Logger) error {
	k8sClient, err := k8s.NewK8SClient(log)
	if err != nil {
		log.Errorf("failed to initalize k8s client, %v", err)
		return err
	}

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
		return fmt.Errorf("failed to create namespace %v, reason: %v", namespace, err)
	}
	log.Infof("Created namesapce: %v", namespace)

	err = generateCertManagerServerCertificate(cmClient, namespace, serverCertSecretName, clusterIssuerName, log)
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
	err = os.WriteFile(CertFileName, []byte(secretData.Data["tls.crt"]), cert.FilePermission)
	if err != nil {
		return fmt.Errorf("failed to create cert file, %v", err)
	}
	err = os.WriteFile(KeyFileName, []byte(secretData.Data["tls.key"]), cert.FilePermission)
	if err != nil {
		return fmt.Errorf("failed to create key file, %v", err)
	}

	return nil
}

func generateCertManagerServerCertificate(
	cmClient *cmclient.Clientset,
	namespace string,
	certName string,
	issuerRefName string,
	log logging.Logger,
) error {
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
