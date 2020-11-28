package recipe_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

func TestDisableRecipe(t *testing.T) {
	t.Run("given an unknown id", func(t *testing.T) {
		s := &recipe.MockRecipeStore{}

		var r recipe.Recipe

		s.On("Get", mock.Anything).Return(r, recipe.ErrRecipeNotFound)

		i := &recipe.MockInventory{}
		b := recipe.Book{DB: s, Inventory: i}

		err := b.UpdateStatus(1, false)

		t.Run("calls store", func(t *testing.T) {
			s.AssertExpectations(t)
		})

		t.Run("should return error", func(t *testing.T) {
			assert.Equal(t, recipe.ErrRecipeNotFound, err)
		})
	})

	t.Run("given a known id", func(t *testing.T) {
		s := &recipe.MockRecipeStore{}

		r := recipe.Recipe{
			ID:          1,
			Ingredients: nil,
			Name:        "test",
			Enabled:     true,
		}

		s.On("Get", mock.Anything).Return(r, nil)

		var r2 = r
		r2.Enabled = false

		s.On("Save", r2).Return(nil)

		b := recipe.Book{DB: s}

		err := b.UpdateStatus(1, false)

		t.Run("calls store", func(t *testing.T) {
			s.AssertExpectations(t)
		})

		t.Run("should disable recipe", func(t *testing.T) {
			assert.NoError(t, err)
			assert.Equal(t, false, r2.Enabled)
		})
	})
}

func TestAddRecipe(t *testing.T) {
	t.Run("should reject empty name", func(t *testing.T) {
		b := recipe.Book{}
		_, err := b.Add("", nil)
		assert.Equal(t, recipe.ErrEmptyName, err)
	})

	t.Run("should reject empty list of ingredients", func(t *testing.T) {
		b := recipe.Book{}
		_, err := b.Add("test", nil)
		assert.Equal(t, recipe.ErrNoIngredients, err)
	})

	t.Run("should reject missing quantity", func(t *testing.T) {
		b := recipe.Book{}
		_, err := b.Add("test", []recipe.Ingredient{{ID: 1, Qty: 0}})
		assert.Equal(t, recipe.ErrQuantityNotProvided, err)
	})

	t.Run("should reject if ingredients are disabled", func(t *testing.T) {
		s := &recipe.MockRecipeStore{}

		i := &recipe.MockInventory{}
		b := recipe.Book{DB: s, Inventory: i}

		var item = inventory.Item{
			ID:      1,
			Enabled: false,
			Name:    "test",
		}

		i.On("Get", 1).Return(item, nil)

		_, err := b.Add("test", []recipe.Ingredient{{ID: 1, Qty: 2}})

		assert.Equal(t, recipe.ErrIgredientDisabled, err)

		s.AssertExpectations(t)
		i.AssertExpectations(t)
	})

	t.Run("should reject if incredients are not present in inventory", func(t *testing.T) {
		s := &recipe.MockRecipeStore{}

		i := &recipe.MockInventory{}
		b := recipe.Book{DB: s, Inventory: i}

		var zeroInventoryItem inventory.Item
		i.On("Get", 1).Return(zeroInventoryItem, inventory.ErrItemNotFound)

		_, err := b.Add("test", []recipe.Ingredient{{ID: 1, Qty: 2}})

		assert.Equal(t, recipe.ErrIgredientNotFound, err)

		s.AssertExpectations(t)
		i.AssertExpectations(t)
	})

	t.Run("should propagate downstream failure", func(t *testing.T) {
		s := &recipe.MockRecipeStore{}
		i := &recipe.MockInventory{}
		b := recipe.Book{DB: s, Inventory: i}

		i.On("Get", 1).Return(inventory.Item{
			ID:      1,
			Enabled: true,
		}, nil)

		var expectedErr = errors.New("could not save")
		s.On("Add", mock.Anything).Return(recipe.ID(0), expectedErr)

		_, err := b.Add("test", []recipe.Ingredient{{ID: 1, Qty: 2}})

		assert.Equal(t, expectedErr, err)

		s.AssertExpectations(t)
		i.AssertExpectations(t)
	})

	t.Run("should succeed with correct name and components", func(t *testing.T) {
		s := &recipe.MockRecipeStore{}
		i := &recipe.MockInventory{}
		b := recipe.Book{DB: s, Inventory: i}

		i.On("Get", 1).Return(inventory.Item{
			ID:      1,
			Enabled: true,
		}, nil)

		s.On("Add", recipe.Recipe{
			Name:        "test",
			Ingredients: []recipe.Ingredient{{ID: 1, Qty: 2}},
			Enabled:     true,
		}).Return(recipe.ID(1), nil)

		recipeID, err := b.Add("test", []recipe.Ingredient{{ID: 1, Qty: 2}})

		assert.NoError(t, err)
		assert.Equal(t, recipe.ID(1), recipeID)

		s.AssertExpectations(t)
		i.AssertExpectations(t)
	})
}
