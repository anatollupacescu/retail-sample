package warehouse_test

import (
	"testing"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/warehouse"

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

	t.Run("should not allow placing outbound order", func(t *testing.T) {
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

		oneTest := warehouse.Item{
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

		carrotsInStock := warehouse.Item{
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

		twoCarrots := warehouse.Item{
			Type: "carrot",
			Qty:  2,
		}
		_, _ = stock.PlaceInbound(twoCarrots)

		threeCabbages := warehouse.Item{
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

		twoCarrots := warehouse.Item{
			Type: "carrot",
			Qty:  2,
		}
		_, _ = stock.PlaceInbound(twoCarrots)

		threeCabbages := warehouse.Item{
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

//test helper
type TestWrapper struct {
	stock warehouse.Stock
}

func stock() TestWrapper {
	return TestWrapper{
		stock: warehouse.NewInMemoryStock(),
	}
}

func (t TestWrapper) with(s string, q int) TestWrapper {
	if err := t.stock.ConfigureInboundType(s); err != nil {
		panic(err)
	}
	item := warehouse.Item{Type: s, Qty: q}
	if _, err := t.stock.PlaceInbound(item); err != nil {
		panic(err)
	}
	return t
}

func (t TestWrapper) build() warehouse.Stock {
	return t.stock
}
