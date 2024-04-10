package captensdk

import (
	"context"
	"fmt"

	v1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	cmclient "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type MTLSClient struct {
	log logging.Logger
}

func NewMTLSClient(log logging.Logger) (*MTLSClient, error) {
	return &MTLSClient{log: log}, nil
}

func (m *MTLSClient) CreateCertificate(name, namespace, issuerRefName string) error {
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
		err = m.generateCertificate(cmClient, namespace, name, issuerRefName, isClient)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MTLSClient) generateCertificate(cmClient *cmclient.Clientset, namespace string, name string, issuerRefName string, isClient bool) error {
	var usages []v1.KeyUsage
	certName := name
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
				CommonName: name,
				Usages:     usages,
			},
		},
		metav1.CreateOptions{},
	)
	if err != nil {
		m.log.Infof("%v Certificate generation failed, reason: %v", certName, err)
	} else {
		m.log.Infof("%v Certificate generation successful", certName)
	}

	return err
}

func (m *MTLSClient) DeleteCertificate(name, namespace string) error {
	kubeconfigPath := m.PrepareFilePath("captenConfig.ConfigDirPath", "captenConfig.KubeConfigFileName")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return errors.WithMessage(err, "error while building kubeconfig")
	}
	cmClient, err := cmclient.NewForConfig(config)
	if err != nil {
		return err
	}

	// Create cert-manager certificate request
	err = cmClient.CertmanagerV1().Certificates(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	return err
}

func (c *MTLSClient) PrepareFilePath(dir, path string) string {
	return fmt.Sprintf("%s%s%s", "c.CurrentDirPath", dir, path)
}

func (c *MTLSClient) PrepareDirPath(dir string) string {
	return fmt.Sprintf("%s%s", "c.CurrentDirPath", dir)
}
