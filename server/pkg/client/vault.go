package client

import (
	"context"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"github.com/kube-tarian/kad/server/pkg/types"
	"log"
	"os"
)

//var (
//	vaultObj      *Vault
//	vaultSyncOnce sync.Once
//)

type Vault struct {
	client *vault.Client
}

func NewVault() (*Vault, error) {
	config := vault.DefaultConfig()
	config.Address = os.Getenv("VAULT_ADDR")
	tlsConfig := vault.TLSConfig{CACert: os.Getenv("VAULT_CACERT")}
	if err := config.ConfigureTLS(&tlsConfig); err != nil {
		log.Fatalf("unable to configure tls %v", err)
	}

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	token, err := getFileContent(os.Getenv("VAULT_TOKEN"))
	if err != nil {
		return nil, fmt.Errorf("failed to get token %w", err)
	}

	fmt.Println("token is ", token)
	client.SetToken(token)
	return &Vault{
		client: client,
	}, nil
}

func getFileContent(fileName string) (string, error) {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to read file %w", err)
	}

	return string(fileContent), nil
}

func (v *Vault) PutCert(secretName, certChain, clientCert, clientKey, customerID string) error {
	secretData := map[string]interface{}{
		types.ClientCertChainFileName: certChain,
		types.ClientKeyFileName:       clientKey,
		types.ClientCertFileName:      clientCert,
	}

	ctx := context.Background()

	// Write a secret
	_, err := v.client.KVv2(secretName).Put(ctx, fmt.Sprintf("cert-%s", customerID), secretData)
	if err != nil {
		return fmt.Errorf("unable to write secret: %w", err)
	}

	log.Println("Secret written successfully.")
	return nil
}

func (v *Vault) GetCert(secretName, customerID string) (map[string]string, error) {
	secret, err := v.client.KVv2(secretName).Get(context.Background(), fmt.Sprintf("cert-%s", customerID))
	if err != nil {
		return nil, fmt.Errorf("unable to read secret: %w", err)
	}

	certMap := make(map[string]string)
	certMap[types.ClientCertChainFileName] = secret.Data[types.ClientCertChainFileName].(string)
	certMap[types.ClientCertFileName] = secret.Data[types.ClientCertFileName].(string)
	certMap[types.ClientKeyFileName] = secret.Data[types.ClientKeyFileName].(string)

	return certMap, nil
}
