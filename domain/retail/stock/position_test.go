package stock_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

func TestProvision(t *testing.T) {
	t.Run("errors when given quantity is negative", func(t *testing.T) {
		st := new(stock.Position)
		err := st.Provision(-12)

		assert.Equal(t, stock.ErrInvalidProvisionQuantity, err)
		assert.Zero(t, st.Qty)
	})

	t.Run("propagates error from save op", func(t *testing.T) {
		expectedErr := errors.New("err")

		db := &stock.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Save", mock.Anything).Return(expectedErr)

		st := &stock.Position{InventoryID: 1, DB: db}
		err := st.Provision(10)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("success", func(t *testing.T) {
		db := &stock.MockDB{}
		defer db.AssertExpectations(t)

		dto := stock.PositionDTO{
			InventoryID: 1,
			Qty:         11,
		}
		db.On("Save", dto).Return(nil)

		st := &stock.Position{InventoryID: 1, Qty: 1, DB: db}
		err := st.Provision(10)

		assert.NoError(t, err)
		assert.Equal(t, 11, st.Qty)
	})
}

func TestExtract(t *testing.T) {
	t.Run("errors when given quantity is negative", func(t *testing.T) {
		st := new(stock.Position)
		err := st.Extract(-12)

		assert.Equal(t, stock.ErrInvalidExtractQuantity, err)
		assert.Zero(t, st.Qty)
	})

	t.Run("errors when existing quantity is not enough", func(t *testing.T) {
		st := &stock.Position{InventoryID: 1, Qty: 1}
		err := st.Extract(2)

		assert.Equal(t, stock.ErrNotEnoughStock, err)
		assert.Equal(t, 1, st.Qty)
	})

	t.Run("propagates error from save op", func(t *testing.T) {
		expectedErr := errors.New("err")

		db := &stock.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Save", mock.Anything).Return(expectedErr)

		st := &stock.Position{InventoryID: 1, Qty: 10, DB: db}
		err := st.Extract(10)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("success", func(t *testing.T) {
		db := &stock.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Save", stock.PositionDTO{InventoryID: 1, Qty: 1}).Return(nil)

		st := &stock.Position{InventoryID: 1, Qty: 11, DB: db}
		err := st.Extract(10)

		assert.NoError(t, err)
		assert.Equal(t, 1, st.Qty)
	})
}
