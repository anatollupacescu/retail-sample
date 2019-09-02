// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	x "."
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockItemTypeStore is a mock of ItemTypeStore interface
type MockItemTypeStore struct {
	ctrl     *gomock.Controller
	recorder *MockItemTypeStoreMockRecorder
}

// MockItemTypeStoreMockRecorder is the mock recorder for MockItemTypeStore
type MockItemTypeStoreMockRecorder struct {
	mock *MockItemTypeStore
}

// NewMockItemTypeStore creates a new mock instance
func NewMockItemTypeStore(ctrl *gomock.Controller) *MockItemTypeStore {
	mock := &MockItemTypeStore{ctrl: ctrl}
	mock.recorder = &MockItemTypeStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockItemTypeStore) EXPECT() *MockItemTypeStoreMockRecorder {
	return m.recorder
}

// Add mocks base method
func (m *MockItemTypeStore) Add(arg0 string) uint64 {
	ret := m.ctrl.Call(m, "Add", arg0)
	ret0, _ := ret[0].(uint64)
	return ret0
}

// Add indicates an expected call of Add
func (mr *MockItemTypeStoreMockRecorder) Add(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockItemTypeStore)(nil).Add), arg0)
}

// Get mocks base method
func (m *MockItemTypeStore) Get(arg0 uint64) x.Entity {
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(x.Entity)
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockItemTypeStoreMockRecorder) Get(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockItemTypeStore)(nil).Get), arg0)
}

// Remove mocks base method
func (m *MockItemTypeStore) Remove(arg0 uint64) {
	m.ctrl.Call(m, "Remove", arg0)
}

// Remove indicates an expected call of Remove
func (mr *MockItemTypeStoreMockRecorder) Remove(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockItemTypeStore)(nil).Remove), arg0)
}

// List mocks base method
func (m *MockItemTypeStore) List() []x.Entity {
	ret := m.ctrl.Call(m, "List")
	ret0, _ := ret[0].([]x.Entity)
	return ret0
}

// List indicates an expected call of List
func (mr *MockItemTypeStoreMockRecorder) List() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockItemTypeStore)(nil).List))
}
