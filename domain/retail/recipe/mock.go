package recipe

import (
	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

type MockRecipeStore struct {
	mock.Mock
}

func (m *MockRecipeStore) Add(r Recipe) (ID, error) {
	args := m.Called(r)
	return args.Get(0).(ID), args.Error(1)
}

func (m *MockRecipeStore) Find(n Name) (*Recipe, error) {
	args := m.Called(n)
	return args.Get(0).(*Recipe), args.Error(1)
}

func (m *MockRecipeStore) Save(r *Recipe) error {
	return m.Called(r).Error(0)
}

//Inventory

type MockInventory struct {
	mock.Mock
}

func (m *MockInventory) Get(i int) (inventory.Item, error) {
	return m.Called(i).Get(0).(inventory.Item), m.Called(i).Error(1)
}
