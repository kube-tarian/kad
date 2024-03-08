package types

type Plugin struct {
	PluginName          string   `json:"pluginName,omitempty"`
	Description         string   `json:"description,omitempty"`
	Category            string   `json:"category,omitempty"`
	Icon                string   `json:"icon,omitempty"`
	ChartName           string   `json:"chartName,omitempty"`
	ChartRepo           string   `json:"chartRepo,omitempty"`
	Versions            []string `json:"versions,omitempty"`
	DefaultNamespace    string   `json:"defaultNamespace,omitempty"`
	PrivilegedNamespace bool     `json:"privilegedNamespace"`
	PluginEndpoint      string   `json:"pluginAccessEndpoint"`
	Capabilities        []string `json:"capabilities,omitempty"`
}
