package client

import (
	"context"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"github.com/kube-tarian/kad/vaultserv/pkg/pb/vaultservpb"
	"log"
	"os"
)

type Vault struct {
	client *vault.Client
}

func NewVault() (*Vault, error) {
	config := vault.DefaultConfig()
	config.Address = os.Getenv("VAULT_ADDR")
	//tlsConfig := vault.TLSConfig{CACert: os.Getenv("VAULT_CACERT")}
	tlsConfig := vault.TLSConfig{Insecure: true}
	if err := config.ConfigureTLS(&tlsConfig); err != nil {
		log.Fatalf("unable to configure tls %v", err)
	}

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	token, err := getFileContent(os.Getenv("VAULT_TOKEN_FILE"))
	if err != nil {
		return nil, fmt.Errorf("failed to get token %w", err)
	}

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

func (v *Vault) PutCredential(secretName string, credDetails *vaultservpb.StoreCredRequest) error {
	secretData := map[string]interface{}{
		"username": credDetails.Username,
		"password": credDetails.Password,
	}

	ctx := context.Background()

	// Write a secret
	_, err := v.client.KVv2(secretName).Put(ctx, credDetails.Credname, secretData)
	if err != nil {
		return fmt.Errorf("unable to write secret: %w", err)
	}

	log.Println("credential written successfully.")
	return nil
}
