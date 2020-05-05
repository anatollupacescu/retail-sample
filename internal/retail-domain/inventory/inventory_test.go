package inventory_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
)

func TestInventory(t *testing.T) {

	t.Run("should have no registered names on creation", func(t *testing.T) {
		mockStore := &inventory.MockStore{}

		var emptyResp []inventory.Item
		mockStore.On("List").Return(emptyResp, nil)

		i := inventory.Inventory{Store: mockStore}
		names, err := i.List()

		assert.NoError(t, err)
		assert.Len(t, names, 0)
	})

	t.Run("should return zero ID for missing name", func(t *testing.T) {
		t.Skip("will return an error")
		mockStore := &inventory.MockStore{}
		i := inventory.Inventory{Store: mockStore}

		mockStore.On("Find", "test").Return(0, nil)

		id, err := i.Find("test")

		assert.NoError(t, err)
		assert.Equal(t, id, 0)
		mockStore.AssertExpectations(t)
	})

	t.Run("should reject empty name", func(t *testing.T) {
		var mockStore inventory.Store
		i := inventory.Inventory{Store: mockStore}
		_, err := i.Add("")

		assert.Equal(t, inventory.ErrEmptyName, err)
	})

	t.Run("should reject duplicate name", func(t *testing.T) {
		milk := "milk"

		mockStore := &inventory.MockStore{}
		i := inventory.Inventory{Store: mockStore}

		mockStore.On("Find", milk).Return(1, nil)

		id, err := i.Add(milk)

		assert.Equal(t, inventory.ErrDuplicateName, err)
		assert.Zero(t, id)
		mockStore.AssertExpectations(t)
	})

	t.Run("should persist when given a valid name", func(t *testing.T) {
		t.SkipNow()

		milk := "milk"

		mockStore := &inventory.MockStore{}
		i := inventory.Inventory{Store: mockStore}

		mockStore.On("Find", milk).Return(0, nil)
		mockStore.On("Add", milk).Return(1, nil)

		id, err := i.Add(milk)

		assert.NoError(t, err)
		assert.Equal(t, id, 1)
		mockStore.AssertExpectations(t)
	})

	t.Run("should log successfully added name", func(t *testing.T) {
		t.SkipNow()

		milk := "milk"

		mockStore := &inventory.MockStore{}
		i := inventory.Inventory{Store: mockStore}

		mockStore.On("Find", milk).Return(0, nil)
		mockStore.On("Add", milk).Return(1, nil)

		id, err := i.Add(milk)

		assert.NoError(t, err)
		assert.Equal(t, id, 1)
		mockStore.AssertExpectations(t)
	})

	t.Run("should provide full list of names", func(t *testing.T) {
		mockStore := &inventory.MockStore{}
		i := inventory.Inventory{Store: mockStore}

		mockStore.On("List").Return([]inventory.Item{
			{
				ID:   1,
				Name: "test",
			},
		}, nil)

		records, err := i.List()

		assert.NoError(t, err)
		assert.Len(t, records, 1)

		mockStore.AssertExpectations(t)
	})

	t.Run("should return ID for correct name", func(t *testing.T) {
		mockStore := &inventory.MockStore{}
		i := inventory.Inventory{Store: mockStore}

		mockStore.On("Find", "test").Return(1, nil)

		id, err := i.Find("test")

		assert.NoError(t, err)
		assert.Equal(t, id, 1)

		mockStore.AssertExpectations(t)
	})
}
