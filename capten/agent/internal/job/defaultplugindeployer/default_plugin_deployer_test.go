package defaultplugindeployer

import (
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/pluginstorepb"
	"github.com/stretchr/testify/assert"
)

func TestDefaultPluginsDeployer_CronSpec(t *testing.T) {
	// Create a mock logger
	logger := logging.NewLogger()

	// Create a mock plugin store
	mockPluginStore := &MockpluginStore{}

	// Create an instance of DefaultPluginsDeployer
	deployer := &DefaultPluginsDeployer{
		pluginStore: mockPluginStore,
		log:         logger,
		frequency:   "@every 10m",
	}

	// Call the CronSpec method
	cronSpec := deployer.CronSpec()

	// Assert that the returned cron spec is correct
	assert.Equal(t, "@every 10m", cronSpec)
}

func TestDefaultPluginsDeployer_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock logger
	logger := logging.NewLogger()

	// Create a mock plugin store
	mockPluginStore := NewMockpluginStore(ctrl)

	// Create an instance of DefaultPluginsDeployer
	deployer := &DefaultPluginsDeployer{
		pluginStore: mockPluginStore,
		log:         logger,
		frequency:   "@every 10m",
	}

	// Mock the SyncPlugins method of the mock plugin store
	mockPluginStore.EXPECT().SyncPlugins(pluginstorepb.StoreType_DEFAULT_STORE).Return(nil).AnyTimes()

	// Mock the GetPlugins method of the mock plugin store
	mockPluginStore.EXPECT().GetPlugins(pluginstorepb.StoreType_DEFAULT_STORE).Return([]*pluginstorepb.Plugin{
		{
			PluginName: "plugin1",
			Versions:   []string{"1.0.0"},
		},
	}, nil)

	mockPluginStore.EXPECT().DeployPlugin(pluginstorepb.StoreType_DEFAULT_STORE, "plugin1", "1.0.0", []byte{}).Return(nil).AnyTimes()

	// Call the Run method
	deployer.Run()
}

func TestDefaultPluginsDeployer_Run_MultiplePlugins(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock logger
	logger := logging.NewLogger()

	// Create a mock plugin store
	mockPluginStore := NewMockpluginStore(ctrl)

	// Create an instance of DefaultPluginsDeployer
	deployer := &DefaultPluginsDeployer{
		pluginStore: mockPluginStore,
		log:         logger,
		frequency:   "@every 10m",
	}

	// Mock the SyncPlugins method of the mock plugin store
	mockPluginStore.EXPECT().SyncPlugins(pluginstorepb.StoreType_DEFAULT_STORE).Return(nil).AnyTimes()

	// Mock the GetPlugins method of the mock plugin store
	mockPluginStore.EXPECT().GetPlugins(pluginstorepb.StoreType_DEFAULT_STORE).Return([]*pluginstorepb.Plugin{
		{
			PluginName: "plugin1",
			Versions:   []string{"1.0.0"},
		},
		{
			PluginName: "plugin2",
			Versions:   []string{"1.0.0"},
		},
	}, nil)

	mockPluginStore.EXPECT().DeployPlugin(pluginstorepb.StoreType_DEFAULT_STORE, "plugin1", "1.0.0", []byte{}).Return(nil).AnyTimes()
	mockPluginStore.EXPECT().DeployPlugin(pluginstorepb.StoreType_DEFAULT_STORE, "plugin2", "1.0.0", []byte{}).Return(nil).AnyTimes()

	// Call the Run method
	deployer.Run()
}

func TestDefaultPluginsDeployer_Run_MultiplePlugins_OnePluginFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock logger
	logger := logging.NewLogger()

	// Create a mock plugin store
	mockPluginStore := NewMockpluginStore(ctrl)

	// Create an instance of DefaultPluginsDeployer
	deployer := &DefaultPluginsDeployer{
		pluginStore: mockPluginStore,
		log:         logger,
		frequency:   "@every 10m",
	}

	// Mock the SyncPlugins method of the mock plugin store
	mockPluginStore.EXPECT().SyncPlugins(pluginstorepb.StoreType_DEFAULT_STORE).Return(nil).AnyTimes()

	// Mock the GetPlugins method of the mock plugin store
	mockPluginStore.EXPECT().GetPlugins(pluginstorepb.StoreType_DEFAULT_STORE).Return([]*pluginstorepb.Plugin{
		{
			PluginName: "plugin1",
			Versions:   []string{"1.0.0"},
		},
		{
			PluginName: "plugin2",
			Versions:   []string{"1.0.0"},
		},
	}, nil)

	mockPluginStore.EXPECT().DeployPlugin(pluginstorepb.StoreType_DEFAULT_STORE, "plugin1", "1.0.0", []byte{}).Return(nil).AnyTimes()
	mockPluginStore.EXPECT().DeployPlugin(pluginstorepb.StoreType_DEFAULT_STORE, "plugin2", "1.0.0", []byte{}).Return(fmt.Errorf("error")).AnyTimes()

	// Call the Run method
	deployer.Run()
}
