package argocd

import "time"

type Configuration struct {
	ServiceURL   string `envconfig:"ARGOCD_SERVICE_URL" default:"argo-cd-argocd-server.argo-cd.svc.cluster.local"`
	IsSSLEnabled bool   `envconfig:"IS_SSL_ENABLED" default:"false"`
	Username     string `envconfig:"USERNAME" default:"amdin"`
	Password     string `envconfig:"ARGOCD_PASSWORD" required:"true"`
}

type ConnectionState struct {
	AttemptedAt time.Time `json:"attemptedAt"`
	Message     string    `json:"message"`
	Status      string    `json:"status"`
}

type Repository struct {
	Project         string          `json:"project"`
	Repo            string          `json:"repo"`
	SSHPrivateKey   string          `json:"sshPrivateKey"`
	Type            string          `json:"type"`
	ConnectionState ConnectionState `json:"connectionState"`
}
