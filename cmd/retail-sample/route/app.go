package route

import (
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/recipe"
	"github.com/anatollupacescu/retail-sample/internal/retail-sample/warehouse"
)

type App struct {
	inventory inventory.Inventory
	recipe    recipe.Book
	stock     warehouse.Stock
}

func newInMemoryApp() App {
	inventryStore := inventory.NewInMemoryStore()
	inventory := inventory.Inventory{Store: &inventryStore}

	recipeStore := recipe.NewInMemoryStore()
	recipeBook := recipe.Book{Store: &recipeStore, Inventory: &inventory}

	return App{
		inventory: inventory,
		recipe:    recipeBook,
		stock: warehouse.Stock{
			Inventory:   inventory,
			RecipeBook:  recipeBook,
			InboundLog:  make(warehouse.InMemoryInboundLog),
			OutboundLog: make(warehouse.InMemoryOutboundLog),
			Data:        make(map[int]int),
		},
	}
}
