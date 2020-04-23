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
	return args.Get(0).(recipe.ID), args.Error(0)
}

func (b *MockRecipeBook) Get(id recipe.ID) recipe.Recipe {
	return b.Called(id).Get(0).(recipe.Recipe)
}

func (b *MockRecipeBook) List() []recipe.Recipe {
	return b.Called().Get(0).([]recipe.Recipe)
}

// inventory

type MockInventory struct {
	mock.Mock
}

func (m *MockInventory) Add(s inventory.Name) (inventory.ID, error) {
	args := m.Called(s)
	return args.Get(0).(inventory.ID), args.Error(1)
}

func (m *MockInventory) List() []inventory.Item {
	args := m.Called()
	return args.Get(0).([]inventory.Item)
}

func (m *MockInventory) Get(s inventory.ID) inventory.Item {
	args := m.Called(s)
	return args.Get(0).(inventory.Item)
}

func (m *MockInventory) Find(s inventory.Name) inventory.ID {
	args := m.Called(s)
	return args.Get(0).(inventory.ID)
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

func (m *MockOrders) Add(oe order.OrderEntry) order.ID {
	return m.Called(oe).Get(0).(order.ID)
}

func (m *MockOrders) List() []order.Order {
	return m.Called().Get(0).([]order.Order)
}
