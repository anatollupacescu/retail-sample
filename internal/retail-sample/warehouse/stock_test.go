package warehouse_test

import (
	"testing"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/warehouse"

	"github.com/stretchr/testify/assert"
)

func TestConfigureItemType(t *testing.T) {

	t.Run("should reject empty type name", func(t *testing.T) {
		stock := stock().build()
		err := stock.ConfigureInboundType("")
		assert.Equal(t, warehouse.ErrInboundNameNotProvided, err)
	})

	t.Run("should succeed for valid type name", func(t *testing.T) {
		stock := stock().build()
		err := stock.ConfigureInboundType("milk")
		assert.NoError(t, err)
	})

	t.Run("newly added types have 0 quantity in stock", func(t *testing.T) {
		stock := stock().build()
		err := stock.ConfigureInboundType("milk")
		assert.NoError(t, err)
		qty, err := stock.Quantity("milk")
		assert.NoError(t, err)
		assert.Equal(t, 0, qty)
	})

	t.Run("should return err when getting quantity for non existent item", func(t *testing.T) {
		stock := stock().build()
		_, err := stock.Quantity("iDoNotExist")
		assert.Equal(t, warehouse.ErrInventoryItemNotFound, err)
	})

	t.Run("should reject duplicate name", func(t *testing.T) {
		stock := stock().build()
		err := stock.ConfigureInboundType("milk")
		assert.NoError(t, err)
		err = stock.ConfigureInboundType("milk")
		assert.Equal(t, warehouse.ErrInboundItemTypeAlreadyConfigured, err)
	})
}

func TestStockWithoutConfiguredItemTypes(t *testing.T) {

	t.Run("can not add stock item with non existent type", func(t *testing.T) {
		stock := stock().build()
		item := warehouse.Item{Type: "milk", Qty: 31}
		_, err := stock.PlaceInbound(item)
		assert.Equal(t, warehouse.ErrInboundItemTypeNotFound, err)
	})
}

func TestStockWithConfiguredItems(t *testing.T) {

	t.Run("should place inbound when item type exists", func(t *testing.T) {
		stock := stock().build()
		_ = stock.ConfigureInboundType("milk")
		item := warehouse.Item{Type: "milk", Qty: 31}
		qty, err := stock.PlaceInbound(item)
		assert.NoError(t, err)
		assert.Equal(t, 31, qty)

		qty, _ = stock.Quantity("milk")
		assert.Equal(t, 31, qty)
	})

	t.Run("should add to inbound log", func(t *testing.T) {
		stock := stock().with("milk", 31).build()

		l := stock.ListInbound()
		assert.Len(t, l, 1)
	})

	t.Run("should increment existing stock levels", func(t *testing.T) {
		stock := stock().with("milk", 31).build()

		item := warehouse.Item{Type: "milk", Qty: 9}
		qty, err := stock.PlaceInbound(item)
		assert.NoError(t, err)
		assert.Equal(t, 40, qty)
	})
}
