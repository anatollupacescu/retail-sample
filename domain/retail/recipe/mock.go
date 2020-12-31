package recipe

import (
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Add(r Recipe) (ID, error) {
	args := m.Called(r)
	return args.Get(0).(ID), args.Error(1)
}

func (m *MockDB) Find(n Name) (*Recipe, error) {
	args := m.Called(n)
	return args.Get(0).(*Recipe), args.Error(1)
}

func (m *MockDB) Save(r *Recipe) error {
	return m.Called(r).Error(0)
}

// Get not covered by the tests of this package
func (rb *MockDB) Get(id ID) (Recipe, error) {
	args := rb.Called(id)
	return args.Get(0).(Recipe), args.Error(1)
}
