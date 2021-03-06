package recipe_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

func TestCreateRecipe(t *testing.T) {
	var (
		db    *recipe.MockDB
		reset = func() {
			db = &recipe.MockDB{}
		}
	)
	t.Run("given empty name", func(t *testing.T) {
		rr := recipe.Recipes{}
		_, err := rr.Create("", nil)
		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, recipe.ErrEmptyName, err)
		})
	})
	t.Run("given no ingredients", func(t *testing.T) {
		rr := recipe.Recipes{}
		_, err := rr.Create("test", nil)
		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, recipe.ErrNoIngredients, err)
		})
	})
	t.Run("given invalid ingredient quantity", func(t *testing.T) {
		rr := recipe.Recipes{}
		_, err := rr.Create("test", []recipe.InventoryItem{{ID: 1, Qty: 0}})
		t.Run("assert error", func(t *testing.T) {
			assert.True(t, errors.Is(err, recipe.ErrQuantityNotProvided))
		})
	})
	t.Run("given duplicate ingredient", func(t *testing.T) {
		reset()

		item := recipe.DTO{ID: 1, Name: "test"}
		db.On("Find", "test").Return(item, nil)

		recipes := recipe.Recipes{DB: db}
		id, err := recipes.Create("test", []recipe.InventoryItem{{ID: 1, Qty: 2}})

		db.AssertExpectations(t)

		t.Run("assert error", func(t *testing.T) {
			assert.Zero(t, id)
			assert.Equal(t, recipe.ErrDuplicateName, err)
		})
	})
	t.Run("given ingredient invalid", func(t *testing.T) {
		reset()

		db.On("Find", "test").Return(recipe.DTO{}, recipe.ErrNotFound)

		mi := &recipe.MockValidator{}
		mi.On("Valid", mock.Anything).Return(false, nil)

		b := recipe.Recipes{DB: db, ItemValidator: mi}
		id, err := b.Create("test", []recipe.InventoryItem{{ID: 1, Qty: 2}})

		db.AssertExpectations(t)
		mi.AssertExpectations(t)

		t.Run("assert error", func(t *testing.T) {
			assert.True(t, errors.Is(err, recipe.ErrIngredientNotValid))
			assert.Zero(t, id)
		})
	})
	t.Run("given fail to check validity", func(t *testing.T) {
		reset()

		db.On("Find", "test").Return(recipe.DTO{}, recipe.ErrNotFound)

		mi := &recipe.MockValidator{}
		expected := errors.New("other")
		mi.On("Valid", mock.Anything).Return(false, expected)

		b := recipe.Recipes{DB: db, ItemValidator: mi}
		id, err := b.Create("test", []recipe.InventoryItem{{ID: 1, Qty: 2}})

		db.AssertExpectations(t)
		mi.AssertExpectations(t)

		t.Run("assert error", func(t *testing.T) {
			assert.True(t, errors.Is(err, recipe.ErrIngredientNotValid))
			assert.Zero(t, id)
		})
	})
	t.Run("given recipe name not unique", func(t *testing.T) {
		reset()

		db.On("Find", "test").Return(recipe.DTO{}, nil)

		recipes := recipe.Recipes{DB: db}
		id, err := recipes.Create("test", []recipe.InventoryItem{{ID: 1, Qty: 2}})

		db.AssertExpectations(t)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, recipe.ErrDuplicateName, err)
			assert.Zero(t, id)
		})
	})
	t.Run("given fail to check for name uniqueness", func(t *testing.T) {
		reset()

		expectedErr := errors.New("test")
		db.On("Find", "test").Return(recipe.DTO{}, expectedErr)

		recipes := recipe.Recipes{DB: db}
		id, err := recipes.Create("test", []recipe.InventoryItem{{ID: 1, Qty: 2}})

		db.AssertExpectations(t)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, expectedErr, err)
			assert.Zero(t, id)
		})
	})
	t.Run("given recipe created", func(t *testing.T) {
		reset()

		db.On("Find", "test").Return(recipe.DTO{}, recipe.ErrNotFound)
		db.On("Add", mock.Anything).Return(3, nil)

		validator := &recipe.MockValidator{}
		validator.On("Valid", mock.Anything).Return(true, nil)

		recipes := recipe.Recipes{DB: db, ItemValidator: validator}
		id, err := recipes.Create("test", []recipe.InventoryItem{{ID: 1, Qty: 2}})

		db.AssertExpectations(t)
		validator.AssertExpectations(t)

		t.Run("assert success", func(t *testing.T) {
			assert.NoError(t, err)
			assert.Equal(t, 3, id)
		})
	})
	t.Run("given fail to save recipe", func(t *testing.T) {
		reset()

		db.On("Find", "test").Return(recipe.DTO{}, recipe.ErrNotFound)
		dbErr := errors.New("expected")
		db.On("Add", mock.Anything).Return(0, dbErr)

		validator := &recipe.MockValidator{}
		validator.On("Valid", mock.Anything).Return(true, nil)

		recipes := recipe.Recipes{DB: db, ItemValidator: validator}
		_, err := recipes.Create("test", []recipe.InventoryItem{{ID: 1, Qty: 2}})

		db.AssertExpectations(t)
		validator.AssertExpectations(t)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, dbErr, err)
		})
	})
}

func TestDisableRecipe(t *testing.T) {
	t.Run("given recipe is disabled", func(t *testing.T) {
		db := &recipe.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Save", mock.Anything).Return(nil)

		r := recipe.Recipe{
			Enabled: true,
			DB:      db,
		}

		err := r.Disable()

		t.Run("assert success", func(t *testing.T) {
			assert.NoError(t, err)
			assert.False(t, r.Enabled)
		})
	})
	t.Run("given fail to disable recipe", func(t *testing.T) {
		db := &recipe.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Save", mock.Anything).Return(recipe.ErrNotFound)

		r := recipe.Recipe{ID: 1, DB: db}

		err := r.Disable()

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, recipe.ErrNotFound, err)
		})
	})
}

func TestEnableRecipe(t *testing.T) {
	t.Run("given recipe is enabled", func(t *testing.T) {
		db := &recipe.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Save", mock.Anything).Return(nil)

		r := recipe.Recipe{
			DB: db,
		}

		err := r.Enable()

		t.Run("assert success", func(t *testing.T) {
			assert.NoError(t, err)
			assert.True(t, r.Enabled)
		})
	})
	t.Run("given fail to enable recipe", func(t *testing.T) {
		db := &recipe.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Save", mock.Anything).Return(recipe.ErrNotFound)

		r := recipe.Recipe{ID: 1, DB: db}

		err := r.Enable()

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, recipe.ErrNotFound, err)
		})
	})
}

func TestValidateRecipe(t *testing.T) {
	var (
		db    *recipe.MockDB
		reset = func() {
			db = &recipe.MockDB{}
		}
		validator = func() recipe.Validator {
			return recipe.Validator{
				Recipes: db,
			}
		}
	)
	t.Run("given recipe not found", func(t *testing.T) {
		reset()

		expected := recipe.ErrNotFound
		db.On("Get", 1).Return(recipe.DTO{}, expected)

		_, err := validator().Valid(1)

		db.AssertExpectations(t)
		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, expected, err)
		})
	})
	t.Run("given fail to check for presence", func(t *testing.T) {
		reset()

		expected := errors.New("test")
		db.On("Get", 1).Return(recipe.DTO{}, expected)

		_, err := validator().Valid(1)

		db.AssertExpectations(t)
		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, expected, err)
		})
	})
	t.Run("given recipe disabled", func(t *testing.T) {
		reset()

		db.On("Get", 1).Return(recipe.DTO{Enabled: false}, nil)

		valid, err := validator().Valid(1)

		db.AssertExpectations(t)
		t.Run("assert error", func(t *testing.T) {
			assert.NoError(t, err)
			assert.False(t, valid)
		})
	})
	t.Run("given recipe validated", func(t *testing.T) {
		reset()

		db.On("Get", 1).Return(recipe.DTO{Enabled: true}, nil)

		_, err := validator().Valid(1)

		db.AssertExpectations(t)
		t.Run("assert success", func(t *testing.T) {
			assert.NoError(t, err)
		})
	})
}
