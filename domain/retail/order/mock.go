package order

import (
	"github.com/stretchr/testify/mock"
)

type MockOrderDB struct {
	mock.Mock
}

func (m *MockOrderDB) Add(i Order) (ID, error) {
	args := m.Called(i)
	return args.Get(0).(ID), args.Error(1)
}
