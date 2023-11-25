package crossplane

type appConfig struct {
	MainAppGitPath string   `json:"mainAppGitPath"`
	ChildAppNames  []string `json:"childAppNames"`
	SynchApp       bool     `json:"synchApp"`
}

type clusterConfig struct {
	MainAppGitPath string `json:"mainAppGitPath"`
	File           string `json:"file"`
}

type crossplanePluginConfig struct {
	TemplateGitRepo          string            `json:"templateGitRepo"`
	CrossplaneConfigSyncPath string            `json:"crossplaneConfigSyncPath"`
	ProviderConfigSyncPath   string            `json:"providerConfigSyncPath"`
	ProviderPackages         map[string]string `json:"providerPackages"`
	ArgoCDApps               []appConfig       `json:"argoCDApps"`
	ClusterEndpointUpdates   clusterConfig     `json:"clusterEndpointUpdates"`
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
        GOOGLE_CLOUD_KEYFILE_JSON="{{ .Data.data.keyfileJSON | toString }}"
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
	crossplaneAzureProviderTemplate = `
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
        AZURE_SUBSCRIPTION_ID="{{ .Data.data.subscriptionID }}"
        AZURE_TENANT_ID="{{ .Data.data.tenantID }}"
        AZURE_CLIENT_ID="{{ .Data.data.clientID }}"
        AZURE_CLIENT_SECRET="{{ .Data.data.clientSecret }}"
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
