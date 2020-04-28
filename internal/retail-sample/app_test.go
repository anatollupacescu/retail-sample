package retailsample_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/order"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"

	retail "github.com/anatollupacescu/retail-sample/internal/retail-sample"

	"github.com/stretchr/testify/assert"
)

func TestPlaceOrder(t *testing.T) {
	t.Skip("needs mocks")

	t.Run("should reject order with non existent recipe id", func(t *testing.T) {
		b := retail.MockRecipeBook{}

		var zeroRecipe recipe.Recipe
		b.On("Get", recipe.ID(1)).Return(zeroRecipe)

		app := retail.App{
			RecipeBook: &b,
		}

		entryID, err := app.PlaceOrder(1, 10)

		assert.Zero(t, entryID)
		assert.Equal(t, retail.ErrRecipeNotFound, err)
		b.AssertExpectations(t)
	})

	t.Run("should reject outbound when not enough stock", func(t *testing.T) {
		b := retail.MockRecipeBook{}

		b.On("Get", recipe.ID(1)).Return(recipe.Recipe{
			Name:        "test",
			Ingredients: []recipe.Ingredient{{ID: 51, Qty: 2}},
		})

		// data := map[int]int{
		// 	51: 9,
		// }

		app := retail.App{
			RecipeBook: &b,
			// Stock:      warehouse.NewStockWithData(data),
		}

		entryID, err := app.PlaceOrder(1, 5)

		assert.Zero(t, entryID)
		assert.Equal(t, retail.ErrNotEnoughStock, err)
		b.AssertExpectations(t)
	})

	t.Run("should update inventory on success", func(t *testing.T) {
		b := retail.MockRecipeBook{}

		b.On("Get", recipe.ID(1)).Return(recipe.Recipe{
			Name:        recipe.Name("test"),
			Ingredients: []recipe.Ingredient{{ID: 51, Qty: 2}},
		})

		ol := retail.MockOrders{}
		ol.On("Add", mock.Anything).Return(order.ID(1))

		// data := map[int]int{
		// 	51: 11,
		// }

		app := retail.App{
			RecipeBook: &b,
			Orders:     &ol,
			// Stock:      warehouse.NewStockWithData(data),
		}

		entryID, err := app.PlaceOrder(1, 5)

		assert.NoError(t, err)
		assert.Equal(t, order.ID(1), entryID)
		// assert.Equal(t, 1, data[51])
		b.AssertExpectations(t)
	})
}

func TestProvision(t *testing.T) {
	t.SkipNow()

	t.Run("should reject stock item with non existent type", func(t *testing.T) {
		i := retail.MockInventory{}
		i.On("Get", inventory.ID(1)).Return(inventory.Item{})

		// data := map[int]int{1: 0}

		app := retail.App{
			Inventory: &i,
			// Stock:     warehouse.NewStockWithData(data),
		}

		_, err := app.Provision([]retail.ProvisionEntry{{ID: 1, Qty: 31}})

		assert.Equal(t, inventory.ErrInventoryItemNotFound, err)
		i.AssertExpectations(t)
	})

	t.Run("should place inbound when item type exists", func(t *testing.T) {
		milk := "milk"
		i := retail.MockInventory{}
		i.On("Get", inventory.ID(51)).Return(inventory.Item{
			ID:   51,
			Name: inventory.Name(milk),
		})

		provisionLog := retail.MockInboundLog{}
		provisionLog.On("Add", mock.Anything, mock.Anything).Times(1)

		// data := map[int]int{
		// 	51: 9,
		// }

		app := retail.App{
			Inventory:    &i,
			ProvisionLog: &provisionLog,
			// Stock:        warehouse.NewStockWithData(data),
		}

		qty, err := app.Provision([]retail.ProvisionEntry{{ID: 51, Qty: 31}})

		assert.NoError(t, err)
		assert.Equal(t, 40, qty)

		provisionLog.AssertExpectations(t)
		i.AssertExpectations(t)
	})
}
