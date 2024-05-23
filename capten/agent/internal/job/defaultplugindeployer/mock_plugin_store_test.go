package defaultplugindeployer

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	"github.com/kube-tarian/kad/capten/common-pkg/pb/pluginstorepb"
)

// MockpluginStore is a mock of pluginStore interface.
type MockpluginStore struct {
	ctrl     *gomock.Controller
	recorder *MockpluginStoreMockRecorder
}

// MockpluginStoreMockRecorder is the mock recorder for MockpluginStore.
type MockpluginStoreMockRecorder struct {
	mock *MockpluginStore
}

// NewMockpluginStore creates a new mock instance.
func NewMockpluginStore(ctrl *gomock.Controller) *MockpluginStore {
	mock := &MockpluginStore{ctrl: ctrl}
	mock.recorder = &MockpluginStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockpluginStore) EXPECT() *MockpluginStoreMockRecorder {
	return m.recorder
}

// ConfigureStore mocks base method.
func (m *MockpluginStore) ConfigureStore(config *pluginstorepb.PluginStoreConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfigureStore", config)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConfigureStore indicates an expected call of ConfigureStore.
func (mr *MockpluginStoreMockRecorder) ConfigureStore(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfigureStore", reflect.TypeOf((*MockpluginStore)(nil).ConfigureStore), config)
}

// DeployPlugin mocks base method.
func (m *MockpluginStore) DeployPlugin(storeType pluginstorepb.StoreType, pluginName, version string, values []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeployPlugin", storeType, pluginName, version, values)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeployPlugin indicates an expected call of DeployPlugin.
func (mr *MockpluginStoreMockRecorder) DeployPlugin(storeType, pluginName, version, values interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeployPlugin", reflect.TypeOf((*MockpluginStore)(nil).DeployPlugin), storeType, pluginName, version, values)
}

// GetPluginData mocks base method.
func (m *MockpluginStore) GetPluginData(storeType pluginstorepb.StoreType, pluginName string) (*pluginstorepb.PluginData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPluginData", storeType, pluginName)
	ret0, _ := ret[0].(*pluginstorepb.PluginData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPluginData indicates an expected call of GetPluginData.
func (mr *MockpluginStoreMockRecorder) GetPluginData(storeType, pluginName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPluginData", reflect.TypeOf((*MockpluginStore)(nil).GetPluginData), storeType, pluginName)
}

// GetPluginValues mocks base method.
func (m *MockpluginStore) GetPluginValues(storeType pluginstorepb.StoreType, pluginName, version string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPluginValues", storeType, pluginName, version)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPluginValues indicates an expected call of GetPluginValues.
func (mr *MockpluginStoreMockRecorder) GetPluginValues(storeType, pluginName, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPluginValues", reflect.TypeOf((*MockpluginStore)(nil).GetPluginValues), storeType, pluginName, version)
}

// GetPlugins mocks base method.
func (m *MockpluginStore) GetPlugins(storeType pluginstorepb.StoreType) ([]*pluginstorepb.Plugin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlugins", storeType)
	ret0, _ := ret[0].([]*pluginstorepb.Plugin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlugins indicates an expected call of GetPlugins.
func (mr *MockpluginStoreMockRecorder) GetPlugins(storeType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlugins", reflect.TypeOf((*MockpluginStore)(nil).GetPlugins), storeType)
}

// GetStoreConfig mocks base method.
func (m *MockpluginStore) GetStoreConfig(storeType pluginstorepb.StoreType) (*pluginstorepb.PluginStoreConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStoreConfig", storeType)
	ret0, _ := ret[0].(*pluginstorepb.PluginStoreConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStoreConfig indicates an expected call of GetStoreConfig.
func (mr *MockpluginStoreMockRecorder) GetStoreConfig(storeType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStoreConfig", reflect.TypeOf((*MockpluginStore)(nil).GetStoreConfig), storeType)
}

// SyncPlugins mocks base method.
func (m *MockpluginStore) SyncPlugins(storeType pluginstorepb.StoreType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncPlugins", storeType)
	ret0, _ := ret[0].(error)
	return ret0
}

// SyncPlugins indicates an expected call of SyncPlugins.
func (mr *MockpluginStoreMockRecorder) SyncPlugins(storeType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncPlugins", reflect.TypeOf((*MockpluginStore)(nil).SyncPlugins), storeType)
}

// UnDeployPlugin mocks base method.
func (m *MockpluginStore) UnDeployPlugin(storeType pluginstorepb.StoreType, pluginName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnDeployPlugin", storeType, pluginName)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnDeployPlugin indicates an expected call of UnDeployPlugin.
func (mr *MockpluginStoreMockRecorder) UnDeployPlugin(storeType, pluginName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnDeployPlugin", reflect.TypeOf((*MockpluginStore)(nil).UnDeployPlugin), storeType, pluginName)
}
