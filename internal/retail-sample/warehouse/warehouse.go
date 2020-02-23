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
	Names() []string
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
	var zeroInventoryItem inventory.Item
	if s.inventory.Get(id) == zeroInventoryItem {
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

func (s *Stock) AddInventoryName(name string) (int, error) {
	id, err := s.inventory.Add(name)

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

func (s *Stock) AddRecipe(name string, ingredients []recipe.Ingredient) error {
	return s.recipeBook.Add(name, ingredients)
}

var (
	ErrRecipeNotFound = errors.New("outbound type not found")
	ErrNotEnoughStock = errors.New("not enough stock")
)

func (s *Stock) PlaceOrder(id int, qty int) error {
	r := s.recipeBook.Get(id)

	ingredients := r.Ingredients

	if ingredients == nil {
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

func (s *Stock) RecipeNames() (r []string) {
	r = append(r, s.recipeBook.Names()...)

	return
}
