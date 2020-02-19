package warehouse

import (
	"errors"
	"time"
)

type ( //outbound

	OutboundConfiguration interface {
		add(OutboundItem)
		list() []OutboundItem
		hasConfig(string) bool
		components(string) []OutboundItemComponent
	}

	OutboundItemComponent struct {
		ItemType string
		Qty      int
	}

	OutboundItem struct {
		Name  string
		Items []OutboundItemComponent
	}
)

type ( //inventory

	Inventory interface {
		setQty(string, int)
		qty(string) int
		addType(string) int
		hasType(string) bool
		types() []ItemConfig
		disable(string)
	}

	Log interface {
		Add(time.Time, ProvisionEntry)
		List() []ProvisionEntry
	}

	OutboundLog interface {
		Add(SoldItem)
		List() []SoldItem
	}

	ItemConfig struct {
		Type     string
		Disabled bool
	}

	ProvisionEntry struct {
		Time time.Time
		Type string
		Qty  int
	}
)

type Stock struct {
	inboundLog            Log
	soldItems             OutboundLog
	inventory             Inventory
	outboundConfiguration OutboundConfiguration
}

func NewStock(log Log, inv Inventory, config OutboundConfiguration, outboundItemLog OutboundLog) Stock {
	return Stock{
		inboundLog:            log,
		soldItems:             outboundItemLog,
		inventory:             inv,
		outboundConfiguration: config,
	}
}

func NewInMemoryStock() Stock {
	zero := 0
	return Stock{
		inboundLog:            make(InMemoryInboundLog),
		soldItems:             make(InMemoryOutboundLog),
		outboundConfiguration: make(InMemoryOutboundConfiguration),
		inventory: &InMemoryInventory{
			data:    make(map[string]int),
			config:  make(map[string]bool),
			ids:     make(map[string]int),
			counter: &zero,
		},
	}
}

var ErrInboundItemTypeNotFound = errors.New("type not found")

func (s Stock) GetType(name string) ItemConfig {
	for _, i := range s.inventory.types() {
		if i.Type == name {
			return i
		}
	}

	return ItemConfig{}
}

func (s Stock) Disable(item string) error {

	if !s.inventory.hasType(item) {
		return ErrInboundItemTypeNotFound
	}

	s.inventory.disable(item)

	return nil
}

func (s Stock) PlaceInbound(item ProvisionEntry) (int, error) {
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

func (s *Stock) ConfigureInboundType(typeName string) (int, error) {
	if len(typeName) == 0 {
		return 0, ErrInboundNameNotProvided
	}

	if s.inventory.hasType(typeName) {
		return 0, ErrInboundItemTypeAlreadyConfigured
	}

	id := s.inventory.addType(typeName)

	return id, nil
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
	types := s.inventory.types()
	for _, t := range types {
		r = append(r, t.Type)
	}
	return
}

func (s Stock) ListInbound() (r []ProvisionEntry) {
	return s.inboundLog.List()
}

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

		if !s.inventory.hasType(item.ItemType) {
			return ErrInboundItemTypeNotFound
		}

		if item.Qty == 0 {
			return ErrOutboundZeroQuantityNotAllowed
		}
	}

	outboundItem := OutboundItem{
		Name:  name,
		Items: items,
	}

	s.outboundConfiguration.add(outboundItem)

	return nil
}

func (s *Stock) OutboundConfigurations() []OutboundItem {
	return s.outboundConfiguration.list()
}

var (
	ErrOutboundTypeNotFound = errors.New("outbound type not found")
	ErrNotEnoughStock       = errors.New("not enough stock")
)

func (s *Stock) PlaceOutbound(typeName string, qty int) error {

	if !s.outboundConfiguration.hasConfig(typeName) {
		return ErrOutboundTypeNotFound
	}

	components := s.outboundConfiguration.components(typeName)

	for _, outboundItem := range components {
		inventoryQty := s.inventory.qty(outboundItem.ItemType)
		if outboundItem.Qty*qty > inventoryQty {
			return ErrNotEnoughStock
		}
	}

	for _, outboundItem := range components {
		inventoryQty := s.inventory.qty(outboundItem.ItemType)
		inventoryQty -= outboundItem.Qty * qty
		s.inventory.setQty(outboundItem.ItemType, inventoryQty)
	}

	s.soldItems.Add(SoldItem{
		Date: time.Now(),
		Name: typeName,
		Qty:  qty,
	})

	return nil
}

type SoldItem struct {
	Date time.Time
	Name string
	Qty  int
}

func (s *Stock) ListOutbound() ([]SoldItem, error) {
	return s.soldItems.List(), nil
}
