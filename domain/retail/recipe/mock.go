package recipe

import (
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Add(r RecipeDTO) (ID, error) {
	args := m.Called(r)
	return args.Get(0).(ID), args.Error(1)
}

func (m *MockDB) Find(n Name) (*RecipeDTO, error) {
	args := m.Called(n)
	return args.Get(0).(*RecipeDTO), args.Error(1)
}

func (m *MockDB) Save(r *RecipeDTO) error {
	return m.Called(r).Error(0)
}

// Get not covered by the tests of this package
func (rb *MockDB) Get(id ID) (RecipeDTO, error) {
	args := rb.Called(id)
	return args.Get(0).(RecipeDTO), args.Error(1)
}
