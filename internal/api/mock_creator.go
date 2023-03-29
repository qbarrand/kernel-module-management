// Code generated by MockGen. DO NOT EDIT.
// Source: creator.go

// Package api is a generated GoMock package.
package api

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1beta1 "github.com/kubernetes-sigs/kernel-module-management/api/v1beta1"
)

// MockModuleLoaderDataCreator is a mock of ModuleLoaderDataCreator interface.
type MockModuleLoaderDataCreator struct {
	ctrl     *gomock.Controller
	recorder *MockModuleLoaderDataCreatorMockRecorder
}

// MockModuleLoaderDataCreatorMockRecorder is the mock recorder for MockModuleLoaderDataCreator.
type MockModuleLoaderDataCreatorMockRecorder struct {
	mock *MockModuleLoaderDataCreator
}

// NewMockModuleLoaderDataCreator creates a new mock instance.
func NewMockModuleLoaderDataCreator(ctrl *gomock.Controller) *MockModuleLoaderDataCreator {
	mock := &MockModuleLoaderDataCreator{ctrl: ctrl}
	mock.recorder = &MockModuleLoaderDataCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModuleLoaderDataCreator) EXPECT() *MockModuleLoaderDataCreatorMockRecorder {
	return m.recorder
}

// FromModule mocks base method.
func (m *MockModuleLoaderDataCreator) FromModule(mod *v1beta1.Module, kernelVersion string) (*ModuleLoaderData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FromModule", mod, kernelVersion)
	ret0, _ := ret[0].(*ModuleLoaderData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FromModule indicates an expected call of FromModule.
func (mr *MockModuleLoaderDataCreatorMockRecorder) FromModule(mod, kernelVersion interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FromModule", reflect.TypeOf((*MockModuleLoaderDataCreator)(nil).FromModule), mod, kernelVersion)
}

// MockmoduleLoaderDataCreatorHelper is a mock of moduleLoaderDataCreatorHelper interface.
type MockmoduleLoaderDataCreatorHelper struct {
	ctrl     *gomock.Controller
	recorder *MockmoduleLoaderDataCreatorHelperMockRecorder
}

// MockmoduleLoaderDataCreatorHelperMockRecorder is the mock recorder for MockmoduleLoaderDataCreatorHelper.
type MockmoduleLoaderDataCreatorHelperMockRecorder struct {
	mock *MockmoduleLoaderDataCreatorHelper
}

// NewMockmoduleLoaderDataCreatorHelper creates a new mock instance.
func NewMockmoduleLoaderDataCreatorHelper(ctrl *gomock.Controller) *MockmoduleLoaderDataCreatorHelper {
	mock := &MockmoduleLoaderDataCreatorHelper{ctrl: ctrl}
	mock.recorder = &MockmoduleLoaderDataCreatorHelperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmoduleLoaderDataCreatorHelper) EXPECT() *MockmoduleLoaderDataCreatorHelperMockRecorder {
	return m.recorder
}

// findKernelMapping mocks base method.
func (m *MockmoduleLoaderDataCreatorHelper) findKernelMapping(mappings []v1beta1.KernelMapping, kernelVersion string) (*v1beta1.KernelMapping, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "findKernelMapping", mappings, kernelVersion)
	ret0, _ := ret[0].(*v1beta1.KernelMapping)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// findKernelMapping indicates an expected call of findKernelMapping.
func (mr *MockmoduleLoaderDataCreatorHelperMockRecorder) findKernelMapping(mappings, kernelVersion interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "findKernelMapping", reflect.TypeOf((*MockmoduleLoaderDataCreatorHelper)(nil).findKernelMapping), mappings, kernelVersion)
}

// prepareModuleLoaderData mocks base method.
func (m *MockmoduleLoaderDataCreatorHelper) prepareModuleLoaderData(mapping *v1beta1.KernelMapping, mod *v1beta1.Module, kernelVersion string) (*ModuleLoaderData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "prepareModuleLoaderData", mapping, mod, kernelVersion)
	ret0, _ := ret[0].(*ModuleLoaderData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// prepareModuleLoaderData indicates an expected call of prepareModuleLoaderData.
func (mr *MockmoduleLoaderDataCreatorHelperMockRecorder) prepareModuleLoaderData(mapping, mod, kernelVersion interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "prepareModuleLoaderData", reflect.TypeOf((*MockmoduleLoaderDataCreatorHelper)(nil).prepareModuleLoaderData), mapping, mod, kernelVersion)
}

// replaceTemplates mocks base method.
func (m *MockmoduleLoaderDataCreatorHelper) replaceTemplates(mld *ModuleLoaderData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "replaceTemplates", mld)
	ret0, _ := ret[0].(error)
	return ret0
}

// replaceTemplates indicates an expected call of replaceTemplates.
func (mr *MockmoduleLoaderDataCreatorHelperMockRecorder) replaceTemplates(mld interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "replaceTemplates", reflect.TypeOf((*MockmoduleLoaderDataCreatorHelper)(nil).replaceTemplates), mld)
}

// MocksectionGetter is a mock of sectionGetter interface.
type MocksectionGetter struct {
	ctrl     *gomock.Controller
	recorder *MocksectionGetterMockRecorder
}

// MocksectionGetterMockRecorder is the mock recorder for MocksectionGetter.
type MocksectionGetterMockRecorder struct {
	mock *MocksectionGetter
}

// NewMocksectionGetter creates a new mock instance.
func NewMocksectionGetter(ctrl *gomock.Controller) *MocksectionGetter {
	mock := &MocksectionGetter{ctrl: ctrl}
	mock.recorder = &MocksectionGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocksectionGetter) EXPECT() *MocksectionGetterMockRecorder {
	return m.recorder
}

// GetRelevantBuild mocks base method.
func (m *MocksectionGetter) GetRelevantBuild(moduleBuild, mappingBuild *v1beta1.Build) *v1beta1.Build {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRelevantBuild", moduleBuild, mappingBuild)
	ret0, _ := ret[0].(*v1beta1.Build)
	return ret0
}

// GetRelevantBuild indicates an expected call of GetRelevantBuild.
func (mr *MocksectionGetterMockRecorder) GetRelevantBuild(moduleBuild, mappingBuild interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRelevantBuild", reflect.TypeOf((*MocksectionGetter)(nil).GetRelevantBuild), moduleBuild, mappingBuild)
}

// GetRelevantSign mocks base method.
func (m *MocksectionGetter) GetRelevantSign(moduleSign, mappingSign *v1beta1.Sign, kernel string) (*v1beta1.Sign, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRelevantSign", moduleSign, mappingSign, kernel)
	ret0, _ := ret[0].(*v1beta1.Sign)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRelevantSign indicates an expected call of GetRelevantSign.
func (mr *MocksectionGetterMockRecorder) GetRelevantSign(moduleSign, mappingSign, kernel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRelevantSign", reflect.TypeOf((*MocksectionGetter)(nil).GetRelevantSign), moduleSign, mappingSign, kernel)
}
