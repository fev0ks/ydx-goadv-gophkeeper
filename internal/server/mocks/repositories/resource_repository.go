// Code generated by MockGen. DO NOT EDIT.
// Source: resource_repository.go

// Package repositories is a generated GoMock package.
package repositories

import (
	context "context"
	reflect "reflect"
	model "ydx-goadv-gophkeeper/internal/server/model"
	enum "ydx-goadv-gophkeeper/pkg/model/enum"

	gomock "github.com/golang/mock/gomock"
)

// MockResourceRepository is a mock of ResourceRepository interface.
type MockResourceRepository struct {
	ctrl     *gomock.Controller
	recorder *MockResourceRepositoryMockRecorder
}

// MockResourceRepositoryMockRecorder is the mock recorder for MockResourceRepository.
type MockResourceRepositoryMockRecorder struct {
	mock *MockResourceRepository
}

// NewMockResourceRepository creates a new mock instance.
func NewMockResourceRepository(ctrl *gomock.Controller) *MockResourceRepository {
	mock := &MockResourceRepository{ctrl: ctrl}
	mock.recorder = &MockResourceRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResourceRepository) EXPECT() *MockResourceRepositoryMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockResourceRepository) Delete(ctx context.Context, resId, userId int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, resId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockResourceRepositoryMockRecorder) Delete(ctx, resId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockResourceRepository)(nil).Delete), ctx, resId, userId)
}

// Get mocks base method.
func (m *MockResourceRepository) Get(ctx context.Context, resId, userId int32) (*model.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, resId, userId)
	ret0, _ := ret[0].(*model.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockResourceRepositoryMockRecorder) Get(ctx, resId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockResourceRepository)(nil).Get), ctx, resId, userId)
}

// GetResDescriptionsByType mocks base method.
func (m *MockResourceRepository) GetResDescriptionsByType(ctx context.Context, userId int32, resType enum.ResourceType) ([]*model.ResourceDescription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResDescriptionsByType", ctx, userId, resType)
	ret0, _ := ret[0].([]*model.ResourceDescription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResDescriptionsByType indicates an expected call of GetResDescriptionsByType.
func (mr *MockResourceRepositoryMockRecorder) GetResDescriptionsByType(ctx, userId, resType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResDescriptionsByType", reflect.TypeOf((*MockResourceRepository)(nil).GetResDescriptionsByType), ctx, userId, resType)
}

// Save mocks base method.
func (m *MockResourceRepository) Save(ctx context.Context, resource *model.Resource) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, resource)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockResourceRepositoryMockRecorder) Save(ctx, resource interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockResourceRepository)(nil).Save), ctx, resource)
}

// Update mocks base method.
func (m *MockResourceRepository) Update(ctx context.Context, resource *model.Resource) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, resource)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockResourceRepositoryMockRecorder) Update(ctx, resource interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockResourceRepository)(nil).Update), ctx, resource)
}
