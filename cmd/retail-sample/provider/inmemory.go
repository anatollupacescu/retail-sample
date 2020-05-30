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

	InMemoryProvider struct{}
)

var memFactory InMemoryProviderFactory

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

	memFactory = InMemoryProviderFactory{
		inventory:  inventory,
		recipeBook: recipeBook,
		orders:     orders,
		stock:      stock,
	}

	return &memFactory
}

func (pf *InMemoryProviderFactory) New() middleware.PersistenceProvider {
	return &InMemoryProvider{}
}

func (_ *InMemoryProviderFactory) Commit(_ middleware.PersistenceProvider) {
	/*no op*/
}

func (_ *InMemoryProviderFactory) Rollback(_ middleware.PersistenceProvider) {
	/*no op*/
}

func (_ *InMemoryProviderFactory) Ping() error {
	/*no op*/
	return nil
}

func (i *InMemoryProvider) Inventory() inventory.Inventory {
	return memFactory.inventory
}

func (pp *InMemoryProvider) RecipeBook() recipe.Book {
	return memFactory.recipeBook
}

func (pp *InMemoryProvider) Orders() order.Orders {
	return memFactory.orders
}

func (pp *InMemoryProvider) Stock() stock.Stock {
	return memFactory.stock
}
