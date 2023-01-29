package fetcher

type PluginRequest struct {
	PluginName string
}

type PluginResponse struct {
	ServiceURL   string
	IsSSLEnabled bool
	Username     string
	Password     string
}

type PluginDetails struct {
	Name        string
	RepoName    string
	RepoURL     string
	ChartName   string
	Namespace   string
	ReleaseName string
	Version     string
}
