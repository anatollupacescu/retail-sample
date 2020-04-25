package inventory

import (
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) Add(s Name) (ID, error) {
	args := m.Called(s)
	return args.Get(0).(ID), args.Error(1)
}

func (m *MockStore) Find(s Name) (ID, error) {
	args := m.Called(s)
	return args.Get(0).(ID), args.Error(1)
}

func (m *MockStore) List() ([]Item, error) {
	args := m.Called()
	return args.Get(0).([]Item), args.Error(1)
}

func (m *MockStore) Get(id ID) (Item, error) {
	args := m.Called()
	return args.Get(0).(Item), args.Error(1)
}
