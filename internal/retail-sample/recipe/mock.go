package recipe

import "github.com/stretchr/testify/mock"

type MockRecipeStore struct {
	mock.Mock
}

func (m *MockRecipeStore) add(r Recipe) (ID, error) {
	args := m.Called(r)
	return args.Get(0).(ID), args.Error(1)
}

func (m *MockRecipeStore) all() []Recipe {
	return m.Called().Get(0).([]Recipe)
}

func (m *MockRecipeStore) get(id ID) Recipe {
	return m.Called(id).Get(0).(Recipe)
}

//Inventory

type MockInventory struct {
	mock.Mock
}

func (m *MockInventory) Get(i int) string {
	return m.Called(i).String(0)
}
