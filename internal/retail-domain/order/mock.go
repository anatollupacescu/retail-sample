package order

import "github.com/stretchr/testify/mock"

type MockOrderStore struct {
	mock.Mock
}

func (m *MockOrderStore) Add(i Order) (ID, error) {
	return m.Called(i).Get(0).(ID), m.Called(i).Error(1)
}

func (m *MockOrderStore) List() ([]Order, error) {
	return m.Called().Get(0).([]Order), m.Called().Error(1)
}
