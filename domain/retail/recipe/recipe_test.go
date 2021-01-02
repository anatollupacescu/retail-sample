package recipe_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

func TestDisableRecipe(t *testing.T) {
	t.Run("errors if save fails", func(t *testing.T) {
		db := &recipe.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Save", mock.Anything).Return(recipe.ErrRecipeNotFound)

		r := recipe.Recipe{ID: 1, DB: db}

		err := r.Disable()

		assert.Equal(t, recipe.ErrRecipeNotFound, err)
	})

	t.Run("succeeds", func(t *testing.T) {
		db := &recipe.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Save", mock.Anything).Return(nil)

		r := recipe.Recipe{
			Enabled: true,
			DB:      db,
		}

		err := r.Disable()

		assert.NoError(t, err)
		assert.False(t, r.Enabled)
	})
}

func TestEnableRecipe(t *testing.T) {
	t.Run("errors if save fails", func(t *testing.T) {
		db := &recipe.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Save", mock.Anything).Return(recipe.ErrRecipeNotFound)

		r := recipe.Recipe{DB: db}

		err := r.Enable()

		assert.Equal(t, recipe.ErrRecipeNotFound, err)
	})

	t.Run("succeeds", func(t *testing.T) {
		db := &recipe.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Save", mock.Anything).Return(nil)

		r := recipe.Recipe{
			Enabled: true,
			DB:      db,
		}

		err := r.Enable()

		assert.NoError(t, err)
		assert.True(t, r.Enabled)
	})
}

func TestAddRecipe(t *testing.T) {
	t.Run("errors on empty name", func(t *testing.T) {
		rr := recipe.Recipes{}
		_, err := rr.Add("", nil)
		assert.Equal(t, recipe.ErrEmptyName, err)
	})

	t.Run("errors on empty list of ingredients", func(t *testing.T) {
		rr := recipe.Recipes{}
		_, err := rr.Add("test", nil)
		assert.Equal(t, recipe.ErrNoIngredients, err)
	})

	t.Run("errors on missing quantity", func(t *testing.T) {
		rr := recipe.Recipes{}
		_, err := rr.Add("test", []recipe.InventoryItem{{ID: 1, Qty: 0}})
		assert.Equal(t, recipe.ErrQuantityNotProvided, err)
	})

	t.Run("errors if ingredient name is taken", func(t *testing.T) {
		db := &recipe.MockDB{}
		defer db.AssertExpectations(t)

		item := &recipe.RecipeDTO{ID: 1, Name: "test"}

		db.On("Find", recipe.Name("test")).Return(item, nil)

		recipes := recipe.Recipes{DB: db}
		id, err := recipes.Add("test", []recipe.InventoryItem{{ID: 1, Qty: 2}})

		assert.Zero(t, id)
		assert.Equal(t, recipe.ErrDuplicateName, err)
	})

	t.Run("errors when fails to check that ingredient name is unique", func(t *testing.T) {
		db := &recipe.MockDB{}
		defer db.AssertExpectations(t)

		expectedErr := errors.New("test")
		var nilRecipe *recipe.RecipeDTO
		db.On("Find", recipe.Name("test")).Return(nilRecipe, expectedErr)

		recipes := recipe.Recipes{DB: db}
		id, err := recipes.Add("test", []recipe.InventoryItem{{ID: 1, Qty: 2}})

		assert.Zero(t, id)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("errors if fails to validate ingredients", func(t *testing.T) {
		db := &recipe.MockDB{}
		defer db.AssertExpectations(t)

		var nilRecipe *recipe.RecipeDTO
		db.On("Find", recipe.Name("test")).Return(nilRecipe, recipe.ErrRecipeNotFound)

		mi := &recipe.MockInventory{}
		defer mi.AssertExpectations(t)

		err := recipe.ErrIgredientDisabled
		mi.On("Validate", mock.Anything).Return(err)

		b := recipe.Recipes{DB: db, Inventory: mi}
		id, err := b.Add("test", []recipe.InventoryItem{{ID: 1, Qty: 2}})

		assert.Zero(t, id)
		assert.Equal(t, recipe.ErrIgredientDisabled, err)
	})

	t.Run("errors when persistence fails", func(t *testing.T) {
		db := &recipe.MockDB{}
		defer db.AssertExpectations(t)

		var nilRecipe *recipe.RecipeDTO
		db.On("Find", recipe.Name("test")).Return(nilRecipe, recipe.ErrRecipeNotFound)

		mi := &recipe.MockInventory{}
		defer mi.AssertExpectations(t)

		mi.On("Validate", mock.Anything).Return(nil)

		var expectedErr = errors.New("could not save")
		db.On("Add", mock.Anything).Return(recipe.ID(0), expectedErr)

		b := recipe.Recipes{DB: db, Inventory: mi}
		id, err := b.Add("test", []recipe.InventoryItem{{ID: 1, Qty: 2}})

		assert.Zero(t, id)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("succeeds", func(t *testing.T) {
		db := &recipe.MockDB{}
		defer db.AssertExpectations(t)

		var nilRecipe *recipe.RecipeDTO
		db.On("Find", recipe.Name("test")).Return(nilRecipe, recipe.ErrRecipeNotFound)

		mi := &recipe.MockInventory{}
		defer mi.AssertExpectations(t)

		mi.On("Validate", mock.Anything).Return(nil)

		add := recipe.RecipeDTO{
			Name:        "test",
			Ingredients: []recipe.InventoryItem{{ID: 1, Qty: 2}},
			Enabled:     true,
		}
		db.On("Add", add).Return(recipe.ID(1), nil)

		b := recipe.Recipes{DB: db, Inventory: mi}
		recipeID, err := b.Add("test", []recipe.InventoryItem{{ID: 1, Qty: 2}})

		assert.NoError(t, err)
		assert.Equal(t, recipe.ID(1), recipeID)
	})
}
