package activities

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kube-tarian/kad/capten/model"
	cp "github.com/otiai10/copy"
	"gopkg.in/yaml.v2"
)

func (hg *HandleGit) configureCrossplaneProvider(ctx context.Context, params *model.UseCase, pathInRepo, token string) error {
	gitConfigPlugin := getCICDPlugin()

	// create a dummy directory for creating all the files
	tempDir, err := os.MkdirTemp(".", "temp*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	reqRepo, err := os.MkdirTemp(hg.config.GitCLoneDir, "clone*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(reqRepo) // clean up

	if err := gitConfigPlugin.Clone(reqRepo, params.RepoURL, token); err != nil {
		return err
	}

	if err := createFiles(
		filepath.Join(tempDir, pathInRepo),
		params,
		hg.pluginConfig.GetPluginMap(CrossPlaneProvider)); err != nil {
		return err
	}

	// copy the contents to the cloned repo and if dir exists, then replace it
	err = cp.Copy(
		filepath.Join(tempDir, pathInRepo),
		filepath.Join(reqRepo, pathInRepo),
		cp.Options{
			OnDirExists: func(src, dest string) cp.DirExistsAction {
				return cp.Replace
			}})
	if err != nil {
		return err
	}

	if err := gitConfigPlugin.Commit(".", "configure CrossplaneProvider for the repo",
		hg.config.GitDefaultCommiterName, hg.config.GitDefaultCommiterEmail); err != nil {
		return err
	}

	localBranchName := branchName + "-" + params.Type
	defaultBranch, err := gitConfigPlugin.GetDefaultBranchName()
	if err != nil {
		return err
	}

	if params.PushToDefaultBranch {
		localBranchName = defaultBranch
	}

	if err := gitConfigPlugin.Push(localBranchName, token); err != nil || params.PushToDefaultBranch {
		return err
	}

	_, err = createPR(ctx, params.RepoURL, branchName+"-"+params.Type, defaultBranch, token)
	if err != nil {
		return err
	}

	return nil
}

func createFiles(dir string, params *model.UseCase, pluginMap map[string]string) error {

	// create controllerConfigs
	for _, providerInfo := range params.CrossplaneProviders {
		cloudType := providerInfo.CloudType
		controllerFile := filepath.Join(dir, fmt.Sprintf("controllerConfig-%s.yaml", cloudType))
		dir := filepath.Dir(controllerFile)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("err while creating directories: %v", dir)
		}
		file, err := os.Create(controllerFile)
		if err != nil {
			return fmt.Errorf("err while creating file for controllerconfig: %v", err)
		}

		secretPath := fmt.Sprintf("generic/CloudProvider/%s", providerInfo.CloudProviderId)
		controllerConfig := fmt.Sprintf(controllerConfig, cloudType, secretPath, secretPath)
		if _, err := file.WriteString(controllerConfig); err != nil {
			return fmt.Errorf("err while writing to controllerconfig: %v", err)
		}
		file.Close()
	}

	// create providers
	for _, providerInfo := range params.CrossplaneProviders {
		cloudType := providerInfo.CloudType
		providerFile := filepath.Join(dir, fmt.Sprintf("provider-%s.yaml", cloudType))
		dir := filepath.Dir(providerFile)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("err while creating directories: %v", dir)
		}
		file, err := os.Create(providerFile)
		if err != nil {
			return fmt.Errorf("err while creating file for provider: %v", err)
		}

		val, found := pluginMap[fmt.Sprintf("%s_package", cloudType)]
		if !found {
			return fmt.Errorf("plugin package not found for cloudType: %s", cloudType)
		}

		override := map[string]any{
			"Package":   val,
			"CloudType": providerInfo.CloudType,
		}
		newMapping, err := replaceTemplateValues(yamlStringToMapping(provider), override)
		if err != nil {
			return fmt.Errorf("err while enriching provider yaml: %v", err)
		}

		buffer := new(bytes.Buffer)
		if err := yaml.NewEncoder(buffer).Encode(newMapping); err != nil {
			return fmt.Errorf("err while encoding newMapping: %v", err)
		}

		if _, err := file.Write(buffer.Bytes()); err != nil {
			return fmt.Errorf("err while writing to controllerconfig: %v", err)
		}
	}
	return nil
}

func yamlStringToMapping(s string) map[string]any {
	// can't marshal directly from string, need to convert to map first
	var initialMapping map[string]any
	err := yaml.NewDecoder(strings.NewReader(s)).Decode(&initialMapping)
	if err != nil {
		log.Println("yamlStringToMapping: err while decoding", err)
		return nil
	}
	return initialMapping
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

// yaml marshalling and unmarshalling doesn't work, err: secret func not found
const controllerConfig = `
apiVersion: pkg.crossplane.io/v1alpha1
kind: ControllerConfig
metadata:
  name: "%s-vault-config"
spec:
  args:
    - --debug
  metadata:
    annotations:
      vault.hashicorp.com/agent-inject: "true"
      vault.hashicorp.com/role: "crossplane-providers"
      vault.hashicorp.com/agent-inject-secret-creds.txt: "%s"
      vault.hashicorp.com/agent-inject-template-creds.txt: |
        {{- with secret "%s" -}}
          [default]
          aws_access_key_id="{{ .access_key }}"
          aws_secret_access_key="{{ .secret_key }}"
        {{- end -}}
`

const provider = `
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-{{.CloudType}}
spec:
  package: "{{.Package}}"
  controllerConfigRef:
    name: "{{.CloudType}}-vault-config"
`
