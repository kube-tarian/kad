package tekton

type Source struct {
	RepoURL        string `json:"repoURL,omitempty"`
	TargetRevision string `json:"targetRevision,omitempty"`
}

type Dest struct {
	Server string `json:"server,omitempty"`
}

type GlobalValues struct {
	ClusterConfigPath string `json:"clusterConfigPath,omitempty"`
}

type DefaultApps struct {
	Name           string `yaml:"name" json:"name,omitempty"`
	ValuesPath     string `yaml:"valuesPath" json:"valuesPath,omitempty"`
	RepoURL        string `yaml:"repoURL" json:"repoURL,omitempty"`
	Namespace      string `yaml:"namespace" json:"namespace,omitempty"`
	Chart          string `yaml:"chart" json:"chart,omitempty"`
	TargetRevision string `yaml:"targetRevision" json:"targetRevision,omitempty"`
}

type DefaultAppList struct {
	DefaultApps []DefaultApps `yaml:"defaultApps"`
}

type TektonPipeline struct {
	Name string `json:"name,omitempty"`
}

type SecretNames struct {
	Name string `json:"name,omitempty"`
}

type TektonConfigValues struct {
	Project         string            `json:"project,omitempty"`
	Src             Source            `json:"source,omitempty"`
	Destination     Dest              `json:"destination,omitempty"`
	TektonPath      string            `json:"tektonPath,omitempty"`
	TektonPipelines *[]TektonPipeline `json:"tektonPipelines,omitempty"`
}

type TektonPieplineConfigValues struct {
	PipelineName      string         `json:"pipelineName,omitempty"`
	IngressDomainName string         `json:"ingressDomainName,omitempty"`
	Namespace         string         `json:"namespace,omitempty"`
	SecretName        *[]SecretNames `json:"secretName,omitempty"`
}
