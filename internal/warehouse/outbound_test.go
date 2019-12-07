package warehouse_test

import (
	warehouse "github.com/anatollupacescu/retail-sample/internal/warehouse"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOutbound(t *testing.T) {
	stock := warehouse.Stock{}
	err := stock.PlaceOutbound("test", 3)
	assert.Equal(t, warehouse.ErrOutboundTypeNotFound, err)
}

func TestOutbound1(t *testing.T) {
	t.Run("can not place outbound order when inventory is empty", func(t *testing.T) {
		stock := warehouse.Stock{}
		stock.ConfigureInboundType("test")
		items := []warehouse.OutboundItemComponent{
			{
				ItemType: "test",
				Qty:      1,
			},
		}
		stock.ConfigureOutbound("goesOut", items)
		err := stock.PlaceOutbound("goesOut", 3)
		assert.Equal(t, warehouse.ErrNotEnoughStock, err)
	})
}

func TestOutbound2(t *testing.T) {
	t.Run("can not place outbound order when inventory is not enough", func(t *testing.T) {
		stock := warehouse.Stock{}
		stock.ConfigureInboundType("test")
		items := []warehouse.OutboundItemComponent{
			{
				ItemType: "test",
				Qty:      2,
			},
		}
		stock.ConfigureOutbound("test", items)

		oneTest := warehouse.InboundItem{
			Type: "test",
			Qty:  1,
		}
		stock.Provision(oneTest)

		err := stock.PlaceOutbound("test", 1)
		assert.Equal(t, warehouse.ErrNotEnoughStock, err)
	})
}

func TestOutbound3(t *testing.T) {
	t.Run("should update stock when outbound succeeds", func(t *testing.T) {
		stock := warehouse.Stock{}
		stock.ConfigureInboundType("carrot")

		items := []warehouse.OutboundItemComponent{
			{
				ItemType: "carrot",
				Qty:      1,
			},
		}
		stock.ConfigureOutbound("salad", items)

		carrotsInStock := warehouse.InboundItem{
			Type: "carrot",
			Qty:  2,
		}
		stock.Provision(carrotsInStock)

		err := stock.PlaceOutbound("salad", 1)
		assert.Nil(t, err)

		qty, _ := stock.Quantity("carrot")
		assert.Equal(t, 1, qty)
	})

	t.Run("should update stock when outbound succeeds for compound items", func(t *testing.T) {
		stock := warehouse.Stock{}
		stock.ConfigureInboundType("carrot")
		stock.ConfigureInboundType("cabbage")

		saladConfig := []warehouse.OutboundItemComponent{
			{
				ItemType: "carrot",
				Qty:      1,
			},
			{
				ItemType: "cabbage",
				Qty: 1,
			},
		}
		stock.ConfigureOutbound("salad", saladConfig)

		twoCarrots := warehouse.InboundItem{
			Type: "carrot",
			Qty:  2,
		}
		stock.Provision(twoCarrots)

		threeCabbages := warehouse.InboundItem{
			Type: "cabbage",
			Qty:  2,
		}
		stock.Provision(threeCabbages)

		err := stock.PlaceOutbound("salad", 1)
		assert.Nil(t, err)

		qty, _ := stock.Quantity("carrot")
		assert.Equal(t, 1, qty)

		qty, _ = stock.Quantity("cabbage")
		assert.Equal(t, 1, qty)
	})
}

func TestConfigureOutbound(t *testing.T) {

	assert := assert.New(t)

	t.Run("empty name not accepted", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureOutbound("", nil)
		assert.Equal(warehouse.ErrNameNotProvided, err)
	})

	t.Run("empty list of outbound items not accepted", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureOutbound("mocha", []warehouse.OutboundItemComponent{})
		assert.Equal(warehouse.ErrItemsNotProvided, err)
	})

	t.Run("can not add finished product that has unknown item types", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureOutbound("mocha", []warehouse.OutboundItemComponent{{
			ItemType: "nope",
			Qty:      1,
		}})
		assert.Equal(warehouse.ErrItemTypeNotFound, err)
	})

	t.Run("can not add finished product that has zero quantity", func(t *testing.T) {
		stock := warehouse.Stock{}
		stock.ConfigureInboundType("milk")
		err := stock.ConfigureOutbound("mocha", []warehouse.OutboundItemComponent{{
			ItemType: warehouse.InboundType("milk"),
			Qty:      0,
		}})
		assert.Equal(warehouse.ErrZeroQuantityNotAllowed, err)
	})

	t.Run("can add finished product", func(t *testing.T) {
		stock := warehouse.Stock{}
		stock.ConfigureInboundType("milk")
		err := stock.ConfigureOutbound("mocha", []warehouse.OutboundItemComponent{{
			ItemType: warehouse.InboundType("milk"),
			Qty:      5,
		}})
		assert.Nil(err)

		list := stock.Outbounds()
		assert.Len(list, 1)

		assert.Len(list[0].Items, 1)
	})
}
