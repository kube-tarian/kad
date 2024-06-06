package defaultplugindeployer

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	agentpb "github.com/kube-tarian/kad/capten/common-pkg/pb/agentpb"
)

// MockAgent is a mock of defaultPluginsDeployer interface.
type MockAgent struct {
	ctrl     *gomock.Controller
	recorder *MockAgentMockRecorder
}

// MockAgentMockRecorder is the mock recorder for MockdefaultPluginsDeployer.
type MockAgentMockRecorder struct {
	mock *MockAgent
}

// NewMockdefaultPluginsDeployer creates a new mock instance.
func NewMockAgent(ctrl *gomock.Controller) *MockAgent {
	mock := &MockAgent{ctrl: ctrl}
	mock.recorder = &MockAgentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAgent) EXPECT() *MockAgentMockRecorder {
	return m.recorder
}

// DeployDefaultApps mocks base method.
func (m *MockAgent) DeployDefaultApps(ctx context.Context, request *agentpb.DeployDefaultAppsRequest) (*agentpb.DeployDefaultAppsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeployDefaultApps", ctx, request)
	ret0, _ := ret[0].(*agentpb.DeployDefaultAppsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeployDefaultApps indicates an expected call of DeployDefaultApps.
func (mr *MockAgentMockRecorder) DeployDefaultApps(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeployDefaultApps", reflect.TypeOf((*MockAgent)(nil).DeployDefaultApps), ctx, request)
}
