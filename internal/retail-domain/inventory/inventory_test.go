package inventory_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
)

func TestAdd(t *testing.T) {
	t.Run("when name is empty", func(t *testing.T) {
		i := inventory.Inventory{}
		_, err := i.Add("")

		t.Run("throws empty name error", func(t *testing.T) {
			assert.Equal(t, inventory.ErrEmptyName, err)
		})
	})

	t.Run("when name is already present", func(t *testing.T) {
		milk := "milk"

		mockStore := &inventory.MockStore{}
		i := inventory.Inventory{Store: mockStore}

		mockStore.On("Find", milk).Return(1, nil)

		id, err := i.Add(milk)

		t.Run("calls the store", func(t *testing.T) {
			mockStore.AssertExpectations(t)
		})

		t.Run("throws duplicate name error", func(t *testing.T) {
			assert.Equal(t, inventory.ErrDuplicateName, err)
			assert.Zero(t, id)
		})
	})

	t.Run("when name is valid", func(t *testing.T) {
		milk := "milk"

		mockStore := &inventory.MockStore{}
		i := inventory.Inventory{Store: mockStore}

		mockStore.On("Find", milk).Return(0, inventory.ErrItemNotFound)
		mockStore.On("Add", milk).Return(1, nil)

		id, err := i.Add(milk)

		t.Run("calls the store", func(t *testing.T) {
			mockStore.AssertExpectations(t)
		})

		t.Run("saves it to store", func(t *testing.T) {
			assert.NoError(t, err)
			assert.Equal(t, 1, id)
		})
	})

	t.Run("when store throws error during find op", func(t *testing.T) {
		milk := "milk"

		mockStore := &inventory.MockStore{}
		i := inventory.Inventory{Store: mockStore}

		expected := errors.New("unknown")
		mockStore.On("Find", milk).Return(0, expected)

		id, err := i.Add(milk)

		t.Run("calls the store", func(t *testing.T) {
			mockStore.AssertExpectations(t)
		})

		t.Run("should propage it to the caller", func(t *testing.T) {
			assert.Equal(t, expected, err)
			assert.Equal(t, 0, id)
		})
	})
}

func TestList(t *testing.T) {
	t.Run("calls store", func(t *testing.T) {})
}
func TestFind(t *testing.T) {
	t.Run("calls store", func(t *testing.T) {})
}
func TestGet(t *testing.T) {
	t.Run("calls store", func(t *testing.T) {})
}
