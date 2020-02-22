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

func (m *MockStore) all() []Item {
	args := m.Called()
	return args.Get(0).([]Item)
}

func (m *MockStore) get(id ID) Item {
	args := m.Called()
	return args.Get(0).(Item)
}
