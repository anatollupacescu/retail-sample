package stock

import (
	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) Quantity(id int) (qty int, err error) {
	args := m.Called(id)
	return args.Int(0), args.Error(1)
}

func (m *MockStore) Provision(a1 int, a2 int) error {
	return m.Called(a1, a2).Error(0)
}

func (m *MockStore) Sell(a1 []recipe.Ingredient, a2 int) error {
	return m.Called(a1, a2).Error(0)
}

type MockInventory struct {
	mock.Mock
}

func (m *MockInventory) Get(id int) (inventory.Item, error) {
	args := m.Called(id)
	return args.Get(0).(inventory.Item), args.Error(1)
}
