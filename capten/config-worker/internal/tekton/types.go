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
