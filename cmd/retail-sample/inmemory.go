package main

import (
	invCmd "github.com/anatollupacescu/retail-sample/cmd/retail-sample/inventory"
	orderCmd "github.com/anatollupacescu/retail-sample/cmd/retail-sample/order"
	recipeCmd "github.com/anatollupacescu/retail-sample/cmd/retail-sample/recipe"
	stockCmd "github.com/anatollupacescu/retail-sample/cmd/retail-sample/stock"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/order"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/stock"

	retail "github.com/anatollupacescu/retail-sample/cmd/retail-sample/types"
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

//nolint:golint,unused,deadcode
func newInMemoryPersistentFactory() *InMemoryProviderFactory {
	invStore := invCmd.NewInMemoryStore()
	recipeStore := recipeCmd.NewInMemoryStore()
	orderStore := orderCmd.NewInMemoryStore()
	stockStore := stockCmd.NewInMemoryStock()
	provisionLog := make(stockCmd.InMemoryProvisionLog)

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

func (pf *InMemoryProviderFactory) New() retail.PersistenceProvider {
	return &InMemoryProvider{}
}

func (_ *InMemoryProviderFactory) Commit(_ retail.PersistenceProvider) {
	/*no op*/
}

func (_ *InMemoryProviderFactory) Rollback(_ retail.PersistenceProvider) {
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
