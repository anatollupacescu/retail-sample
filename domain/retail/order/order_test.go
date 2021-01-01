package order_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/domain/retail/order"
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

func TestAdd(t *testing.T) {
	t.Run("errors when quantity is zero", func(t *testing.T) {
		orders := order.Orders{}

		_, err := orders.Add(1, 0)

		assert.Equal(t, order.ErrInvalidQuantity, err)
	})

	t.Run("errors when quantity is negative", func(t *testing.T) {
		orders := order.Orders{}

		_, err := orders.Add(1, -1)

		assert.Equal(t, order.ErrInvalidQuantity, err)
	})

	t.Run("errors when get recipe fails", func(t *testing.T) {
		recipeDB := &recipe.MockDB{}
		defer recipeDB.AssertExpectations(t)

		expectedErr := errors.New("not found")
		recipeDB.On("Get", recipe.ID(1)).Return(recipe.RecipeDTO{}, expectedErr)

		orders := order.Orders{Recipes: recipeDB}

		receivedID, err := orders.Add(1, 1)

		assert.Equal(t, expectedErr, err)
		assert.Equal(t, order.ID(0), receivedID)
	})

	t.Run("errors when recipe is disabled", func(t *testing.T) {
		var recipeID = recipe.ID(1)

		var r = recipe.RecipeDTO{
			ID:      recipeID,
			Enabled: false,
		}

		recipeDB := &recipe.MockDB{}
		defer recipeDB.AssertExpectations(t)

		recipeDB.On("Get", recipeID).Return(r, nil)

		orders := order.Orders{Recipes: recipeDB}

		receivedID, err := orders.Add(1, 1)

		assert.Equal(t, order.ErrInvalidRecipe, err)
		assert.Equal(t, order.ID(0), receivedID)
	})

	t.Run("errors when extracting from stock position fails", func(t *testing.T) {
		var recipeID = recipe.ID(1)

		var r = recipe.RecipeDTO{
			ID: recipeID,
			Ingredients: []recipe.InventoryItem{
				{ID: 1},
			},
			Name:    "test",
			Enabled: true,
		}

		recipeDB := &recipe.MockDB{}
		defer recipeDB.AssertExpectations(t)

		recipeDB.On("Get", recipeID).Return(r, nil)

		expectedErr := errors.New("expected")

		stockDB := &order.MockStock{}
		defer stockDB.AssertExpectations(t)

		stockDB.On("Extract", mock.Anything, mock.Anything).Return(expectedErr)

		orders := order.Orders{Recipes: recipeDB, Stock: stockDB}

		receivedID, err := orders.Add(1, 1)

		assert.Equal(t, expectedErr, err)
		assert.Zero(t, receivedID)
	})

	t.Run("errors when save to DB fails", func(t *testing.T) {
		var recipeID = recipe.ID(1)

		var r = recipe.RecipeDTO{
			ID: recipeID,
			Ingredients: []recipe.InventoryItem{
				{ID: 1},
			},
			Name:    "test",
			Enabled: true,
		}

		recipeDB := &recipe.MockDB{}
		defer recipeDB.AssertExpectations(t)

		recipeDB.On("Get", recipeID).Return(r, nil)

		stockDB := &order.MockStock{}
		defer stockDB.AssertExpectations(t)

		stockDB.On("Extract", mock.Anything, mock.Anything).Return(nil)

		db := &order.MockDB{}
		defer db.AssertExpectations(t)

		var dbErr = errors.New("test")
		db.On("Add", mock.Anything).Return(order.ID(0), dbErr)

		orders := order.Orders{DB: db, Recipes: recipeDB, Stock: stockDB}

		receivedID, err := orders.Add(1, 1)

		assert.Equal(t, dbErr, err)
		assert.Zero(t, receivedID)
	})

	t.Run("when selling ingredients succeeds", func(t *testing.T) {
		var recipeID = recipe.ID(1)

		var r = recipe.RecipeDTO{
			ID: recipeID,
			Ingredients: []recipe.InventoryItem{
				{ID: 1},
			},
			Name:    "test",
			Enabled: true,
		}

		recipeDB := &recipe.MockDB{}
		defer recipeDB.AssertExpectations(t)

		recipeDB.On("Get", recipeID).Return(r, nil)

		stockDB := &order.MockStock{}
		defer stockDB.AssertExpectations(t)

		stockDB.On("Extract", mock.Anything, mock.Anything).Return(nil)

		db := &order.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Add", mock.Anything).Return(order.ID(1), nil)

		orders := order.Orders{DB: db, Recipes: recipeDB, Stock: stockDB}

		receivedID, err := orders.Add(1, 1)

		assert.NoError(t, err)
		assert.Equal(t, order.ID(1), receivedID)
	})
}
