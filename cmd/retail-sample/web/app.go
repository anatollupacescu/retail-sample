package web

import (
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/order"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/recipe"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/warehouse"
)

type App struct {
	inventory inventory.Inventory
	recipe    recipe.Book
	orders    order.Orders
	stock     warehouse.Stock
}

func newInMemoryApp() App {
	inventryStore := inventory.NewInMemoryStore()
	inventory := inventory.Inventory{Store: &inventryStore}

	recipeStore := recipe.NewInMemoryStore()
	recipeBook := recipe.Book{Store: &recipeStore, Inventory: &inventory}

	provisionLog := make(warehouse.InMemoryProvisionLog)
	orderStore := order.NewInMemoryStore()
	orders := order.Orders{Store: orderStore}
	stock := warehouse.NewStock(inventory, recipeBook, provisionLog, orders)

	return App{
		inventory: inventory,
		recipe:    recipeBook,
		orders:    orders,
		stock:     stock,
	}
}
