package warehouse_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	warehouse "github.com/anatollupacescu/retail-sample/internal/retail-sample"
)

func TestDoOutbound(t *testing.T) {

	assert := assert.New(t)

	t.Run("outbound product subtracts from stock", func(t *testing.T) {
		assert.Fail("quick break?")
	})
}

func TestConfigureOutbount(t *testing.T) {

	assert := assert.New(t)

	t.Run("empty name not accepted", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.AddFinishedProduct("", nil)
		assert.Equal(warehouse.ErrNameNotProvided, err)
	})

	t.Run("empty list of outbound items not accepted", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.AddFinishedProduct("mocha", []warehouse.OutboundItem{})
		assert.Equal(warehouse.ErrItemsNotProvided, err)
	})

	t.Run("can not add finished product that has unknown item types", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.AddFinishedProduct("mocha", []warehouse.OutboundItem{{
			ItemType: warehouse.ItemType("nope"),
			Qty:      1,
		}})
		assert.Equal(warehouse.ErrItemTypeNotFound, err)
	})

	t.Run("can not add finished product that has zero quantity", func(t *testing.T) {
		stock := warehouse.Stock{}
		stock.AddType("milk")
		err := stock.AddFinishedProduct("mocha", []warehouse.OutboundItem{{
			ItemType: warehouse.ItemType("milk"),
			Qty:      0,
		}})
		assert.Equal(warehouse.ErrZeroQuantityNotAllowed, err)
	})

	t.Run("can add finished product", func(t *testing.T) {
		stock := warehouse.Stock{}
		stock.AddType("milk")
		err := stock.AddFinishedProduct("mocha", []warehouse.OutboundItem{{
			ItemType: warehouse.ItemType("milk"),
			Qty:      5,
		}})
		assert.Nil(err)

		list := stock.FinishedProducts()
		assert.Len(list, 1)

		assert.Len(list[0].Items, 1)
	})
}

func TestInbound(t *testing.T) {

	assert := assert.New(t)

	t.Run("can add type", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.AddType("milk")
		assert.NoError(err)
	})

	t.Run("type must be unique", func(t *testing.T) {
		stock := warehouse.Stock{}
		stock.AddType("milk")
		err := stock.AddType("milk")
		assert.Equal(warehouse.ErrItemTypePresent, err)
	})

	t.Run("can add stock item with existent type", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.AddType("milk")
		item := warehouse.InboundItem{Type: "milk", Qty: 31}
		qty, err := stock.Add(item)
		assert.NoError(err)
		assert.Equal(31, qty)
	})

	t.Run("can not add stock item with non existent type", func(t *testing.T) {
		stock := warehouse.Stock{}
		item := warehouse.InboundItem{Type: "milk", Qty: 31}
		_, err := stock.Add(item)
		assert.Equal(warehouse.ErrItemTypeNotFound, err)
	})

	t.Run("stock quantity accumulates", func(t *testing.T) {
		stock := warehouse.Stock{}
		stock.AddType("milk")
		item := warehouse.InboundItem{Type: "milk", Qty: 31}
		stock.Add(item)
		item = warehouse.InboundItem{Type: "milk", Qty: 9}
		qty, _ := stock.Add(item)
		assert.Equal(40, qty)
	})
}
