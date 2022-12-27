package argocd

type Configuration struct {
	ServiceURL   string `envconfig:"ARGOCD_SERVICE_URL" default:"localhost:9081"`
	IsSSLEnabled bool   `envconfig:"IS_SSL_ENABLED" default:"false"`
	Username     string `envconfig:"USERNAME" default:"amdin"`
	Password     string `envconfig:"PASSWORD" required:"true"`
}
