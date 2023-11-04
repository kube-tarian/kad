package activities

type PluginConfigExtractor struct {
	tektonPluginData     tektonPluginDS
	crossPlanePluginData crossplanePluginDS
}

func NewPluginExtractor(tektonFileName, crossplaneFilename string) (*PluginConfigExtractor, error) {
	pluginInfo, err := ReadTektonPluginConfig(tektonFileName)
	if err != nil {
		return nil, err
	}

	cpluginInfo, err := ReadCrossPlanePluginConfig(crossplaneFilename)
	if err != nil {
		return nil, err
	}

	return &PluginConfigExtractor{
		tektonPluginData:     pluginInfo,
		crossPlanePluginData: cpluginInfo,
	}, nil
}

func (pc *PluginConfigExtractor) tektonGetGitRepo() string {
	return pc.tektonPluginData[GitRepo]
}

func (pc *PluginConfigExtractor) tektonGetGitConfigPath() string {
	return pc.tektonPluginData[GitConfigPath]
}

func (pc *PluginConfigExtractor) tektonGetConfigMainApp() string {
	return pc.tektonPluginData[ConfigMainApp]
}

func (pc *PluginConfigExtractor) crossplaneGetGitRepo() string {
	return pc.crossPlanePluginData[GitRepo]
}

func (pc *PluginConfigExtractor) crossplaneGetGitConfigPath() string {
	return pc.crossPlanePluginData[GitConfigPath]
}

func (pc *PluginConfigExtractor) crossplaneGetConfigMainApp() string {
	return pc.crossPlanePluginData[ConfigMainApp]
}

func (pc *PluginConfigExtractor) GetPluginMap(plugin string) map[string]string {
	switch plugin {
	case Tekton:
		return pc.tektonPluginData
	case CrossPlane:
		return pc.crossPlanePluginData
	default:
		return map[string]string{}
	}
}
