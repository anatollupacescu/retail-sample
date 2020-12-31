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

		invDB := &stock.MockInventory{}
		defer invDB.AssertExpectations(t)

		invDB.On("Get", mock.Anything).Return(inventory.ItemDTO{}, expectedErr)

		st := &stock.Stock{InventoryDB: invDB}
		err := st.Provision(1, 1)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("propagates error from store", func(t *testing.T) {
		expectedErr := errors.New("err")

		db := &stock.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Provision", mock.Anything, mock.Anything).Return(expectedErr)

		invDB := &stock.MockInventory{}
		defer invDB.AssertExpectations(t)

		invDB.On("Get", mock.Anything).Return(inventory.ItemDTO{}, nil)

		st := &stock.Stock{DB: db, InventoryDB: invDB}
		err := st.Provision(1, 1)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("success", func(t *testing.T) {
		db := &stock.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Provision", mock.Anything, mock.Anything).Return(nil)

		invDB := &stock.MockInventory{}
		defer invDB.AssertExpectations(t)

		invDB.On("Get", mock.Anything).Return(inventory.ItemDTO{}, nil)

		s := &stock.Stock{
			DB:          db,
			InventoryDB: invDB,
		}

		err := s.Provision(1, 5)

		assert.NoError(t, err)
	})
}

func TestSell(t *testing.T) {
	t.Run("propagates error from checking quantity", func(t *testing.T) {
		db := &stock.MockDB{}
		defer db.AssertExpectations(t)

		s := stock.Stock{DB: db}

		expectedErr := errors.New("test")
		db.On("Quantity", mock.Anything).Return(0, expectedErr)

		ii := []recipe.Ingredient{{
			ID: 1,
		}}

		err := s.Sell(ii, 1)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("returns error when quantity not enough", func(t *testing.T) {
		db := &stock.MockDB{}
		defer db.AssertExpectations(t)

		s := stock.Stock{DB: db}

		db.On("Quantity", mock.Anything).Return(1, nil)

		ii := []recipe.Ingredient{{
			ID:  1,
			Qty: 1,
		}}

		err := s.Sell(ii, 2)

		assert.Equal(t, stock.ErrNotEnoughStock, err)
	})

	t.Run("propagates error from sell operation", func(t *testing.T) {
		db := &stock.MockDB{}
		defer db.AssertExpectations(t)

		s := stock.Stock{DB: db}

		db.On("Quantity", mock.Anything).Return(2, nil)

		expectedErr := errors.New("expected")
		db.On("Sell", mock.Anything, mock.Anything).Return(expectedErr)

		ii := []recipe.Ingredient{{
			ID:  1,
			Qty: 1,
		}}

		err := s.Sell(ii, 2)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("success", func(t *testing.T) {
		db := &stock.MockDB{}
		defer db.AssertExpectations(t)

		s := stock.Stock{DB: db}

		db.On("Quantity", mock.Anything).Return(2, nil)
		db.On("Sell", mock.Anything, mock.Anything).Return(nil)

		ii := []recipe.Ingredient{{
			ID:  1,
			Qty: 1,
		}}

		err := s.Sell(ii, 2)

		assert.NoError(t, err)
	})
}
