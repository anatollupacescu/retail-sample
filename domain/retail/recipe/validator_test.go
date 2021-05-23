package recipe_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

func TestValidator(t *testing.T) {
	t.Run("given a bad recipe", func(t *testing.T) {
		db := &recipe.MockDB{}
		defer db.AssertExpectations(t)

		expected := errors.New("test")

		db.On("Get", 1).Return(recipe.DTO{}, expected)
		v := recipe.Validator{
			Recipes: db,
		}

		_, err := v.Valid(1)

		t.Run("assert is not valid", func(t *testing.T) {
			assert.Equal(t, expected, err)
		})
	})

	t.Run("given a disabled recipe", func(t *testing.T) {
		db := &recipe.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Get", 1).Return(recipe.DTO{Enabled: false}, nil)

		v := recipe.Validator{
			Recipes: db,
		}

		valid, err := v.Valid(1)

		t.Run("assert is not valid", func(t *testing.T) {
			assert.NoError(t, err)
			assert.False(t, valid)
		})
	})

	t.Run("given a valid recipe", func(t *testing.T) {
		db := &recipe.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Get", 1).Return(recipe.DTO{Enabled: true}, nil)

		v := recipe.Validator{
			Recipes: db,
		}

		valid, err := v.Valid(1)

		t.Run("assert is valid", func(t *testing.T) {
			assert.NoError(t, err)
			assert.True(t, valid)
		})
	})
}
