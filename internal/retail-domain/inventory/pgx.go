package inventory

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PgxPersistentStore struct {
	pool *pgxpool.Pool
}

func NewPersistentStore(pool *pgxpool.Pool) PgxPersistentStore {
	return PgxPersistentStore{pool: pool}
}

func (ps *PgxPersistentStore) add(n Name) ID {
	conn, err := ps.pool.Acquire(context.Background())

	defer conn.Release()

	if err != nil {
		log.Fatal(err)
	}

	t, err := conn.Begin(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	var id int32
	err = t.QueryRow(context.Background(), "insert into inventory(name) values($1) returning id", n).Scan(&id)

	if err != nil {
		_ = t.Rollback(context.Background())
		log.Fatal(err)
	}

	_ = t.Commit(context.Background())

	return ID(id)
}

func (ps *PgxPersistentStore) find(n Name) ID {
	conn, err := ps.pool.Acquire(context.Background())

	defer conn.Release()

	if err != nil {
		log.Fatal(err)
	}

	t, err := conn.Begin(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	rows, err := t.Query(context.Background(), "select id from inventory where name = $1", n)

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}

		return ID(id)
	}

	return ID(0)
}

func (ps *PgxPersistentStore) get(id ID) Item {
	conn, err := ps.pool.Acquire(context.Background())

	defer conn.Release()

	if err != nil {
		log.Fatal(err)
	}

	t, err := conn.Begin(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	rows, err := t.Query(context.Background(), "select name from inventory where id = $1", id)

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		return Item{
			ID:   id,
			Name: Name(name),
		}
	}

	_ = t.Commit(context.Background())

	return Item{}
}

func (ps *PgxPersistentStore) all() (items []Item) {
	conn, err := ps.pool.Acquire(context.Background())

	defer conn.Release()

	if err != nil {
		log.Fatal(err)
	}

	t, err := conn.Begin(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	rows, err := t.Query(context.Background(), "select id, name from inventory")

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var id int32
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		items = append(items, Item{
			ID:   ID(id),
			Name: Name(name),
		})
	}

	_ = t.Commit(context.Background())

	return
}
