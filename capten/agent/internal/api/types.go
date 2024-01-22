package api

const (
	vaultStore = `
apiVersion: external-secrets.io/v1beta1
kind: SecretStore
metadata:
  name: vault-store
  namespace: %s
spec:
  provider:
    vault:
      server: "%s"
      path: "secret"
      version: "v2"
      auth:
        tokenSecretRef:
          name: "%s"
          key: "token"
          `
	natsVaultExternalSecret = `
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: vault-nats-external
  namespace: %s
spec:
  refreshInterval: "10s"
  secretStoreRef:
    name: vault-store
    kind: SecretStore
  target:
    name: vault-nats-secret
  data:
  - secretKey: credentials
    remoteRef:
      key: %s
      property: nats
  `
)
