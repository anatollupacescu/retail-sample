package stock_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

func TestProvision(t *testing.T) {
	t.Run("propagates error from inventory", func(t *testing.T) {
		var expectedErr = inventory.ErrItemNotFound

		inv := &stock.MockInventory{}
		defer inv.AssertExpectations(t)

		inv.On("Get", mock.Anything).Return(inventory.Item{}, expectedErr)

		st := &stock.Stock{InventoryDB: inv}
		err := st.Provision(1, 1)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("propagates error from store", func(t *testing.T) {
		expectedErr := errors.New("err")

		store := &stock.MockStore{}
		defer store.AssertExpectations(t)

		store.On("Provision", mock.Anything, mock.Anything).Return(expectedErr)

		inv := &stock.MockInventory{}
		defer inv.AssertExpectations(t)

		inv.On("Get", mock.Anything).Return(inventory.Item{}, nil)

		st := &stock.Stock{DB: store, InventoryDB: inv}
		err := st.Provision(1, 1)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("success", func(t *testing.T) {
		store := &stock.MockStore{}
		defer store.AssertExpectations(t)

		store.On("Provision", mock.Anything, mock.Anything).Return(nil)

		inv := &stock.MockInventory{}
		defer inv.AssertExpectations(t)

		inv.On("Get", mock.Anything).Return(inventory.Item{}, nil)

		st := &stock.Stock{
			DB:          store,
			InventoryDB: inv,
		}

		err := st.Provision(1, 5)

		assert.NoError(t, err)
	})
}

func TestSell(t *testing.T) {
	t.Run("propagates error from checking quantity", func(t *testing.T) {
		store := &stock.MockStore{}
		defer store.AssertExpectations(t)

		s := stock.Stock{DB: store}

		expectedErr := errors.New("test")
		store.On("Quantity", mock.Anything).Return(0, expectedErr)

		ii := []recipe.Ingredient{{
			ID: 1,
		}}

		err := s.Sell(ii, 1)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("returns error when quantity not enough", func(t *testing.T) {
		store := &stock.MockStore{}
		defer store.AssertExpectations(t)

		s := stock.Stock{DB: store}

		store.On("Quantity", mock.Anything).Return(1, nil)

		ii := []recipe.Ingredient{{
			ID:  1,
			Qty: 1,
		}}

		err := s.Sell(ii, 2)

		assert.Equal(t, stock.ErrNotEnoughStock, err)
	})

	t.Run("propagates error from sell operation", func(t *testing.T) {
		store := &stock.MockStore{}
		defer store.AssertExpectations(t)

		s := stock.Stock{DB: store}

		store.On("Quantity", mock.Anything).Return(2, nil)

		expectedErr := errors.New("expected")
		store.On("Sell", mock.Anything, mock.Anything).Return(expectedErr)

		ii := []recipe.Ingredient{{
			ID:  1,
			Qty: 1,
		}}

		err := s.Sell(ii, 2)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("success", func(t *testing.T) {
		store := &stock.MockStore{}
		defer store.AssertExpectations(t)

		s := stock.Stock{DB: store}

		store.On("Quantity", mock.Anything).Return(2, nil)
		store.On("Sell", mock.Anything, mock.Anything).Return(nil)

		ii := []recipe.Ingredient{{
			ID:  1,
			Qty: 1,
		}}

		err := s.Sell(ii, 2)

		assert.NoError(t, err)
	})
}
