package warehouse_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/recipe"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/warehouse"

	"github.com/stretchr/testify/assert"
)

func TestPlaceOrder(t *testing.T) {
	t.Run("should reject order with non existent id", func(t *testing.T) {
		b := warehouse.MockRecipeBoook{}

		var zeroRecipe recipe.Recipe
		b.On("Get", recipe.ID(1)).Return(zeroRecipe)

		stock := warehouse.NewStock(nil, nil, &b, nil)

		err := stock.PlaceOrder(1, 10)

		assert.Equal(t, warehouse.ErrRecipeNotFound, err)
		b.AssertExpectations(t)
	})

	t.Run("should reject outbound when not enough stock", func(t *testing.T) {
		b := warehouse.MockRecipeBoook{}

		b.On("Get", recipe.ID(1)).Return(recipe.Recipe{
			Name:        "test",
			Ingredients: []recipe.Ingredient{{ID: 51, Qty: 2}},
		})

		data := map[int]int{
			51: 9,
		}
		stock := warehouse.NewStockWithData(nil, nil, &b, nil, data)

		err := stock.PlaceOrder(1, 5)

		assert.Equal(t, warehouse.ErrNotEnoughStock, err)
		b.AssertExpectations(t)
	})

	t.Run("should update inventory on success", func(t *testing.T) {
		b := warehouse.MockRecipeBoook{}

		b.On("Get", recipe.ID(1)).Return(recipe.Recipe{
			Name:        recipe.Name("test"),
			Ingredients: []recipe.Ingredient{{ID: 51, Qty: 2}},
		})

		ol := warehouse.MockOutboundLog{}
		ol.On("Add", mock.Anything).Times(1)

		data := map[int]int{
			51: 11,
		}
		stock := warehouse.NewStockWithData(nil, nil, &b, &ol, data)

		err := stock.PlaceOrder(1, 5)

		assert.NoError(t, err)

		assert.Equal(t, 1, data[51])
		b.AssertExpectations(t)
	})
}

func TestProvision(t *testing.T) {
	t.Run("should reject stock item with non existent type", func(t *testing.T) {
		i := warehouse.MockInventory{}
		i.On("Get", inventory.ID(1)).Return(inventory.Item{})

		data := map[int]int{1: 0}
		stock := warehouse.NewStockWithData(nil, &i, nil, nil, data)

		_, err := stock.Provision(1, 31)

		assert.Equal(t, warehouse.ErrInventoryItemNotFound, err)
		i.AssertExpectations(t)
	})

	t.Run("should place inbound when item type exists", func(t *testing.T) {
		milk := "milk"
		i := warehouse.MockInventory{}
		i.On("Get", inventory.ID(51)).Return(inventory.Item{
			ID:   51,
			Name: inventory.Name(milk),
		})

		inboundLog := warehouse.MockInboundLog{}
		inboundLog.On("Add", mock.Anything, mock.Anything).Times(1)

		data := map[int]int{
			51: 9,
		}
		stock := warehouse.NewStockWithData(&inboundLog, &i, nil, nil, data)

		qty, err := stock.Provision(51, 31)

		assert.NoError(t, err)
		assert.Equal(t, 40, qty)

		inboundLog.AssertExpectations(t)
		i.AssertExpectations(t)
	})
}
