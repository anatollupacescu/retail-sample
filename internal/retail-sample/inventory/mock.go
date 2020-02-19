package inventory

import (
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) add(s Name) ID {
	args := m.Called(s)
	return args.Get(0).(ID)
}

func (m *MockStore) find(s Name) ID {
	args := m.Called(s)
	return args.Get(0).(ID)
}

func (m *MockStore) all() []Record {
	args := m.Called()
	return args.Get(0).([]Record)
}
