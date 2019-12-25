package warehouse_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/warehouse"

	"github.com/stretchr/testify/assert"
)

func TestPlaceOutbound(t *testing.T) {

	t.Run("should reject outbound order with non existent type", func(t *testing.T) {
		chips := "chips"
		config := warehouse.MockOutboundConfig{}
		config.On("hasConfig", chips).Return(false)

		stock := warehouse.NewStock(nil, nil, &config)

		err := stock.PlaceOutbound(chips, 0)

		assert.Equal(t, warehouse.ErrOutboundTypeNotFound, err)
		config.AssertExpectations(t)
	})

	t.Run("should reject outbound when not enough stock", func(t *testing.T) {
		chips := "chips"
		potato := "potato"

		config := warehouse.MockOutboundConfig{}
		config.On("hasConfig", chips).Return(true)
		config.On("components", chips).Return([]warehouse.OutboundItemComponent{{
			ItemType: potato,
			Qty:      2,
		}})

		inventory := warehouse.MockInventory{}
		inventory.On("qty", potato).Return(3)

		stock := warehouse.NewStock(nil, &inventory, &config)

		err := stock.PlaceOutbound(chips, 2)

		assert.Equal(t, warehouse.ErrNotEnoughStock, err)
		inventory.AssertExpectations(t)
	})

	t.Run("should update inventory on success", func(t *testing.T) {
		chips := "chips"
		potato := "potato"

		config := warehouse.MockOutboundConfig{}
		config.On("hasConfig", chips).Return(true)
		config.On("components", chips).Return([]warehouse.OutboundItemComponent{{
			ItemType: potato,
			Qty:      2,
		}})

		inventory := warehouse.MockInventory{}
		inventory.On("qty", potato).Return(5)
		inventory.On("setQty", potato, 1).Times(1)

		stock := warehouse.NewStock(nil, &inventory, &config)

		err := stock.PlaceOutbound(chips, 2)

		assert.NoError(t, err)
		inventory.AssertExpectations(t)
	})
}

func TestConfigureOutbound(t *testing.T) {

	t.Run("should reject empty name", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureOutbound("", nil)
		assert.Equal(t, warehouse.ErrOutboundNameNotProvided, err)
	})

	t.Run("should reject empty list of outbound items", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureOutbound("mocha", []warehouse.OutboundItemComponent{})
		assert.Equal(t, warehouse.ErrOutboundItemsNotProvided, err)
	})

	t.Run("should reject when request has unknown item types", func(t *testing.T) {
		milk := "milk"
		inventory := warehouse.MockInventory{}
		inventory.On("hasType", milk).Return(false)

		stock := warehouse.NewStock(nil, &inventory, nil)

		err := stock.ConfigureOutbound("mocha", []warehouse.OutboundItemComponent{{
			ItemType: milk,
			Qty:      1,
		}})

		assert.Equal(t, warehouse.ErrInboundItemTypeNotFound, err)
		inventory.AssertExpectations(t)
	})

	t.Run("should reject when request has zero quantity", func(t *testing.T) {
		milk := "milk"
		inventory := warehouse.MockInventory{}
		inventory.On("hasType", milk).Return(true)

		stock := warehouse.NewStock(nil, &inventory, nil)

		err := stock.ConfigureOutbound("mocha", []warehouse.OutboundItemComponent{{
			ItemType: "milk",
			Qty:      0,
		}})

		assert.Equal(t, warehouse.ErrOutboundZeroQuantityNotAllowed, err)
		inventory.AssertExpectations(t)
	})

	t.Run("should accept when configured correctly", func(t *testing.T) {
		milk := "milk"
		inventory := warehouse.MockInventory{}
		inventory.On("hasType", milk).Return(true)

		outboundConfig := warehouse.MockOutboundConfig{}
		outboundConfig.On("add", mock.AnythingOfType("warehouse.OutboundItem")).Times(1)

		stock := warehouse.NewStock(nil, &inventory, &outboundConfig)

		err := stock.ConfigureOutbound("mocha", []warehouse.OutboundItemComponent{{
			ItemType: "milk",
			Qty:      5,
		}})
		assert.Nil(t, err)

		inventory.AssertExpectations(t)
		outboundConfig.AssertExpectations(t)
	})
}

func TestConfigureInboundType(t *testing.T) {

	t.Run("should reject empty inbound type name", func(t *testing.T) {
		stock := warehouse.Stock{}
		err := stock.ConfigureInboundType("")
		assert.Equal(t, warehouse.ErrInboundNameNotProvided, err)
	})

	t.Run("should reject duplicate inbound type name", func(t *testing.T) {
		milk := "milk"
		inventory := warehouse.MockInventory{}
		inventory.On("hasType", milk).Return(true)

		stock := warehouse.NewStock(nil, &inventory, nil)

		err := stock.ConfigureInboundType(milk)

		assert.Equal(t, warehouse.ErrInboundItemTypeAlreadyConfigured, err)
		inventory.AssertExpectations(t)
	})

	t.Run("should succeed given a valid inbound type name", func(t *testing.T) {
		milk := "milk"
		inventory := warehouse.MockInventory{}
		inventory.On("hasType", milk).Return(false)
		inventory.On("addType", milk).Times(1)

		stock := warehouse.NewStock(nil, &inventory, nil)

		err := stock.ConfigureInboundType(milk)

		assert.NoError(t, err)
		inventory.AssertExpectations(t)
	})

	t.Run("should return err when getting quantity for non existent item", func(t *testing.T) {
		milk := "milk"
		inventory := warehouse.MockInventory{}
		inventory.On("hasType", milk).Return(false)

		stock := warehouse.NewStock(nil, &inventory, nil)

		_, err := stock.Quantity(milk)

		assert.Equal(t, warehouse.ErrInventoryItemNotFound, err)
		inventory.AssertExpectations(t)
	})
}

func TestPlaceInbound(t *testing.T) {

	t.Run("should reject stock item with non existent type", func(t *testing.T) {
		milk := "milk"
		inventory := warehouse.MockInventory{}
		inventory.On("hasType", milk).Return(false)

		stock := warehouse.NewStock(nil, &inventory, nil)

		item := warehouse.Item{Type: milk, Qty: 31}
		_, err := stock.PlaceInbound(item)

		assert.Equal(t, warehouse.ErrInboundItemTypeNotFound, err)
		inventory.AssertExpectations(t)
	})

	t.Run("should place inbound when item type exists", func(t *testing.T) {
		milk := "milk"
		inventory := warehouse.MockInventory{}
		inventory.On("hasType", milk).Return(true)
		inventory.On("qty", milk).Return(9)
		inventory.On("setQty", milk, 40).Times(1)

		inboundLog := warehouse.MockInboundLog{}
		inboundLog.On("Add", mock.Anything, mock.Anything).Times(1)
		stock := warehouse.NewStock(&inboundLog, &inventory, nil)

		item := warehouse.Item{Type: milk, Qty: 31}
		qty, err := stock.PlaceInbound(item)

		assert.NoError(t, err)
		assert.Equal(t, 40, qty)

		inboundLog.AssertExpectations(t)
		inventory.AssertExpectations(t)
	})
}
