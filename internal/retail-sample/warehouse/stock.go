package warehouse

import (
	"errors"
	"time"
)

type (

	Inventory interface {
		setQty(string, int)
		qty(string) int
		addType(string)
		hasType(string) bool
		types() []string
	}

	OutboundConfiguration interface {
		add(OutboundItem)
		list() []OutboundItem
		hasConfig(string) bool
		components(string) []OutboundItemComponent
	}

	Stock struct {
		inboundLog            Log
		inventory             Inventory
		outboundConfiguration OutboundConfiguration
	}
)

func NewInMemoryStock() Stock {
	return Stock{
		inboundLog:            make(InMemoryInboundLog),
		inventory:             make(InMemoryInventory),
		outboundConfiguration: make(InMemoryOutboundConfiguration),
	}
}

var ErrInboundItemTypeNotFound = errors.New("type not found")

func (s Stock) PlaceInbound(item Item) (int, error) {
	if !s.inventory.hasType(item.Type) {
		return 0, ErrInboundItemTypeNotFound
	}

	newQty := s.inventory.qty(item.Type) + item.Qty

	s.inventory.setQty(item.Type, newQty)

	s.inboundLog.Add(time.Now(), item)

	return newQty, nil
}

var (
	ErrInboundItemTypeAlreadyConfigured = errors.New("item type already present")
	ErrInboundNameNotProvided           = errors.New("name not provided")
)

func (s *Stock) ConfigureInboundType(typeName string) error {
	if len(typeName) == 0 {
		return ErrInboundNameNotProvided
	}

	if s.inventory.hasType(typeName) {
		return ErrInboundItemTypeAlreadyConfigured
	}

	s.inventory.addType(typeName)

	return nil
}

var ErrInventoryItemNotFound = errors.New("inventory item not found")

func (s Stock) Quantity(typeName string) (int, error) {

	if !s.inventory.hasType(typeName) {
		return 0, ErrInventoryItemNotFound
	}

	qty := s.inventory.qty(typeName)

	return qty, nil
}

func (s Stock) ItemTypes() (r []string) {
	return s.inventory.types()
}

func (s Stock) ListInbound() (r []Item) {
	return s.inboundLog.List()
}
