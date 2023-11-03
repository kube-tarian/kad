package activities

const (
	branchName   = "capten-template-bot"
	gitUrlSuffix = ".git"
)

// plugin structure
const (
	Tekton        = "tekton"
	CrossPlane    = "crossplane"
	GitRepo       = "git_repo"
	GitConfigPath = "git_config_path"
	ConfigMainApp = "configure_main_app"
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
      vault.hashicorp.com/role: "crossplane-providers"
      vault.hashicorp.com/agent-inject-secret-creds.txt: "%s"
      vault.hashicorp.com/agent-inject-template-creds.txt: |
        {{- with secret "%s" -}}
          [default]
          aws_access_key_id="{{ .access_key }}"
          aws_secret_access_key="{{ .secret_key }}"
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
