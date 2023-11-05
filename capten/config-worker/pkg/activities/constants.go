package activities

const (
	branchName   = "capten-template-bot"
	gitUrlSuffix = ".git"
	accessToken  = "accessToken"
)

const crossplaneProviderTemplate = `
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
      vault.hashicorp.com/agent-inject-secret-creds.txt: "%s"
      vault.hashicorp.com/agent-inject-template-creds.txt: |
        {{- with secret "%s" -}}
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
