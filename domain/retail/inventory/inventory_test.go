package inventory_test

import (
	"errors"
	"testing"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateInventoryItem(t *testing.T) {
	var ( // set up
		name string

		db *inventory.MockDB
		i  inventory.Collection

		reset = func() {
			name = "milk"
			db = &inventory.MockDB{}
			i = inventory.Collection{DB: db}
		}

		givenNameIsEmpty = func() {
			name = ""
		}

		create = func() (id int, err error) {
			id, err = i.Create(name)
			db.AssertExpectations(t)
			return
		}

		givenNonUniqueName = func() {
			db.On("Find", "milk").Return(1, nil)
		}

		uniqunessCheckErr = errors.New("unknown")

		givenErrorCheckingForUniqueness = func() {
			db.On("Find", "milk").Return(0, uniqunessCheckErr)
		}

		saveItemErr = errors.New("db")

		givenErrorSavingItem = func() {
			db.On("Find", "milk").Return(0, inventory.ErrNotFound)
			db.On("Add", "milk").Return(0, saveItemErr)
		}

		givenCanSaveItem = func() {
			db.On("Find", "milk").Return(0, inventory.ErrNotFound)
			db.On("Add", "milk").Return(1, nil)
		}
	)

	{ // tests
		reset()
		givenNameIsEmpty()
		_, err := create()
		assert.Equal(t, inventory.ErrEmptyName, err)

		reset()
		givenErrorCheckingForUniqueness()
		id, err := create()
		assert.Equal(t, uniqunessCheckErr, err)
		assert.Zero(t, id)

		reset()
		givenNonUniqueName()
		id, err = create()
		assert.Equal(t, inventory.ErrDuplicateName, err)
		assert.Zero(t, id)

		reset()
		givenErrorSavingItem()
		id, err = create()
		assert.Equal(t, saveItemErr, err)
		assert.Zero(t, id)

		reset()
		givenCanSaveItem()
		id, err = create()
		assert.NoError(t, err)
		assert.Equal(t, 1, id)
	}
}

func TestDisableItem(t *testing.T) {
	var (
		db   *inventory.MockDB
		item inventory.Item
	)

	var reset = func() {
		db = &inventory.MockDB{}
		item = inventory.Item{DB: db, Enabled: true}
	}

	t.Run("given item is saved", func(t *testing.T) {
		reset()

		db.On("Save", mock.Anything).Return(nil)

		err := item.Disable()

		t.Run("assert success", func(t *testing.T) {
			assert.NoError(t, err)
			assert.False(t, item.Enabled)
			db.AssertExpectations(t)
		})
	})
	t.Run("given failed to save item", func(t *testing.T) {
		reset()

		expectedErr := errors.New("test")
		db.On("Save", mock.Anything).Return(expectedErr)

		err := item.Disable()

		t.Run("assert error", func(t *testing.T) {
			assert.Equal(t, expectedErr, err)
			assert.True(t, item.Enabled)
			db.AssertExpectations(t)
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
	t.Run("given a non existent item", func(t *testing.T) {
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
	t.Run("given item is enabled", func(t *testing.T) {
		db := &inventory.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Get", 1).Return(inventory.DTO{Enabled: true}, nil)

		v := inventory.Validator{
			Inventory: db,
		}

		err := v.Validate(1)

		t.Run("assert returns valid", func(t *testing.T) {
			assert.NoError(t, err)
		})
	})
	t.Run("given item is not enabled", func(t *testing.T) {
		db := &inventory.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Get", 1).Return(inventory.DTO{Enabled: false}, nil)

		v := inventory.Validator{
			Inventory: db,
		}

		err := v.Validate(1)

		t.Run("assert returns invalid", func(t *testing.T) {
			assert.Equal(t, inventory.ErrItemDisabled, err)
		})
	})
}
