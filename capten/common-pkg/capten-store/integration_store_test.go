package captenstore

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/agentpb"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/captenpluginspb"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/clusterpluginspb"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/pluginstorepb"
	dbinit "github.com/kube-tarian/kad/capten/common-pkg/postgres/db-init"
	"github.com/kube-tarian/kad/capten/model"
	"github.com/stretchr/testify/assert"
)

var log = logging.NewLogger()

func setEnvConfig() {
	os.Setenv("PG_DB_HOST", "127.0.0.1")
	os.Setenv("PG_DB_PORT", "5432")
	os.Setenv("PG_DB_SERVICE_USERNAME", "capten")
	os.Setenv("PG_DB_SERVICE_USERPASSWORD", "test123!")
	os.Setenv("PG_DB_ADMIN_PASSWORD", "mysecretpassword")
	os.Setenv("PG_DB_DEFAULT_NAME", "postgres")
	os.Setenv("PG_DB_NAME", "captendb")
	os.Setenv("PG_SOURCE_URI", "file://../../database/postgres/migrations")
}

func TestDBCreate(t *testing.T) {
	setEnvConfig()
	if err := dbinit.CreatedDatabase(log); err != nil {
		t.Error(err)
	}

	if err := dbinit.RunMigrations(dbinit.UP); err != nil {
		t.Error(err)
	}
}

func TestGitProjects(t *testing.T) {
	setEnvConfig()

	store, err := NewStore(log)
	if !assert.Nil(t, err) {
		t.Logf("store create error, %v", err)
		return
	}

	err = store.UpsertGitProject(&captenpluginspb.GitProject{
		ProjectUrl: "https://github.com/kube-tarian/kad",
		Labels:     []string{"env", "dev"},
	})
	if !assert.Nil(t, err) {
		t.Logf("add git provider error, %v", err)
		return
	}

	err = store.UpsertGitProject(&captenpluginspb.GitProject{
		ProjectUrl: "https://github.com/kube-tarian/kad",
		Labels:     []string{"env2", "dev"},
	})
	if !assert.Nil(t, err) {
		t.Logf("add git provider error, %v", err)
		return
	}

	labelProjects, err := store.GetGitProjectsByLabels([]string{"env2"})
	if !assert.Nil(t, err) {
		t.Logf("get label git provider error, %v", err)
		return
	}
	assert.Equal(t, len(labelProjects), 1)
	for _, provider := range labelProjects {
		t.Logf("label provider: %v", provider)
	}

	err = store.UpsertGitProject(&captenpluginspb.GitProject{
		Id:         labelProjects[0].Id,
		ProjectUrl: "https://github.com/kube-tarian/kad",
		Labels:     []string{"entest"},
	})
	if !assert.Nil(t, err) {
		t.Logf("add git provider error, %v", err)
		return
	}

	labelProjects, err = store.GetGitProjectsByLabels([]string{"entest"})
	if !assert.Nil(t, err) {
		t.Logf("get label git provider error, %v", err)
		return
	}
	assert.Equal(t, len(labelProjects), 1)

	project, err := store.GetGitProjectForID(labelProjects[0].Id)
	if !assert.Nil(t, err) {
		t.Logf("get git provider error, %v", err)
		return
	}
	assert.Equal(t, project.ProjectUrl, "https://github.com/kube-tarian/kad")
	assert.Equal(t, project.Labels, []string{"entest"})

	allProviders, err := store.GetGitProjects()
	if !assert.Nil(t, err) {
		t.Logf("get all git provider error, %v", err)
		return
	}
	for _, provider := range allProviders {
		t.Logf("all provider: %v", provider)
	}

	err = store.DeleteGitProjectById(labelProjects[0].Id)
	if !assert.Nil(t, err) {
		t.Logf("delete git provider error, %v", err)
		return
	}

	allProviders, err = store.GetGitProjects()
	if !assert.Nil(t, err) {
		t.Logf("get all git provider error, %v", err)
		return
	}
	assert.Equal(t, len(allProviders), 1)
	err = store.DeleteGitProjectById(allProviders[0].Id)
	if !assert.Nil(t, err) {
		t.Logf("delete git provider error, %v", err)
		return
	}
}

func TestContainerRegistry(t *testing.T) {
	setEnvConfig()

	store, err := NewStore(log)
	if !assert.Nil(t, err) {
		t.Logf("store create error, %v", err)
		return
	}

	err = store.UpsertContainerRegistry(&captenpluginspb.ContainerRegistry{
		RegistryUrl: "https://github.com/kube-tarian/kad",
		Labels:      []string{"env", "dev"},
	})
	if !assert.Nil(t, err) {
		t.Logf("add registry error, %v", err)
		return
	}

	err = store.UpsertContainerRegistry(&captenpluginspb.ContainerRegistry{
		RegistryUrl: "https://github.com/kube-tarian/kad",
		Labels:      []string{"env2", "dev"},
	})
	if !assert.Nil(t, err) {
		t.Logf("add registry error, %v", err)
		return
	}

	labelProjects, err := store.GetContainerRegistriesByLabels([]string{"env2"})
	if !assert.Nil(t, err) {
		t.Logf("get label registry error, %v", err)
		return
	}
	assert.Equal(t, len(labelProjects), 1)
	for _, provider := range labelProjects {
		t.Logf("label registry: %v", provider)
	}

	err = store.UpsertContainerRegistry(&captenpluginspb.ContainerRegistry{
		Id:          labelProjects[0].Id,
		RegistryUrl: "https://github.com/kube-tarian/kad",
		Labels:      []string{"entest"},
	})
	if !assert.Nil(t, err) {
		t.Logf("add registry error, %v", err)
		return
	}

	labelProjects, err = store.GetContainerRegistriesByLabels([]string{"entest"})
	if !assert.Nil(t, err) {
		t.Logf("get label registry error, %v", err)
		return
	}
	assert.Equal(t, len(labelProjects), 1)

	project, err := store.GetContainerRegistryForID(labelProjects[0].Id)
	if !assert.Nil(t, err) {
		t.Logf("get registry error, %v", err)
		return
	}
	assert.Equal(t, project.RegistryUrl, "https://github.com/kube-tarian/kad")
	assert.Equal(t, project.Labels, []string{"entest"})

	allProviders, err := store.GetContainerRegistries()
	if !assert.Nil(t, err) {
		t.Logf("get all registry error, %v", err)
		return
	}
	for _, provider := range allProviders {
		t.Logf("all registry: %v", provider)
	}

	err = store.DeleteContainerRegistryById(labelProjects[0].Id)
	if !assert.Nil(t, err) {
		t.Logf("delete registry error, %v", err)
		return
	}

	allProviders, err = store.GetContainerRegistries()
	if !assert.Nil(t, err) {
		t.Logf("get all registry error, %v", err)
		return
	}
	assert.Equal(t, len(allProviders), 1)
	err = store.DeleteContainerRegistryById(allProviders[0].Id)
	if !assert.Nil(t, err) {
		t.Logf("delete registry error, %v", err)
		return
	}
}

func TestCloudProviders(t *testing.T) {
	setEnvConfig()

	store, err := NewStore(log)
	if !assert.Nil(t, err) {
		t.Logf("store create error, %v", err)
		return
	}

	err = store.UpsertCloudProvider(&captenpluginspb.CloudProvider{
		CloudType: "aws",
		Labels:    []string{"env", "dev"},
	})
	if !assert.Nil(t, err) {
		t.Logf("add cloud provider error, %v", err)
		return
	}

	err = store.UpsertCloudProvider(&captenpluginspb.CloudProvider{
		CloudType: "azure",
		Labels:    []string{"env2", "dev"},
	})
	if !assert.Nil(t, err) {
		t.Logf("add cloud provider error, %v", err)
		return
	}

	labelProviders, err := store.GetCloudProvidersByLabels([]string{"env2"})
	if !assert.Nil(t, err) {
		t.Logf("get label cloud provider error, %v", err)
		return
	}
	assert.Equal(t, len(labelProviders), 1)
	for _, provider := range labelProviders {
		t.Logf("label provider: %v", provider)
		assert.Contains(t, provider.Labels, "env2")
	}

	providers, err := store.GetCloudProviders()
	assert.Nil(t, err)
	if !assert.Nil(t, err) {
		t.Logf("get cloud provider error, %v", err)
		return
	}
	assert.Equal(t, len(providers), 2)
	for _, provider := range providers {
		t.Logf("all provider: %v", provider)
	}

	labelCloudTypeProviders, err := store.GetCloudProvidersByLabelsAndCloudType([]string{"env"}, "aws")
	if !assert.Nil(t, err) {
		t.Logf("get label cloud provider error, %v", err)
		return
	}
	assert.Equal(t, len(labelCloudTypeProviders), 1)

	for _, provider := range labelCloudTypeProviders {
		assert.Equal(t, provider.CloudType, "aws")
		assert.Contains(t, provider.Labels, "env")
	}

	for _, provider := range providers {
		idProvider1, err := store.GetCloudProviderForID(provider.Id)
		if !assert.Nil(t, err) {
			t.Logf("get cloud provider error, %v", err)
			return
		}
		assert.Equal(t, provider.Id, idProvider1.Id)

		err = store.DeleteCloudProviderById(provider.Id)
		if !assert.Nil(t, err) {
			t.Logf("delete cloud provider error, %v", err)
			return
		}
	}
	providers, err = store.GetCloudProviders()
	assert.Nil(t, err)
	if !assert.Nil(t, err) {
		t.Logf("get cloud provider error, %v", err)
		return
	}
	assert.Equal(t, len(providers), 0)
}

func TestAppConfig(t *testing.T) {
	setEnvConfig()

	store, err := NewStore(log)
	if !assert.Nil(t, err) {
		t.Logf("store create error, %v", err)
		return
	}

	err = store.UpsertAppConfig(&agentpb.SyncAppData{
		Config: &agentpb.AppConfig{
			ReleaseName:         "release1",
			Version:             "1.0.0",
			Category:            "cat1",
			Description:         "app1 description",
			ChartName:           "app1 chart",
			RepoURL:             "url1",
			Namespace:           "namespace1",
			PrivilegedNamespace: true,
			Icon:                []byte("icon1"),
			UiEndpoint:          "uiendpoint1",
			UiModuleEndpoint:    "uimoduleendpoint1",
			InstallStatus:       "install status1",
			RuntimeStatus:       "runtime status1",
			DefualtApp:          true,
			PluginName:          "plugin1",
			ApiEndpoint:         "apiendpoint1",
			PluginStoreType:     1,
		},
		Values: &agentpb.AppValues{
			OverrideValues: []byte("override values1"),
			LaunchUIValues: []byte("launch ui values1"),
			TemplateValues: []byte("template values1"),
		},
	})
	if !assert.Nil(t, err) {
		t.Logf("add cloud provider error, %v", err)
		return
	}

	apps, err := store.GetAllApps()
	assert.Nil(t, err)
	if !assert.Nil(t, err) {
		t.Logf("get apps error, %v", err)
		return
	}
	assert.Equal(t, len(apps), 1)
	for _, app := range apps {
		assert.Equal(t, app.Config.ReleaseName, "release1")
		assert.Equal(t, app.Config.Category, "cat1")
		assert.Equal(t, app.Config.Description, "app1 description")
		assert.Equal(t, app.Config.Icon, []byte("icon1"))
		assert.Equal(t, app.Config.ChartName, "app1 chart")
		assert.Equal(t, app.Config.RepoURL, "url1")
		assert.Equal(t, app.Config.Namespace, "namespace1")
		assert.Equal(t, app.Config.PrivilegedNamespace, true)
		assert.Equal(t, app.Config.UiEndpoint, "uiendpoint1")
		assert.Equal(t, app.Config.UiModuleEndpoint, "uimoduleendpoint1")
		assert.Equal(t, app.Config.InstallStatus, "install status1")
		assert.Equal(t, app.Config.DefualtApp, true)
		assert.Equal(t, app.Config.PluginName, "plugin1")
	}

	err = store.DeleteAppConfigByReleaseName("release1")
	if !assert.Nil(t, err) {
		t.Logf("get apps error, %v", err)
		return
	}

	apps, err = store.GetAllApps()
	if !assert.Nil(t, err) {
		t.Logf("get apps error, %v", err)
		return
	}
	assert.Equal(t, len(apps), 0)
}

func TestClusterPluginConfig(t *testing.T) {
	setEnvConfig()

	store, err := NewStore(log)
	if !assert.Nil(t, err) {
		t.Logf("store create error, %v", err)
		return
	}

	err = store.UpsertClusterPluginConfig(&clusterpluginspb.Plugin{
		PluginName:          "plugin1",
		StoreType:           1,
		Category:            "cat1",
		Capabilities:        []string{"cap1"},
		Description:         "desc1",
		Icon:                []byte("icon1"),
		ChartName:           "chart1",
		ChartRepo:           "repo1",
		DefaultNamespace:    "namespace1",
		PrivilegedNamespace: true,
		ApiEndpoint:         "apiendpoint1",
		UiEndpoint:          "uiendpoint1",
		Version:             "1.0.0",
		InstallStatus:       "install status1",
		Values:              []byte("values1"),
		OverrideValues:      []byte("override values1"),
	})
	if !assert.Nil(t, err) {
		t.Logf("add cluster plugin config error, %v", err)
		return
	}

	clusterPluginConfigs, err := store.GetAllClusterPluginConfigs()
	if !assert.Nil(t, err) {
		t.Logf("get cluster plugin configs error, %v", err)
		return
	}
	assert.Equal(t, len(clusterPluginConfigs), 1)

	clusterPluginConfig, err := store.GetClusterPluginConfig("plugin1")
	if !assert.Nil(t, err) {
		t.Logf("get cluster plugin configs error, %v", err)
		return
	}
	assert.Equal(t, clusterPluginConfig.PluginName, "plugin1")
	assert.Equal(t, clusterPluginConfig.Category, "cat1")
	assert.Equal(t, clusterPluginConfig.Capabilities, []string{"cap1"})
	assert.Equal(t, clusterPluginConfig.Description, "desc1")
	assert.Equal(t, clusterPluginConfig.Icon, []byte("icon1"))
	assert.Equal(t, clusterPluginConfig.ChartName, "chart1")
	assert.Equal(t, clusterPluginConfig.ChartRepo, "repo1")
	assert.Equal(t, clusterPluginConfig.DefaultNamespace, "namespace1")
	assert.Equal(t, clusterPluginConfig.PrivilegedNamespace, true)
	assert.Equal(t, clusterPluginConfig.ApiEndpoint, "apiendpoint1")
	assert.Equal(t, clusterPluginConfig.UiEndpoint, "uiendpoint1")
	assert.Equal(t, clusterPluginConfig.Version, "1.0.0")
	assert.Equal(t, clusterPluginConfig.InstallStatus, "install status1")
	assert.Equal(t, clusterPluginConfig.Values, []byte("values1"))
	assert.Equal(t, clusterPluginConfig.OverrideValues, []byte("override values1"))

	err = store.DeleteClusterPluginConfig("plugin1")
	if !assert.Nil(t, err) {
		t.Logf("delete cluster plugin config error, %v", err)
		return
	}

	clusterPluginConfigs, err = store.GetAllClusterPluginConfigs()
	if !assert.Nil(t, err) {
		t.Logf("get cluster plugin configs error, %v", err)
		return
	}
	assert.Equal(t, len(clusterPluginConfigs), 0)
}

func TestCrossplaneProvider(t *testing.T) {
	setEnvConfig()

	store, err := NewStore(log)
	if !assert.Nil(t, err) {
		t.Logf("store create error, %v", err)
		return
	}

	err = store.UpsertCrossplaneProvider(&model.CrossplaneProvider{
		ProviderName:    "provider1",
		CloudType:       "aws",
		Status:          "status1",
		CloudProviderId: "cloud provider id1",
	})
	if !assert.Nil(t, err) {
		t.Logf("add crossplane provider error, %v", err)
		return
	}

	crossplaneProviders, err := store.GetCrossplaneProviders()
	if !assert.Nil(t, err) {
		t.Logf("get crossplane providers error, %v", err)
		return
	}
	assert.Equal(t, len(crossplaneProviders), 1)

	err = store.UpsertCrossplaneProvider(&model.CrossplaneProvider{
		Id:              crossplaneProviders[0].Id,
		ProviderName:    "provider1",
		CloudType:       "aws",
		Status:          "status2",
		CloudProviderId: "cloud provider id1",
	})
	if !assert.Nil(t, err) {
		t.Logf("add crossplane provider error, %v", err)
		return
	}

	crossplaneProvider, err := store.GetCrossplanProviderByCloudType("aws")
	if !assert.Nil(t, err) {
		t.Logf("get crossplane providers error, %v", err)
		return
	}
	assert.Equal(t, crossplaneProvider.ProviderName, "provider1")
	assert.Equal(t, crossplaneProvider.CloudType, "aws")
	assert.Equal(t, crossplaneProvider.Status, "status2")
	assert.Equal(t, crossplaneProvider.CloudProviderId, "cloud provider id1")

	crossplaneProvider, err = store.GetCrossplanProviderById(crossplaneProviders[0].Id)
	if !assert.Nil(t, err) {
		t.Logf("get crossplane providers error, %v", err)
		return
	}
	assert.Equal(t, crossplaneProvider.ProviderName, "provider1")
	assert.Equal(t, crossplaneProvider.CloudType, "aws")
	assert.Equal(t, crossplaneProvider.Status, "status2")
	assert.Equal(t, crossplaneProvider.CloudProviderId, "cloud provider id1")

	err = store.DeleteCrossplaneProviderById(crossplaneProviders[0].Id)
	if !assert.Nil(t, err) {
		t.Logf("delete crossplane provider error, %v", err)
		return
	}

	crossplaneProviders, err = store.GetCrossplaneProviders()
	if !assert.Nil(t, err) {
		t.Logf("get crossplane providers error, %v", err)
		return
	}
	assert.Equal(t, len(crossplaneProviders), 0)
}

func TestTektonProject(t *testing.T) {
	setEnvConfig()

	store, err := NewStore(log)
	if !assert.Nil(t, err) {
		t.Logf("store create error, %v", err)
		return
	}

	err = store.UpsertGitProject(&captenpluginspb.GitProject{
		ProjectUrl: "https://github.com/kube-tarian/kad",
		Labels:     []string{"tekton", "dev"},
	})
	if !assert.Nil(t, err) {
		t.Logf("add git provider error, %v", err)
		return
	}

	projects, err := store.GetGitProjectsByLabels([]string{"tekton"})
	if !assert.Nil(t, err) {
		t.Logf("add git project error, %v", err)
		return
	}

	err = store.UpsertTektonProject(&model.TektonProject{
		Id:            "1",
		GitProjectId:  projects[0].Id,
		GitProjectUrl: "git url1",
		Status:        "status1",
	})
	if !assert.Nil(t, err) {
		t.Logf("add tekton project error, %v", err)
		return
	}

	tektonProject, err := store.GetTektonProject()
	if !assert.Nil(t, err) {
		t.Logf("get tekton projects error, %v", err)
		return
	}
	assert.Equal(t, tektonProject.GitProjectId, projects[0].Id)
	assert.Equal(t, tektonProject.GitProjectUrl, "https://github.com/kube-tarian/kad")

	err = store.UpsertTektonProject(&model.TektonProject{
		Id:            "1",
		GitProjectId:  projects[0].Id,
		GitProjectUrl: "git url1",
		Status:        "status1",
	})
	if !assert.Nil(t, err) {
		t.Logf("add tekton project error, %v", err)
		return
	}

	tektonProject, err = store.GetTektonProject()
	if !assert.Nil(t, err) {
		t.Logf("get tekton projects error, %v", err)
		return
	}
	assert.Equal(t, tektonProject.GitProjectId, projects[0].Id)
	assert.Equal(t, tektonProject.GitProjectUrl, "https://github.com/kube-tarian/kad")
	assert.Equal(t, tektonProject.Status, "status1")

	err = store.DeleteTektonProject("1")
	if !assert.Nil(t, err) {
		t.Logf("delete tekton project error, %v", err)
		return
	}

	err = store.DeleteGitProjectById(tektonProject.GitProjectId)
	assert.Nil(t, err)
	assert.Equal(t, tektonProject.GitProjectId, projects[0].Id)
}

func TestCrossplaneProject(t *testing.T) {
	setEnvConfig()

	store, err := NewStore(log)
	if !assert.Nil(t, err) {
		t.Logf("store create error, %v", err)
		return
	}

	err = store.UpsertGitProject(&captenpluginspb.GitProject{
		ProjectUrl: "https://github.com/kube-tarian/kad",
		Labels:     []string{"crossplane", "dev"},
	})
	if !assert.Nil(t, err) {
		t.Logf("add git provider error, %v", err)
		return
	}

	projects, err := store.GetGitProjectsByLabels([]string{"crossplane"})
	if !assert.Nil(t, err) {
		t.Logf("add git project error, %v", err)
		return
	}

	err = store.UpsertCrossplaneProject(&model.CrossplaneProject{
		Id:            "1",
		GitProjectId:  projects[0].Id,
		GitProjectUrl: "git url1",
		Status:        "status1",
	})
	if !assert.Nil(t, err) {
		t.Logf("add crossplane project error, %v", err)
		return
	}

	tektonProject, err := store.GetCrossplaneProject()
	if !assert.Nil(t, err) {
		t.Logf("get crossplane projects error, %v", err)
		return
	}
	assert.Equal(t, tektonProject.GitProjectId, projects[0].Id)
	assert.Equal(t, tektonProject.GitProjectUrl, "https://github.com/kube-tarian/kad")

	err = store.UpsertCrossplaneProject(&model.CrossplaneProject{
		Id:            "1",
		GitProjectId:  projects[0].Id,
		GitProjectUrl: "git url1",
		Status:        "status1",
	})
	if !assert.Nil(t, err) {
		t.Logf("add tekton project error, %v", err)
		return
	}

	tektonProject, err = store.GetCrossplaneProject()
	if !assert.Nil(t, err) {
		t.Logf("get crossplane projects error, %v", err)
		return
	}
	assert.Equal(t, tektonProject.GitProjectId, projects[0].Id)
	assert.Equal(t, tektonProject.GitProjectUrl, "https://github.com/kube-tarian/kad")
	assert.Equal(t, tektonProject.Status, "status1")

	err = store.DeleteCrossplaneProject("1")
	assert.Nil(t, err)

	err = store.DeleteGitProjectById(tektonProject.GitProjectId)
	assert.Nil(t, err)
	assert.Equal(t, tektonProject.GitProjectId, projects[0].Id)
}

func TestManagedCluster(t *testing.T) {
	setEnvConfig()

	store, err := NewStore(log)
	if !assert.Nil(t, err) {
		t.Logf("store create error, %v", err)
		return
	}

	err = store.UpsertManagedCluster(&captenpluginspb.ManagedCluster{
		ClusterName:         "cluster1",
		ClusterEndpoint:     "endpoint1",
		ClusterDeployStatus: "deploy1",
		AppDeployStatus:     "deploy2",
	})
	if !assert.Nil(t, err) {
		t.Logf("add managed cluster error, %v", err)
		return
	}

	managedClusters, err := store.GetManagedClusters()
	if !assert.Nil(t, err) {
		t.Logf("get managed cluster error, %v", err)
		return
	}
	assert.Equal(t, len(managedClusters), 1)
	assert.Equal(t, managedClusters[0].ClusterName, "cluster1")

	err = store.UpsertManagedCluster(&captenpluginspb.ManagedCluster{
		Id:                  managedClusters[0].Id,
		ClusterName:         "cluster2",
		ClusterEndpoint:     "endpoint2",
		ClusterDeployStatus: "deploy1",
		AppDeployStatus:     "deploy3",
	})
	if !assert.Nil(t, err) {
		t.Logf("add managed cluster error, %v", err)
		return
	}
	err = store.UpsertManagedCluster(&captenpluginspb.ManagedCluster{
		Id:                  managedClusters[0].Id,
		ClusterName:         "cluster2",
		ClusterEndpoint:     "endpoint2",
		ClusterDeployStatus: "deploy2",
		AppDeployStatus:     "deploy3",
	})
	if !assert.Nil(t, err) {
		t.Logf("add managed cluster error, %v", err)
		return
	}

	managedCluster, err := store.GetManagedClusterForID(managedClusters[0].Id)
	if !assert.Nil(t, err) {
		t.Logf("get managed cluster error, %v", err)
		return
	}
	assert.Equal(t, managedCluster.ClusterName, "cluster2")
	assert.Equal(t, managedCluster.Id, managedClusters[0].Id)
	assert.Equal(t, managedCluster.ClusterEndpoint, "endpoint2")
	assert.Equal(t, managedCluster.ClusterDeployStatus, "deploy2")
	assert.Equal(t, managedCluster.AppDeployStatus, "deploy3")

	err = store.DeleteManagedClusterById(managedClusters[0].Id)
	if !assert.Nil(t, err) {
		t.Logf("delete managed cluster error, %v", err)
		return
	}
}

func TestPluginStoreConfig(t *testing.T) {
	setEnvConfig()

	store, err := NewStore(log)
	if !assert.Nil(t, err) {
		t.Logf("store create error, %v", err)
		return
	}

	gitProjectId := uuid.New().String()
	err = store.UpsertPluginStoreConfig(&pluginstorepb.PluginStoreConfig{
		StoreType:     pluginstorepb.StoreType_CENTRAL_STORE,
		GitProjectId:  gitProjectId,
		GitProjectURL: "url1",
	})
	if !assert.Nil(t, err) {
		t.Logf("add plugin store config error, %v", err)
		return
	}

	config, err := store.GetPluginStoreConfig(pluginstorepb.StoreType_CENTRAL_STORE)
	if !assert.Nil(t, err) {
		t.Logf("get plugin store config error, %v", err)
		return
	}
	assert.Equal(t, config.GitProjectId, gitProjectId)
	assert.Equal(t, config.GitProjectURL, "url1")

	err = store.UpsertPluginStoreConfig(&pluginstorepb.PluginStoreConfig{
		StoreType:     pluginstorepb.StoreType_CENTRAL_STORE,
		GitProjectId:  gitProjectId,
		GitProjectURL: "url2",
	})
	if !assert.Nil(t, err) {
		t.Logf("add plugin store config error, %v", err)
		return
	}

	config, err = store.GetPluginStoreConfig(pluginstorepb.StoreType_CENTRAL_STORE)
	if !assert.Nil(t, err) {
		t.Logf("get plugin store config error, %v", err)
		return
	}
	assert.Equal(t, config.GitProjectURL, "url2")

	err = store.UpsertPluginStoreConfig(&pluginstorepb.PluginStoreConfig{
		StoreType:     pluginstorepb.StoreType_LOCAL_STORE,
		GitProjectId:  gitProjectId,
		GitProjectURL: "url2",
	})
	if !assert.Nil(t, err) {
		t.Logf("add plugin store config error, %v", err)
		return
	}

	config, err = store.GetPluginStoreConfig(pluginstorepb.StoreType_LOCAL_STORE)
	if !assert.Nil(t, err) {
		t.Logf("get plugin store config error, %v", err)
		return
	}
	assert.Equal(t, config.GitProjectURL, "url2")

	err = store.DeletePluginStoreConfig(pluginstorepb.StoreType_CENTRAL_STORE)
	if !assert.Nil(t, err) {
		t.Logf("delete plugin store config error, %v", err)
		return
	}

	err = store.DeletePluginStoreConfig(pluginstorepb.StoreType_LOCAL_STORE)
	if !assert.Nil(t, err) {
		t.Logf("delete plugin store config error, %v", err)
		return
	}

	_, err = store.GetPluginStoreConfig(pluginstorepb.StoreType_CENTRAL_STORE)
	assert.Error(t, err)
}

func TestPluginStoreData(t *testing.T) {
	setEnvConfig()

	store, err := NewStore(log)
	if !assert.Nil(t, err) {
		t.Logf("store create error, %v", err)
		return
	}

	gitProjectId := uuid.New().String()
	err = store.UpsertPluginStoreData(gitProjectId, &pluginstorepb.PluginData{
		StoreType:   pluginstorepb.StoreType_LOCAL_STORE,
		PluginName:  "plugin1",
		Description: "value1",
		Category:    "category1",
		Versions:    []string{"v1", "v2"},
		Icon:        []byte("icon1"),
	})
	if !assert.Nil(t, err) {
		t.Logf("add plugin store data error, %v", err)
		return
	}

	err = store.UpsertPluginStoreData(gitProjectId, &pluginstorepb.PluginData{
		StoreType:   pluginstorepb.StoreType_CENTRAL_STORE,
		PluginName:  "plugin1",
		Description: "value1",
		Category:    "category2",
		Versions:    []string{"v1", "v2"},
		Icon:        []byte("icon1"),
	})
	if !assert.Nil(t, err) {
		t.Logf("add plugin store data error, %v", err)
		return
	}

	pluginData, err := store.GetPluginStoreData(pluginstorepb.StoreType_LOCAL_STORE, gitProjectId, "plugin1")
	if !assert.Nil(t, err) {
		t.Logf("get plugin store data error, %v", err)
		return
	}

	assert.Equal(t, pluginData.StoreType, pluginstorepb.StoreType_LOCAL_STORE)
	assert.Equal(t, pluginData.PluginName, "plugin1")
	assert.Equal(t, pluginData.Description, "value1")
	assert.Equal(t, pluginData.Category, "category1")
	assert.Equal(t, pluginData.Versions, []string{"v1", "v2"})
	assert.Equal(t, pluginData.Icon, []byte("icon1"))

	pluginData, err = store.GetPluginStoreData(pluginstorepb.StoreType_CENTRAL_STORE, gitProjectId, "plugin1")
	if !assert.Nil(t, err) {
		t.Logf("get plugin store data error, %v", err)
		return
	}
	assert.Equal(t, pluginData.StoreType, pluginstorepb.StoreType_CENTRAL_STORE)
	assert.Equal(t, pluginData.PluginName, "plugin1")
	assert.Equal(t, pluginData.Category, "category2")
	assert.Equal(t, pluginData.Icon, []byte("icon1"))

	plugins, err := store.GetPlugins(gitProjectId)
	if !assert.Nil(t, err) {
		t.Logf("get plugins error, %v", err)
		return
	}
	assert.Equal(t, len(plugins), 2)
	assert.Equal(t, plugins[0].PluginName, "plugin1")
	assert.Equal(t, pluginData.Category, "category2")

	err = store.DeletePluginStoreData(pluginstorepb.StoreType_CENTRAL_STORE, gitProjectId, "plugin1")
	if !assert.Nil(t, err) {
		t.Logf("delete plugin store data error, %v", err)
		return
	}

	err = store.DeletePluginStoreData(pluginstorepb.StoreType_LOCAL_STORE, gitProjectId, "plugin1")
	if !assert.Nil(t, err) {
		t.Logf("delete plugin store data error, %v", err)
		return
	}

	_, err = store.GetPluginStoreData(pluginstorepb.StoreType_LOCAL_STORE, gitProjectId, "plugin1")
	assert.Error(t, err)
}
