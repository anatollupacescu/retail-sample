package warehouse

import (
	"errors"
)

type (
	OutboundType string

	OutboundItemComponent struct {
		ItemType ItemType
		Qty      int
	}

	OutboundItem struct {
		name  string
		Items []OutboundItemComponent
	}
)

var (
	ErrOutboundNameNotProvided        = errors.New("name not provided")
	ErrOutboundItemsNotProvided       = errors.New("items not provided")
	ErrOutboundZeroQuantityNotAllowed = errors.New("zero quantity not allowed")
)

func (s *Stock) ConfigureOutbound(name string, items []OutboundItemComponent) error {

	if len(name) == 0 {
		return ErrOutboundNameNotProvided
	}

	if len(items) == 0 {
		return ErrOutboundItemsNotProvided
	}

	for _, item := range items {
		if !s.hasType(item.ItemType) {
			return ErrInboundItemTypeNotFound
		}
		if item.Qty == 0 {
			return ErrOutboundZeroQuantityNotAllowed
		}
	}

	outboundItem := OutboundItem{
		name:  name,
		Items: items,
	}

	if s.outboundConfiguration == nil {
		s.outboundConfiguration = make(map[OutboundType]OutboundItem)
	}
	s.outboundConfiguration[OutboundType(name)] = outboundItem

	return nil
}

func (s *Stock) Outbounds() (result []OutboundItem) {
	for _, v := range s.outboundConfiguration {
		result = append(result, v)
	}
	return
}

var (
	ErrOutboundTypeNotFound = errors.New("outbound type not found")
	ErrNotEnoughStock       = errors.New("not enough stock")
)

func (s *Stock) PlaceOutbound(typeName string, qty int) error {

	outboundType := OutboundType(typeName)

	config, found := s.outboundConfiguration[outboundType]

	if !found {
		return ErrOutboundTypeNotFound
	}

	for _, outboundItem := range config.Items {
		inventoryQty := s.inventory[outboundItem.ItemType]
		if outboundItem.Qty*qty > inventoryQty {
			return ErrNotEnoughStock
		}
	}

	for _, outboundItem := range config.Items {
		inventoryQty := s.inventory[outboundItem.ItemType]
		inventoryQty -= outboundItem.Qty * qty
		s.inventory[outboundItem.ItemType] = inventoryQty
	}

	return nil
}
