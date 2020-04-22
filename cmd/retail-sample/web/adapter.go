package web

import (
	"context"
	"log"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/persistence"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/order"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
	retail "github.com/anatollupacescu/retail-sample/internal/retail-sample"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type WebApp struct {
	retail.App
}

func (w *WebApp) IsHealthy() bool {
	return true
}

func NewApp() WebApp {
	config, err := pgxpool.ParseConfig("postgres://docker:docker@localhost:5432/retail?pool_max_conns=10")

	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)

	if err != nil {
		log.Fatal(err)
	}

	inventoryStore := persistence.PgxInventoryStore{DB: pool}
	inventory := inventory.Inventory{Store: &inventoryStore}

	orderStore := persistence.PgxOrderStore{DB: pool}
	orders := order.Orders{Store: &orderStore}

	provisionLog := make(retail.InMemoryProvisionLog)

	recipeStore := recipe.NewInMemoryStore()
	recipeBook := recipe.Book{Store: &recipeStore, Inventory: &inventory}

	stock := retail.NewInMemoryStock()

	app := retail.App{
		PersistentProviderFactory: newFactory(pool),
		Inventory:                 inventory,
		Orders:                    orders,
		ProvisionLog:              provisionLog,
		RecipeBook:                recipeBook,
		Stock:                     stock,
	}

	return WebApp{
		App: app,
	}
}

type (
	PgxProviderFactory struct {
		pool *pgxpool.Pool
	}

	PgxTransactionalProvider struct {
		tx pgx.Tx
	}
)

func (pf *PgxProviderFactory) Begin() retail.PersistenceProvider {
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

func newFactory(pool *pgxpool.Pool) *PgxProviderFactory {
	return &PgxProviderFactory{pool: pool}
}

func (pp *PgxTransactionalProvider) Inventory() retail.Inventory {
	store := &persistence.PgxInventoryStore{DB: pp.tx}
	return inventory.Inventory{Store: store}
}

func (pp *PgxTransactionalProvider) RecipeBook() retail.RecipeBook {
	recipeStore := &persistence.PgxRecipeStore{DB: pp.tx}
	inventory := pp.Inventory()
	return recipe.Book{Store: recipeStore, Inventory: inventory}
}

func (pp *PgxTransactionalProvider) Orders() retail.Orders {
	orderStore := persistence.PgxOrderStore{DB: pp.tx}
	return order.Orders{Store: &orderStore}
}

//no store since are not domain
func (pp *PgxTransactionalProvider) Stock() retail.Stock {
	return &persistence.PgxStock{DB: pp.tx}
}

func (pp *PgxTransactionalProvider) ProvisionLog() retail.ProvisionLog {
	return &persistence.PgxProvisionLog{DB: pp.tx}
}
