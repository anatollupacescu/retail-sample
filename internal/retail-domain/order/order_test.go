package order

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrder(t *testing.T) {
	t.Run("can add", func(t *testing.T) {
		var expectedID = ID(1)

		var orderEntry = OrderEntry{
			RecipeID: 2,
			Qty:      1,
		}

		mockStore := &MockOrderStore{}
		mockStore.On("add", mock.AnythingOfType("Order")).Return(expectedID)

		orders := Orders{
			Store: mockStore,
		}

		receivedID := orders.Add(orderEntry)

		mockStore.AssertExpectations(t)
		assert.Equal(t, receivedID, expectedID)
	})

	t.Run("can list", func(t *testing.T) {
		var expectedOrders = []Order{{
			OrderEntry: OrderEntry{
				RecipeID: 1,
				Qty:      3,
			},
			Date: time.Now(),
			ID:   ID(1),
		}}

		mockStore := &MockOrderStore{}
		mockStore.On("all").Return(expectedOrders)

		orders := Orders{
			Store: mockStore,
		}

		receivedOrders := orders.All()

		mockStore.AssertExpectations(t)
		assert.Equal(t, expectedOrders, receivedOrders)
	})

	t.Skip("can get single order")
}
