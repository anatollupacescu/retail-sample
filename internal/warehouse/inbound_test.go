package warehouse_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	warehouse "github.com/anatollupacescu/retail-sample/internal/warehouse"
)

func TestInbound(t *testing.T) {

	asrt := assert.New(t)

	t.Run("can add type", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureInboundType("milk")
		asrt.NoError(err)
	})

	t.Run("newly added types have 0 quantity in stock", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureInboundType("milk")
		asrt.NoError(err)
		qty, err := stock.Quantity("milk")
		asrt.NoError(err)
		asrt.Equal(0, qty)
	})

	t.Run("type must be unique", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureInboundType("milk")
		asrt.NoError(err)
		err = stock.ConfigureInboundType("milk")
		asrt.Equal(warehouse.ErrItemTypePresent, err)
	})

	t.Run("can add stock item with existent type", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureInboundType("milk")
		item := warehouse.InboundItem{Type: "milk", Qty: 31}
		qty, err := stock.Provision(item)
		asrt.NoError(err)
		asrt.Equal(31, qty)

		qty, _ = stock.Quantity("milk")
		asrt.Equal(31, qty)
	})

	t.Run("can not add stock item with non existent type", func(t *testing.T) {
		stock := warehouse.Stock{}
		item := warehouse.InboundItem{Type: "milk", Qty: 31}
		_, err := stock.Provision(item)
		asrt.Equal(warehouse.ErrItemTypeNotFound, err)
	})

	t.Run("stock quantity accumulates", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureInboundType("milk")
		asrt.NoError(err)
		item := warehouse.InboundItem{Type: "milk", Qty: 31}
		_, err = stock.Provision(item)
		asrt.NoError(err)
		item = warehouse.InboundItem{Type: "milk", Qty: 9}
		qty, err := stock.Provision(item)
		asrt.NoError(err)
		asrt.Equal(40, qty)
	})
}
