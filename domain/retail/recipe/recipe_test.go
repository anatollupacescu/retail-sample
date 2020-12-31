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
	t.Run("propagates error from persistence", func(t *testing.T) {
		s := &recipe.MockDB{}
		defer s.AssertExpectations(t)

		s.On("Save", mock.Anything).Return(recipe.ErrRecipeNotFound)

		r := recipe.Recipe{ID: 1, DB: s}

		err := r.Disable()

		assert.Equal(t, recipe.ErrRecipeNotFound, err)
	})

	t.Run("succeeds", func(t *testing.T) {
		s := &recipe.MockDB{}
		defer s.AssertExpectations(t)

		s.On("Save", mock.Anything).Return(nil)

		r := recipe.Recipe{
			Enabled: true,
			DB:      s,
		}

		err := r.Disable()

		assert.NoError(t, err)
		assert.False(t, r.Enabled)
	})
}

func TestEnableRecipe(t *testing.T) {
	t.Run("propagates error from persistence", func(t *testing.T) {
		s := &recipe.MockDB{}
		defer s.AssertExpectations(t)

		s.On("Save", mock.Anything).Return(recipe.ErrRecipeNotFound)

		r := recipe.Recipe{DB: s}

		err := r.Disable()

		assert.Equal(t, recipe.ErrRecipeNotFound, err)
	})

	t.Run("succeeds", func(t *testing.T) {
		s := &recipe.MockDB{}
		defer s.AssertExpectations(t)

		s.On("Save", mock.Anything).Return(nil)

		r := recipe.Recipe{
			Enabled: true,
			DB:      s,
		}

		err := r.Enable()

		assert.NoError(t, err)
		assert.True(t, r.Enabled)
	})
}

func TestAddRecipe(t *testing.T) {
	t.Run("errors on empty name", func(t *testing.T) {
		b := recipe.Collection{}
		_, err := b.Add("", nil)
		assert.Equal(t, recipe.ErrEmptyName, err)
	})

	t.Run("errors on empty list of ingredients", func(t *testing.T) {
		b := recipe.Collection{}
		_, err := b.Add("test", nil)
		assert.Equal(t, recipe.ErrNoIngredients, err)
	})

	t.Run("errors on missing quantity", func(t *testing.T) {
		b := recipe.Collection{}
		_, err := b.Add("test", []recipe.Ingredient{{ID: 1, Qty: 0}})
		assert.Equal(t, recipe.ErrQuantityNotProvided, err)
	})

	t.Run("errors if ingredient name is taken", func(t *testing.T) {
		s := &recipe.MockDB{}
		defer s.AssertExpectations(t)

		item := &recipe.Recipe{ID: 1, Name: "test"}

		s.On("Find", recipe.Name("test")).Return(item, nil)

		recipes := recipe.Collection{DB: s}
		id, err := recipes.Add("test", []recipe.Ingredient{{ID: 1, Qty: 2}})

		assert.Zero(t, id)
		assert.Equal(t, recipe.ErrDuplicateName, err)
	})

	t.Run("errors when fails to check ingredient name", func(t *testing.T) {
		s := &recipe.MockDB{}
		defer s.AssertExpectations(t)

		expectedErr := errors.New("test")
		var nilRecipe *recipe.Recipe
		s.On("Find", recipe.Name("test")).Return(nilRecipe, expectedErr)

		recipes := recipe.Collection{DB: s}
		id, err := recipes.Add("test", []recipe.Ingredient{{ID: 1, Qty: 2}})

		assert.Zero(t, id)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("errors if an ingredient is disabled", func(t *testing.T) {
		s := &recipe.MockDB{}
		defer s.AssertExpectations(t)

		i := &inventory.MockDB{}
		defer i.AssertExpectations(t)

		var nilRecipe *recipe.Recipe
		s.On("Find", recipe.Name("test")).Return(nilRecipe, recipe.ErrRecipeNotFound)

		var item = inventory.ItemDTO{
			ID:      1,
			Enabled: false,
			Name:    "test",
		}

		i.On("Get", 1).Return(item, nil)

		b := recipe.Collection{DB: s, Inventory: i}
		id, err := b.Add("test", []recipe.Ingredient{{ID: 1, Qty: 2}})

		assert.Zero(t, id)
		assert.Equal(t, recipe.ErrIgredientDisabled, err)
	})

	t.Run("errors if an incredient is not found", func(t *testing.T) {
		s := &recipe.MockDB{}
		defer s.AssertExpectations(t)

		i := &inventory.MockDB{}
		defer i.AssertExpectations(t)

		var nilRecipe *recipe.Recipe
		s.On("Find", recipe.Name("test")).Return(nilRecipe, recipe.ErrRecipeNotFound)

		var zeroInventoryItem inventory.ItemDTO
		i.On("Get", 1).Return(zeroInventoryItem, inventory.ErrItemNotFound)

		b := recipe.Collection{DB: s, Inventory: i}
		id, err := b.Add("test", []recipe.Ingredient{{ID: 1, Qty: 2}})

		assert.Zero(t, id)
		assert.Equal(t, recipe.ErrIgredientNotFound, err)
	})

	t.Run("errors when inventory call fails", func(t *testing.T) {
		s := &recipe.MockDB{}
		defer s.AssertExpectations(t)

		i := &inventory.MockDB{}
		defer i.AssertExpectations(t)

		var nilRecipe *recipe.Recipe
		s.On("Find", recipe.Name("test")).Return(nilRecipe, recipe.ErrRecipeNotFound)

		expectedErr := errors.New("test")
		i.On("Get", 1).Return(inventory.ItemDTO{}, expectedErr)

		b := recipe.Collection{Inventory: i, DB: s}
		id, err := b.Add("test", []recipe.Ingredient{{ID: 1, Qty: 2}})

		assert.Zero(t, id)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("errors when persistence fails", func(t *testing.T) {
		s := &recipe.MockDB{}
		defer s.AssertExpectations(t)

		i := &inventory.MockDB{}
		defer i.AssertExpectations(t)

		var nilRecipe *recipe.Recipe
		s.On("Find", recipe.Name("test")).Return(nilRecipe, recipe.ErrRecipeNotFound)

		item := inventory.ItemDTO{
			ID:      1,
			Enabled: true,
		}
		i.On("Get", 1).Return(item, nil)

		var expectedErr = errors.New("could not save")
		s.On("Add", mock.Anything).Return(recipe.ID(0), expectedErr)

		b := recipe.Collection{DB: s, Inventory: i}
		id, err := b.Add("test", []recipe.Ingredient{{ID: 1, Qty: 2}})

		assert.Zero(t, id)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("succeeds", func(t *testing.T) {
		db := &recipe.MockDB{}
		defer db.AssertExpectations(t)

		invDB := &inventory.MockDB{}
		defer invDB.AssertExpectations(t)

		var nilRecipe *recipe.Recipe
		db.On("Find", recipe.Name("test")).Return(nilRecipe, recipe.ErrRecipeNotFound)

		item := inventory.ItemDTO{
			ID:      1,
			Enabled: true,
		}
		invDB.On("Get", 1).Return(item, nil)

		db.On("Add", recipe.Recipe{
			Name:        "test",
			Ingredients: []recipe.Ingredient{{ID: 1, Qty: 2}},
			Enabled:     true,
		}).Return(recipe.ID(1), nil)

		b := recipe.Collection{DB: db, Inventory: invDB}
		recipeID, err := b.Add("test", []recipe.Ingredient{{ID: 1, Qty: 2}})

		assert.NoError(t, err)
		assert.Equal(t, recipe.ID(1), recipeID)
	})
}
