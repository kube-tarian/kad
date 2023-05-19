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

// type Vault struct {
// 	client *api.Client
// }

func (v Vault) GenerateUnsealKeysFromVault() ([]string, string, error) {

	res := &api.InitRequest{
		SecretThreshold: 2,
		SecretShares:    3,
	}
	unsealKeys := []string{}

	key, err := v.client.Sys().Init(res)
	if err != nil {
		fmt.Println("Error while initializing ", err)
	}
	for _, key := range key.Keys {
		fmt.Println("Key is ", key)

		unsealKeys = append(unsealKeys, key)
	}

	rootToken := key.RootToken
	fmt.Println("Root Token is ", rootToken)
	fmt.Print("Unsealed keys are generated")
	return unsealKeys, rootToken, err
}
func (v Vault) Unseal(keys []string) error {
	flag := true
	for _, key := range keys {

		_, err := v.client.Sys().Unseal(key)
		if err != nil {
			flag = false
			fmt.Println("Error while unsealing", err)
			return err
		}

	}
	if flag {
		fmt.Println("Unsealed")
	}
	return nil

}
func (v Vault) Storekeys() []string {

	var values []string

	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Error while creating clientset", err)
	}

	namespace := "default"       // Namespace where you want to create the Secret
	secretName := "vault-secret" // Name of the Secret

	// Retrieve the secret with the given name
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		// Secret not found, create a new one
		if err.Error() == "secrets \""+secretName+"\" not found" {
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
				fmt.Printf("Failed to create secret: %v\n", err)

			}

			fmt.Printf("Secret '%s' created in namespace '%s'\n", createdSecret.Name, createdSecret.Namespace)
		} else {
			// Error occurred while retrieving secret
			fmt.Printf("Failed to retrieve secret: %v\n", err)

		}
	} else {
		// Secret found
		secret2, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
		if err != nil {
			fmt.Println("Error while getting secret", err)
		}

		for key, value := range secret2.Data {

			// Use the secret value as needed
			fmt.Printf("Retrieved value for key %s: %s\n", key, value)
			keys := string(value)
			values = append(values, keys)
			fmt.Println("Key is ", keys)
		}

		fmt.Printf("Secret '%s' found in namespace '%s'\n", secret.Name, secret.Namespace)
		// Use the secret as needed
	}
	return values
}
func startMonitoringService() {

	client, err := api.NewClient(&api.Config{
		Address: "http://127.0.0.1:8200",
	})

	if err != nil {
		fmt.Println("Error while connecting to client", err)
	}

	vault := Vault{
		client: client,
	}

	status, err := vault.client.Sys().SealStatus()
	if err != nil {
		fmt.Println("Error while sealing", err)
	}
	for {
		if !status.Initialized {
			fmt.Println("Vault server is not initialized")
			//	keys,_,_:=vault.GenerateUnsealKeysFromVault()
			keys := vault.Storekeys()
			err := vault.Unseal(keys)
			if err != nil {
				fmt.Println("Error is", err)
			}

			fmt.Printf("Key is Stored Successfully")

		}

		if status.Sealed {
			log.Printf("Vault server is sealed")
			keys := vault.Storekeys()
			vault.Unseal(keys)
		}
		// Sleep for a given interval before the next check
		time.Sleep(1 * time.Minute)
	}

}

//func main() {
// 	go startMonitoringService()

// 	// Sleep indefinitely to keep the service running
// 	select {}
// }
