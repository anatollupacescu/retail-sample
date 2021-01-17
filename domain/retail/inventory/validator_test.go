package inventory_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

func TestValidator(t *testing.T) {
	t.Run("given a non existent inventory item", func(t *testing.T) {
		db := &inventory.MockDB{}
		defer db.AssertExpectations(t)

		expected := errors.New("test")

		db.On("Get", 1).Return(inventory.DTO{}, expected)

		v := inventory.Validator{
			Inventory: db,
		}

		err := v.Validate(1)

		t.Run("assert not valid", func(t *testing.T) {
			assert.Equal(t, expected, err)
		})
	})

	t.Run("given a disabled inventory item", func(t *testing.T) {
		db := &inventory.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Get", 1).Return(inventory.DTO{Enabled: false}, nil)

		v := inventory.Validator{
			Inventory: db,
		}

		err := v.Validate(1)

		t.Run("assert not valid", func(t *testing.T) {
			assert.Equal(t, inventory.ErrItemDisabled, err)
		})
	})

	t.Run("given a valid inventory item", func(t *testing.T) {
		db := &inventory.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Get", 1).Return(inventory.DTO{Enabled: true}, nil)

		v := inventory.Validator{
			Inventory: db,
		}

		err := v.Validate(1)

		t.Run("assert is valid", func(t *testing.T) {
			assert.NoError(t, err)
		})
	})
}
