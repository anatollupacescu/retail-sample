package provider

import (
	"context"
	"log"

	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

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
	PgxProviderFactory struct {
		pool *pgxpool.Pool
	}

	PgxTransactionalProvider struct {
		tx pgx.Tx
	}
)

func newPersistenceFactory(dbConn string) middleware.PersistenceProviderFactory {
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

func (pf *PgxProviderFactory) New() middleware.PersistenceProvider {
	tx, err := pf.pool.Begin(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	return &PgxTransactionalProvider{
		tx: tx,
	}
}

func (pf *PgxProviderFactory) Ping() error {
	tx, err := pf.pool.Begin(context.Background())

	if err != nil {
		return err
	}

	if _, err := tx.Exec(context.Background(), "SELECT true"); err != nil {
		return err
	}

	return nil
}

func (pf *PgxProviderFactory) Commit(pp middleware.PersistenceProvider) {
	provider := pp.(*PgxTransactionalProvider)
	err := provider.tx.Commit(context.Background())

	if err != nil {
		log.Printf("commit: %v", err)
	}
}

func (pf *PgxProviderFactory) Rollback(pp middleware.PersistenceProvider) {
	provider := pp.(*PgxTransactionalProvider)
	err := provider.tx.Rollback(context.Background())

	if err != nil {
		log.Printf("rollback: %v", err)
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
