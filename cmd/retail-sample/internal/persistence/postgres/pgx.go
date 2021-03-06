package persistence

import (
	"context"
	"log"
	"strings"

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

	if _, err = tx.Exec(context.Background(), "SELECT true"); err != nil {
		return err
	}

	err = tx.Commit(context.Background())

	return err
}
