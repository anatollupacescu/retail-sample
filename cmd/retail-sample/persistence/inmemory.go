package persistence

import (
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"
	"github.com/anatollupacescu/retail-sample/persistence/inmemory"

	"github.com/anatollupacescu/retail-sample/domain/retail-sample/inventory"
	"github.com/anatollupacescu/retail-sample/domain/retail-sample/order"
	"github.com/anatollupacescu/retail-sample/domain/retail-sample/recipe"
	"github.com/anatollupacescu/retail-sample/domain/retail-sample/stock"
)

type (
	InMemoryProviderFactory struct {
		inventory  inventory.Inventory
		recipeBook recipe.Book
		orders     order.Orders
		stock      stock.Stock
	}

	InMemoryProvider struct {
		factory *InMemoryProviderFactory
	}
)

func New() *InMemoryProviderFactory {
	invStore := inmemory.NewInventory()
	recipeStore := inmemory.NewRecipe()
	orderStore := inmemory.NewOrder()
	stockStore := inmemory.NewStock()
	provisionLog := inmemory.NewProvisionLog()

	inventory := inventory.Inventory{Store: &invStore}

	recipeBook := recipe.Book{
		Store:     &recipeStore,
		Inventory: inventory,
	}

	stock := stock.Stock{
		Store:        stockStore,
		Inventory:    inventory,
		ProvisionLog: provisionLog,
	}

	orders := order.Orders{
		Store:      orderStore,
		RecipeBook: recipeBook,
		Stock:      stock,
	}

	return &InMemoryProviderFactory{
		inventory:  inventory,
		recipeBook: recipeBook,
		orders:     orders,
		stock:      stock,
	}
}

func (f *InMemoryProviderFactory) New() middleware.PersistenceProvider {
	return &InMemoryProvider{factory: f}
}

func (*InMemoryProviderFactory) Commit(middleware.PersistenceProvider) {
	/*no op*/
}

func (*InMemoryProviderFactory) Rollback(middleware.PersistenceProvider) {
	/*no op*/
}

func (*InMemoryProviderFactory) Ping() error {
	/*no op*/
	return nil
}

func (p *InMemoryProvider) Inventory() inventory.Inventory {
	return p.factory.inventory
}

func (p *InMemoryProvider) RecipeBook() recipe.Book {
	return p.factory.recipeBook
}

func (p *InMemoryProvider) Orders() order.Orders {
	return p.factory.orders
}

func (p *InMemoryProvider) Stock() stock.Stock {
	return p.factory.stock
}
