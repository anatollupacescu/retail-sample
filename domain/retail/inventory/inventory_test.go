package inventory_test

import (
	"errors"
	"testing"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateItem(t *testing.T) {
	t.Run("given empty name", func(t *testing.T) {
		i := inventory.Collection{}
		_, err := i.Create("")

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, inventory.ErrEmptyName, err)
		})
	})
	t.Run("given an item with a non unique name", func(t *testing.T) {
		db := &inventory.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Find", "milk").Return(1, nil)

		i := inventory.Collection{DB: db}
		id, err := i.Create("milk")

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, inventory.ErrDuplicateName, err)
			assert.Zero(t, id)
		})
	})
	t.Run("given failure to check for uniqueness", func(t *testing.T) {
		db := &inventory.MockDB{}
		db.AssertExpectations(t)

		expected := errors.New("unknown")
		db.On("Find", "milk").Return(0, expected)

		i := inventory.Collection{DB: db}
		id, err := i.Create("milk")

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, expected, err)
			assert.Equal(t, 0, id)
		})
	})
	t.Run("given valid item is not created", func(t *testing.T) {
		db := &inventory.MockDB{}
		defer db.AssertExpectations(t)

		expectedErr := errors.New("test")
		db.On("Find", "milk").Return(0, inventory.ErrNotFound)
		db.On("Add", "milk").Return(1, expectedErr)

		i := inventory.Collection{DB: db}
		id, err := i.Create("milk")

		t.Run("assert error", func(t *testing.T) {
			assert.Zero(t, id)
			assert.Equal(t, expectedErr, err)
		})
	})
	t.Run("given valid item", func(t *testing.T) {
		db := &inventory.MockDB{}
		db.AssertExpectations(t)

		db.On("Find", "milk").Return(0, inventory.ErrNotFound)
		db.On("Add", "milk").Return(1, nil)

		i := inventory.Collection{DB: db}
		id, err := i.Create("milk")

		t.Run("assert item is created", func(t *testing.T) {
			assert.NoError(t, err)
			assert.Equal(t, 1, id)
		})
	})
}

func TestDisableItem(t *testing.T) {
	t.Run("given item is disabled", func(t *testing.T) {
		db := &inventory.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Save", mock.Anything).Return(nil)

		var item = inventory.Item{DB: db}

		err := item.Disable()

		t.Run("assert success", func(t *testing.T) {
			assert.NoError(t, err)
			assert.False(t, item.Enabled)
		})
	})
	t.Run("given failed to disable item", func(t *testing.T) {
		db := &inventory.MockDB{}
		defer db.AssertExpectations(t)

		expectedErr := errors.New("test")
		db.On("Save", mock.Anything).Return(expectedErr)

		var item = inventory.Item{Enabled: true, DB: db}

		err := item.Disable()

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, expectedErr, err)
			assert.True(t, item.Enabled)
		})
	})
}

func TestEnableItem(t *testing.T) {
	t.Run("given item is enabled", func(t *testing.T) {
		db := &inventory.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Save", mock.Anything).Return(nil)

		var item = inventory.Item{DB: db}

		err := item.Enable()

		t.Run("assert success", func(t *testing.T) {
			assert.NoError(t, err)
			assert.True(t, item.Enabled)
		})
	})
	t.Run("given failed to enable item", func(t *testing.T) {
		db := &inventory.MockDB{}
		defer db.AssertExpectations(t)

		expectedErr := errors.New("test")
		db.On("Save", mock.Anything).Return(expectedErr)

		var item = inventory.Item{DB: db}

		err := item.Enable()

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, expectedErr, err)
			assert.False(t, item.Enabled)
		})
	})
}

func TestValidateItem(t *testing.T) {
	t.Run("given item is not found", func(t *testing.T) {
		db := &inventory.MockDB{}
		defer db.AssertExpectations(t)

		expected := errors.New("test")

		db.On("Get", 1).Return(inventory.DTO{}, expected)

		v := inventory.Validator{
			Inventory: db,
		}

		err := v.Validate(1)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, expected, err)
		})
	})
	t.Run("given item enabled status is false", func(t *testing.T) {
		db := &inventory.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Get", 1).Return(inventory.DTO{Enabled: false}, nil)

		v := inventory.Validator{
			Inventory: db,
		}

		err := v.Validate(1)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, inventory.ErrItemDisabled, err)
		})
	})
	t.Run("given item is valid", func(t *testing.T) {
		db := &inventory.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Get", 1).Return(inventory.DTO{Enabled: true}, nil)

		v := inventory.Validator{
			Inventory: db,
		}

		err := v.Validate(1)

		t.Run("assert success", func(t *testing.T) {
			assert.NoError(t, err)
		})
	})
}
