package inventory

import (
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) Add(s Name) ID {
	args := m.Called(s)
	return args.Get(0).(ID)
}

func (m *MockStore) Find(s Name) ID {
	args := m.Called(s)
	return args.Get(0).(ID)
}

func (m *MockStore) List() []Item {
	args := m.Called()
	return args.Get(0).([]Item)
}

func (m *MockStore) Get(id ID) Item {
	args := m.Called()
	return args.Get(0).(Item)
}
