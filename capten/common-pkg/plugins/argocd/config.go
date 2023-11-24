package argocd

import "time"

type Configuration struct {
	ServiceURL   string `envconfig:"ARGOCD_SERVICE_URL" default:"argo-cd-argocd-server.argo-cd.svc.cluster.local"`
	IsSSLEnabled bool   `envconfig:"IS_SSL_ENABLED" default:"false"`
	Username     string `envconfig:"USERNAME" default:"admin"`
	Password     string `envconfig:"ARGOCD_PASSWORD" required:"true"`
}

type ConnectionState struct {
	AttemptedAt time.Time `json:"AttemptedAt"`
	Message     string    `json:"Message"`
	Status      string    `json:"Status"`
}

type Repository struct {
	Project               string          `json:"Project"`
	Repo                  string          `json:"Repo"`
	Username              string          `json:"Username"`
	Password              string          `json:"Password"`
	Type                  string          `json:"Type"`
	Insecure              bool            `json:"Insecure"`
	EnableLFS             bool            `json:"EnableLFS"`
	InsecureIgnoreHostKey bool            `json:"InsecureIgnoreHostKey"`
	ConnectionState       ConnectionState `json:"ConnectionState"`
}

type TLSClientConfig struct {
	// Insecure specifies that the server should be accessed without verifying the TLS certificate. For testing only.
	Insecure bool `json:"insecure" `
	// ServerName is passed to the server for SNI and is used in the client to check server
	// certificates against. If ServerName is empty, the hostname used to contact the
	// server is used.
	ServerName string `json:"serverName,omitempty" `
	// CertData holds PEM-encoded bytes (typically read from a client certificate file).
	// CertData takes precedence over CertFile
	CertData []byte `json:"certData,omitempty" `
	// KeyData holds PEM-encoded bytes (typically read from a client certificate key file).
	// KeyData takes precedence over KeyFile
	KeyData []byte `json:"keyData,omitempty" `
	// CAData holds PEM-encoded bytes (typically read from a root certificates bundle).
	// CAData takes precedence over CAFile
	CAData []byte `json:"caData,omitempty" `
}

type ClusterConfig struct {
	BearerToken string `json:"bearerToken,omitempty"`
	// TLSClientConfig contains settings to enable transport layer security
	TLSClientConfig `json:"tlsClientConfig"`
}

type Cluster struct {
	Server string        `json:"server"`
	Name   string        `json:"name"`
	Config ClusterConfig `json:"config"`
}
