package stock_test

import (
	"testing"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/stock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProvision(t *testing.T) {
	t.Run("propates error from the store", func(t *testing.T) {
		mockStore := &MockStore{}

		expectedErr := errors.New("err")
		mockStore.On("Provision", mock.Anything, mock.Anything).Return(0, expectedErr)

		st := &stock.Stock{Store: mockStore}

		entries := []stock.StockProvisionEntry{{
			ID:  1,
			Qty: 1,
		}}

		qty, err := st.Provision(entries)

		assert.Zero(t, qty)
		assert.Equal(t, expectedErr, err)
		mockStore.AssertExpectations(t)
	})

	t.Run("returns empty map for empty input", func(t *testing.T) {
		mockStore := &MockStore{}

		st := &stock.Stock{Store: mockStore}

		var entries []stock.StockProvisionEntry

		qty, err := st.Provision(entries)

		assert.Empty(t, qty)
		assert.Nil(t, err)

		mockStore.AssertExpectations(t)
	})

	t.Run("calls store", func(t *testing.T) {
		mockStore := &MockStore{}

		mockStore.On("Provision", mock.Anything, mock.Anything).Return(10, nil)

		st := &stock.Stock{Store: mockStore}

		entries := []stock.StockProvisionEntry{{
			ID:  1,
			Qty: 5,
		}}

		qty, err := st.Provision(entries)

		assert.Nil(t, err)
		assert.NotNil(t, qty)
		assert.Equal(t, qty[1], 10)
		mockStore.AssertExpectations(t)
	})
}

func TestQuantity(t *testing.T) {
	t.Run("propates from the store", func(t *testing.T) {
		mockStore := &MockStore{}

		expectedErr := errors.New("err")
		mockStore.On("Quantity", mock.Anything).Return(0, expectedErr)

		st := &stock.Stock{Store: mockStore}
		qty, err := st.Quantity(1)

		assert.Zero(t, qty)
		assert.Equal(t, expectedErr, err)
		mockStore.AssertExpectations(t)
	})
}

func TestCurrentStock(t *testing.T) {
	t.Run("propates error from the store", func(t *testing.T) {
		mockStore := &MockStore{}

		expectedErr := errors.New("err")
		mockStore.On("Quantity", mock.Anything).Return(0, expectedErr)

		st := &stock.Stock{Store: mockStore}
		oneItem := []inventory.Item{{
			ID:   inventory.ID(1),
			Name: inventory.Name("test"),
		}}

		sp, err := st.CurrentStock(oneItem)

		assert.Nil(t, sp)
		assert.Equal(t, expectedErr, err)
		mockStore.AssertExpectations(t)
	})

	t.Run("ignores 'item not found' errors", func(t *testing.T) {
		mockStore := &MockStore{}

		expectedErr := stock.ErrItemNotFound
		mockStore.On("Quantity", mock.Anything).Return(10, expectedErr)

		st := &stock.Stock{Store: mockStore}

		oneItem := []inventory.Item{{
			ID:   inventory.ID(1),
			Name: inventory.Name("test"),
		}}

		sp, err := st.CurrentStock(oneItem)

		assert.Nil(t, err)
		assert.Empty(t, sp)
		mockStore.AssertExpectations(t)
	})

	t.Run("returns the quantity for each item", func(t *testing.T) {
		mockStore := &MockStore{}

		mockStore.On("Quantity", mock.Anything).Return(10, nil)

		st := &stock.Stock{Store: mockStore}
		oneItem := []inventory.Item{{
			ID:   inventory.ID(1),
			Name: inventory.Name("test"),
		}}

		sp, err := st.CurrentStock(oneItem)

		assert.Nil(t, err)

		expectedResult := []stock.StockPosition{{
			ID:   1,
			Name: "test",
			Qty:  10,
		}}

		assert.Equal(t, expectedResult, sp)
		mockStore.AssertExpectations(t)
	})
}

func TestSell(t *testing.T) {
	t.Run("propagates error from store", func(t *testing.T) {
		store := &MockStore{}

		st := stock.Stock{Store: store}

		expectedErr := errors.New("test")

		store.On("Quantity", mock.Anything).Return(0, expectedErr)

		ii := []recipe.Ingredient{{
			ID: 1,
		}}

		err := st.Sell(ii, 1)

		assert.Equal(t, expectedErr, err)
		store.AssertExpectations(t)
	})

	t.Run("returns error when quantity not enough", func(t *testing.T) {
		store := &MockStore{}

		st := stock.Stock{Store: store}

		store.On("Quantity", mock.Anything).Return(1, nil)

		ii := []recipe.Ingredient{{
			ID:  1,
			Qty: 1,
		}}

		err := st.Sell(ii, 2)

		assert.Equal(t, stock.ErrNotEnoughStock, err)
		store.AssertExpectations(t)
	})

	t.Run("calls store to sell items", func(t *testing.T) {
		store := &MockStore{}

		st := stock.Stock{Store: store}

		store.On("Quantity", mock.Anything).Return(2, nil)

		expectedErr := errors.New("expected")
		store.On("Sell", mock.Anything, mock.Anything).Return(expectedErr)

		ii := []recipe.Ingredient{{
			ID:  1,
			Qty: 1,
		}}

		err := st.Sell(ii, 2)

		assert.Equal(t, expectedErr, err)
		store.AssertExpectations(t)
	})
}

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
