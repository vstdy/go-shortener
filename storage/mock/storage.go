// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package storagemock is a generated GoMock package.
package storagemock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	model "github.com/vstdy/go-shortener/model"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AddURLs mocks base method.
func (m *MockStorage) AddURLs(ctx context.Context, objs []model.URL) ([]model.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddURLs", ctx, objs)
	ret0, _ := ret[0].([]model.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddURLs indicates an expected call of AddURLs.
func (mr *MockStorageMockRecorder) AddURLs(ctx, objs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddURLs", reflect.TypeOf((*MockStorage)(nil).AddURLs), ctx, objs)
}

// Close mocks base method.
func (m *MockStorage) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorage)(nil).Close))
}

// GetURL mocks base method.
func (m *MockStorage) GetURL(ctx context.Context, urlID int) (model.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURL", ctx, urlID)
	ret0, _ := ret[0].(model.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURL indicates an expected call of GetURL.
func (mr *MockStorageMockRecorder) GetURL(ctx, urlID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURL", reflect.TypeOf((*MockStorage)(nil).GetURL), ctx, urlID)
}

// GetUsersURLs mocks base method.
func (m *MockStorage) GetUsersURLs(ctx context.Context, userID uuid.UUID) ([]model.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersURLs", ctx, userID)
	ret0, _ := ret[0].([]model.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersURLs indicates an expected call of GetUsersURLs.
func (mr *MockStorageMockRecorder) GetUsersURLs(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersURLs", reflect.TypeOf((*MockStorage)(nil).GetUsersURLs), ctx, userID)
}

// HasURL mocks base method.
func (m *MockStorage) HasURL(ctx context.Context, urlID int) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasURL", ctx, urlID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasURL indicates an expected call of HasURL.
func (mr *MockStorageMockRecorder) HasURL(ctx, urlID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasURL", reflect.TypeOf((*MockStorage)(nil).HasURL), ctx, urlID)
}

// Ping mocks base method.
func (m *MockStorage) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockStorageMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockStorage)(nil).Ping))
}

// RemoveUsersURLs mocks base method.
func (m *MockStorage) RemoveUsersURLs(ctx context.Context, objs []model.URL) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUsersURLs", ctx, objs)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUsersURLs indicates an expected call of RemoveUsersURLs.
func (mr *MockStorageMockRecorder) RemoveUsersURLs(ctx, objs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUsersURLs", reflect.TypeOf((*MockStorage)(nil).RemoveUsersURLs), ctx, objs)
}
