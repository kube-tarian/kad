package crossplane

type appConfig struct {
	MainAppGitPath string   `json:"mainAppGitPath"`
	ChildAppNames  []string `json:"childAppNames"`
	SynchApp       bool     `json:"synchApp"`
}

type clusterUpdateConfig struct {
	MainAppGitPath              string                  `json:"mainAppGitPath"`
	ClusterValuesFile           string                  `json:"clusterValuesFile"`
	DefaultAppListFile          string                  `json:"defaultAppListFile"`
	DefaultAppValuesPath        string                  `json:"defaultAppValuesPath"`
	ClusterDefaultAppValuesPath string                  `json:"clusterDefaultAppValuesPath"`
	AppRoleTokenVaultPaths      []string                `json:"appRoleTokenVaultPaths"`
	ExternalSecrets             []clusterExternalSecret `json:"externalSecrets"`
}

type providerUpdateConfig struct {
	SyncAppPath string `json:"syncAppPath"`
}

type CrossplanePluginConfig struct {
	TemplateGitRepo          string               `json:"templateGitRepo"`
	CrossplaneConfigSyncPath string               `json:"crossplaneConfigSyncPath"`
	ProviderConfigSyncPath   string               `json:"providerConfigSyncPath"`
	ProviderPackages         map[string]string    `json:"providerPackages"`
	ArgoCDApps               []appConfig          `json:"argoCDApps"`
	ClusterEndpointUpdates   clusterUpdateConfig  `json:"clusterUpdateConfig"`
	ProviderEndpointUpdates  providerUpdateConfig `json:"providerUpdate"`
}

type secretData struct {
	SecretKey  string `json:"secretKey"`
	SecretPath string `json:"secretPath"`
}

type clusterExternalSecret struct {
	Namespace    string       `json:"namespace"`
	SecretName   string       `json:"secretName"`
	SecretType   string       `json:"secretType"`
	VaultSecrets []secretData `json:"vaultSecrets"`
}

const (
	crossplaneAWSProviderTemplate = `
apiVersion: pkg.crossplane.io/v1alpha1
kind: ControllerConfig
metadata:
  name: "%s-vault-config"
spec:
  args:
    - --debug
  metadata:
    annotations:
      vault.hashicorp.com/agent-inject: "true"
      vault.hashicorp.com/role: "vault-role-crossplane"
      vault.hashicorp.com/agent-inject-secret-creds.txt: "secret/%s"
      vault.hashicorp.com/agent-inject-template-creds.txt: |
        {{- with secret "secret/%s" -}}
          [default]
          aws_access_key_id="{{ .Data.data.accessKey }}"
          aws_secret_access_key="{{ .Data.data.secretKey }}"
        {{- end -}}
---
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-%s
spec:
  package: "%s"
  controllerConfigRef:
    name: "%s-vault-config"
`
)

const (
	crossplaneGCPProviderTemplate = `
apiVersion: pkg.crossplane.io/v1alpha1
kind: ControllerConfig
metadata:
  name: "%s-vault-config"
spec:
  metadata:
    annotations:
      vault.hashicorp.com/agent-inject: "true"
      vault.hashicorp.com/role: "vault-role-crossplane"
      vault.hashicorp.com/agent-inject-secret-creds.txt: "secret/%s"
      vault.hashicorp.com/agent-inject-template-creds.txt: |
        {{- with secret "secret/%s" -}}
         {{ .Data.data | toJSON }}
        {{- end -}}
---
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-%s
spec:
  package: "%s"
  controllerConfigRef:
    name: "%s-vault-config"`
)

const (
	crossplaneAzureProviderTemplate = `
  apiVersion: pkg.crossplane.io/v1
  kind: Provider
  metadata:
    name: provider-%s
  spec:
    package: "%s"
`
)
