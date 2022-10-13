// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nhost/hasura-storage/controller (interfaces: MetadataStorage)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	controller "github.com/nhost/hasura-storage/controller"
)

// MockMetadataStorage is a mock of MetadataStorage interface.
type MockMetadataStorage struct {
	ctrl     *gomock.Controller
	recorder *MockMetadataStorageMockRecorder
}

// MockMetadataStorageMockRecorder is the mock recorder for MockMetadataStorage.
type MockMetadataStorageMockRecorder struct {
	mock *MockMetadataStorage
}

// NewMockMetadataStorage creates a new mock instance.
func NewMockMetadataStorage(ctrl *gomock.Controller) *MockMetadataStorage {
	mock := &MockMetadataStorage{ctrl: ctrl}
	mock.recorder = &MockMetadataStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetadataStorage) EXPECT() *MockMetadataStorageMockRecorder {
	return m.recorder
}

// DeleteFileByID mocks base method.
func (m *MockMetadataStorage) DeleteFileByID(arg0 context.Context, arg1 string, arg2 http.Header) *controller.APIError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFileByID", arg0, arg1, arg2)
	ret0, _ := ret[0].(*controller.APIError)
	return ret0
}

// DeleteFileByID indicates an expected call of DeleteFileByID.
func (mr *MockMetadataStorageMockRecorder) DeleteFileByID(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFileByID", reflect.TypeOf((*MockMetadataStorage)(nil).DeleteFileByID), arg0, arg1, arg2)
}

// GetBucketByID mocks base method.
func (m *MockMetadataStorage) GetBucketByID(arg0 context.Context, arg1 string, arg2 http.Header) (controller.BucketMetadata, *controller.APIError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBucketByID", arg0, arg1, arg2)
	ret0, _ := ret[0].(controller.BucketMetadata)
	ret1, _ := ret[1].(*controller.APIError)
	return ret0, ret1
}

// GetBucketByID indicates an expected call of GetBucketByID.
func (mr *MockMetadataStorageMockRecorder) GetBucketByID(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBucketByID", reflect.TypeOf((*MockMetadataStorage)(nil).GetBucketByID), arg0, arg1, arg2)
}

// GetFileByID mocks base method.
func (m *MockMetadataStorage) GetFileByID(arg0 context.Context, arg1 string, arg2 http.Header) (controller.FileMetadata, *controller.APIError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileByID", arg0, arg1, arg2)
	ret0, _ := ret[0].(controller.FileMetadata)
	ret1, _ := ret[1].(*controller.APIError)
	return ret0, ret1
}

// GetFileByID indicates an expected call of GetFileByID.
func (mr *MockMetadataStorageMockRecorder) GetFileByID(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileByID", reflect.TypeOf((*MockMetadataStorage)(nil).GetFileByID), arg0, arg1, arg2)
}

// InitializeFile mocks base method.
func (m *MockMetadataStorage) InitializeFile(arg0 context.Context, arg1 string, arg2 http.Header) *controller.APIError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitializeFile", arg0, arg1, arg2)
	ret0, _ := ret[0].(*controller.APIError)
	return ret0
}

// InitializeFile indicates an expected call of InitializeFile.
func (mr *MockMetadataStorageMockRecorder) InitializeFile(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitializeFile", reflect.TypeOf((*MockMetadataStorage)(nil).InitializeFile), arg0, arg1, arg2)
}

// ListFiles mocks base method.
func (m *MockMetadataStorage) ListFiles(arg0 context.Context, arg1 http.Header) ([]controller.FileSummary, *controller.APIError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFiles", arg0, arg1)
	ret0, _ := ret[0].([]controller.FileSummary)
	ret1, _ := ret[1].(*controller.APIError)
	return ret0, ret1
}

// ListFiles indicates an expected call of ListFiles.
func (mr *MockMetadataStorageMockRecorder) ListFiles(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFiles", reflect.TypeOf((*MockMetadataStorage)(nil).ListFiles), arg0, arg1)
}

// PopulateMetadata mocks base method.
func (m *MockMetadataStorage) PopulateMetadata(arg0 context.Context, arg1, arg2 string, arg3 int64, arg4, arg5 string, arg6 bool, arg7 string, arg8 http.Header) (controller.FileMetadata, *controller.APIError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PopulateMetadata", arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8)
	ret0, _ := ret[0].(controller.FileMetadata)
	ret1, _ := ret[1].(*controller.APIError)
	return ret0, ret1
}

// PopulateMetadata indicates an expected call of PopulateMetadata.
func (mr *MockMetadataStorageMockRecorder) PopulateMetadata(arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PopulateMetadata", reflect.TypeOf((*MockMetadataStorage)(nil).PopulateMetadata), arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8)
}

// SetIsUploaded mocks base method.
func (m *MockMetadataStorage) SetIsUploaded(arg0 context.Context, arg1 string, arg2 bool, arg3 http.Header) *controller.APIError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetIsUploaded", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*controller.APIError)
	return ret0
}

// SetIsUploaded indicates an expected call of SetIsUploaded.
func (mr *MockMetadataStorageMockRecorder) SetIsUploaded(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetIsUploaded", reflect.TypeOf((*MockMetadataStorage)(nil).SetIsUploaded), arg0, arg1, arg2, arg3)
}