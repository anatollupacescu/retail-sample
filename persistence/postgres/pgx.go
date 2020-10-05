package persistence

import (
	"context"
	"log"

	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"

	"github.com/anatollupacescu/retail-sample/domain/retail-sample/inventory"
	"github.com/anatollupacescu/retail-sample/domain/retail-sample/order"
	"github.com/anatollupacescu/retail-sample/domain/retail-sample/recipe"
	"github.com/anatollupacescu/retail-sample/domain/retail-sample/stock"
)

type (
	PgxProviderFactory struct {
		pool *pgxpool.Pool
	}

	PgxTransactionalProvider struct {
		tx pgx.Tx
	}
)

func NewPersistenceFactory(dbConn string) middleware.PersistenceProviderFactory {
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
	store := &InventoryPgxStore{DB: pp.tx}
	return inventory.Inventory{Store: store}
}

func (pp *PgxTransactionalProvider) RecipeBook() recipe.Book {
	recipeStore := &RecipePgxStore{DB: pp.tx}
	inventory := pp.Inventory()

	return recipe.Book{Store: recipeStore, Inventory: inventory}
}

func (pp *PgxTransactionalProvider) Orders() order.Orders {
	orderStore := &OrderPgxStore{DB: pp.tx}
	recipeBook := pp.RecipeBook()
	stock := pp.Stock()

	return order.Orders{
		Store:      orderStore,
		RecipeBook: recipeBook,
		Stock:      stock,
	}
}

func (pp *PgxTransactionalProvider) Stock() stock.Stock {
	store := &StockPgxStore{DB: pp.tx}
	provisionLog := &PgxProvisionLog{DB: pp.tx}
	inventory := pp.Inventory()

	return stock.Stock{
		Store:        store,
		Inventory:    inventory,
		ProvisionLog: provisionLog,
	}
}
