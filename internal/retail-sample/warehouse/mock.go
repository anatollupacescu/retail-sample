package warehouse

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/recipe"
)

//recipe book

type MockRecipeBoook struct {
	mock.Mock
}

func (b *MockRecipeBoook) Add(name string, is []recipe.Ingredient) error {
	return b.Called(name, is).Error(0)
}

func (b *MockRecipeBoook) Get(id int) recipe.Recipe {
	return b.Called(id).Get(0).(recipe.Recipe)
}

// inventory

type MockInventory struct {
	mock.Mock
}

func (m *MockInventory) Add(s string) (int, error) {
	args := m.Called(s)
	return args.Int(0), args.Error(1)
}

func (m *MockInventory) All() []inventory.Record {
	args := m.Called()
	return args.Get(0).([]inventory.Record)
}

func (m *MockInventory) Get(s int) string {
	args := m.Called(s)
	return args.String(0)
}

func (m *MockInventory) Find(s string) int {
	args := m.Called(s)
	return args.Int(0)
}

//inbound log

type MockInboundLog struct {
	mock.Mock
}

func (m *MockInboundLog) Add(t time.Time, i ProvisionEntry) {
	_ = m.Called(t, i)
}

func (m *MockInboundLog) List() []ProvisionEntry {
	return nil
}

//outbound log

type MockOutboundLog struct {
	mock.Mock
}

func (m *MockOutboundLog) Add(i SoldItem) {
	_ = m.Called(i)
}

func (m *MockOutboundLog) List() []SoldItem {
	return nil
}
