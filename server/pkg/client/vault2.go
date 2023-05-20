package client

import (
	"fmt"
	"log"
	"time"

	"context"

	"github.com/hashicorp/vault/api"

	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func (v *Vault) Config() *kubernetes.Clientset {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		v.log.Errorf("Error while creating config %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		v.log.Errorf("Error while creating clientset %v", err)
	}
	return clientset
}

func (v *Vault) GenerateUnsealKeysFromVault() ([]string, string, error) {

	res := &api.InitRequest{
		SecretThreshold: 2,
		SecretShares:    3,
	}
	unsealKeys := []string{}

	key, err := v.client.Sys().Init(res)
	if err != nil {
		v.log.Error("Error while initializing ", err)
	}
	for _, key := range key.Keys {

		unsealKeys = append(unsealKeys, key)
	}

	rootToken := key.RootToken
	v.log.Info("Unsealed keys are generated")
	return unsealKeys, rootToken, err
}
func (v *Vault) Unseal(keys []string) error {

	for _, key := range keys {

		_, err := v.client.Sys().Unseal(key)
		if err != nil {

			v.log.Error("Error while unsealing", err)
			return err
		}

	}

	return nil
}
func (v *Vault) Storekeys(nameSpace string, SecretName string) []string {
	clientset := v.Config()
	var values []string
	namespace := nameSpace   // Namespace where you want to create the Secret
	secretName := SecretName // Name of the Secret
	//generate unseal keys
	unsealKeys, _, _ := v.GenerateUnsealKeysFromVault()

	stringData := make(map[string]string)
	for i, value := range unsealKeys {
		key := fmt.Sprintf("key%d", i+1)
		stringData[key] = value
		values = append(values, value)
	}
	newSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		StringData: stringData,
	}
	createdSecret, err := clientset.CoreV1().Secrets(namespace).Create(context.TODO(), newSecret, metav1.CreateOptions{})
	if err != nil {
		v.log.Error("Failed to create secret: %v\n", err)
	}
	v.log.Info("Secret '%s' created in namespace '%s'\n", createdSecret.Name, createdSecret.Namespace)
	return values
}

func (v *Vault) RetrieveKeys(nameSpace string, SecretName string) []string {

	var values []string
	clientset := v.Config()
	namespace := nameSpace // Namespace where you want to create the Secret
	secretName := SecretName

	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		v.log.Error("Error while getting secret", err)
	}

	for key, value := range secret.Data {
		// Use the secret value as needed
		v.log.Infof("Retrieved value for key %s: %s\n", key, value)
		keys := string(value)
		values = append(values, keys)
	}

	v.log.Infof("Secret '%s' found in namespace '%s'\n", secret.Name, secret.Namespace)
	// Use the secret as needed
	return values
}

func StartMonitoringService() {
	vault, err := NewVault()
	if err != nil {
		log.Fatal("Error while connecting to vault", err)
	}

	status, err := vault.client.Sys().SealStatus()
	if err != nil {
		log.Fatalf("Error while checking seal status %v", err)
	}
	for {
		if !status.Initialized {
			keys := vault.Storekeys("default", "vault-server")
			err := vault.Unseal(keys)
			if err != nil {
				log.Fatal("Error while unsealing the keys", err)
			}

		}

		if status.Sealed {
			keys := vault.RetrieveKeys("default", "vault-secret")
			vault.Unseal(keys)
		}
		// Sleep for a given interval before the next check
		time.Sleep(1 * time.Minute)
	}

}

// func main() {
// 	go startMonitoringService()

// 	// Sleep indefinitely to keep the service running
// 	select {}
// }
