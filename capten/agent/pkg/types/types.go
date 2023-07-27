package types

// Todo: These structs along with proto file should be a part of common package

type LaunchUIConfig struct {
	RedirectURL string `yaml:"RedirectURL"`
}

type Override struct {
	LaunchUIConfig LaunchUIConfig `yaml:"LaunchUIConfig"`
	LaunchUIValues map[string]any `yaml:"LaunchUIValues"`
	Values         map[string]any `yaml:"Values"`
}

type AppConfig struct {
	Name        string    `yaml:"Name"`
	ChartName   string    `yaml:"ChartName"`
	RepoName    string    `yaml:"RepoName"`
	RepoURL     string    `yaml:"RepoURL"`
	Namespace   string    `yaml:"Namespace"`
	ReleaseName string    `yaml:"ReleaseName"`
	Version     string    `yaml:"Version"`
	Override    *Override `yaml:"Override"`

	CreateNamespace     *bool `yaml:"CreateNamespace"`
	PrivilegedNamespace *bool `yaml:"PrivilegedNamespace"`
}
