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
		mockStore.On("List").Return(emptyResp)

		i := inventory.Inventory{Store: mockStore}
		names := i.List()
		assert.Len(t, names, 0)
	})

	t.Run("should return zero ID for missing name", func(t *testing.T) {
		mockStore := &inventory.MockStore{}
		i := inventory.Inventory{Store: mockStore}

		mockStore.On("Find", inventory.Name("test")).Return(inventory.ID(0))

		id := i.Find("test")

		assert.Equal(t, id, inventory.ID(0))
		mockStore.AssertExpectations(t)
	})

	t.Run("should reject empty name", func(t *testing.T) {
		var mockStore inventory.Store
		i := inventory.Inventory{Store: mockStore}
		_, err := i.Add("")

		assert.Equal(t, inventory.ErrEmptyName, err)
	})

	t.Run("should reject duplicate name", func(t *testing.T) {
		milk := inventory.Name("milk")

		mockStore := &inventory.MockStore{}
		i := inventory.Inventory{Store: mockStore}

		mockStore.On("Find", inventory.Name(milk)).Return(inventory.ID(1))

		id, err := i.Add(milk)

		assert.Equal(t, inventory.ErrDuplicateName, err)
		assert.Zero(t, id)
		mockStore.AssertExpectations(t)
	})

	t.Run("should persist when given a valid name", func(t *testing.T) {
		milk := inventory.Name("milk")

		mockStore := &inventory.MockStore{}
		i := inventory.Inventory{Store: mockStore}

		mockStore.On("Find", inventory.Name(milk)).Return(inventory.ID(0))
		mockStore.On("Add", inventory.Name(milk)).Return(inventory.ID(1))

		id, err := i.Add(milk)

		assert.NoError(t, err)
		assert.Equal(t, id, inventory.ID(1))
		mockStore.AssertExpectations(t)
	})

	t.Run("should log successfully added name", func(t *testing.T) {
		milk := inventory.Name("milk")

		mockStore := &inventory.MockStore{}
		i := inventory.Inventory{Store: mockStore}

		mockStore.On("Find", inventory.Name(milk)).Return(inventory.ID(0))
		mockStore.On("Add", inventory.Name(milk)).Return(inventory.ID(1))

		id, err := i.Add(milk)

		assert.NoError(t, err)
		assert.Equal(t, id, inventory.ID(1))
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
		})

		records := i.List()
		assert.Len(t, records, 1)

		mockStore.AssertExpectations(t)
	})

	t.Run("should return ID for correct name", func(t *testing.T) {
		mockStore := &inventory.MockStore{}
		i := inventory.Inventory{Store: mockStore}

		mockStore.On("Find", inventory.Name("test")).Return(inventory.ID(1))

		id := i.Find("test")
		assert.Equal(t, id, inventory.ID(1))

		mockStore.AssertExpectations(t)
	})
}
