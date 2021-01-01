package order

import (
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Add(i OrderDTO) (ID, error) {
	args := m.Called(i)
	return args.Get(0).(ID), args.Error(1)
}

type MockStock struct {
	mock.Mock
}

func (m *MockStock) Extract(id, qty int) error {
	return m.Called(id, qty).Error(0)
}
