package warehouse

import (
	"errors"
	"time"

	"github.com/anatollupacescu/retail-sample/internal/retail-sample/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/recipe"
)

type Inventory interface {
	Add(inventory.Name) (inventory.ID, error)
	All() []inventory.Item
	Get(inventory.ID) inventory.Item
	Find(inventory.Name) inventory.ID
}

type RecipeBook interface {
	Add(recipe.Name, []recipe.Ingredient) (recipe.ID, error)
	Get(recipe.ID) recipe.Recipe
	Names() []recipe.Name
}

type ( //log
	ProvisionEntry struct {
		Time time.Time
		ID   int
		Qty  int
	}

	InboundLog interface {
		Add(ProvisionEntry)
		List() []ProvisionEntry
	}

	OrderLogEntry struct {
		Date time.Time
		Name string
		Qty  int
	}

	OutboundLog interface {
		Add(OrderLogEntry)
		List() []OrderLogEntry
	}
)

type Position struct {
	ID   int
	Name string
	Qty  int
}

func (s Stock) CurrentState() (ps []Position) {
	for _, item := range s.Inventory.All() {
		itemID := int(item.ID)
		qty := s.Quantity(itemID)
		ps = append(ps, Position{
			ID:   itemID,
			Name: string(item.Name),
			Qty:  qty,
		})
	}

	return
}

var ErrInventoryItemNotFound = errors.New("inventory item not found")

func (s Stock) Provision(id, qty int) (int, error) {
	var zeroInventoryItem inventory.Item

	itemID := inventory.ID(id)

	if s.Inventory.Get(itemID) == zeroInventoryItem {
		return 0, ErrInventoryItemNotFound
	}

	newQty := s.Data[id] + qty

	s.Data[id] = newQty

	s.InboundLog.Add(ProvisionEntry{
		ID:  id,
		Qty: qty,
	})

	return newQty, nil
}

func (s Stock) Quantity(id int) int {
	return s.Data[id]
}

func (s Stock) ProvisionLog() (r []ProvisionEntry) {
	r = append(r, s.InboundLog.List()...)

	return
}

var (
	ErrRecipeNotFound = errors.New("outbound type not found")
	ErrNotEnoughStock = errors.New("not enough stock")
)

func (s *Stock) PlaceOrder(id int, qty int) error {
	recipeID := recipe.ID(id)
	r := s.RecipeBook.Get(recipeID)

	ingredients := r.Ingredients

	if ingredients == nil {
		return ErrRecipeNotFound
	}

	for _, i := range ingredients {
		presentQty := s.Data[i.ID]
		requestedQty := i.Qty * qty
		if requestedQty > presentQty {
			return ErrNotEnoughStock
		}
	}

	for _, i := range ingredients {
		presentQty := s.Data[i.ID]
		requestedQty := i.Qty * qty
		s.Data[i.ID] = presentQty - requestedQty
	}

	s.OutboundLog.Add(OrderLogEntry{
		Date: time.Now(),
		Qty:  qty,
	})

	return nil
}

func (s *Stock) OrderLog() (r []OrderLogEntry) {
	r = append(r, s.OutboundLog.List()...)

	return
}
