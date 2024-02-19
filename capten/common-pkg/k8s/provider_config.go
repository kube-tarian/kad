package k8s

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v2"
)

type ProviderSpec struct {
	Package             string `yaml:"package,omitempty"`
	ControllerConfigRef struct {
		Name string `yaml:"name,omitempty"`
	} `yaml:"controllerConfigRef,omitempty"`
}

type ProviderConfig struct {
	APIVersion string `yaml:"apiVersion,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`
	Kind       string `yaml:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`
	Metadata   struct {
		Name string `yaml:"name,omitempty" protobuf:"bytes,2,opt,name=name"`
	} `yaml:"metadata,omitempty"`
	Spec ProviderSpec `yaml:"spec,omitempty"`
}

type ControllerConfigSpec struct {
	Args     []string `yaml:"args,omitempty"`
	Metadata struct {
		Annotations map[string]string `yaml:"annotations,omitempty"`
	} `yaml:"metadata,omitempty"`
}

type ControllerConfig struct {
	APIVersion string `yaml:"apiVersion,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`
	Kind       string `yaml:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`
	Metadata   struct {
		Name string `yaml:"name,omitempty" protobuf:"bytes,2,opt,name=name"`
	} `yaml:"metadata,omitempty"`
	Spec ControllerConfigSpec `yaml:"spec,omitempty"`
}

func (k *K8SClient) CreateAWSProviderConfig(ctx context.Context, cloudType, secretPath, pkg string) (string, error) {

	controlerConfig := ControllerConfig{
		APIVersion: "pkg.crossplane.io/v1alpha1",
		Kind:       "ControllerConfig",
	}
	controlerConfig.Metadata.Name = fmt.Sprintf("%s-vault-config", cloudType)
	controlerConfig.Spec.Args = []string{"--debug"}
	controlerConfig.Spec.Metadata.Annotations = map[string]string{
		"vault.hashicorp.com/agent-inject":                  "true",
		"vault.hashicorp.com/role":                          "vault-role-crossplane",
		"vault.hashicorp.com/agent-inject-secret-creds.txt": fmt.Sprintf("secret/%s", secretPath),
		"vault.hashicorp.com/agent-inject-template-creds.txt": fmt.Sprintf(`{{- with secret "secret/%s" -}}
        [default]
        aws_access_key_id="{{ .Data.data.accessKey }}"
        aws_secret_access_key="{{ .Data.data.secretKey }}"
      {{- end -}}`, secretPath),
	}

	providerConfig := ProviderConfig{
		APIVersion: "pkg.crossplane.io/v1",
		Kind:       "Provider",
	}
	providerConfig.Metadata.Name = fmt.Sprintf("provider-%s", cloudType)
	providerConfig.Spec.Package = pkg
	providerConfig.Spec.ControllerConfigRef.Name = fmt.Sprintf("%s-vault-config", cloudType)

	controlerConfigByte, err := yaml.Marshal(&controlerConfig)
	if err != nil {
		return "", err
	}

	providerConfigByte, err := yaml.Marshal(&providerConfig)
	if err != nil {
		return "", err
	}

	return string(controlerConfigByte) + "---\n" + string(providerConfigByte), nil
}

func (k *K8SClient) CreateGCPProviderConfig(ctx context.Context, cloudType, secretPath, pkg string) (string, error) {

	controlerConfig := ControllerConfig{
		APIVersion: "pkg.crossplane.io/v1alpha1",
		Kind:       "ControllerConfig",
	}
	controlerConfig.Metadata.Name = fmt.Sprintf("%s-vault-config", cloudType)
	controlerConfig.Spec.Args = []string{"--debug"}
	controlerConfig.Spec.Metadata.Annotations = map[string]string{
		"vault.hashicorp.com/agent-inject":                  "true",
		"vault.hashicorp.com/role":                          "vault-role-crossplane",
		"vault.hashicorp.com/agent-inject-secret-creds.txt": fmt.Sprintf("secret/%s", secretPath),
		"vault.hashicorp.com/agent-inject-template-creds.txt": fmt.Sprintf(`{{- with secret "secret/%s" -}}
        {{ .Data.data | toJSON }}
      {{- end -}}`, secretPath),
	}

	providerConfig := ProviderConfig{
		APIVersion: "pkg.crossplane.io/v1",
		Kind:       "Provider",
	}
	providerConfig.Metadata.Name = fmt.Sprintf("provider-%s", cloudType)
	providerConfig.Spec.Package = pkg
	providerConfig.Spec.ControllerConfigRef.Name = fmt.Sprintf("%s-vault-config", cloudType)

	controlerConfigByte, err := yaml.Marshal(&controlerConfig)
	if err != nil {
		return "", err
	}

	providerConfigByte, err := yaml.Marshal(&providerConfig)
	if err != nil {
		return "", err
	}

	return string(controlerConfigByte) + "---\n" + string(providerConfigByte), nil
}

func (k *K8SClient) CreateAzureProviderConfig(ctx context.Context, cloudType, secretPath, pkg string) (string, error) {

	controlerConfig := ControllerConfig{
		APIVersion: "pkg.crossplane.io/v1alpha1",
		Kind:       "ControllerConfig",
	}
	controlerConfig.Metadata.Name = fmt.Sprintf("%s-vault-config", cloudType)
	controlerConfig.Spec.Args = []string{"--debug"}
	controlerConfig.Spec.Metadata.Annotations = map[string]string{
		"vault.hashicorp.com/agent-inject":                  "true",
		"vault.hashicorp.com/role":                          "vault-role-crossplane",
		"vault.hashicorp.com/agent-inject-secret-creds.txt": fmt.Sprintf("secret/%s", secretPath),
		"vault.hashicorp.com/agent-inject-template-creds.txt": fmt.Sprintf(`{{- with secret "secret/%s" -}}
       [default]
       AZURE_SUBSCRIPTION_ID="{{ .Data.data.subscriptionID }}"
       AZURE_TENANT_ID="{{ .Data.data.tenantID }}"
       AZURE_CLIENT_ID="{{ .Data.data.clientID }}"
       AZURE_CLIENT_SECRET="{{ .Data.data.clientSecret }}"
      {{- end -}}`, secretPath),
	}

	providerConfig := ProviderConfig{
		APIVersion: "pkg.crossplane.io/v1",
		Kind:       "Provider",
	}
	providerConfig.Metadata.Name = fmt.Sprintf("provider-%s", cloudType)
	providerConfig.Spec.Package = pkg
	providerConfig.Spec.ControllerConfigRef.Name = fmt.Sprintf("%s-vault-config", cloudType)

	controlerConfigByte, err := yaml.Marshal(&controlerConfig)
	if err != nil {
		return "", err
	}

	providerConfigByte, err := yaml.Marshal(&providerConfig)
	if err != nil {
		return "", err
	}

	return string(controlerConfigByte) + "---\n" + string(providerConfigByte), nil
}
