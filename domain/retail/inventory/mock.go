package inventory

import (
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Save(i *ItemDTO) error {
	args := m.Called(i)
	return args.Error(0)
}

func (m *MockDB) Add(s string) (int, error) {
	args := m.Called(s)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockDB) Find(s string) (int, error) {
	args := m.Called(s)
	return args.Get(0).(int), args.Error(1)
}
