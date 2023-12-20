package argocd

import (
	"time"
)

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
	Upsert                bool            `json:"Upsert"`
}
