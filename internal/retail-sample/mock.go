package retailsample

import (
	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/order"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
)

//recipe book

type MockRecipeBook struct {
	mock.Mock
}

func (b *MockRecipeBook) Add(name recipe.Name, is []recipe.Ingredient) (recipe.ID, error) {
	args := b.Called(name, is)
	return args.Get(0).(recipe.ID), args.Error(1)
}

func (b *MockRecipeBook) Get(id recipe.ID) (recipe.Recipe, error) {
	args := b.Called(id)
	return args.Get(0).(recipe.Recipe), args.Error(1)
}

func (b *MockRecipeBook) List() ([]recipe.Recipe, error) {
	args := b.Called()
	return args.Get(0).([]recipe.Recipe), args.Error(1)
}

// inventory

type MockInventory struct {
	mock.Mock
}

func (m *MockInventory) Add(s inventory.Name) (inventory.ID, error) {
	args := m.Called(s)
	return args.Get(0).(inventory.ID), args.Error(1)
}

func (m *MockInventory) List() ([]inventory.Item, error) {
	args := m.Called()
	return args.Get(0).([]inventory.Item), args.Error(1)
}

func (m *MockInventory) Get(s inventory.ID) (inventory.Item, error) {
	args := m.Called(s)
	return args.Get(0).(inventory.Item), args.Error(1)
}

func (m *MockInventory) Find(s inventory.Name) (inventory.ID, error) {
	args := m.Called(s)
	return args.Get(0).(inventory.ID), args.Error(1)
}

//inbound log

type MockInboundLog struct {
	mock.Mock
}

func (m *MockInboundLog) Add(i ProvisionEntry) {
	_ = m.Called(i)
}

func (m *MockInboundLog) List() []ProvisionEntry {
	return nil
}

type MockOrders struct {
	mock.Mock
}

func (m *MockOrders) Add(oe order.OrderEntry) (order.ID, error) {
	args := m.Called(oe)
	return args.Get(0).(order.ID), args.Error(1)
}

func (m *MockOrders) List() ([]order.Order, error) {
	args := m.Called()
	return args.Get(0).([]order.Order), args.Error(1)
}
