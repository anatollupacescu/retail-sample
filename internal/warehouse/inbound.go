package warehouse

import (
	"errors"
	"time"
)

type (
	InboundType string

	Stock struct {
		inboundLog            map[time.Time]InboundItem
		inventory             map[InboundType]int
		outboundConfiguration map[OutboundType]OutboundItem
	}

	InboundItem struct {
		Type InboundType
		Qty  int
	}
)

func NewStock() Stock {
	return Stock{
		inboundLog:            make(map[time.Time]InboundItem),
		inventory:             make(map[InboundType]int),
		outboundConfiguration: make(map[OutboundType]OutboundItem),
	}
}

var ErrInboundItemTypeNotFound = errors.New("type not found")

func (s Stock) PlaceInbound(item InboundItem) (int, error) {

	if !s.hasType(item.Type) {
		return 0, ErrInboundItemTypeNotFound
	}

	currentQty := s.inventory[item.Type] + item.Qty

	s.inventory[item.Type] = currentQty

	addLogEntry(s, item)

	return s.inventory[item.Type], nil
}

func addLogEntry(s Stock, item InboundItem) {
	s.inboundLog[time.Now()] = item
}

var (
	ErrInboundItemTypeDuplicated = errors.New("item type already present")
	ErrInboundNameNotProvided    = errors.New("name not provided")
)

func (s *Stock) ConfigureInboundType(typeName string) error {

	if len(typeName) == 0 {
		return ErrInboundNameNotProvided
	}

	typeToAdd := InboundType(typeName)

	if s.hasType(typeToAdd) {
		return ErrInboundItemTypeDuplicated
	}

	s.inventory[typeToAdd] = 0

	return nil
}

func (s Stock) hasType(itemType InboundType) bool {
	_, found := s.inventory[itemType]
	return found
}

func (s Stock) Quantity(typeName string) (int, error) {

	itemType := InboundType(typeName)
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

func (s Stock) ListInbound() (r []InboundItem) {
	for _, item := range s.inboundLog {
		r = append(r, item)
	}
	return
}
