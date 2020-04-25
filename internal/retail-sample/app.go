package retailsample

import (
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/order"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
	"github.com/pkg/errors"
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
	List() ([]inventory.Item, error)
	Get(inventory.ID) (inventory.Item, error)
	Find(inventory.Name) (inventory.ID, error)
}

type RecipeBook interface {
	Add(recipe.Name, []recipe.Ingredient) (recipe.ID, error)
	Get(recipe.ID) (recipe.Recipe, error)
	List() ([]recipe.Recipe, error)
}

type Orders interface {
	Add(order.OrderEntry) (order.ID, error)
	List() ([]order.Order, error)
}

type ( //app specific
	ProvisionEntry struct {
		ID  int
		Qty int
	}

	ProvisionLog interface {
		Add(ProvisionEntry)
		List() []ProvisionEntry
	}

	Stock interface {
		Quantity(int) (int, error)
		Provision(int, int) (int, error)
		Sell([]recipe.Ingredient, int) error
	}
)

func (a App) CurrentStock() (ps []StockPosition, err error) {
	provider := a.PersistentProviderFactory.Begin()
	defer a.PersistentProviderFactory.Commit(provider)

	stock := &stock{store: provider.Stock()}
	inv := provider.Inventory()

	items, err := inv.List()

	if err != nil {
		return ps, err
	}

	return stock.CurrentStock(items), nil
}

var ErrInventoryItemNotFound = errors.New("inventory item not found")

func (a App) Provision(in []ProvisionEntry) (newQty map[int]int, err error) {
	provider := a.PersistentProviderFactory.Begin()

	defer func() {
		if err != nil {
			a.PersistentProviderFactory.Rollback(provider)
			return
		}
		a.PersistentProviderFactory.Commit(provider)
	}()

	inv := provider.Inventory()

	//>garbage
	id := 0
	qty := 0
	//<garbage

	itemID := inventory.ID(id)

	if _, err = inv.Get(itemID); err != nil {
		return nil, err
	}

	stock := provider.Stock()

	newQty = make(map[int]int, 0)

	for _, v := range in {
		_, err = stock.Provision(v.ID, v.Qty)
	}

	provisionLog := provider.ProvisionLog()

	provisionLog.Add(ProvisionEntry{
		ID:  id,
		Qty: qty,
	})

	return newQty, nil
}

func (a App) Quantity(id int) (int, error) {
	return a.Stock.Quantity(id)
}

func (a App) GetProvisionLog() (r []ProvisionEntry) {
	r = append(r, a.ProvisionLog.List()...)

	return
}

var (
	BusinessErr = errors.New("business")

	ErrRecipeNotFound = errors.New("outbound type not found")
	ErrNotEnoughStock = errors.New("not enough stock")
)

func (a App) PlaceOrder(id int, qty int) (orderID order.ID, err error) {
	recipeID := recipe.ID(id)

	provider := a.PersistentProviderFactory.Begin()

	defer func() {
		if err != nil {
			a.PersistentProviderFactory.Rollback(provider)
			return
		}
		a.PersistentProviderFactory.Commit(provider)
	}()

	recipeBook := provider.RecipeBook()

	r, _ := recipeBook.Get(recipeID)

	ingredients := r.Ingredients

	//TODO check the error instead of nil
	if ingredients == nil {
		return 0, errors.Wrap(BusinessErr, ErrRecipeNotFound.Error())
	}

	stock := &stock{store: provider.Stock()}

	if err := stock.Sell(ingredients, qty); err != nil {
		return 0, err
	}

	orders := provider.Orders()

	return orders.Add(order.OrderEntry{
		RecipeID: id,
		Qty:      qty,
	})
}
