package middleware

import (
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/order"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/stock"
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
