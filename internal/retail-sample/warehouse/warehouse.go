package warehouse

import (
	"errors"
	"time"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/recipe"
)

type Inventory interface {
	Add(string) (int, error)
	All() []inventory.Item
	Get(int) inventory.Item
	Find(string) int
}

type RecipeBook interface {
	Add(string, []recipe.Ingredient) error
	Get(int) recipe.Recipe
}

type ( //log
	InboundLog interface {
		Add(ProvisionEntry)
		List() []ProvisionEntry
	}

	OutboundLog interface {
		Add(OrderLogEntry)
		List() []OrderLogEntry
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

func isPresent(i Inventory, id int) bool {
	return i.Get(id) != inventory.Item{}
}

type Position struct {
	Name string
	Qty  int
}

func (s Stock) CurrentState() (ps []Position) {
	for _, item := range s.inventory.All() {
		itemID := int(item.ID)
		qty := s.Quantity(itemID)
		ps = append(ps, Position{
			Name: string(item.Name),
			Qty:  qty,
		})
	}

	return
}

var ErrInventoryItemNotFound = errors.New("inventory item not found")

type ProvisionEntry struct {
	Time time.Time
	ID   int
	Qty  int
}

func (s Stock) Provision(id, qty int) (int, error) {
	if !isPresent(s.inventory, id) {
		return 0, ErrInventoryItemNotFound
	}

	newQty := s.data[id] + qty

	s.data[id] = newQty

	s.inboundLog.Add(ProvisionEntry{
		ID:  id,
		Qty: qty,
	})

	return newQty, nil
}

func (s *Stock) AddInventoryName(typeName string) (int, error) {
	id, err := s.inventory.Add(typeName)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s Stock) Quantity(id int) int {
	return s.data[id]
}

func (s Stock) InventoryItems() (r []inventory.Item) {
	r = append(r, s.inventory.All()...)

	return
}

func (s Stock) ProvisionLog() (r []ProvisionEntry) {
	r = append(r, s.inboundLog.List()...)

	return
}

var (
	ErrOutboundNameNotProvided = errors.New("name not provided")
	ErrIngredientsNotProvided  = errors.New("items not provided")
	ErrZeroQuantityNotAllowed  = errors.New("zero quantity not allowed")
)

func (s *Stock) AddRecipe(name string, items []recipe.Ingredient) error {
	if len(name) == 0 {
		return ErrOutboundNameNotProvided
	}

	if len(items) == 0 {
		return ErrIngredientsNotProvided
	}

	for _, item := range items {

		if !isPresent(s.inventory, item.ID) {
			return ErrInventoryItemNotFound
		}

		if item.Qty == 0 {
			return ErrZeroQuantityNotAllowed
		}
	}

	return s.recipeBook.Add(name, nil)
}

var (
	ErrRecipeNotFound = errors.New("outbound type not found")
	ErrNotEnoughStock = errors.New("not enough stock")
)

func (s *Stock) PlaceOrder(id int, qty int) error {
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

	s.outboundLog.Add(OrderLogEntry{
		Date: time.Now(),
		Qty:  qty,
	})

	return nil
}

type OrderLogEntry struct {
	Date time.Time
	Name string
	Qty  int
}

func (s *Stock) OrderLog() (r []OrderLogEntry) {
	r = append(r, s.outboundLog.List()...)

	return
}
