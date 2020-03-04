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

	ProvisionLog interface {
		Add(ProvisionEntry)
		List() []ProvisionEntry
	}

	OrderLogEntry struct {
		RecipeID int
		Date     time.Time
		Qty      int
	}

	OrderLog interface {
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
	for _, item := range s.inventory.All() {
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

	if s.inventory.Get(itemID) == zeroInventoryItem {
		return 0, ErrInventoryItemNotFound
	}

	newQty := s.data[id] + qty

	s.data[id] = newQty

	s.provisionLog.Add(ProvisionEntry{
		ID:  id,
		Qty: qty,
	})

	return newQty, nil
}

func (s Stock) Quantity(id int) int {
	return s.data[id]
}

func (s Stock) ProvisionLog() (r []ProvisionEntry) {
	r = append(r, s.provisionLog.List()...)

	return
}

var (
	ErrRecipeNotFound = errors.New("outbound type not found")
	ErrNotEnoughStock = errors.New("not enough stock")
)

func (s *Stock) PlaceOrder(id int, qty int) error {
	recipeID := recipe.ID(id)
	r := s.recipeBook.Get(recipeID)

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

	s.orderLog.Add(OrderLogEntry{
		RecipeID: id,
		Date:     time.Now(),
		Qty:      qty,
	})

	return nil
}

func (s *Stock) OrderLog() (r []OrderLogEntry) {
	r = append(r, s.orderLog.List()...)

	return
}
