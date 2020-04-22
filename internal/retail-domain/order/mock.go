package order

import "github.com/stretchr/testify/mock"

type MockOrderStore struct {
	mock.Mock
}

func (m *MockOrderStore) Add(i Order) ID {
	return m.Called(i).Get(0).(ID)
}

func (m *MockOrderStore) List() []Order {
	return m.Called().Get(0).([]Order)
}
