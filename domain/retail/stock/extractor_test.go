package stock_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

func TestExtractor(t *testing.T) {
	t.Run("given a non existent recipe", func(t *testing.T) {
		recipes := &stock.MockRecipeDB{}
		defer recipes.AssertExpectations(t)

		recipes.On("Get", 1).Return(recipe.DTO{}, recipe.ErrNotFound)

		e := stock.Extractor{
			Recipes: recipes,
		}

		err := e.Extract(1, 1)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, recipe.ErrNotFound, err)
		})
	})

	t.Run("given a recipe with a non existent ingredient", func(t *testing.T) {
		recipes := &stock.MockRecipeDB{}
		defer recipes.AssertExpectations(t)

		ingredients := []recipe.InventoryItem{{
			ID:  1,
			Qty: 1,
		}}

		recipes.On("Get", 1).Return(recipe.DTO{Ingredients: ingredients}, nil)

		db := &stock.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Get", 1).Return(stock.PositionDTO{}, stock.ErrPositionNotFound)

		e := stock.Extractor{
			Recipes: recipes,
			Stock:   db,
		}

		err := e.Extract(1, 1)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, stock.ErrNotEnoughStock, err)
		})
	})

	t.Run("given a recipe with a bad ingredient", func(t *testing.T) {
		recipes := &stock.MockRecipeDB{}
		defer recipes.AssertExpectations(t)

		ingredients := []recipe.InventoryItem{{
			ID:  1,
			Qty: 1,
		}}

		recipes.On("Get", 1).Return(recipe.DTO{Ingredients: ingredients}, nil)

		db := &stock.MockDB{}
		defer db.AssertExpectations(t)

		expectedErr := errors.New("unexpected")

		db.On("Get", 1).Return(stock.PositionDTO{}, expectedErr)

		e := stock.Extractor{
			Recipes: recipes,
			Stock:   db,
		}

		err := e.Extract(1, 1)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("given a recipe with an ingredient we can not extract", func(t *testing.T) {
		recipes := &stock.MockRecipeDB{}
		defer recipes.AssertExpectations(t)

		ingredients := []recipe.InventoryItem{{
			ID:  1,
			Qty: 1,
		}}

		recipes.On("Get", 1).Return(recipe.DTO{Ingredients: ingredients}, nil)

		db := &stock.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Get", 1).Return(stock.PositionDTO{Qty: 1}, nil)

		expectedErr := errors.New("unexpected")

		db.On("Save", mock.Anything).Return(expectedErr)

		e := stock.Extractor{
			Recipes: recipes,
			Stock:   db,
		}

		err := e.Extract(1, 1)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, expectedErr, err)
		})
	})

	t.Run("given a recipe with two good ingredients", func(t *testing.T) {
		recipes := &stock.MockRecipeDB{}
		defer recipes.AssertExpectations(t)

		ingredients := []recipe.InventoryItem{{
			ID:  1,
			Qty: 2,
		}, {
			ID:  2,
			Qty: 3,
		}}

		recipes.On("Get", 1).Return(recipe.DTO{Ingredients: ingredients}, nil)

		db := &stock.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Get", 1).Return(stock.PositionDTO{InventoryID: 1, Qty: 4}, nil)
		db.On("Get", 2).Return(stock.PositionDTO{InventoryID: 2, Qty: 6}, nil)

		db.On("Save", stock.PositionDTO{InventoryID: 1, Qty: 0}).Return(nil)
		db.On("Save", stock.PositionDTO{InventoryID: 2, Qty: 0}).Return(nil)

		e := stock.Extractor{
			Recipes: recipes,
			Stock:   db,
		}

		err := e.Extract(1, 2)

		t.Run("assert both positions are updated", func(t *testing.T) {
			assert.NoError(t, err)
		})
	})
}
