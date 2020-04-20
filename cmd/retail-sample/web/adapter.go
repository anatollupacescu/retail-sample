package web

import (
	"context"
	"log"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/order"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
	retail "github.com/anatollupacescu/retail-sample/internal/retail-sample"
	"github.com/jackc/pgx/v4/pgxpool"
)

type WebApp struct {
	retail.App
}

func NewInMemoryApp() WebApp {
	config, err := pgxpool.ParseConfig("postgres://docker:docker@localhost:5432/retail?pool_max_conns=10")

	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)

	if err != nil {
		log.Fatal(err)
	}

	inventoryStore := inventory.NewPersistentStore(pool)
	inventory := inventory.Inventory{Store: &inventoryStore}

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
