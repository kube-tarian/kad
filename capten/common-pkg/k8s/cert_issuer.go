package k8s

import (
	"context"
	"fmt"

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

func CreateOrUpdateClusterIssuer(clusterCAIssuer string, k8sclient *K8SClient, forceUpdate bool) (*cert.CertificatesData, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, errors.WithMessage(err, "error while building kubeconfig")
	}

	cmClient, err := cmclient.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	certsData, err := cert.GenerateRootCerts()
	if err != nil {
		return nil, err
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
	if err != nil {
		if k8serrors.IsNotFound(err) {
			err = CreateOrUpdateClusterCAIssuerSecret(k8sclient, certsData.RootCert.CertData, certsData.RootKey.KeyData, certsData.CaChainCertData)
			if err != nil {
				return nil, fmt.Errorf("failed to create/update CA Issuer Secret: %v", err)
			}

			result, err := cmClient.CertmanagerV1().ClusterIssuers().Create(context.Background(), issuer, metav1.CreateOptions{})
			if err != nil {
				return nil, errors.WithMessage(err, "error in creating cert issuer")
			}
			log.Debugf("ClusterIssuer %s created successfully", result.Name)
			return certsData, nil
		} else if k8serrors.IsAlreadyExists(err) {
			secret, err := k8sclient.GetSecretObject(cert.CertManagerNamespace, cert.ClusterCACertSecretName)
			if err != nil {
				log.Errorf("Failed to read secert %s, %v", cert.ClusterCACertSecretName, err)
				return nil, err
			}
			certsData.CaChainCertData = secret.Data["ca.crt"]
			certsData.RootCert.CertData = secret.Data[corev1.TLSCertKey]
			certsData.RootKey.KeyData = secret.Data[corev1.TLSPrivateKeyKey]
			return certsData, nil
		}
		return nil, err
	}

	if forceUpdate {
		err = CreateOrUpdateClusterCAIssuerSecret(k8sclient, certsData.RootCert.CertData, certsData.RootKey.KeyData, certsData.CaChainCertData)
		if err != nil {
			return nil, fmt.Errorf("failed to create/update CA Issuer Secret: %v", err)
		}

		serverIssuer.Spec.IssuerConfig.CA.SecretName = cert.ClusterCACertSecretName
		issuerClient := cmClient.CertmanagerV1().ClusterIssuers()
		_, err := issuerClient.Update(context.TODO(), serverIssuer, metav1.UpdateOptions{})
		if err != nil {
			return nil, errors.WithMessage(err, "error while updating cluster issuer")
		}
	}
	log.Debugf("ClusterIssuer %s updated successfully", issuer.Name)
	return certsData, nil
}

func CreateOrUpdateClusterCAIssuerSecret(k8sClient *K8SClient, caCertData, caKeyData, caCertChainData []byte) error {
	// Create the Secret object
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cert.ClusterCACertSecretName,
			Namespace: cert.CertManagerNamespace,
		},
		Data: map[string][]byte{
			corev1.TLSCertKey:       caCertData,
			corev1.TLSPrivateKeyKey: caKeyData,
			"ca.crt":                caCertChainData,
		},
		Type: corev1.SecretTypeTLS,
	}
	return k8sClient.CreateOrUpdateSecretObject(context.TODO(), secret)
}
