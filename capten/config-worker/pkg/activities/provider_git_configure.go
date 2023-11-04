package activities

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	model2 "github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/model"
	cp "github.com/otiai10/copy"
)

func (hg *HandleGit) configureCrossplaneProvider(ctx context.Context, params *model.CrossplaneUseCase, dirPathInRepo, token string) error {
	gitConfigPlugin := getCICDPlugin()

	// create a dummy directory for creating all the files
	tempDir, err := os.MkdirTemp(hg.config.GitCLoneDir, "temp*")
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

	if err := createProviderConfigs(
		filepath.Join(tempDir, dirPathInRepo),
		params,
		hg.pluginConfig.GetPluginMap(CrossPlane)); err != nil {
		return err
	}

	// copy the contents to the cloned repo and if dir exists, then replace it
	err = cp.Copy(
		filepath.Join(tempDir, dirPathInRepo),
		filepath.Join(reqRepo, dirPathInRepo),
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

	if hg.config.CreatePr {
		_, err = createPR(ctx, params.RepoURL, branchName+"-"+params.Type, defaultBranch, token)
		if err != nil {
			return err
		}
	}

	return nil
}

func createProviderConfigs(dir string, params *model.CrossplaneUseCase, pluginMap map[string]string) error {
	for _, provider := range params.CrossplaneProviders {
		cloudType := provider.CloudType
		providerConfigString, err := createProviderCrdString(provider, params, pluginMap)
		if err != nil {
			return fmt.Errorf("createProviderConfigs: err createProviderCrdString: %v", err)
		}

		// create and write to files
		providerFile := filepath.Join(dir, fmt.Sprintf("%s-provider.yaml", cloudType))

		dir := filepath.Dir(providerFile)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("err while creating directories: %v", dir)
		}

		file, err := os.Create(providerFile)
		if err != nil {
			return fmt.Errorf("err while creating file for provider: %v", err)
		}

		if _, err := file.WriteString(providerConfigString); err != nil {
			return fmt.Errorf("err while writing to controllerconfig: %v", err)
		}

		file.Close()
	}
	return nil
}

func createProviderCrdString(provider model2.CrossplaneProvider, params *model.CrossplaneUseCase, pluginMap map[string]string) (string, error) {
	cloudType := strings.ToLower(provider.CloudType)
	pkg, found := pluginMap[fmt.Sprintf("%s_package", cloudType)]
	if !found {
		return "", fmt.Errorf("plugin package not found for cloudType: %s", cloudType)
	}

	secretPath := fmt.Sprintf("generic/CloudProvider/%s", provider.CloudProviderId)
	providerConfigString := fmt.Sprintf(
		crossplaneProviderTemplate,
		cloudType, secretPath, secretPath,
		cloudType, pkg, cloudType,
	)
	return providerConfigString, nil
}
