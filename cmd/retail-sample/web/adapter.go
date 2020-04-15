package web

import (
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/order"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
	retail "github.com/anatollupacescu/retail-sample/internal/retail-sample"
)

type WebApp struct {
	retail.App
}

func NewInMemoryApp() WebApp {
	inventryStore := inventory.NewInMemoryStore()
	inventory := inventory.Inventory{Store: &inventryStore}

	orderStore := order.NewInMemoryStore()
	orders := order.Orders{Store: orderStore}

	provisionLog := make(retail.InMemoryProvisionLog)

	recipeStore := recipe.NewInMemoryStore()
	recipeBook := recipe.Book{Store: &recipeStore, Inventory: &inventory}

	stock := retail.NewInMemoryStock()

	app := retail.App{
		Inventory:    inventory,
		Orders:       orders,
		ProvisionLog: provisionLog,
		RecipeBook:   recipeBook,
		Stock:        stock,
	}

	return WebApp{
		App: app,
	}
}
