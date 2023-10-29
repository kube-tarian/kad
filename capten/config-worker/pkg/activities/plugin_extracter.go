package activities

type PluginConfigExtractor struct {
	pluginData pluginDS
}

func NewPluginExtractor(fileName string) (*PluginConfigExtractor, error) {
	pluginInfo, err := ReadPluginConfig(fileName)
	if err != nil {
		return nil, err
	}

	return &PluginConfigExtractor{pluginData: pluginInfo}, nil
}

func (pc *PluginConfigExtractor) getGitRepo(appName string) string {
	return pc.pluginData[appName][GitRepo]
}

func (pc *PluginConfigExtractor) getGitConfigPath(appName string) string {
	return pc.pluginData[appName][GitConfigPath]
}

func (pc *PluginConfigExtractor) getConfigMainApp(appName string) string {
	return pc.pluginData[appName][ConfigMainApp]
}
