package argocd

type Configuration struct {
	ServiceURL   string `envconfig:"ARGOCD_SERVICE_URL" default:"https://localhost:9080"`
	IsSSLEnabled bool   `envconfig:"IS_SSL_ENABLED" default:"false"`
	Password     string `envconfig:"PASSWORD" required:"true"`
}
