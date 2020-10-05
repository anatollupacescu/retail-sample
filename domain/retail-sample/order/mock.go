package order

import (
	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/domain/retail-sample/recipe"
)

type MockOrderStore struct {
	mock.Mock
}

func (m *MockOrderStore) Add(i Order) (ID, error) {
	args := m.Called(i)
	return args.Get(0).(ID), args.Error(1)
}

func (m *MockOrderStore) List() ([]Order, error) {
	args := m.Called()
	return args.Get(0).([]Order), args.Error(1)
}

func (m *MockOrderStore) Get(ID) (Order, error) {
	args := m.Called()
	return args.Get(0).(Order), args.Error(1)
}

type MockRecipeBook struct {
	mock.Mock
}

func (rb *MockRecipeBook) Get(id recipe.ID) (recipe.Recipe, error) {
	args := rb.Called(id)
	return args.Get(0).(recipe.Recipe), args.Error(1)
}

type MockStock struct {
	mock.Mock
}

func (m *MockStock) Sell(ingredients []recipe.Ingredient, qty int) error {
	args := m.Called(ingredients, qty)
	return args.Error(0)
}
