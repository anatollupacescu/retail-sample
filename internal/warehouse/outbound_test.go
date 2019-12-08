package warehouse_test

import (
	"github.com/anatollupacescu/retail-sample/internal/warehouse"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInventoryTypesAreNotConfigured(t *testing.T) {
	t.Run("should not allow placing outbound order", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.PlaceOutbound("test", 3)
		assert.Equal(t, warehouse.ErrOutboundTypeNotFound, err)
	})
}

func TestInventoryIsEmpty(t *testing.T) {
	t.Run("not allw placing outbound order", func(t *testing.T) {
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

func TestInventoryDoesNotContainEnoughItems(t *testing.T) {
	t.Run("can not place outbound order", func(t *testing.T) {
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

func TestStockHasEnoughItems(t *testing.T) {
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
				Qty:      1,
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

	t.Run("can place outbound for multiple items", func(t *testing.T) {
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
				Qty:      1,
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

		err := stock.PlaceOutbound("salad", 2)
		assert.Nil(t, err)

		qty, _ := stock.Quantity("carrot")
		assert.Equal(t, 0, qty)

		qty, _ = stock.Quantity("cabbage")
		assert.Equal(t, 0, qty)
	})
}

func TestConfigureOutbound(t *testing.T) {

	assert := assert.New(t)

	t.Run("should reject empty name", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureOutbound("", nil)
		assert.Equal(warehouse.ErrOutboundNameNotProvided, err)
	})

	t.Run("should reject empty list of outbound items", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureOutbound("mocha", []warehouse.OutboundItemComponent{})
		assert.Equal(warehouse.ErrOutboundItemsNotProvided, err)
	})

	t.Run("should reject when request has unknown item types", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureOutbound("mocha", []warehouse.OutboundItemComponent{{
			ItemType: "nope",
			Qty:      1,
		}})
		assert.Equal(warehouse.ErrInboundItemTypeNotFound, err)
	})

	t.Run("should reject when request has zero quantity", func(t *testing.T) {
		stock := warehouse.Stock{}
		stock.ConfigureInboundType("milk")
		err := stock.ConfigureOutbound("mocha", []warehouse.OutboundItemComponent{{
			ItemType: warehouse.InboundType("milk"),
			Qty:      0,
		}})
		assert.Equal(warehouse.ErrOutboundZeroQuantityNotAllowed, err)
	})

	t.Run("should accept when configured correctly", func(t *testing.T) {
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
