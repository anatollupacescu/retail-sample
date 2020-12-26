package order_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/domain/retail/order"
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

func TestPlaceOrder(t *testing.T) {
	t.Run("when quantity is invalid", func(t *testing.T) {
		orders := order.Orders{}

		_, err := orders.PlaceOrder(1, 0)

		t.Run("propagates error", func(t *testing.T) {
			assert.Equal(t, order.ErrInvalidQuantity, err)
		})
	})

	t.Run("when recipe not found", func(t *testing.T) {
		var recipeID = recipe.ID(1)

		var zeroRecipe recipe.Recipe

		expectedErr := errors.New("not found")

		recipeBook := &order.MockRecipeBook{}
		recipeBook.On("Get", recipeID).Return(zeroRecipe, expectedErr)

		orders := order.Orders{RecipeBook: recipeBook}

		receivedID, err := orders.PlaceOrder(1, 1)

		t.Run("calls recipe book", func(t *testing.T) {
			recipeBook.AssertExpectations(t)
		})

		t.Run("propagates the error", func(t *testing.T) {
			assert.Equal(t, expectedErr, err)
			assert.Equal(t, order.ID(0), receivedID)
		})
	})

	t.Run("when recipe is disabled", func(t *testing.T) {
		var recipeID = recipe.ID(1)

		var r = recipe.Recipe{
			ID:      recipeID,
			Enabled: false,
		}

		recipeBook := &order.MockRecipeBook{}
		recipeBook.On("Get", recipeID).Return(r, nil)

		orders := order.Orders{RecipeBook: recipeBook}

		receivedID, err := orders.PlaceOrder(1, 1)

		t.Run("calls recipe book", func(t *testing.T) {
			recipeBook.AssertExpectations(t)
		})

		t.Run("propagates the error", func(t *testing.T) {
			assert.Equal(t, order.ErrInvalidRecipe, err)
			assert.Equal(t, order.ID(0), receivedID)
		})
	})

	t.Run("when selling ingredients fails", func(t *testing.T) {
		var recipeID = recipe.ID(1)

		var r = recipe.Recipe{
			ID:          recipeID,
			Ingredients: []recipe.Ingredient{},
			Name:        "test",
			Enabled:     true,
		}

		recipeBook := &order.MockRecipeBook{}
		recipeBook.On("Get", recipeID).Return(r, nil)

		expectedErr := errors.New("expected")

		mockStock := &order.MockStock{}
		mockStock.On("Sell", mock.Anything, 1).Return(expectedErr)

		orders := order.Orders{RecipeBook: recipeBook, Stock: mockStock}

		receivedID, err := orders.PlaceOrder(1, 1)

		t.Run("makes the expected calls", func(t *testing.T) {
			recipeBook.AssertExpectations(t)
			mockStock.AssertExpectations(t)
		})

		t.Run("propagates the error", func(t *testing.T) {
			assert.Equal(t, expectedErr, err)
			assert.Equal(t, order.ID(0), receivedID)
		})
	})

	t.Run("when call to store fails", func(t *testing.T) {
		var recipeID = recipe.ID(1)

		var r = recipe.Recipe{
			ID:          recipeID,
			Ingredients: []recipe.Ingredient{},
			Name:        "test",
			Enabled:     true,
		}

		recipeBook := &order.MockRecipeBook{}
		recipeBook.On("Get", recipeID).Return(r, nil)

		mockStock := &order.MockStock{}
		mockStock.On("Sell", mock.Anything, 1).Return(nil)

		var zeroOrderID order.ID
		expectedErr := errors.New("expected")

		store := &order.MockOrderStore{}
		store.On("Add", mock.Anything).Return(zeroOrderID, expectedErr)

		orders := order.Orders{DB: store, RecipeBook: recipeBook, Stock: mockStock}

		receivedID, err := orders.PlaceOrder(1, 1)

		t.Run("makes the expected calls", func(t *testing.T) {
			recipeBook.AssertExpectations(t)
			mockStock.AssertExpectations(t)
			store.AssertExpectations(t)
		})

		t.Run("throws error", func(t *testing.T) {
			assert.Equal(t, expectedErr, err)
			assert.Zero(t, receivedID)
		})
	})

	t.Run("when selling ingredients succeeds", func(t *testing.T) {
		var recipeID = recipe.ID(1)

		var r = recipe.Recipe{
			ID:          recipeID,
			Ingredients: []recipe.Ingredient{},
			Name:        "test",
			Enabled:     true,
		}

		recipeBook := &order.MockRecipeBook{}
		recipeBook.On("Get", recipeID).Return(r, nil)

		mockStock := &order.MockStock{}
		mockStock.On("Sell", mock.Anything, 1).Return(nil)

		store := &order.MockOrderStore{}
		store.On("Add", mock.Anything).Return(order.ID(1), nil)

		orders := order.Orders{DB: store, RecipeBook: recipeBook, Stock: mockStock}

		receivedID, err := orders.PlaceOrder(1, 1)

		t.Run("makes the expected calls", func(t *testing.T) {
			recipeBook.AssertExpectations(t)
			mockStock.AssertExpectations(t)
			store.AssertExpectations(t)
		})

		t.Run("does not throw error", func(t *testing.T) {
			assert.NoError(t, err)
			assert.Equal(t, order.ID(1), receivedID)
		})
	})
}
