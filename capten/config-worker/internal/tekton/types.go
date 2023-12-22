package tekton

type appConfig struct {
	MainAppGitPath string `json:"mainAppGitPath"`
	SynchApp       bool   `json:"synchApp"`
}

type pipelineSyncUpdate struct {
	MainAppValues  string `json:"mainAppValues"`
	PipelineValues string `json:"pipelineValues"`
}

type tektonPluginConfig struct {
	TemplateGitRepo               string             `json:"templateGitRepo"`
	PipelineClusterConfigSyncPath string             `json:"pipelineClusterConfigSyncPath"`
	PipelineConfigSyncPath        string             `json:"pipelineConfigSyncPath"`
	TektonProject                 string             `json:"tektonProject"`
	TektonPipelinePath            string             `json:"TektonPipelinePath"`
	ArgoCDApps                    []appConfig        `json:"argoCDApps"`
	PipelineSyncUpdate            pipelineSyncUpdate `json:"pipelineSyncUpdate"`
}

type DockerConfigEntry struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty" datapolicy:"password"`
	Email    string `json:"email,omitempty"`
	Auth     string `json:"auth,omitempty" datapolicy:"token"`
}

type DockerConfig map[string]DockerConfigEntry

type DockerConfigJSON struct {
	Auths DockerConfig `json:"auths" datapolicy:"token"`
	// +optional
	HttpHeaders map[string]string `json:"HttpHeaders,omitempty" datapolicy:"token"`
}
