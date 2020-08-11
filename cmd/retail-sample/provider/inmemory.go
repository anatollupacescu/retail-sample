package provider

import (
	invCmd "github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/inventory"
	orderCmd "github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/order"
	recipeCmd "github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/recipe"
	stockCmd "github.com/anatollupacescu/retail-sample/cmd/retail-sample/app/stock"
	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/order"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/stock"
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

func newInMemoryPersistentFactory() *InMemoryProviderFactory {
	invStore := invCmd.NewInMemoryStore()
	recipeStore := recipeCmd.NewInMemoryStore()
	orderStore := orderCmd.NewInMemoryStore()
	stockStore := stockCmd.NewInMemoryStock()
	provisionLog := stockCmd.NewInMemoryProvisionLog()

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
