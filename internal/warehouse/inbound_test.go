package warehouse_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anatollupacescu/retail-sample/internal/warehouse"
)

func TestConfigureItemType(t *testing.T) {

	t.Run("should reject empty name", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureInboundType("")
		assert.Equal(t, warehouse.ErrInboundNameNotProvided, err)
	})

	t.Run("can add type", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureInboundType("milk")
		assert.NoError(t, err)
	})

	t.Run("newly added types have 0 quantity in stock", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureInboundType("milk")
		assert.NoError(t, err)
		qty, err := stock.Quantity("milk")
		assert.NoError(t, err)
		assert.Equal(t, 0, qty)
	})

	t.Run("should reject duplicate name", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureInboundType("milk")
		assert.NoError(t, err)
		err = stock.ConfigureInboundType("milk")
		assert.Equal(t, warehouse.ErrInboundItemTypeDuplicated, err)
	})
}

func TestStockWithoutConfiguredItemTypes(t *testing.T) {

	t.Run("can not add stock item with non existent type", func(t *testing.T) {
		stock := warehouse.Stock{}
		item := warehouse.InboundItem{Type: "milk", Qty: 31}
		_, err := stock.Provision(item)
		assert.Equal(t, warehouse.ErrInboundItemTypeNotFound, err)
	})
}

func TestStockWithConfiguredItems(t *testing.T) {

	t.Run("should succeed when item type exists", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureInboundType("milk")
		item := warehouse.InboundItem{Type: "milk", Qty: 31}
		qty, err := stock.Provision(item)
		assert.NoError(t, err)
		assert.Equal(t, 31, qty)

		qty, _ = stock.Quantity("milk")
		assert.Equal(t, 31, qty)
	})

	t.Run("should increment existing stock levels", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureInboundType("milk")
		assert.NoError(t, err)

		item := warehouse.InboundItem{Type: "milk", Qty: 31}
		_, err = stock.Provision(item)
		assert.NoError(t, err)

		item = warehouse.InboundItem{Type: "milk", Qty: 9}
		qty, err := stock.Provision(item)
		assert.NoError(t, err)
		assert.Equal(t, 40, qty)
	})
}
