package stock

import (
	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) Quantity(id int) (qty int, err error) {
	args := m.Called(id)
	return args.Int(0), args.Error(1)
}

func (m *MockStore) Provision(a1 int, a2 int) (int, error) {
	args := m.Called(a1, a2)
	return args.Int(0), args.Error(1)
}

func (m *MockStore) Sell(a1 []recipe.Ingredient, a2 int) error {
	args := m.Called(a1, a2)
	return args.Error(0)
}

type MockProvisionLog struct {
	mock.Mock
}

func (m *MockProvisionLog) List() ([]ProvisionEntry, error) {
	args := m.Called()
	return args.Get(0).([]ProvisionEntry), args.Error(1)
}

func (m *MockProvisionLog) Add(id, qty int) (int, error) {
	args := m.Called(id, qty)
	return args.Int(0), args.Error(1)
}

func (m *MockProvisionLog) Get(id int) (ProvisionEntry, error) {
	args := m.Called(id)
	return args.Get(0).(ProvisionEntry), args.Error(1)
}

type MockInventory struct {
	mock.Mock
}

func (m *MockInventory) List() ([]inventory.Item, error) {
	args := m.Called()
	return args.Get(0).([]inventory.Item), args.Error(1)
}

func (m *MockInventory) Get(id int) (inventory.Item, error) {
	args := m.Called(id)
	return args.Get(0).(inventory.Item), args.Error(1)
}
