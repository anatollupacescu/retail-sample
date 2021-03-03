package stock_test

import (
	"testing"

	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProvision(t *testing.T) {
	var (
		db *stock.MockDB

		reset = func() {
			db = &stock.MockDB{}
		}
		expectOnSave = func(err error) {
			db.On("Save", mock.Anything).Return(err)
		}
	)
	t.Run("given quantity is negative", func(t *testing.T) {
		st := new(stock.Position)
		err := st.Provision(-12)
		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, stock.ErrInvalidProvisionQuantity, err)
			assert.Zero(t, st.Qty)
		})
	})
	t.Run("given position updated", func(t *testing.T) {
		reset()
		expectOnSave(nil)

		st := &stock.Position{InventoryID: 1, Qty: 1, DB: db}
		err := st.Provision(10)

		t.Run("assert success", func(t *testing.T) {
			assert.NoError(t, err)
			assert.Equal(t, 11, st.Qty)
			db.AssertExpectations(t)
		})
	})
	t.Run("given failure to update quantity", func(t *testing.T) {
		reset()

		expectedErr := errors.New("err")
		expectOnSave(expectedErr)

		st := &stock.Position{Qty: 4, DB: db}
		err := st.Provision(10)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, expectedErr, err)
			assert.Equal(t, 4, st.Qty)
			db.AssertExpectations(t)
		})
	})
}

func TestExtract(t *testing.T) {
	var (
		recipeDB *stock.MockRecipeDB
		stockDB  *stock.MockDB
		reset    = func() {
			recipeDB = &stock.MockRecipeDB{}
			stockDB = &stock.MockDB{}
		}
		newExtractor = func() stock.Extractor {
			return stock.Extractor{
				Recipes: recipeDB,
				Stock:   stockDB,
			}
		}
		expectRecipeErr = func(id int, err error) {
			recipeDB.On("Get", id).Return(recipe.DTO{}, err)
		}
		expectRecipeInvalid = func(id int) {
			ingredients := []recipe.InventoryItem{{ID: id, Qty: 1}}
			recipeDB.On("Get", 1).Return(recipe.DTO{Enabled: false, Ingredients: ingredients}, nil)
		}
		expectRecipeValid = func(id int) {
			ingredients := []recipe.InventoryItem{{ID: id, Qty: 1}}
			recipeDB.On("Get", 1).Return(recipe.DTO{Enabled: true, Ingredients: ingredients}, nil)
		}
		expectGetStockErr = func(err error) {
			stockDB.On("Get", mock.Anything).Return(stock.PositionDTO{}, err)
		}
		expectGetStockOK = func(qty int) {
			stockDB.On("Get", mock.Anything).Return(stock.PositionDTO{InventoryID: 1, Qty: qty}, nil)
		}
		expectSavePosOK = func(qty int) {
			stockDB.On("Save", stock.PositionDTO{InventoryID: 1, Qty: qty}).Return(nil)
		}
		expectSaveErr = func(err error) {
			stockDB.On("Save", mock.Anything).Return(err)
		}
	)
	t.Run("given recipe not found", func(t *testing.T) {
		reset()

		id := 1
		expectRecipeErr(id, recipe.ErrNotFound)

		err := newExtractor().Extract(id, 1)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, recipe.ErrNotFound, err)
			recipeDB.AssertExpectations(t)
			stockDB.AssertExpectations(t)
		})
	})
	t.Run("given recipe is invalid", func(t *testing.T) {
		reset()

		id := 1
		expectRecipeInvalid(id)

		err := newExtractor().Extract(id, 1)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, stock.ErrRecipeDisabled, err)
			recipeDB.AssertExpectations(t)
			stockDB.AssertExpectations(t)
		})
	})
	t.Run("given failure to check recipe", func(t *testing.T) {
		reset()

		id := 1
		expectedErr := errors.New("db err")
		expectRecipeErr(id, expectedErr)

		err := newExtractor().Extract(id, 1)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, expectedErr, err)
			recipeDB.AssertExpectations(t)
			stockDB.AssertExpectations(t)
		})
	})
	t.Run("given item not present in stock", func(t *testing.T) {
		reset()

		id := 1
		expectRecipeValid(id)
		expectGetStockErr(stock.ErrPositionNotFound)

		err := newExtractor().Extract(id, 1)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, stock.ErrNotEnoughStock, err)
			recipeDB.AssertExpectations(t)
			stockDB.AssertExpectations(t)
		})
	})
	t.Run("given failure to check for stock item", func(t *testing.T) {
		reset()

		id := 1
		expectRecipeValid(id)
		expectedErr := errors.New("db err")
		expectGetStockErr(expectedErr)

		err := newExtractor().Extract(id, 1)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, expectedErr, err)
			recipeDB.AssertExpectations(t)
			stockDB.AssertExpectations(t)
		})
	})
	t.Run("given quantity is negative", func(t *testing.T) {
		reset()

		id := 1
		expectRecipeValid(id)
		qty := 23
		expectGetStockOK(qty)

		err := newExtractor().Extract(id, -1)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, stock.ErrInvalidExtractQuantity, err)
			recipeDB.AssertExpectations(t)
			stockDB.AssertExpectations(t)
		})
	})
	t.Run("given not enough stock for item", func(t *testing.T) {
		reset()

		id := 1
		expectRecipeValid(id)
		expectGetStockOK(0)

		err := newExtractor().Extract(id, 1)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, stock.ErrNotEnoughStock, err)
			recipeDB.AssertExpectations(t)
			stockDB.AssertExpectations(t)
		})
	})
	t.Run("given position updated", func(t *testing.T) {
		reset()

		id := 1
		expectRecipeValid(id)
		expectGetStockOK(3)
		expectSavePosOK(2) // stock value (3) substract (1)

		err := newExtractor().Extract(id, 1)

		t.Run("assert success", func(t *testing.T) {
			assert.NoError(t, err)
			recipeDB.AssertExpectations(t)
			stockDB.AssertExpectations(t)
		})
	})
	t.Run("given failure to update position quantity", func(t *testing.T) {
		reset()

		id := 1
		expectRecipeValid(id)
		expectGetStockOK(3)
		expectedErr := errors.New("db err")
		expectSaveErr(expectedErr)

		err := newExtractor().Extract(id, 1)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, expectedErr, err)
			recipeDB.AssertExpectations(t)
			stockDB.AssertExpectations(t)
		})
	})
}
