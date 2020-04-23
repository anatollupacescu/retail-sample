package retailsample

import (
	"errors"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/order"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
)

//Facade/Use cases
type App struct {
	//domain
	Inventory  Inventory
	Orders     Orders
	RecipeBook RecipeBook

	//app specific
	ProvisionLog ProvisionLog
	Stock        Stock

	//unit-of-work
	PersistentProviderFactory PersistenceProviderFactory
}

type PersistenceProviderFactory interface {
	Begin() PersistenceProvider
	Commit(PersistenceProvider)
	Rollback(PersistenceProvider)
}

type PersistenceProvider interface {
	Inventory() Inventory
	Stock() Stock
	ProvisionLog() ProvisionLog
	RecipeBook() RecipeBook
	Orders() Orders
}

type Inventory interface {
	Add(inventory.Name) (inventory.ID, error)
	List() []inventory.Item
	Get(inventory.ID) inventory.Item
	Find(inventory.Name) inventory.ID
}

type RecipeBook interface {
	Add(recipe.Name, []recipe.Ingredient) (recipe.ID, error)
	Get(recipe.ID) recipe.Recipe
	List() []recipe.Recipe
}

type Orders interface {
	Add(order.OrderEntry) order.ID
	List() []order.Order
}

type ( //provision log
	ProvisionEntry struct {
		ID  int
		Qty int
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
	for _, item := range a.Inventory.List() {
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

	provider := a.PersistentProviderFactory.Begin()
	defer a.PersistentProviderFactory.Commit(provider)

	inv := provider.Inventory()

	itemID := inventory.ID(id)

	if inv.Get(itemID) == zeroInventoryItem {
		return 0, ErrInventoryItemNotFound
	}

	stock := provider.Stock()

	newQty := stock.Provision(id, qty)

	provisionLog := provider.ProvisionLog()

	provisionLog.Add(ProvisionEntry{
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

var (
	ErrRecipeNotFound = errors.New("outbound type not found")
	ErrNotEnoughStock = errors.New("not enough stock")
)

func (a App) PlaceOrder(id int, qty int) (order.ID, error) {
	recipeID := recipe.ID(id)

	provider := a.PersistentProviderFactory.Begin()

	recipeBook := provider.RecipeBook()

	r := recipeBook.Get(recipeID)

	ingredients := r.Ingredients

	if ingredients == nil {
		return 0, ErrRecipeNotFound
	}

	stock := provider.Stock()

	if err := stock.Sell(ingredients, qty); err != nil {
		switch err {
		case ErrNotEnoughStock:
			return 0, ErrNotEnoughStock
		default:
			panic("unexpected error")
		}
	}

	orders := provider.Orders()

	entryID := orders.Add(order.OrderEntry{
		RecipeID: id,
		Qty:      qty,
	})

	a.PersistentProviderFactory.Commit(provider)

	return entryID, nil
}
