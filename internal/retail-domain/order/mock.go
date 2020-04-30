package order

import "github.com/stretchr/testify/mock"

type MockOrderStore struct {
	mock.Mock
}

func (m *MockOrderStore) Add(i Order) (ID, error) {
	args := m.Called(i)
	return args.Get(0).(ID), args.Error(1)
}

func (m *MockOrderStore) List() ([]Order, error) {
	args := m.Called()
	return args.Get(0).([]Order), args.Error(1)
}

func (m *MockOrderStore) Get(ID) (Order, error) {
	args := m.Called()
	return args.Get(0).(Order), args.Error(1)
}
