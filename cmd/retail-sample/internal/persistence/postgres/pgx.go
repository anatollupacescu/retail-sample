package persistence

import (
	"context"
	"log"
	"strings"

	inventory "github.com/anatollupacescu/retail-sample/domain/retail/inventory"
	order "github.com/anatollupacescu/retail-sample/domain/retail/order"
	recipe "github.com/anatollupacescu/retail-sample/domain/retail/recipe"
	"github.com/anatollupacescu/retail-sample/domain/retail/stock"

	v4 "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

type TX struct {
	v4.Tx
}

func (t TX) Commit(ctx context.Context) error {
	return t.Tx.Commit(ctx)
}

func (t TX) Rollback(ctx context.Context) error {
	return t.Tx.Rollback(ctx)
}

func (t TX) Inventory() inventory.Inventory {
	return inventory.New(&InventoryPgxStore{DB: t})
}

func (t TX) Orders() order.Orders {
	s := &OrderPgxStore{DB: t}
	rb := t.Recipe()
	stock := t.Stock()

	return order.New(s, rb, stock)
}

func (t TX) Stock() stock.Stock {
	store := &StockPgxStore{DB: t}
	inventory := t.Inventory()
	log := &PgxProvisionLog{DB: t}

	return stock.New(store, inventory, log)
}

func (t TX) Recipe() recipe.Book {
	store := &RecipePgxStore{DB: t}
	inv := t.Inventory()

	return recipe.New(store, inv)
}

func (db DB) Begin(ctx context.Context) (TX, error) {
	tx, err := db.Pool.Begin(ctx)

	if err != nil {
		return TX{}, err
	}

	wrapped := TX{
		Tx: tx,
	}

	return wrapped, err
}

func NewPersistenceFactory(ctx context.Context, dbConn string) *DB {
	dbConn = strings.TrimSpace(dbConn)
	config, err := pgxpool.ParseConfig(dbConn)

	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.ConnectConfig(ctx, config)

	if err != nil {
		log.Fatal(err)
	}

	db := DB{pool}

	if err = Ping(&db); err != nil {
		log.Fatal(err)
	}

	return &db
}

func Ping(db *DB) error {
	tx, err := db.Begin(context.Background())

	if err != nil {
		return err
	}

	if _, err := tx.Exec(context.Background(), "SELECT true"); err != nil {
		return err
	}

	err = tx.Commit(context.Background())

	return err
}
