package retailsample

import (
	"errors"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/order"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/stock"
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
	NewLogger                 LoggerFactory

	counter *int32
}

type Logger interface {
	Log(keyvals ...interface{}) error
}

type LoggerFactory func() Logger

type PersistenceProviderFactory interface {
	Begin() PersistenceProvider
	Commit(PersistenceProvider)
	Rollback(PersistenceProvider)
}

type Stock interface {
	CurrentStock([]inventory.Item) ([]stock.StockPosition, error)
	Quantity(itemID int) (int, error)
	Sell([]recipe.Ingredient, int) error
	Provision([]stock.StockProvisionEntry) (map[int]int, error)
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
	ProvisionLogEntry struct {
		ID     int
		OldQty int
		NewQty int
	}

	ProvisionLog interface {
		Add(ProvisionEntry) error
		List() ([]ProvisionEntry, error)
	}
)

func (a App) ListInventoryItems() ([]inventory.Item, error) {
	logger := a.NewLogger()

	logger.Log("msg", "[list inventory items] enter method")
	defer logger.Log("msg", "[list inventory items] exit method")

	provider := a.PersistentProviderFactory.Begin()
	defer a.PersistentProviderFactory.Commit(provider)

	inv := provider.Inventory()

	res, err := inv.List()

	if err != nil {
		logger.Log("msg", "[list inventory items]", "error", err)
	}

	return res, err
}

func (a App) AddToInventory(name string) (newID inventory.ID, err error) {
	logger := a.NewLogger()

	logger.Log("msg", "[add inventory item] enter method")
	defer logger.Log("msg", "[add inventory item] exit method")

	provider := a.PersistentProviderFactory.Begin()

	defer func() {
		if err != nil {
			logger.Log("msg", "[provision stock] rollback")
			a.PersistentProviderFactory.Rollback(provider)
			return
		}
		logger.Log("msg", "[provision stock] commit")
		a.PersistentProviderFactory.Commit(provider)
	}()

	inv := provider.Inventory()

	newID, err = inv.Add(inventory.Name(name))

	if err != nil {
		logger.Log("msg", "[add inventory item]", "error", err)
	}

	return
}

func (a App) CurrentStock() (currentStock []stock.StockPosition, err error) {
	logger := a.NewLogger()

	logger.Log("msg", "[current stock] enter method")
	defer logger.Log("msg", "[current stock] exit method")

	provider := a.PersistentProviderFactory.Begin()
	defer a.PersistentProviderFactory.Commit(provider)

	stock := provider.Stock()
	inv := provider.Inventory()

	items, err := inv.List()

	if err != nil {
		logger.Log("msg", "[current stock] fetch inventory items", "error", err)
		return
	}

	currentStock, err = stock.CurrentStock(items)

	if err != nil {
		logger.Log("msg", "[current stock] fetch current stock", "error", err)
	}

	return
}

func (a App) Provision(in []ProvisionEntry) (updatedQtys map[int]int, err error) {
	logger := a.NewLogger()

	logger.Log("msg", "[provision stock] enter method")
	defer logger.Log("msg", "[provision stock] exit method")

	provider := a.PersistentProviderFactory.Begin()

	defer func() {
		if err != nil {
			logger.Log("msg", "[provision stock] rollback")
			a.PersistentProviderFactory.Rollback(provider)
			return
		}
		logger.Log("msg", "[provision stock] commit")
		a.PersistentProviderFactory.Commit(provider)
	}()

	inv := provider.Inventory()

	spes := make([]stock.StockProvisionEntry, 0)

	for _, i := range in {
		id := i.ID

		itemID := inventory.ID(id)

		if _, err = inv.Get(itemID); err != nil {
			logger.Log("msg", "[provision stock] check inventory item exists", "id", id, "error", err)
			return nil, err
		}

		qty := i.Qty

		spes = append(spes, stock.StockProvisionEntry{
			ID:  id,
			Qty: qty,
		})
	}

	st := provider.Stock()

	updatedQtys, err = st.Provision(spes)

	if err != nil {
		logger.Log("msg", "[provision stock] provision item", "error", err)
		return nil, err
	}

	provisionLog := provider.ProvisionLog()

	for id, qty := range updatedQtys {
		entry := ProvisionEntry{
			ID:  id,
			Qty: qty,
		}

		if err := provisionLog.Add(entry); err != nil {
			logger.Log("msg", "[provision stock] log provision entry", "error", err)
			break
		}
	}

	return updatedQtys, nil
}

func (a App) Quantity(id int) (qty int, err error) {
	logger := a.NewLogger()

	logger.Log("msg", "[get stock quantity] enter method")
	defer logger.Log("msg", "[get stock quantity] exit method")

	provider := a.PersistentProviderFactory.Begin()
	defer a.PersistentProviderFactory.Commit(provider)

	qty, err = provider.Stock().Quantity(id)

	if err != nil {
		logger.Log("msg", "[get stock quantity] fetch from store", "error", err)
	}

	return qty, err
}

func (a App) GetProvisionLog() (r []ProvisionEntry, err error) {
	list, err := a.ProvisionLog.List()

	if err != nil {
		return nil, err
	}

	r = append(r, list...)

	return r, err
}

var (
	ErrRecipeNotFound = errors.New("outbound type not found")
	ErrNotEnoughStock = errors.New("not enough stock")
)

func (a App) PlaceOrder(id int, qty int) (orderID order.ID, err error) {
	logger := a.NewLogger()

	logger.Log("msg", "[place order] enter method")
	defer logger.Log("msg", "[place order] exit method")

	recipeID := recipe.ID(id)

	provider := a.PersistentProviderFactory.Begin()

	defer func() {
		if err != nil {
			logger.Log("msg", "[place order] rollback")
			a.PersistentProviderFactory.Rollback(provider)
			return
		}
		logger.Log("msg", "[place order] commit")
		a.PersistentProviderFactory.Commit(provider)
	}()

	recipeBook := provider.RecipeBook()

	r, err := recipeBook.Get(recipeID)

	if err != nil {
		logger.Log("msg", "[place order] get recipe by id", "id", recipeID, "error", err)
		return order.ID(0), err
	}

	ingredients := r.Ingredients

	stock := provider.Stock()

	if err := stock.Sell(ingredients, qty); err != nil {
		return 0, err
	}

	orders := provider.Orders()

	orderID, err = orders.Add(order.OrderEntry{
		RecipeID: id,
		Qty:      qty,
	})

	if err != nil {
		logger.Log("msg", "[place order] save new order", "error", err)
		return order.ID(0), err
	}

	return orderID, nil
}
