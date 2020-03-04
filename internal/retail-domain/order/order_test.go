package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	t.Run("can add", func(t *testing.T) {
		var expectedID = ID(1)

		var orderEntry = OrderEntry{
			Qty: 1,
		}

		mockStore := &MockOrderStore{}

		mockStore.On("add", orderEntry).Return(expectedID)

		orders := Orders{
			Store: mockStore,
		}

		receivedID := orders.Add(orderEntry)

		mockStore.AssertExpectations(t)
		assert.Equal(t, receivedID, expectedID)
	})
}
