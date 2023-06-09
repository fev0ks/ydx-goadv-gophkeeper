// Code generated by MockGen. DO NOT EDIT.
// Source: file_service.go

// Package services is a generated GoMock package.
package services

import (
	os "os"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFileService is a mock of FileService interface.
type MockFileService struct {
	ctrl     *gomock.Controller
	recorder *MockFileServiceMockRecorder
}

// MockFileServiceMockRecorder is the mock recorder for MockFileService.
type MockFileServiceMockRecorder struct {
	mock *MockFileService
}

// NewMockFileService creates a new mock instance.
func NewMockFileService(ctrl *gomock.Controller) *MockFileService {
	mock := &MockFileService{ctrl: ctrl}
	mock.recorder = &MockFileServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileService) EXPECT() *MockFileServiceMockRecorder {
	return m.recorder
}

// ReadFile mocks base method.
func (m *MockFileService) ReadFile(path string, errCh chan error) (chan []byte, os.FileInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadFile", path, errCh)
	ret0, _ := ret[0].(chan []byte)
	ret1, _ := ret[1].(os.FileInfo)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ReadFile indicates an expected call of ReadFile.
func (mr *MockFileServiceMockRecorder) ReadFile(path, errCh interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadFile", reflect.TypeOf((*MockFileService)(nil).ReadFile), path, errCh)
}

// SaveFile mocks base method.
func (m *MockFileService) SaveFile(path string, chunks chan []byte) (chan error, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveFile", path, chunks)
	ret0, _ := ret[0].(chan error)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveFile indicates an expected call of SaveFile.
func (mr *MockFileServiceMockRecorder) SaveFile(path, chunks interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveFile", reflect.TypeOf((*MockFileService)(nil).SaveFile), path, chunks)
}
