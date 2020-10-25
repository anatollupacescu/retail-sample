package inventory

import (
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) Update(i Item) error {
	args := m.Called(i)
	return args.Error(0)
}

func (m *MockStore) Add(s string) (int, error) {
	args := m.Called(s)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockStore) Find(s string) (int, error) {
	args := m.Called(s)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockStore) List() ([]Item, error) {
	args := m.Called()
	return args.Get(0).([]Item), args.Error(1)
}

func (m *MockStore) Get(id int) (Item, error) {
	args := m.Called(id)
	return args.Get(0).(Item), args.Error(1)
}
