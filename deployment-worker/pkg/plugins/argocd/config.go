package argocd

type Configuration struct {
	ServiceURL string `envconfig:"ARGOCD_SERVICE_URL"`
}
