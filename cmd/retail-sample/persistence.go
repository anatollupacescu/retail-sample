package main

import (
	"context"
	"log"

	invCmd "github.com/anatollupacescu/retail-sample/cmd/retail-sample/inventory"
	orderCmd "github.com/anatollupacescu/retail-sample/cmd/retail-sample/order"
	recipeCmd "github.com/anatollupacescu/retail-sample/cmd/retail-sample/recipe"
	stockCmd "github.com/anatollupacescu/retail-sample/cmd/retail-sample/stock"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/order"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/stock"

	retail "github.com/anatollupacescu/retail-sample/cmd/retail-sample/types"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	kitlog "github.com/go-kit/kit/log"
)

func newPersistentFactory(logger kitlog.Logger, dbConn string) *PgxProviderFactory {
	config, err := pgxpool.ParseConfig(dbConn)

	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)

	if err != nil {
		log.Fatal(err)
	}

	return &PgxProviderFactory{pool: pool}
}

type (
	PgxProviderFactory struct {
		pool *pgxpool.Pool
	}

	PgxTransactionalProvider struct {
		tx pgx.Tx
	}
)

func (pf *PgxProviderFactory) New() retail.PersistenceProvider {
	tx, err := pf.pool.Begin(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	return &PgxTransactionalProvider{
		tx: tx,
	}
}

func (pf *PgxProviderFactory) Commit(pp retail.PersistenceProvider) {
	provider := pp.(*PgxTransactionalProvider)
	if err := provider.tx.Commit(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func (pf *PgxProviderFactory) Rollback(pp retail.PersistenceProvider) {
	provider := pp.(*PgxTransactionalProvider)
	if err := provider.tx.Rollback(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func (pp *PgxTransactionalProvider) Inventory() inventory.Inventory {
	store := &invCmd.PgxStore{DB: pp.tx}
	return inventory.Inventory{Store: store}
}

func (pp *PgxTransactionalProvider) RecipeBook() recipe.Book {
	recipeStore := &recipeCmd.PgxStore{DB: pp.tx}
	inventory := pp.Inventory()
	return recipe.Book{Store: recipeStore, Inventory: inventory}
}

func (pp *PgxTransactionalProvider) Orders() order.Orders {
	orderStore := &orderCmd.PgxStore{DB: pp.tx}
	recipeBook := pp.RecipeBook()
	stock := pp.Stock()
	return order.Orders{
		Store:      orderStore,
		RecipeBook: recipeBook,
		Stock:      stock,
	}
}

func (pp *PgxTransactionalProvider) Stock() stock.Stock {
	store := &stockCmd.PgxStore{DB: pp.tx}
	provisionLog := &stockCmd.PgxProvisionLog{DB: pp.tx}
	inventory := pp.Inventory()

	return stock.Stock{
		Store:        store,
		Inventory:    inventory,
		ProvisionLog: provisionLog,
	}
}
