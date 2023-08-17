package agent

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestMergeValues(t *testing.T) {
	assert := require.New(t)

	var mapping map[string]any
	err := yaml.Unmarshal([]byte(launchUiYaml), &mapping)
	assert.Nil(err)

	values := map[string]any{
		"ClientId":     "abc_id",
		"ClientSecret": "abc_secret",
		"OryURL":       "abc_url",
	}
	replaced, err := replaceTemplateValues(mapping, values)
	_ = replaced
	assert.Nil(err)

	m1 := map[any]any{"a": 2, "c": map[any]any{"c_0": 1}}
	m2 := map[any]any{"a": 3, "b": 44, "c": map[any]any{"c_0": 11, "c_1": 12}}
	want := map[any]any{"a": 3, "b": 44, "c": map[any]any{"c_0": 11, "c_1": 12}}
	got := mergeRecursive(m1, m2)
	assert.Equal(want, got)

	var finalMapping map[any]any
	err = yaml.Unmarshal([]byte(overrideYaml), &finalMapping)
	assert.Nil(err)

	finalMapping = mergeRecursive(finalMapping, convertKey(replaced))

	_, err = yaml.Marshal(finalMapping)
	assert.Nil(err)

}

const launchUiYaml = `
grafana:
  sso:
    clientId: "{{.ClientId}}"
    clientSecret: "{{.ClientSecret}}"
    oryAdress: "{{.OryURL}}"
    something: "else_again"`

const overrideYaml = `
grafana:
  enabled: true
  sso:
    something: else`
