package warehouse_test

import (
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/warehouse"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInventoryTypesAreNotConfigured(t *testing.T) {

	t.Run("should not allow placing outbound order", func(t *testing.T) {
		stock := stock().build()
		err := stock.PlaceOutbound("test", 3)
		assert.Equal(t, warehouse.ErrOutboundTypeNotFound, err)
	})
}

func TestInventoryIsEmpty(t *testing.T) {

	t.Run("not allow placing outbound order", func(t *testing.T) {
		stock := stock().build()
		_ = stock.ConfigureInboundType("test")
		items := []warehouse.OutboundItemComponent{
			{
				ItemType: "test",
				Qty:      1,
			},
		}
		_ = stock.ConfigureOutbound("goesOut", items)
		err := stock.PlaceOutbound("goesOut", 3)
		assert.Equal(t, warehouse.ErrNotEnoughStock, err)
	})
}

func TestInventoryDoesNotContainEnoughItems(t *testing.T) {

	t.Run("should not allow placing outbound order", func(t *testing.T) {
		stock := stock().build()
		_ = stock.ConfigureInboundType("test")
		items := []warehouse.OutboundItemComponent{
			{
				ItemType: "test",
				Qty:      2,
			},
		}
		_ = stock.ConfigureOutbound("test", items)

		oneTest := warehouse.InboundItem{
			Type: "test",
			Qty:  1,
		}
		_, _ = stock.PlaceInbound(oneTest)

		err := stock.PlaceOutbound("test", 1)
		assert.Equal(t, warehouse.ErrNotEnoughStock, err)
	})
}

func TestStockHasEnoughItems(t *testing.T) {

	t.Run("should update stock when outbound succeeds", func(t *testing.T) {
		stock := stock().build()
		_ = stock.ConfigureInboundType("carrot")

		items := []warehouse.OutboundItemComponent{
			{
				ItemType: "carrot",
				Qty:      1,
			},
		}
		_ = stock.ConfigureOutbound("salad", items)

		carrotsInStock := warehouse.InboundItem{
			Type: "carrot",
			Qty:  2,
		}
		_, _ = stock.PlaceInbound(carrotsInStock)

		err := stock.PlaceOutbound("salad", 1)
		assert.Nil(t, err)

		qty, _ := stock.Quantity("carrot")
		assert.Equal(t, 1, qty)
	})

	t.Run("should update stock when outbound succeeds for compound items", func(t *testing.T) {
		stock := stock().build()
		_ = stock.ConfigureInboundType("carrot")
		_ = stock.ConfigureInboundType("cabbage")

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
		_ = stock.ConfigureOutbound("salad", saladConfig)

		twoCarrots := warehouse.InboundItem{
			Type: "carrot",
			Qty:  2,
		}
		_, _ = stock.PlaceInbound(twoCarrots)

		threeCabbages := warehouse.InboundItem{
			Type: "cabbage",
			Qty:  2,
		}
		_, _ = stock.PlaceInbound(threeCabbages)

		err := stock.PlaceOutbound("salad", 1)
		assert.Nil(t, err)

		qty, _ := stock.Quantity("carrot")
		assert.Equal(t, 1, qty)

		qty, _ = stock.Quantity("cabbage")
		assert.Equal(t, 1, qty)
	})

	t.Run("can place outbound for multiple items", func(t *testing.T) {
		stock := stock().build()
		_ = stock.ConfigureInboundType("carrot")
		_ = stock.ConfigureInboundType("cabbage")

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
		_ = stock.ConfigureOutbound("salad", saladConfig)

		twoCarrots := warehouse.InboundItem{
			Type: "carrot",
			Qty:  2,
		}
		_, _ = stock.PlaceInbound(twoCarrots)

		threeCabbages := warehouse.InboundItem{
			Type: "cabbage",
			Qty:  2,
		}
		_, _ = stock.PlaceInbound(threeCabbages)

		err := stock.PlaceOutbound("salad", 2)
		assert.Nil(t, err)

		qty, _ := stock.Quantity("carrot")
		assert.Equal(t, 0, qty)

		qty, _ = stock.Quantity("cabbage")
		assert.Equal(t, 0, qty)
	})
}

func TestConfigureOutbound(t *testing.T) {

	t.Run("should reject empty name", func(t *testing.T) {
		stock := stock().build()
		err := stock.ConfigureOutbound("", nil)
		assert.Equal(t, warehouse.ErrOutboundNameNotProvided, err)
	})

	t.Run("should reject empty list of outbound items", func(t *testing.T) {
		stock := stock().build()
		err := stock.ConfigureOutbound("mocha", []warehouse.OutboundItemComponent{})
		assert.Equal(t, warehouse.ErrOutboundItemsNotProvided, err)
	})

	t.Run("should reject when request has unknown item types", func(t *testing.T) {
		stock := stock().build()
		err := stock.ConfigureOutbound("mocha", []warehouse.OutboundItemComponent{{
			ItemType: "nope",
			Qty:      1,
		}})
		assert.Equal(t, warehouse.ErrInboundItemTypeNotFound, err)
	})

	t.Run("should reject when request has zero quantity", func(t *testing.T) {
		stock := stock().build()
		_ = stock.ConfigureInboundType("milk")
		err := stock.ConfigureOutbound("mocha", []warehouse.OutboundItemComponent{{
			ItemType: "milk",
			Qty:      0,
		}})
		assert.Equal(t, warehouse.ErrOutboundZeroQuantityNotAllowed, err)
	})

	t.Run("should accept when configured correctly", func(t *testing.T) {
		stock := stock().build()
		_ = stock.ConfigureInboundType("milk")
		err := stock.ConfigureOutbound("mocha", []warehouse.OutboundItemComponent{{
			ItemType: "milk",
			Qty:      5,
		}})
		assert.Nil(t, err)

		list := stock.Outbounds()
		assert.Len(t, list, 1)
		assert.Len(t, list[0].Items, 1)
	})
}
