// Code generated by MockGen. DO NOT EDIT.
// Source: manager.go

// Package sign is a generated GoMock package.
package sign

import (
	context "context"
	reflect "reflect"

	api "github.com/kubernetes-sigs/kernel-module-management/internal/api"
	utils "github.com/kubernetes-sigs/kernel-module-management/internal/utils"
	gomock "go.uber.org/mock/gomock"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MockSignManager is a mock of SignManager interface.
type MockSignManager struct {
	ctrl     *gomock.Controller
	recorder *MockSignManagerMockRecorder
}

// MockSignManagerMockRecorder is the mock recorder for MockSignManager.
type MockSignManagerMockRecorder struct {
	mock *MockSignManager
}

// NewMockSignManager creates a new mock instance.
func NewMockSignManager(ctrl *gomock.Controller) *MockSignManager {
	mock := &MockSignManager{ctrl: ctrl}
	mock.recorder = &MockSignManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSignManager) EXPECT() *MockSignManagerMockRecorder {
	return m.recorder
}

// GarbageCollect mocks base method.
func (m *MockSignManager) GarbageCollect(ctx context.Context, modName, namespace string, owner v1.Object) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GarbageCollect", ctx, modName, namespace, owner)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GarbageCollect indicates an expected call of GarbageCollect.
func (mr *MockSignManagerMockRecorder) GarbageCollect(ctx, modName, namespace, owner interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GarbageCollect", reflect.TypeOf((*MockSignManager)(nil).GarbageCollect), ctx, modName, namespace, owner)
}

// ShouldSync mocks base method.
func (m *MockSignManager) ShouldSync(ctx context.Context, mld *api.ModuleLoaderData) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShouldSync", ctx, mld)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShouldSync indicates an expected call of ShouldSync.
func (mr *MockSignManagerMockRecorder) ShouldSync(ctx, mld interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShouldSync", reflect.TypeOf((*MockSignManager)(nil).ShouldSync), ctx, mld)
}

// Sync mocks base method.
func (m *MockSignManager) Sync(ctx context.Context, mld *api.ModuleLoaderData, imageToSign string, pushImage bool, owner v1.Object) (utils.Status, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sync", ctx, mld, imageToSign, pushImage, owner)
	ret0, _ := ret[0].(utils.Status)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sync indicates an expected call of Sync.
func (mr *MockSignManagerMockRecorder) Sync(ctx, mld, imageToSign, pushImage, owner interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*MockSignManager)(nil).Sync), ctx, mld, imageToSign, pushImage, owner)
}
