package warehouse

import (
	"errors"
	"time"
)

type (
	InboundType string

	InboundLog interface {
		Add(time.Time, InboundItem)
		List() []InboundItem
	}

	inMemoryInboundLog map[time.Time]InboundItem

	Inventory             map[InboundType]int
	OutboundConfiguration map[OutboundType]OutboundItem

	Stock struct {
		inboundLog            InboundLog
		inventory             Inventory
		outboundConfiguration OutboundConfiguration
	}

	InboundItem struct {
		Type InboundType
		Qty  int
	}
)

func NewStock() Stock {
	return Stock{
		inboundLog:            make(inMemoryInboundLog),
		inventory:             make(map[InboundType]int),
		outboundConfiguration: make(map[OutboundType]OutboundItem),
	}
}

func (i inMemoryInboundLog) Add(k time.Time, v InboundItem) {
	i[k] = v
}

func (i inMemoryInboundLog) List() (r []InboundItem) {
	for _, v := range i {
		r = append(r, v)
	}
	return
}

var ErrInboundItemTypeNotFound = errors.New("type not found")

func (s Stock) PlaceInbound(item InboundItem) (int, error) {

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
	for _, item := range s.inboundLog.List() {
		r = append(r, item)
	}
	return
}
