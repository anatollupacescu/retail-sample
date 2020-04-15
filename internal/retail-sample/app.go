package retailsample

import (
	"errors"
	"time"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/order"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
)

type App struct {
	//domain
	Inventory  Inventory
	Orders     Orders
	RecipeBook RecipeBook
	//app specific
	ProvisionLog ProvisionLog
	Stock        Stock
}

type Inventory interface {
	Add(inventory.Name) (inventory.ID, error)
	All() []inventory.Item
	Get(inventory.ID) inventory.Item
	Find(inventory.Name) inventory.ID
}

type RecipeBook interface {
	Add(recipe.Name, []recipe.Ingredient) (recipe.ID, error)
	Get(recipe.ID) recipe.Recipe
	All() []recipe.Recipe
}

type Orders interface {
	Add(order.OrderEntry) order.ID
	All() []order.Order
}

type ( //provision log
	ProvisionEntry struct {
		Time time.Time
		ID   int
		Qty  int
	}

	ProvisionLog interface {
		Add(ProvisionEntry)
		List() []ProvisionEntry
	}
)

type (
	Stock interface {
		Quantity(int) int
		Provision(int, int) int
		Sell([]recipe.Ingredient, int) error
	}

	StockPosition struct {
		ID   int
		Name string
		Qty  int
	}
)

func (a App) CurrentStock() (ps []StockPosition) {
	for _, item := range a.Inventory.All() {
		itemID := int(item.ID)
		qty := a.Stock.Quantity(itemID)
		ps = append(ps, StockPosition{
			ID:   itemID,
			Name: string(item.Name),
			Qty:  qty,
		})
	}

	return
}

var ErrInventoryItemNotFound = errors.New("inventory item not found")

func (a App) Provision(id, qty int) (int, error) {
	var zeroInventoryItem inventory.Item

	itemID := inventory.ID(id)

	if a.Inventory.Get(itemID) == zeroInventoryItem {
		return 0, ErrInventoryItemNotFound
	}

	//TODO should provision qty=0?

	newQty := a.Stock.Provision(id, qty)

	a.ProvisionLog.Add(ProvisionEntry{
		ID:  id,
		Qty: qty,
	})

	return newQty, nil
}

func (a App) Quantity(id int) int {
	return a.Stock.Quantity(id)
}

func (a App) GetProvisionLog() (r []ProvisionEntry) {
	r = append(r, a.ProvisionLog.List()...)

	return
}

var ErrRecipeNotFound = errors.New("outbound type not found")

func (a App) PlaceOrder(id int, qty int) (order.ID, error) {
	recipeID := recipe.ID(id)
	r := a.RecipeBook.Get(recipeID)

	ingredients := r.Ingredients

	if ingredients == nil {
		return 0, ErrRecipeNotFound
	}

	//TODO zero qty?

	if err := a.Stock.Sell(ingredients, qty); err != nil {
		switch err {
		case ErrNotEnoughStock:
			return 0, ErrNotEnoughStock
		default:
			panic("unexpected error")
		}
	}

	entryID := a.Orders.Add(order.OrderEntry{
		RecipeID: id,
		Qty:      qty,
	})

	return entryID, nil
}
