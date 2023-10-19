package argocd

type Configuration struct {
	ServiceURL   string `envconfig:"ARGOCD_SERVICE_URL" default:"https://argocd.demo.optimizor.app"`
	IsSSLEnabled bool   `envconfig:"IS_SSL_ENABLED" default:"false"`
	Username     string `envconfig:"USERNAME" default:"amdin"`
	Password     string `envconfig:"ARGOCD_PASSWORD" required:"true"`
}
