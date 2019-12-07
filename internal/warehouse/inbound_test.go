package warehouse_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	warehouse "github.com/anatollupacescu/retail-sample/internal/warehouse"
)

func TestInbound(t *testing.T) {

	assert := assert.New(t)

	t.Run("can add type", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureInboundType("milk")
		assert.NoError(err)
	})

	t.Run("type must be unique", func(t *testing.T) {
		stock := warehouse.Stock{}
		stock.ConfigureInboundType("milk")
		err := stock.ConfigureInboundType("milk")
		assert.Equal(warehouse.ErrItemTypePresent, err)
	})

	t.Run("can add stock item with existent type", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureInboundType("milk")
		item := warehouse.InboundItem{Type: "milk", Qty: 31}
		qty, err := stock.Provision(item)
		assert.NoError(err)
		assert.Equal(31, qty)

		qty, _ = stock.Quantity("milk")
		assert.Equal(31, qty)
	})

	t.Run("can not add stock item with non existent type", func(t *testing.T) {
		stock := warehouse.Stock{}
		item := warehouse.InboundItem{Type: "milk", Qty: 31}
		_, err := stock.Provision(item)
		assert.Equal(warehouse.ErrItemTypeNotFound, err)
	})

	t.Run("stock quantity accumulates", func(t *testing.T) {
		stock := warehouse.Stock{}
		stock.ConfigureInboundType("milk")
		item := warehouse.InboundItem{Type: "milk", Qty: 31}
		stock.Provision(item)
		item = warehouse.InboundItem{Type: "milk", Qty: 9}
		qty, _ := stock.Provision(item)
		assert.Equal(40, qty)
	})
}
