package warehouse

import (
	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/recipe"
)

//recipe book

type MockRecipeBoook struct {
	mock.Mock
}

func (b *MockRecipeBoook) Add(name recipe.Name, is []recipe.Ingredient) (recipe.ID, error) {
	args := b.Called(name, is)
	return args.Get(0).(recipe.ID), args.Error(0)
}

func (b *MockRecipeBoook) Get(id recipe.ID) recipe.Recipe {
	return b.Called(id).Get(0).(recipe.Recipe)
}

func (b *MockRecipeBoook) Names() []recipe.Name {
	return b.Called().Get(0).([]recipe.Name)
}

// inventory

type MockInventory struct {
	mock.Mock
}

func (m *MockInventory) Add(s inventory.Name) (inventory.ID, error) {
	args := m.Called(s)
	return args.Get(0).(inventory.ID), args.Error(1)
}

func (m *MockInventory) All() []inventory.Item {
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

//outbound log

type MockOutboundLog struct {
	mock.Mock
}

func (m *MockOutboundLog) Add(i OrderLogEntry) {
	_ = m.Called(i)
}

func (m *MockOutboundLog) List() []OrderLogEntry {
	return nil
}
