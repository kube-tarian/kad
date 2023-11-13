package api

import (
	"bytes"
	"context"
	"html/template"
	"reflect"
	"strings"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/agent/internal/pb/agentpb"
	"github.com/kube-tarian/kad/capten/common-pkg/credential"
	"github.com/kube-tarian/kad/capten/model"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func populateTemplateValues(appConfig *agentpb.SyncAppData, newOverrideValues, launchUiValues []byte, log logging.Logger) (
	*agentpb.SyncAppData, []byte, error) {

	newAppConfig := *appConfig

	// apply new overrideValues on top of the existing ones
	baseOverrideValuesMapping := map[string]any{}
	yaml.Unmarshal(appConfig.Values.OverrideValues, &baseOverrideValuesMapping)

	if len(newOverrideValues) > 0 {
		newOverrideMapping := map[string]any{}
		if err := yaml.Unmarshal(newOverrideValues, &newOverrideMapping); err != nil {
			log.Errorf("failed to unmarshall override for release_name: %s err: %v", appConfig.Config.ReleaseName, err)
			return nil, nil, err
		}
		for k, v := range newOverrideMapping {
			baseOverrideValuesMapping[k] = v
		}
	}

	finalOverrideMappingBytes, err := yaml.Marshal(baseOverrideValuesMapping)
	if err != nil {
		log.Errorf("failed to marshal for release_name: %s err: %v", appConfig.Config.ReleaseName, err)
		return nil, nil, err
	}
	newAppConfig.Values.OverrideValues = finalOverrideMappingBytes

	// replace template with new override values from the request
	populatedTemplateValuesMapping, err := deriveTemplateValuesMapping(finalOverrideMappingBytes, appConfig.Values.TemplateValues)
	if err != nil {
		log.Errorf("failed to derive template values, err: %v", err)
		return nil, nil, err
	}

	launchUiMapping := map[string]any{}

	if len(launchUiValues) > 0 {
		// replace launchUiMapping with the new launchUi values from the request
		// also append override values to that
		launchUiValues = enrichBytesMapping(launchUiValues, finalOverrideMappingBytes)
		launchUiMapping, err = deriveTemplateValuesMapping(launchUiValues, appConfig.Values.LaunchUIValues)
		if err != nil {
			log.Errorf("failed to deriveTemplateValuesMapping, release:%s err: %v", appConfig.Config.ReleaseName, err)
			return nil, nil, err
		}
	}

	// merge final set of values together
	finalTemplateValuesMapping := mergeRecursive(convertKey(populatedTemplateValuesMapping), convertKey(launchUiMapping))
	marshaledOverrideValues, err := yaml.Marshal(finalTemplateValuesMapping)
	if err != nil {
		log.Errorf("failed to Marshal finalOverrideValuesMapping, release:%s err: %v", appConfig.Config.ReleaseName, err)
		return nil, nil, err
	}

	return &newAppConfig, marshaledOverrideValues, nil
}

func getAppLaunchSSOvalues(releaseName string) ([]byte, error) {
	cid, csecret, err := credential.GetAppOauthCredential(context.TODO(), releaseName)
	if err != nil && strings.Contains(err.Error(), "secret not found") {
		// no secret was found so in that case, no sso values need to be returned
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	ssoOverwriteMapping := map[string]any{
		"ClientId":     cid,
		"ClientSecret": csecret,
	}
	return yaml.Marshal(ssoOverwriteMapping)
}

func enrichBytesMapping(base, override []byte) []byte {
	// used to enrich the base mapping with any additional values
	baseMap, overrideMap := map[string]any{}, map[string]any{}
	if err := yaml.Unmarshal(base, &baseMap); err != nil {
		return base
	}
	if err := yaml.Unmarshal(override, &overrideMap); err != nil {
		return base
	}
	for k, v := range overrideMap {
		if _, found := baseMap[k]; !found {
			baseMap[k] = v
		}
	}
	out, err := yaml.Marshal(baseMap)
	if err != nil {
		return base
	}
	return out
}

func replaceTemplateValues(templateData, values map[string]any) (transformedData map[string]any, err error) {
	yamlData, err := yaml.Marshal(templateData)
	if err != nil {
		return
	}

	tmpl, err := template.New("templateVal").Parse(string(yamlData))
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, values)
	if err != nil {
		return
	}

	transformedData = map[string]any{}
	err = yaml.Unmarshal(buf.Bytes(), &transformedData)
	if err != nil {
		return
	}
	return
}

// merge map[any]T and map[any]T where T => map[any]T | any
func mergeRecursive(original, override map[any]any) map[any]any {
	if override == nil {
		return original
	}
	if original == nil {
		original = map[any]any{}
	}
	for k, v := range override {
		// case 1: value not found in original
		if _, found := original[k]; !found {
			original[k] = v
			continue
		}

		// case 2: both are not maps
		if reflect.TypeOf(original[k]).Kind() != reflect.Map &&
			reflect.TypeOf(v).Kind() != reflect.Map {
			original[k] = v
			continue
		}

		// case 3: both are maps and v is not nil
		if reflect.TypeOf(v) != nil {
			original[k] = mergeRecursive(
				original[k].(map[any]any),
				v.(map[any]any),
			)
		}

	}
	return original
}

func convertKey(m map[string]any) map[any]any {
	ret := map[any]any{}
	for k, v := range m {
		ret[k] = v
	}
	return ret
}

func executeTemplateValuesTemplate(data []byte, values map[string]any) (transformedData []byte, err error) {
	if len(data) == 0 {
		return
	}

	tmpl, err := template.New("templateVal").Parse(string(data))
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, values)
	if err != nil {
		return
	}

	transformedData = buf.Bytes()
	return
}

func executeStringTemplateValues(data string, values []byte) (transformedData string, err error) {
	if len(data) == 0 {
		return
	}

	tmpl, err := template.New("templateVal").Parse(data)
	if err != nil {
		return
	}

	mapValues := map[string]any{}
	if err = yaml.Unmarshal(values, &mapValues); err != nil {
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, mapValues)
	if err != nil {
		return
	}

	transformedData = string(buf.Bytes())
	return
}

func deriveTemplateValuesMapping(overrideValues, templateValues []byte) (map[string]any, error) {
	templateValues, err := deriveTemplateValues(overrideValues, templateValues)
	if err != nil {
		return nil, err
	}

	templateValuesMapping := map[string]any{}
	if err := yaml.Unmarshal(templateValues, &templateValuesMapping); err != nil {
		return nil, errors.WithMessagef(err, "failed to Unmarshal template values")
	}
	return templateValuesMapping, nil
}

func deriveTemplateValues(overrideValues, templateValues []byte) ([]byte, error) {
	overrideValuesMapping := map[string]any{}
	if err := yaml.Unmarshal(overrideValues, &overrideValuesMapping); err != nil {
		return nil, errors.WithMessagef(err, "failed to Unmarshal override values")
	}

	templateValues, err := executeTemplateValuesTemplate(templateValues, overrideValuesMapping)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to exeute template values update")
	}
	return templateValues, nil
}

func prepareAppDeployRequestFromSyncApp(data *agentpb.SyncAppData, values []byte) *model.ApplicationInstallRequest {
	return &model.ApplicationInstallRequest{
		PluginName:     "helm",
		RepoName:       data.Config.RepoName,
		RepoURL:        data.Config.RepoURL,
		ChartName:      data.Config.ChartName,
		Namespace:      data.Config.Namespace,
		ReleaseName:    data.Config.ReleaseName,
		Version:        data.Config.Version,
		ClusterName:    "capten",
		OverrideValues: string(values),
		Timeout:        10,
	}
}
