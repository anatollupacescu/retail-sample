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

		stock := warehouse.NewStock(nil, nil, &config, nil)

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

		mockInv := warehouse.MockInventory{}
		mockInv.On("Get", potato).Return(3)

		stock := warehouse.NewStock(nil, &mockInv, &config, nil)

		err := stock.PlaceOutbound(chips, 2)

		assert.Equal(t, warehouse.ErrNotEnoughStock, err)
		mockInv.AssertExpectations(t)
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

		mockInv := warehouse.MockInventory{}
		mockInv.On("Get", potato).Return(5).Times(2)

		outboundLog := warehouse.MockOutboundLog{}
		outboundLog.On("Add", mock.Anything).Times(1)

		data := map[int]int{
			5: 5,
		}
		stock := warehouse.NewStockWithData(nil, &mockInv, &config, &outboundLog, data)

		err := stock.PlaceOutbound(chips, 2)

		assert.NoError(t, err)
		config.AssertExpectations(t)
		mockInv.AssertExpectations(t)
		outboundLog.AssertExpectations(t)
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
		mockInv := warehouse.MockInventory{}
		mockInv.On("Get", milk).Return(0)

		stock := warehouse.NewStock(nil, &mockInv, nil, nil)

		err := stock.ConfigureOutbound("mocha", []warehouse.OutboundItemComponent{{
			ItemType: milk,
			Qty:      1,
		}})

		assert.Equal(t, warehouse.ErrInboundItemTypeNotFound, err)
		mockInv.AssertExpectations(t)
	})

	t.Run("should reject when request has zero quantity", func(t *testing.T) {
		milk := "milk"
		mockInv := warehouse.MockInventory{}
		mockInv.On("Get", milk).Return(1)

		stock := warehouse.NewStock(nil, &mockInv, nil, nil)

		err := stock.ConfigureOutbound("mocha", []warehouse.OutboundItemComponent{{
			ItemType: "milk",
			Qty:      0,
		}})

		assert.Equal(t, warehouse.ErrOutboundZeroQuantityNotAllowed, err)
		mockInv.AssertExpectations(t)
	})

	t.Run("should accept when configured correctly", func(t *testing.T) {
		milk := "milk"
		mockInv := warehouse.MockInventory{}
		mockInv.On("Get", milk).Return(1)

		outboundConfig := warehouse.MockOutboundConfig{}
		outboundConfig.On("add", mock.AnythingOfType("warehouse.OutboundItem")).Times(1)

		stock := warehouse.NewStock(nil, &mockInv, &outboundConfig, nil)

		err := stock.ConfigureOutbound("mocha", []warehouse.OutboundItemComponent{{
			ItemType: "milk",
			Qty:      5,
		}})
		assert.Nil(t, err)

		mockInv.AssertExpectations(t)
		outboundConfig.AssertExpectations(t)
	})
}

func TestPlaceInbound(t *testing.T) {

	t.Run("should reject stock item with non existent type", func(t *testing.T) {
		milk := "milk"
		mockInv := warehouse.MockInventory{}
		mockInv.On("Get", milk).Return(0)

		stock := warehouse.NewStock(nil, &mockInv, nil, nil)

		item := warehouse.ProvisionEntry{
			Type: milk,
			Qty:  31,
		}

		_, err := stock.PlaceInbound(item)

		assert.Equal(t, warehouse.ErrInboundItemTypeNotFound, err)
		mockInv.AssertExpectations(t)
	})

	t.Run("should place inbound when item type exists", func(t *testing.T) {
		milk := "milk"
		mockInv := warehouse.MockInventory{}
		mockInv.On("Get", milk).Return(1)

		inboundLog := warehouse.MockInboundLog{}
		inboundLog.On("Add", mock.Anything, mock.Anything).Times(1)

		data := map[int]int{
			1: 9,
		}
		stock := warehouse.NewStockWithData(&inboundLog, &mockInv, nil, nil, data)

		item := warehouse.ProvisionEntry{
			Type: milk,
			Qty:  31,
		}

		qty, err := stock.PlaceInbound(item)

		assert.NoError(t, err)
		assert.Equal(t, 40, qty)

		inboundLog.AssertExpectations(t)
		mockInv.AssertExpectations(t)
	})
}
