package k8s

import (
	"context"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/cert"
	"github.com/pkg/errors"

	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmclient "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

var log = logging.NewLogger()

func CreateOrUpdateClusterIssuer(clusterCAIssuer string) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return errors.WithMessage(err, "error while building kubeconfig")
	}

	cmClient, err := cmclient.NewForConfig(config)
	if err != nil {
		return err
	}

	issuer := &certmanagerv1.ClusterIssuer{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterCAIssuer,
		},
		Spec: certmanagerv1.IssuerSpec{
			IssuerConfig: certmanagerv1.IssuerConfig{
				CA: &certmanagerv1.CAIssuer{
					SecretName: cert.ClusterCACertSecretName,
				},
			},
		},
	}

	serverIssuer, err := cmClient.CertmanagerV1().ClusterIssuers().Get(context.Background(), issuer.Name, metav1.GetOptions{})
	if err != nil && k8serrors.IsNotFound(err) {
		result, err := cmClient.CertmanagerV1().ClusterIssuers().Create(context.Background(), issuer, metav1.CreateOptions{})
		if err != nil {
			return errors.WithMessage(err, "error in creating cert issuer")
		}
		log.Debugf("ClusterIssuer %s created successfully", result.Name)
		return nil
	}

	serverIssuer.Spec.IssuerConfig.CA.SecretName = cert.ClusterCACertSecretName
	issuerClient := cmClient.CertmanagerV1().ClusterIssuers()
	result, err := issuerClient.Update(context.TODO(), serverIssuer, metav1.UpdateOptions{})
	if err != nil {
		return errors.WithMessage(err, "error while updating cluster issuer")
	}
	log.Debugf("ClusterIssuer %s updated successfully", result.Name)
	return nil
}

func CreateOrUpdateClusterCAIssuerSecret(k8sClient *K8SClient, interCACertData, interCAKeyData, caCertChainData []byte) error {
	// Create the Secret object
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cert.ClusterCACertSecretName,
			Namespace: cert.CertManagerNamespace,
		},
		Data: map[string][]byte{
			corev1.TLSCertKey:       interCACertData,
			corev1.TLSPrivateKeyKey: interCAKeyData,
			"ca.crt":                caCertChainData,
		},
		Type: corev1.SecretTypeTLS,
	}
	return k8sClient.CreateOrUpdateSecretObject(context.TODO(), secret)
}
