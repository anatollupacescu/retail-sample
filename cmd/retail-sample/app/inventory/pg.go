package inventory

import (
	"context"

	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
)

var ErrDB = errors.New("postgres")

type PgxDB interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
}

type PgxStore struct {
	DB PgxDB
}

func (ps *PgxStore) Update(i inventory.Item) error {
	tag, err := ps.DB.Exec(context.Background(), "update inventory set enabled=$1 and name=$2 where id=$3", i.Enabled, i.Name, i.ID)

	if err != nil {
		return errors.Wrapf(ErrDB, "update inventory item: %v", err)
	}

	if tag.RowsAffected() != 1 {
		return inventory.ErrItemNotFound
	}

	return nil
}

func (ps *PgxStore) Add(n string) (int, error) {
	var id int32
	err := ps.DB.QueryRow(context.Background(), "insert into inventory(name, enabled) values($1, true) returning id", n).Scan(&id)

	if err != nil {
		return 0, errors.Wrapf(ErrDB, "add inventory item: %v", err)
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
		return 0, errors.Wrapf(ErrDB, "find inventory item id: %v", err)
	}

	return id, nil
}

func (ps *PgxStore) Get(id int) (inventory.Item, error) {
	var (
		name    string
		enabled bool
	)

	sql := `select 
						name, enabled 
					from 
						inventory 
					where id = $1`

	err := ps.DB.QueryRow(context.Background(), sql, id).Scan(&name, &enabled)

	var zeroItem inventory.Item

	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return zeroItem, inventory.ErrItemNotFound
	default:
		return zeroItem, errors.Wrapf(ErrDB, "get inventory item by id: %v", err)
	}

	return inventory.Item{
		ID:      id,
		Name:    name,
		Enabled: enabled,
	}, nil
}

func (ps *PgxStore) List() (items []inventory.Item, err error) {
	rows, err := ps.DB.Query(context.Background(), "select id, name, enabled from inventory")

	if err != nil {
		return nil, errors.Wrapf(ErrDB, "list inventory: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id      int32
			name    string
			enabled bool
		)

		if err := rows.Scan(&id, &name, &enabled); err != nil {
			return nil, errors.Wrapf(ErrDB, "scan inventory: %v", err)
		}

		items = append(items, inventory.Item{
			ID:      int(id),
			Name:    name,
			Enabled: enabled,
		})
	}

	return
}
