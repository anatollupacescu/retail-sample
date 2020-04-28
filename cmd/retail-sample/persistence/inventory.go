package persistence

import (
	"context"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type PgxInventoryStore struct {
	DB PgxDB
}

func (ps *PgxInventoryStore) Add(n inventory.Name) (inventory.ID, error) {
	var id int32
	err := ps.DB.QueryRow(context.Background(), "insert into inventory(name) values($1) returning id", n).Scan(&id)

	if err != nil {
		return inventory.ID(0), errors.Wrapf(DBErr, "add inventory item: %v", err)
	}

	return inventory.ID(id), nil
}

func (ps *PgxInventoryStore) Find(n inventory.Name) (inventory.ID, error) {
	var id int
	err := ps.DB.QueryRow(context.Background(), "select id from inventory where name = $1", n).Scan(&id)

	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return inventory.ID(0), inventory.ErrStoreItemNotFound
	default:
		return inventory.ID(0), errors.Wrapf(DBErr, "find inventory item id: %v", err)
	}

	return inventory.ID(id), nil
}

func (ps *PgxInventoryStore) Get(id inventory.ID) (inventory.Item, error) {
	var name string
	err := ps.DB.QueryRow(context.Background(), "select name from inventory where id = $1", id).Scan(&name)

	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return inventory.Item{}, inventory.ErrStoreItemNotFound
	default:
		return inventory.Item{}, errors.Wrapf(DBErr, "get inventory item by id: %v", err)
	}

	return inventory.Item{
		ID:   id,
		Name: inventory.Name(name),
	}, nil
}

func (ps *PgxInventoryStore) List() (items []inventory.Item, err error) {
	rows, err := ps.DB.Query(context.Background(), "select id, name from inventory")

	if err != nil {
		return nil, errors.Wrapf(DBErr, "list inventory: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int32
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, errors.Wrapf(DBErr, "scan inventory: %v", err)
		}
		items = append(items, inventory.Item{
			ID:   inventory.ID(id),
			Name: inventory.Name(name),
		})
	}

	return
}
