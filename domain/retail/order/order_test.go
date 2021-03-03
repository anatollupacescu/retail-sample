package order_test

import (
	"errors"
	"testing"

	"github.com/anatollupacescu/retail-sample/domain/retail/order"
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOrder(t *testing.T) {
	t.Run("given quantity is invalid", func(t *testing.T) {
		// test for zero and negative value
		var err error
		t.Run("errors when quantity is zero", func(t *testing.T) {
			orders := order.Orders{}

			_, err = orders.Create(1, 0)

			t.Run("assert error", func(t *testing.T) {
				assert.Equal(t, order.ErrInvalidQuantity, err)
			})
		})

		t.Run("errors when quantity is negative", func(t *testing.T) {
			orders := order.Orders{}

			_, err = orders.Create(1, -1)

			t.Run("assert error", func(t *testing.T) {
				assert.Equal(t, order.ErrInvalidQuantity, err)
			})
		})
	})
	t.Run("given recipe is not valid", func(t *testing.T) {
		recipeDB := &order.MockRecipe{}
		defer recipeDB.AssertExpectations(t)

		expectedErr := recipe.ErrDisabled
		recipeDB.On("Valid", 1).Return(expectedErr)

		orders := order.Orders{Recipes: recipeDB}

		receivedID, err := orders.Create(1, 1)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, expectedErr, err)
			assert.Equal(t, 0, receivedID)
		})
	})
	t.Run("given fail to check recipe validity", func(t *testing.T) {
		recipe := &order.MockRecipe{}
		defer recipe.AssertExpectations(t)

		expectedErr := errors.New("not found")
		recipe.On("Valid", 1).Return(expectedErr)

		orders := order.Orders{Recipes: recipe}

		receivedID, err := orders.Create(1, 1)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, expectedErr, err)
			assert.Equal(t, 0, receivedID)
		})
	})
	t.Run("given stock is updated", func(t *testing.T) {
		var recipeID = 1

		recipeDB := &order.MockRecipe{}
		defer recipeDB.AssertExpectations(t)

		recipeDB.On("Valid", recipeID).Return(nil)

		stockDB := &order.MockStock{}
		defer stockDB.AssertExpectations(t)

		stockDB.On("Extract", mock.Anything, mock.Anything).Return(nil)

		db := &order.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Add", mock.Anything).Return(1, nil)

		orders := order.Orders{DB: db, Recipes: recipeDB, Stock: stockDB}

		receivedID, err := orders.Create(1, 1)

		t.Run("assert success", func(t *testing.T) {
			assert.NoError(t, err)
			assert.Equal(t, 1, receivedID)
		})
	})
	t.Run("given fail to update stock", func(t *testing.T) {
		var recipeID = 1

		recipe := &order.MockRecipe{}
		defer recipe.AssertExpectations(t)

		recipe.On("Valid", recipeID).Return(nil)

		expectedErr := errors.New("expected")

		stock := &order.MockStock{}
		defer stock.AssertExpectations(t)

		stock.On("Extract", mock.Anything, mock.Anything).Return(expectedErr)

		orders := order.Orders{Recipes: recipe, Stock: stock}

		receivedID, err := orders.Create(1, 1)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, expectedErr, err)
			assert.Zero(t, receivedID)
		})
	})
	t.Run("given order created", func(t *testing.T) {
		var recipeID = 1

		recipeDB := &order.MockRecipe{}
		defer recipeDB.AssertExpectations(t)

		recipeDB.On("Valid", recipeID).Return(nil)

		stockDB := &order.MockStock{}
		defer stockDB.AssertExpectations(t)

		stockDB.On("Extract", mock.Anything, mock.Anything).Return(nil)

		db := &order.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Add", mock.Anything).Return(0, nil)

		orders := order.Orders{DB: db, Recipes: recipeDB, Stock: stockDB}

		receivedID, err := orders.Create(1, 1)

		t.Run("assert success", func(t *testing.T) {
			assert.NoError(t, err)
			assert.Zero(t, receivedID)
		})
	})
	t.Run("given fail to save order", func(t *testing.T) {
		var recipeID = 1

		recipeDB := &order.MockRecipe{}
		defer recipeDB.AssertExpectations(t)

		recipeDB.On("Valid", recipeID).Return(nil)

		stockDB := &order.MockStock{}
		defer stockDB.AssertExpectations(t)

		stockDB.On("Extract", mock.Anything, mock.Anything).Return(nil)

		db := &order.MockDB{}
		defer db.AssertExpectations(t)

		var dbErr = errors.New("test")
		db.On("Add", mock.Anything).Return(0, dbErr)

		orders := order.Orders{DB: db, Recipes: recipeDB, Stock: stockDB}

		receivedID, err := orders.Create(1, 1)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, dbErr, err)
			assert.Zero(t, receivedID)
		})
	})
}
