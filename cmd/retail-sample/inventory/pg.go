package inventory

import (
	"context"

	"github.com/pkg/errors"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
)

var DBErr = errors.New("postgres")

type PgxDB interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type PgxStore struct {
	DB PgxDB
}

func (ps *PgxStore) Add(n string) (int, error) {
	var id int32
	err := ps.DB.QueryRow(context.Background(), "insert into inventory(name) values($1) returning id", n).Scan(&id)

	if err != nil {
		return 0, errors.Wrapf(DBErr, "add inventory item: %v", err)
	}

	return int(id), nil
}

func (ps *PgxStore) Find(n string) (int, error) {
	var id int
	err := ps.DB.QueryRow(context.Background(), "select id from inventory where name = $1", n).Scan(&id)

	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return 0, inventory.ErrItemNotFound
	default:
		return 0, errors.Wrapf(DBErr, "find inventory item id: %v", err)
	}

	return id, nil
}

func (ps *PgxStore) Get(id int) (inventory.Item, error) {
	var name string
	err := ps.DB.QueryRow(context.Background(), "select name from inventory where id = $1", id).Scan(&name)

	var zeroItem inventory.Item

	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return zeroItem, inventory.ErrItemNotFound
	default:
		return zeroItem, errors.Wrapf(DBErr, "get inventory item by id: %v", err)
	}

	return inventory.Item{
		ID:   id,
		Name: string(name),
	}, nil
}

func (ps *PgxStore) List() (items []inventory.Item, err error) {
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
			ID:   int(id),
			Name: string(name),
		})
	}

	return
}
