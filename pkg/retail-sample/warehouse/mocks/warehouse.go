// Code generated by MockGen. DO NOT EDIT.
// Source: warehouse.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockItemStore is a mock of ItemStore interface
type MockItemStore struct {
	ctrl     *gomock.Controller
	recorder *MockItemStoreMockRecorder
}

// MockItemStoreMockRecorder is the mock recorder for MockItemStore
type MockItemStoreMockRecorder struct {
	mock *MockItemStore
}

// NewMockItemStore creates a new mock instance
func NewMockItemStore(ctrl *gomock.Controller) *MockItemStore {
	mock := &MockItemStore{ctrl: ctrl}
	mock.recorder = &MockItemStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockItemStore) EXPECT() *MockItemStoreMockRecorder {
	return m.recorder
}

// Add mocks base method
func (m *MockItemStore) Add(arg0 uint64, arg1 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Add", arg0, arg1)
}

// Add indicates an expected call of Add
func (mr *MockItemStoreMockRecorder) Add(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockItemStore)(nil).Add), arg0, arg1)
}

// Update mocks base method
func (m *MockItemStore) Update(arg0 uint64, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockItemStoreMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockItemStore)(nil).Update), arg0, arg1)
}

// Get mocks base method
func (m *MockItemStore) Get(arg0 uint64) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockItemStoreMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockItemStore)(nil).Get), arg0)
}