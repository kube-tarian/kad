package types

type App struct {
	Name            string                 `json:"Name"`
	ChartName       string                 `json:"ChartName"`
	RepoName        string                 `json:"RepoName"`
	RepoURL         string                 `json:"RepoURL"`
	Namespace       string                 `json:"Namespace"`
	ReleaseName     string                 `json:"ReleaseName"`
	Version         string                 `json:"Version"`
	CreateNamespace bool                   `json:"CreateNamespace"`
	Override        map[string]interface{} `json:"Override"`
}
