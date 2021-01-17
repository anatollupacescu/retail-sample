package inventory_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

func TestEnable(t *testing.T) {
	t.Run("propagates error leaving the status unchanged", func(t *testing.T) {
		db := &inventory.MockDB{}
		defer db.AssertExpectations(t)

		expectedErr := errors.New("test")
		db.On("Save", mock.Anything).Return(expectedErr)

		var item = inventory.Item{DB: db}

		err := item.Enable()

		assert.Equal(t, expectedErr, err)
		assert.False(t, item.Enabled)
	})

	t.Run("success", func(t *testing.T) {
		db := &inventory.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Save", mock.Anything).Return(nil)

		var item = inventory.Item{DB: db}

		err := item.Enable()

		assert.NoError(t, err)
		assert.True(t, item.Enabled)
	})

}

func TestDisable(t *testing.T) {
	t.Run("propagates error leaving the status unchanged", func(t *testing.T) {
		db := &inventory.MockDB{}
		defer db.AssertExpectations(t)

		expectedErr := errors.New("test")
		db.On("Save", mock.Anything).Return(expectedErr)

		var item = inventory.Item{Enabled: true, DB: db}

		err := item.Disable()

		assert.Equal(t, expectedErr, err)
		assert.True(t, item.Enabled)
	})

	t.Run("success", func(t *testing.T) {
		db := &inventory.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Save", mock.Anything).Return(nil)

		var item = inventory.Item{DB: db}

		err := item.Disable()

		assert.NoError(t, err)
		assert.False(t, item.Enabled)
	})
}

func TestAdd(t *testing.T) {
	t.Run("errors when name is empty", func(t *testing.T) {
		i := inventory.Collection{}
		_, err := i.Create("")

		assert.Equal(t, inventory.ErrEmptyName, err)
	})

	t.Run("error when name is already present", func(t *testing.T) {
		db := &inventory.MockDB{}
		defer db.AssertExpectations(t)

		db.On("Find", "milk").Return(1, nil)

		i := inventory.Collection{DB: db}
		id, err := i.Create("milk")

		assert.Equal(t, inventory.ErrDuplicateName, err)
		assert.Zero(t, id)
	})

	t.Run("can add", func(t *testing.T) {
		db := &inventory.MockDB{}
		db.AssertExpectations(t)

		db.On("Find", "milk").Return(0, inventory.ErrNotFound)
		db.On("Add", "milk").Return(1, nil)

		i := inventory.Collection{DB: db}
		id, err := i.Create("milk")

		assert.NoError(t, err)
		assert.Equal(t, 1, id)
	})

	t.Run("propagates error during find op", func(t *testing.T) {
		db := &inventory.MockDB{}
		db.AssertExpectations(t)

		expected := errors.New("unknown")
		db.On("Find", "milk").Return(0, expected)

		i := inventory.Collection{DB: db}
		id, err := i.Create("milk")

		assert.Equal(t, expected, err)
		assert.Equal(t, 0, id)
	})

	t.Run("propagates error from save op", func(t *testing.T) {
		db := &inventory.MockDB{}
		defer db.AssertExpectations(t)

		expectedErr := errors.New("test")
		db.On("Find", "milk").Return(0, inventory.ErrNotFound)
		db.On("Add", "milk").Return(1, expectedErr)

		i := inventory.Collection{DB: db}
		id, err := i.Create("milk")

		assert.Zero(t, id)
		assert.Equal(t, expectedErr, err)
	})
}
