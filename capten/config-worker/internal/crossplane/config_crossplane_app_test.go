package crossplane

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kube-tarian/kad/capten/model"
	model2 "github.com/kube-tarian/kad/capten/model"
	cp "github.com/otiai10/copy"
)

func TestFileValuesReplace(t *testing.T) {
	dir := t.TempDir()

	path := filepath.Join(dir, "test.yaml")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := file.WriteString("https://github.com/intelops/capten-templates.git"); err != nil {
		t.Fatal(err)
	}
	file.Close()

	if err := replaceCaptenUrls(dir, "src", "replaced"); err != nil {
		t.Fatal(err)
	}

	readBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(readBytes) != "replaced" {
		t.Fail()
	}
}

func TestCreatefiles(t *testing.T) {
	createDirAndFiles(t)
}

func createDirAndFiles(t *testing.T) string {
	dir := t.TempDir()
	configMap := map[string]string{
		"aws":   "dummy.io/crossplane-contrib/provider-aws:v1.1.1",
		"azure": "dummy.io/crossplane-contrib/provider-azure:v1.1.1",
	}
	pathInRepo := "configs"
	retDir := dir
	dir = filepath.Join(dir, pathInRepo)
	params := &model.CrossplaneUseCase{CrossplaneProviders: dummyProviderInfo()}
	app := CrossPlaneApp{pluginConfig: &CrossplanePluginConfig{ProviderPackages: configMap}}
	if err := app.createProviderConfigs(dir, params); err != nil {
		t.Fatal(err)
	}
	return retDir
}

func TestOverrideFile(t *testing.T) {
	pathInRepo := "configs"

	reqRepo := createDirAndFiles(t)
	tempDir := createDirAndFiles(t)

	err := cp.Copy(
		filepath.Join(tempDir, pathInRepo),
		filepath.Join(reqRepo, pathInRepo),
		cp.Options{
			OnDirExists: func(src, dest string) cp.DirExistsAction {
				return cp.Replace
			}})
	if err != nil {
		t.Fatal(err)
	}
}

func dummyProviderInfo() (ret []model2.CrossplaneProvider) {

	ret = append(ret, model2.CrossplaneProvider{
		CloudType:       "aws",
		CloudProviderId: "aws-cp-id1",
	})
	ret = append(ret, model2.CrossplaneProvider{
		CloudType:       "azure",
		CloudProviderId: "azure-cp-id1",
	})
	return

}
