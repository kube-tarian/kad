package agent

import (
	"log"
	"strings"
	"testing"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestPopulateTemplateValues(t *testing.T) {
	assert := require.New(t)
	logger := logging.NewLogger()
	_ = logger

	appConfig := &agentpb.SyncAppData{
		Config: &agentpb.AppConfig{ReleaseName: "release"},
		Values: &agentpb.AppValues{
			OverrideValues: yamlStringToByte(overrideTemplate),
			LaunchUIValues: yamlStringToByte(launchUiTemplate),
			TemplateValues: yamlStringToByte(totalTemplate),
		},
	}
	_ = appConfig

	overrideRequest := createDummyOverrideValuesRequestBytes()
	launchUiRequest := createDummyLaunchUiValuesRequestBytes()
	assert.True(len(overrideRequest) > 0, "expected overrideRequest to be populated")
	assert.True(len(launchUiRequest) > 0, "expected launchUiRequest to be populated")

	_, marshalled, err := PopulateTemplateValues(appConfig, overrideRequest, launchUiRequest, logger)

	assert.True(strings.Contains(string(marshalled), "capten.intelops.launchUI"))
	assert.True(strings.Contains(string(marshalled), "capten.intelops.override"))
	assert.Nil(err)
}

func TestPopulateTemplateValuesWithNoLaunchValues(t *testing.T) {
	assert := require.New(t)
	logger := logging.NewLogger()
	_ = logger

	appConfig := &agentpb.SyncAppData{
		Config: &agentpb.AppConfig{ReleaseName: "release"},
		Values: &agentpb.AppValues{
			OverrideValues: yamlStringToByte(overrideTemplate),
			LaunchUIValues: yamlStringToByte(launchUiTemplate),
			TemplateValues: yamlStringToByte(totalTemplate),
		},
	}
	_ = appConfig

	overrideRequest := createDummyOverrideValuesRequestBytes()
	assert.True(len(overrideRequest) > 0, "expected overrideRequest to be populated")

	_, marshalled, err := PopulateTemplateValues(appConfig, overrideRequest, nil, logger)

	assert.True(!strings.Contains(string(marshalled), "capten.intelops.launchUI"))
	assert.True(strings.Contains(string(marshalled), "capten.intelops.override"))
	assert.Nil(err)
}

func createDummyOverrideValuesRequestBytes() []byte {
	const overrideTemplate = `
DomainName: "capten.intelops.override"
`
	byt := yamlStringToByte(overrideTemplate)
	return byt

}

func createDummyLaunchUiValuesRequestBytes() []byte {
	const launchUiTemplate = `
DomainName: "capten.intelops.launchUI"
OAuthBaseURL: "capten.base.intelops.launchUI"
ClientId: "some_client_id"
ClientSecret: "some_client_secret"
`
	byt := yamlStringToByte(launchUiTemplate)
	return byt
}

func yamlStringToByte(s string) []byte {
	// can't marshal directly from string, need to convert to map first
	var initialMapping map[string]any
	err := yaml.NewDecoder(strings.NewReader(s)).Decode(&initialMapping)
	if err != nil {
		log.Println("err while decoding", err)
		return nil
	}
	out, err := yaml.Marshal(initialMapping)
	if err != nil {
		log.Println("err while marshalling", err)
		return nil
	}
	return out
}

func byteToMap(byt []byte) map[string]any {
	var initialMapping map[string]any
	err := yaml.Unmarshal(byt, &initialMapping)
	if err != nil {
		log.Println("err while Unmarshal", err)
		return nil
	}
	return initialMapping
}

const launchUiTemplate = `
grafana:
  grafana.ini:
    server:
      root_url: https://grafana.{{.DomainName}}.app/
    auth.generic_oauth:
      allow_assign_grafana_admin: true
      allow_sign_up: true
      api_url: '{{.OAuthBaseURL}}/userinfo'
      auth_url: '{{.OAuthBaseURL}}/oauth2/auth'
      client_id: '{{.ClientId}}'
      client_secret: '{{.ClientSecret}}'
      token_url: '{{.OAuthBaseURL}}/oauth2/token'
`

const totalTemplate = `
alertmanager:
  alertmanagerSpec:
    alertmanagerConfigMatcherStrategy:
      type: None
grafana:
  grafana.ini:
    server:
      root_url: https://grafana.{{.DomainName}}.app/
  ingress:
    enabled: true
    hosts:
    - grafana.{{.DomainName}}.app
    tls:
    - hosts:
      - grafana.{{.DomainName}}.app
`

const overrideTemplate = `
DomainName: "{{.DomainName}}"
`
