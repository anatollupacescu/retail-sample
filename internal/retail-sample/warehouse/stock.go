package warehouse

import (
	"errors"
	"time"
)

type (
	Inventory             map[ItemType]int
	OutboundConfiguration map[OutboundType]OutboundItem

	Stock struct {
		inboundLog            Log
		inventory             Inventory
		outboundConfiguration OutboundConfiguration
	}
)

func NewStock() Stock {
	return Stock{
		inboundLog:            make(InMemoryInboundLog),
		inventory:             make(map[ItemType]int),
		outboundConfiguration: make(map[OutboundType]OutboundItem),
	}
}

var ErrInboundItemTypeNotFound = errors.New("type not found")

func (s Stock) PlaceInbound(item Item) (int, error) {
	if !s.hasType(item.Type) {
		return 0, ErrInboundItemTypeNotFound
	}
	currentQty := s.inventory[item.Type] + item.Qty
	s.inventory[item.Type] = currentQty
	s.inboundLog.Add(time.Now(), item)
	return currentQty, nil
}

var (
	ErrInboundItemTypeDuplicated = errors.New("item type already present")
	ErrInboundNameNotProvided    = errors.New("name not provided")
)

func (s *Stock) ConfigureInboundType(typeName string) error {
	if len(typeName) == 0 {
		return ErrInboundNameNotProvided
	}

	typeToAdd := ItemType(typeName)

	if s.hasType(typeToAdd) {
		return ErrInboundItemTypeDuplicated
	}

	s.inventory[typeToAdd] = 0

	return nil
}

func (s Stock) hasType(itemType ItemType) bool {
	_, found := s.inventory[itemType]
	return found
}

func (s Stock) Quantity(typeName string) (int, error) {
	itemType := ItemType(typeName)
	qty, found := s.inventory[itemType]

	if !found {
		return 0, ErrInboundItemTypeDuplicated
	}

	return qty, nil
}

func (s Stock) ItemTypes() (r []string) {
	for key := range s.inventory {
		r = append(r, string(key))
	}
	return
}

func (s Stock) ListInbound() (r []Item) {
	return s.inboundLog.List()
}
