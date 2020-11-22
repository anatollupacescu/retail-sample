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
	t.Run("given a non existent inventory item", func(t *testing.T) {
		mockStore := &stock.MockStore{}

		var expectedErr = inventory.ErrItemNotFound

		inv := &stock.MockInventory{}
		inv.On("Get", mock.Anything).Return(inventory.Item{}, inventory.ErrItemNotFound)

		st := &stock.Stock{Inventory: inv}

		qty, err := st.Provision(1, 1)

		t.Run("returns error", func(t *testing.T) {
			mockStore.AssertExpectations(t)

			assert.Zero(t, qty)
			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("given that store returns error", func(t *testing.T) {
		mockStore := &stock.MockStore{}

		expectedErr := errors.New("err")
		mockStore.On("Provision", mock.Anything, mock.Anything).Return(0, expectedErr)

		inv := &stock.MockInventory{}
		inv.On("Get", mock.Anything).Return(inventory.Item{}, nil)

		st := &stock.Stock{Store: mockStore, Inventory: inv}

		qty, err := st.Provision(1, 1)

		t.Run("error is propagated", func(t *testing.T) {
			mockStore.AssertExpectations(t)
			assert.Zero(t, qty)
			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("given that provision log returns error", func(t *testing.T) {
		mockStore := &stock.MockStore{}
		mockStore.On("Provision", mock.Anything, mock.Anything).Return(10, nil)

		inv := &stock.MockInventory{}
		inv.On("Get", mock.Anything).Return(inventory.Item{}, nil)

		provisionLog := &stock.MockProvisionLog{}

		var expectedErr = errors.New("expected")
		provisionLog.On("Add", mock.Anything, mock.Anything).Return(0, expectedErr)

		st := &stock.Stock{
			Store:        mockStore,
			Inventory:    inv,
			ProvisionLog: provisionLog,
		}

		qty, err := st.Provision(1, 5)

		t.Run("error is propagated", func(t *testing.T) {
			mockStore.AssertExpectations(t)
			assert.Equal(t, expectedErr, err)
			assert.Zero(t, qty)
		})
	})

	t.Run("given that all calls succeed", func(t *testing.T) {
		mockStore := &stock.MockStore{}
		mockStore.On("Provision", mock.Anything, mock.Anything).Return(10, nil)

		inv := &stock.MockInventory{}
		inv.On("Get", mock.Anything).Return(inventory.Item{}, nil)

		provisionLog := &stock.MockProvisionLog{}
		provisionLog.On("Add", mock.Anything, mock.Anything).Return(1, nil)

		st := &stock.Stock{
			Store:        mockStore,
			Inventory:    inv,
			ProvisionLog: provisionLog,
		}

		id, err := st.Provision(1, 5)

		t.Run("return provision entry id", func(t *testing.T) {
			mockStore.AssertExpectations(t)
			assert.Nil(t, err)
			assert.Equal(t, 1, id)
		})
	})
}

func TestSell(t *testing.T) {
	t.Run("propagates error from store", func(t *testing.T) {
		store := &stock.MockStore{}

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
		store := &stock.MockStore{}

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
		store := &stock.MockStore{}

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
