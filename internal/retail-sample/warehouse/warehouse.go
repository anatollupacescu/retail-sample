package warehouse

import (
	"errors"
	"time"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/recipe"
)

type Inventory interface {
	Add(string) (int, error)
	All() []inventory.Record
	Get(int) string
	Find(string) int
}

type RecipeBook interface {
	Add(string, []recipe.Ingredient) error
	Get(int) recipe.Recipe
}

type ( //log
	InboundLog interface {
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
		ID   int
		Qty  int
	}
)

type Stock struct {
	inventory   Inventory
	inboundLog  InboundLog
	outboundLog OutboundLog
	recipeBook  RecipeBook
	data        map[int]int
}

func NewStock(log InboundLog, inv Inventory, config RecipeBook, outboundItemLog OutboundLog) Stock {
	return Stock{
		inboundLog:  log,
		outboundLog: outboundItemLog,
		inventory:   inv,
		recipeBook:  config,
	}
}

func NewStockWithData(log InboundLog, inv Inventory, config RecipeBook, outboundItemLog OutboundLog, d map[int]int) Stock {
	return Stock{
		inboundLog:  log,
		outboundLog: outboundItemLog,
		inventory:   inv,
		recipeBook:  config,
		data:        d,
	}
}

var ErrInventoryNameNotFound = errors.New("name not found")

func (s Stock) GetType(_ string) ItemConfig {
	return ItemConfig{}
}

var zeroValueName = ""

func isPresent(i Inventory, id int) bool {
	return i.Get(id) != zeroValueName
}

func (s Stock) PlaceInbound(item ProvisionEntry) (int, error) {
	if !isPresent(s.inventory, item.ID) {
		return 0, ErrInventoryNameNotFound
	}

	id := item.ID

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

var zeroValueID = 0

func (s Stock) Quantity(typeName string) (int, error) {
	id := s.inventory.Find(typeName)

	if id == zeroValueID {
		return 0, ErrInventoryItemNotFound
	}

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

func (s *Stock) ConfigureOutbound(name string, items []recipe.Ingredient) error {
	if len(name) == 0 {
		return ErrOutboundNameNotProvided
	}

	if len(items) == 0 {
		return ErrOutboundItemsNotProvided
	}

	for _, item := range items {

		if !isPresent(s.inventory, item.ID) {
			return ErrInventoryNameNotFound
		}

		if item.Qty == 0 {
			return ErrOutboundZeroQuantityNotAllowed
		}
	}

	return s.recipeBook.Add(name, nil)
}

// func (s *Stock) OutboundConfigurations() []OutboundItem {
// return s.outboundConfiguration.list()
// }

var (
	ErrRecipeNotFound = errors.New("outbound type not found")
	ErrNotEnoughStock = errors.New("not enough stock")
)

var zeroValueRecipe = recipe.Recipe{}

func (s *Stock) PlaceOutbound(id int, qty int) error {
	r := s.recipeBook.Get(id)

	ingredients := r.Ingredients

	if r.Ingredients == nil {
		return ErrRecipeNotFound
	}

	for _, i := range ingredients {
		presentQty := s.data[i.ID]
		requestedQty := i.Qty * qty
		if requestedQty > presentQty {
			return ErrNotEnoughStock
		}
	}

	for _, i := range ingredients {
		presentQty := s.data[i.ID]
		requestedQty := i.Qty * qty
		s.data[i.ID] = presentQty - requestedQty
	}

	s.outboundLog.Add(SoldItem{
		Date: time.Now(),
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
	return s.outboundLog.List(), nil
}
