package middleware

import (
	"github.com/anatollupacescu/retail-sample/domain/retail-sample/inventory"
	"github.com/anatollupacescu/retail-sample/domain/retail-sample/order"
	"github.com/anatollupacescu/retail-sample/domain/retail-sample/recipe"
	"github.com/anatollupacescu/retail-sample/domain/retail-sample/stock"
)

type (
	PersistenceProviderFactory interface {
		New() PersistenceProvider
		Commit(PersistenceProvider)
		Rollback(PersistenceProvider)
		Ping() error
	}

	PersistenceProvider interface {
		Inventory() inventory.Inventory
		Stock() stock.Stock
		RecipeBook() recipe.Book
		Orders() order.Orders
	}
)
