package warehouse

import (
	"errors"
	"time"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/inventory"
)

type ( //outbound

	OutboundConfiguration interface {
		add(OutboundItem)
		list() []OutboundItem
		hasConfig(string) bool
		components(string) []OutboundItemComponent
	}

	OutboundItemComponent struct {
		Name string
		Qty  int
	}

	OutboundItem struct {
		Name  string
		Items []OutboundItemComponent
	}
)

type ( //inventory

	Inventory interface {
		Add(string) (int, error)
		All() []inventory.Record
		Get(int) string
		Find(string) int
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
	inventory             Inventory
	inboundLog            Log
	soldItems             OutboundLog
	outboundConfiguration OutboundConfiguration
	data                  map[int]int
}

func NewStock(log Log, inv Inventory, config OutboundConfiguration, outboundItemLog OutboundLog) Stock {
	return Stock{
		inboundLog:            log,
		soldItems:             outboundItemLog,
		inventory:             inv,
		outboundConfiguration: config,
	}
}

func NewStockWithData(log Log, inv Inventory, config OutboundConfiguration, outboundItemLog OutboundLog, d map[int]int) Stock {
	return Stock{
		inboundLog:            log,
		soldItems:             outboundItemLog,
		inventory:             inv,
		outboundConfiguration: config,
		data:                  d,
	}
}

func NewInMemoryStock() Stock {
	inMemoryStore := inventory.NewInMemoryStore()
	return Stock{
		inboundLog:            make(InMemoryInboundLog),
		soldItems:             make(InMemoryOutboundLog),
		outboundConfiguration: make(InMemoryOutboundConfiguration),
		inventory:             inventory.NewInventory(inMemoryStore),
	}
}

var ErrInboundItemTypeNotFound = errors.New("type not found")

func (s Stock) GetType(name string) ItemConfig {
	for _, i := range s.inventory.All() {
		if string(i.Name) == name {
			return ItemConfig{
				Type:     name,
				Disabled: false,
			}
		}
	}

	return ItemConfig{}
}

var zeroID = 0

func isPresent(i Inventory, s string) bool {
	return i.Find(s) != zeroID
}

func (s Stock) Disable(item string) error {

	if !isPresent(s.inventory, item) {
		return ErrInboundItemTypeNotFound
	}

	// s.inventory.disable(item)

	return nil
}

func (s Stock) PlaceInbound(item ProvisionEntry) (int, error) {
	if !isPresent(s.inventory, item.Type) {
		return 0, ErrInboundItemTypeNotFound
	}

	id := s.inventory.Find(item.Type)

	newQty := s.data[id] + item.Qty

	s.data[id] = newQty

	s.inboundLog.Add(time.Now(), item)

	return newQty, nil
}

func (s *Stock) ConfigureInboundType(typeName string) (int, error) {
	id, err := s.inventory.Add(typeName)

	if err != nil {
		return 0, err
	}

	return id, nil
}

var ErrInventoryItemNotFound = errors.New("inventory item not found")

func (s Stock) Quantity(typeName string) (int, error) {
	if !isPresent(s.inventory, typeName) {
		return 0, ErrInventoryItemNotFound
	}

	id := s.inventory.Find(typeName)

	qty := s.data[id]

	return qty, nil
}

func (s Stock) ItemTypes() (r []string) {
	types := s.inventory.All()
	for _, t := range types {
		r = append(r, string(t.Name))
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

		if !isPresent(s.inventory, item.Name) {
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
		id := s.inventory.Find(outboundItem.Name)
		inventoryQty := s.data[id]
		if outboundItem.Qty*qty > inventoryQty {
			return ErrNotEnoughStock
		}
	}

	for _, outboundItem := range components {
		id := s.inventory.Find(outboundItem.Name)
		inventoryQty := s.data[id]
		inventoryQty -= outboundItem.Qty * qty
		s.data[id] = inventoryQty
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
