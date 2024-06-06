package defaultplugindeployer

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/agentpb"
	"github.com/stretchr/testify/assert"
)

func TestDefaultPluginsDeployer_CronSpec(t *testing.T) {
	// Create a mock logger
	logger := logging.NewLogger()

	// Create a mock plugin store
	mockAgent := &MockAgent{}

	// Create an instance of DefaultPluginsDeployer
	deployer := &DefaultPluginsDeployer{
		log:       logger,
		frequency: "@every 10m",
		agent:     mockAgent,
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
	mockAgent := NewMockAgent(ctrl)

	// Create an instance of DefaultPluginsDeployer
	deployer := &DefaultPluginsDeployer{
		log:       logger,
		frequency: "@every 10m",
		agent:     mockAgent,
	}

	// Mock the SyncPlugins method of the mock plugin store
	mockAgent.EXPECT().DeployDefaultApps(context.TODO(), &agentpb.DeployDefaultAppsRequest{Upgrade: false}).Return(&agentpb.DeployDefaultAppsResponse{
		Status:        agentpb.StatusCode_OK,
		StatusMessage: "",
	}, nil).AnyTimes()

	// Call the Run method
	deployer.Run()
}

func TestDefaultPluginsDeployer_RunWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock logger
	logger := logging.NewLogger()

	// Create a mock plugin store
	mockAgent := NewMockAgent(ctrl)

	// Create an instance of DefaultPluginsDeployer
	deployer := &DefaultPluginsDeployer{
		log:       logger,
		frequency: "@every 10m",
		agent:     mockAgent,
	}

	// Mock the SyncPlugins method of the mock plugin store
	mockAgent.EXPECT().DeployDefaultApps(context.TODO(), &agentpb.DeployDefaultAppsRequest{Upgrade: false}).Return(&agentpb.DeployDefaultAppsResponse{
		Status:        agentpb.StatusCode_INTERNRAL_ERROR,
		StatusMessage: "failed",
	}, nil).AnyTimes()

	// Call the Run method
	deployer.Run()
}
