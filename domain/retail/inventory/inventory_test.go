package inventory_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

func TestCreateInventoryItem(t *testing.T) {
	var (
		db *inventory.MockDB
		i  inventory.Collection
	)

	var reset = func() {
		db = &inventory.MockDB{}
		i = inventory.Collection{DB: db}
	}

	t.Run("given name empty", func(t *testing.T) {
		reset()

		id, err := i.Create("")
		db.AssertExpectations(t)

		t.Run("assert error", func(t *testing.T) {
			assert.Zero(t, id)
			assert.Equal(t, inventory.ErrEmptyName, err)
		})
	})
	t.Run("given name non unique", func(t *testing.T) {
		reset()

		db.On("Find", "milk").Return(1, nil)
		id, err := i.Create("milk")
		db.AssertExpectations(t)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, inventory.ErrDuplicateName, err)
			assert.Zero(t, id)
		})
	})
	t.Run("given fail to check for uniqueness", func(t *testing.T) {
		reset()

		var uniqunessCheckErr = errors.New("unknown")

		db.On("Find", "milk").Return(0, uniqunessCheckErr)
		id, err := i.Create("milk")
		db.AssertExpectations(t)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, uniqunessCheckErr, err)
			assert.Zero(t, id)
		})
	})
	t.Run("given item is saved", func(t *testing.T) {
		reset()

		db.On("Find", "milk").Return(0, inventory.ErrNotFound)
		db.On("Add", "milk").Return(1, nil)

		id, err := i.Create("milk")
		db.AssertExpectations(t)

		t.Run("assert success", func(t *testing.T) {
			assert.NoError(t, err)
			assert.Equal(t, 1, id)
		})
	})
	t.Run("given fail to saving the item", func(t *testing.T) {
		reset()

		var saveItemErr = errors.New("db")

		db.On("Find", "milk").Return(0, inventory.ErrNotFound)
		db.On("Add", "milk").Return(0, saveItemErr)

		id, err := i.Create("milk")
		db.AssertExpectations(t)

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, saveItemErr, err)
			assert.Zero(t, id)
		})
	})
}

func TestDisableItem(t *testing.T) {
	var db *inventory.MockDB

	var reset = func() {
		db = &inventory.MockDB{}
	}

	t.Run("given item is disabled", func(t *testing.T) {
		reset()

		db.On("Save", mock.Anything).Return(nil)

		item := inventory.Item{DB: db, Enabled: true}
		err := item.Disable()

		db.AssertExpectations(t)
		t.Run("assert success", func(t *testing.T) {
			assert.NoError(t, err)
			assert.False(t, item.Enabled)
		})
	})
	t.Run("given fail to disable item", func(t *testing.T) {
		reset()

		expectedErr := errors.New("test")
		db.On("Save", mock.Anything).Return(expectedErr)

		item := inventory.Item{DB: db, Enabled: true}
		err := item.Disable()

		db.AssertExpectations(t)
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
	t.Run("given fail to enable item", func(t *testing.T) {
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
	t.Run("given item does not exist", func(t *testing.T) {
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
	t.Run("given fail to check for presence", func(t *testing.T) {
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
	t.Run("given item valid", func(t *testing.T) {
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
