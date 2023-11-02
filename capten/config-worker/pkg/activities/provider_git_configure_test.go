package activities

import (
	"path/filepath"
	"testing"

	model2 "github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/model"
	cp "github.com/otiai10/copy"
)

func TestCreatefiles(t *testing.T) {
	createDirAndFiles(t)
}

func createDirAndFiles(t *testing.T) string {
	dir := t.TempDir()
	configMap := map[string]string{
		"aws_package":   "xpkg.upbound.io/crossplane-contrib/provider-aws:v0.44.0",
		"azure_package": "xpkg.upbound.io/crossplane-contrib/provider-azure:v0.44.0",
	}
	pathInRepo := "configs"
	retDir := dir
	dir = filepath.Join(dir, pathInRepo)
	params := &model.UseCase{CrossplaneProviders: dummyProviderInfo()}
	if err := createFiles(dir, params, configMap); err != nil {
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
