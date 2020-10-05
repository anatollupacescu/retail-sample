package recipe

import (
	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/domain/retail-sample/inventory"
)

type MockRecipeStore struct {
	mock.Mock
}

func (m *MockRecipeStore) Add(r Recipe) (ID, error) {
	args := m.Called(r)
	return args.Get(0).(ID), args.Error(1)
}

func (m *MockRecipeStore) List() ([]Recipe, error) {
	args := m.Called()
	return args.Get(0).([]Recipe), args.Error(1)
}

func (m *MockRecipeStore) Get(id ID) (Recipe, error) {
	args := m.Called(id)
	return args.Get(0).(Recipe), args.Error(1)
}

func (m *MockRecipeStore) Save(r Recipe) error {
	return m.Called(r).Error(0)
}

//Inventory

type MockInventory struct {
	mock.Mock
}

func (m *MockInventory) Get(i int) (inventory.Item, error) {
	return m.Called(i).Get(0).(inventory.Item), m.Called(i).Error(1)
}
