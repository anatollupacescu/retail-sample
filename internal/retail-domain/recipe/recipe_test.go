package recipe_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
)

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

	t.Run("should return error if incredients are no present in inventory", func(t *testing.T) {
		s := &recipe.MockRecipeStore{}

		i := &recipe.MockInventory{}
		b := recipe.Book{Store: s, Inventory: i}

		var zeroInventoryItem inventory.Item
		i.On("Get", inventory.ID(1)).Return(zeroInventoryItem)

		_, err := b.Add("test", []recipe.Ingredient{{ID: 1, Qty: 2}})

		assert.Equal(t, recipe.ErrIgredientNotFound, err)

		s.AssertExpectations(t)
		i.AssertExpectations(t)
	})

	t.Run("should propagate downstream failure", func(t *testing.T) {
		s := &recipe.MockRecipeStore{}
		i := &recipe.MockInventory{}
		b := recipe.Book{Store: s, Inventory: i}

		i.On("Get", inventory.ID(1)).Return(inventory.Item{
			ID: 1,
		})

		var expectedErr = errors.New("could not save")
		s.On("add", mock.Anything).Return(recipe.ID(0), expectedErr)

		_, err := b.Add("test", []recipe.Ingredient{{ID: 1, Qty: 2}})

		assert.Equal(t, expectedErr, err)

		s.AssertExpectations(t)
		i.AssertExpectations(t)
	})

	t.Run("should succeed with correct name and components", func(t *testing.T) {
		s := &recipe.MockRecipeStore{}
		i := &recipe.MockInventory{}
		b := recipe.Book{Store: s, Inventory: i}

		i.On("Get", inventory.ID(1)).Return(inventory.Item{
			ID: 1,
		})
		s.On("add", recipe.Recipe{
			Name:        "test",
			Ingredients: []recipe.Ingredient{{ID: 1, Qty: 2}},
		}).Return(recipe.ID(1), nil)

		recipeID, err := b.Add("test", []recipe.Ingredient{{ID: 1, Qty: 2}})

		assert.NoError(t, err)
		assert.Equal(t, recipe.ID(1), recipeID)

		s.AssertExpectations(t)
		i.AssertExpectations(t)
	})
}

func TestGetRecipe(t *testing.T) {

	t.Run("should return zero value for non existent", func(t *testing.T) {
		s := &recipe.MockRecipeStore{}
		b := recipe.Book{Store: s}

		var zeroValueRecipe = recipe.Recipe{}
		s.On("get", recipe.ID(1)).Return(zeroValueRecipe)

		r := b.Get(1)
		assert.Equal(t, r, zeroValueRecipe)

		s.AssertExpectations(t)
	})

	t.Run("should fetch stored recipe for valid id", func(t *testing.T) {
		s := &recipe.MockRecipeStore{}
		b := recipe.Book{Store: s}

		var foundRecipe = recipe.Recipe{
			Name: "test",
			Ingredients: []recipe.Ingredient{{
				ID:  1,
				Qty: 2,
			}},
		}
		s.On("get", recipe.ID(1)).Return(foundRecipe)

		r := b.Get(1)
		assert.Equal(t, r, foundRecipe)

		s.AssertExpectations(t)
	})
}

func TestGetRecipeNames(t *testing.T) {
	s := &recipe.MockRecipeStore{}
	b := recipe.Book{Store: s}

	s.On("all").Return([]recipe.Recipe{{
		Name: recipe.Name("glintwine"),
	}})

	r := b.Names()
	assert.Equal(t, r, []recipe.Name{recipe.Name("glintwine")})

	s.AssertExpectations(t)
}
