package captensdk

import (
	"context"
	"fmt"

	v1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	cmclient "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"github.com/kube-tarian/kad/capten/deployment-worker/internal/k8sops"
	"github.com/pkg/errors"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type MTLSClient struct {
	log logging.Logger
}

func NewMTLSClient(log logging.Logger) (*MTLSClient, error) {
	return &MTLSClient{log: log}, nil
}

func (m *MTLSClient) CreateCertificates(certName, namespace, issuerRefName, cmName string, k8sClient *k8s.K8SClient) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return errors.WithMessage(err, "error while building kubeconfig")
	}
	cmClient, err := cmclient.NewForConfig(config)
	if err != nil {
		return err
	}

	// Create cert-manager certificate for client and server
	for _, isClient := range []bool{true, false} {
		err = m.generateCertificate(cmClient, namespace, certName, issuerRefName, isClient)
		if err != nil {
			return err
		}
	}

	data := map[string]string{}
	data["client-certificate"] = certName + "-client-mtls-capten-sdk"
	data["server-certificate"] = certName + "-server-mtls-capten-sdk"
	k8sops.CreateUpdateConfigmap(context.TODO(), m.log, namespace, cmName, data, k8sClient)
	return nil
}

func (m *MTLSClient) generateCertificate(cmClient *cmclient.Clientset, namespace string, certName string, issuerRefName string, isClient bool) error {
	var usages []v1.KeyUsage
	if isClient {
		usages = []v1.KeyUsage{v1.UsageDigitalSignature, v1.UsageKeyEncipherment, v1.UsageClientAuth}
		certName = certName + "-client-mtls-capten-sdk"
	} else {
		usages = []v1.KeyUsage{v1.UsageDigitalSignature, v1.UsageKeyEncipherment, v1.UsageServerAuth}
		certName = certName + "-server-mtls-capten-sdk"
	}

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
		m.log.Infof("%v Certificate already exists", certName)
		return nil
	}
	if err != nil {
		m.log.Infof("%v Certificate generation failed, reason: %v", certName, err)
	} else {
		m.log.Infof("%v Certificate generation successful", certName)
	}

	return err
}

func (m *MTLSClient) DeleteCertificate(name, namespace string) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return errors.WithMessage(err, "error while building kubeconfig")
	}
	cmClient, err := cmclient.NewForConfig(config)
	if err != nil {
		return err
	}

	// Create cert-manager certificate request
	err = cmClient.CertmanagerV1().Certificates(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if k8serror.IsNotFound(err) {
		m.log.Infof("%v Certificate not found", name)
		return nil
	}
	return err
}

func (c *MTLSClient) PrepareFilePath(dir, path string) string {
	return fmt.Sprintf("%s%s%s", "c.CurrentDirPath", dir, path)
}

func (c *MTLSClient) PrepareDirPath(dir string) string {
	return fmt.Sprintf("%s%s", "c.CurrentDirPath", dir)
}
