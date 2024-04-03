package k8s

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v2"
)

type ObjectMeta struct {
	Name      string `yaml:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	Namespace string `yaml:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`
}

type SecretStoreRef struct {
	Name string `yaml:"name"`
	Kind string `yaml:"kind,omitempty"`
}

type ExternalSecretTargetTemplate struct {
	Type string `yaml:"type,omitempty"`
}

type ExternalSecretTarget struct {
	Name     string                       `yaml:"name,omitempty"`
	Template ExternalSecretTargetTemplate `yaml:"template,omitempty"`
}

type ExternalSecretData struct {
	SecretKey string                      `yaml:"secretKey"`
	RemoteRef ExternalSecretDataRemoteRef `yaml:"remoteRef"`
}

type ExternalSecretDataRemoteRef struct {
	Key      string `yaml:"key"`
	Property string `yaml:"property,omitempty"`
}

type ExternalSecretSpec struct {
	SecretStoreRef  SecretStoreRef       `yaml:"secretStoreRef,omitempty"`
	Target          ExternalSecretTarget `yaml:"target,omitempty"`
	RefreshInterval string               `yaml:"refreshInterval,omitempty"`
	Data            []ExternalSecretData `yaml:"data,omitempty"`
}

type ExternalSecret struct {
	Kind       string     `yaml:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`
	APIVersion string     `yaml:"apiVersion,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`
	Metadata   ObjectMeta `yaml:"metadata,omitempty"`

	Spec ExternalSecretSpec `yaml:"spec,omitempty"`
}

type SecretStoreSpec struct {
	Provider        *SecretStoreProvider `yaml:"provider"`
	RefreshInterval int                  `yaml:"refreshInterval,omitempty"`
}

type SecretKeySelector struct {
	Name string `yaml:"name,omitempty"`
	Key  string `yaml:"key,omitempty"`
}

type VaultAuth struct {
	TokenSecretRef *SecretKeySelector `yaml:"tokenSecretRef,omitempty"`
}

type VaultProvider struct {
	Auth    VaultAuth `yaml:"auth"`
	Server  string    `yaml:"server"`
	Path    string    `yaml:"path,omitempty"`
	Version string    `yaml:"version"`
}

type SecretStoreProvider struct {
	Vault *VaultProvider `yaml:"vault,omitempty"`
}

type SecretStore struct {
	Kind       string     `yaml:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`
	APIVersion string     `yaml:"apiVersion,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`
	Metadata   ObjectMeta `yaml:"metadata,omitempty"`

	Spec SecretStoreSpec `yaml:"spec,omitempty"`
}

func (k *K8SClient) CreateOrUpdateSecretStore(ctx context.Context, secretStoreName, namespace, vaultAddr,
	tokenSecretName, tokenSecretKey string) (err error) {
	secretStore := SecretStore{
		APIVersion: "external-secrets.io/v1beta1",
		Kind:       "SecretStore",
		Metadata: ObjectMeta{
			Name:      secretStoreName,
			Namespace: namespace,
		},
		Spec: SecretStoreSpec{
			RefreshInterval: 10,
			Provider: &SecretStoreProvider{
				Vault: &VaultProvider{
					Server:  vaultAddr,
					Path:    "secret",
					Version: "v2",
					Auth: VaultAuth{
						TokenSecretRef: &SecretKeySelector{
							Key:  tokenSecretKey,
							Name: tokenSecretName,
						},
					},
				},
			},
		},
	}

	secretStoreData, err := yaml.Marshal(&secretStore)

	if err != nil {
		return
	}
	_, _, err = k.DynamicClient.CreateResource(ctx, []byte(secretStoreData))
	if err != nil {
		err = fmt.Errorf("failed to create cluter vault token secret %s/%s, %v", namespace, secretStoreName, err)
		return
	}
	return
}

func (k *K8SClient) CreateOrUpdateExternalSecret(ctx context.Context, externalSecretName, namespace,
	secretStoreRefName, secretName, secretType string, vaultKeyPathdata map[string]string) (err error) {
	secretKeysData := []ExternalSecretData{}
	for key, path := range vaultKeyPathdata {
		secretKeyData := ExternalSecretData{
			SecretKey: key,
			RemoteRef: ExternalSecretDataRemoteRef{
				Key:      path,
				Property: key,
			},
		}

		secretKeysData = append(secretKeysData, secretKeyData)
	}
	externalSecret := ExternalSecret{
		APIVersion: "external-secrets.io/v1beta1",
		Kind:       "ExternalSecret",
		Metadata: ObjectMeta{
			Name:      externalSecretName,
			Namespace: namespace,
		},
		Spec: ExternalSecretSpec{
			RefreshInterval: "10s",
			Target: ExternalSecretTarget{
				Name:     secretName,
				Template: ExternalSecretTargetTemplate{Type: secretType}},
			SecretStoreRef: SecretStoreRef{
				Name: secretStoreRefName,
				Kind: "SecretStore",
			},
			Data: secretKeysData,
		},
	}

	externalSecretData, err := yaml.Marshal(&externalSecret)
	if err != nil {
		return
	}

	_, _, err = k.DynamicClient.CreateResource(ctx, []byte(externalSecretData))
	if err != nil {
		err = fmt.Errorf("failed to create vault external secret %s/%s, %v", namespace, externalSecretName, err)
		return
	}
	return
}
