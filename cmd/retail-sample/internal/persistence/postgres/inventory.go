package persistence

import (
	"context"

	pgx "github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

type InventoryPgxStore struct {
	DB pgx.Tx
}

func (ps *InventoryPgxStore) Save(i inventory.DTO) error {
	tag, err := ps.DB.Exec(context.Background(), "update inventory set enabled=$1 and name=$2 where id=$3", i.Enabled, i.Name, i.ID)

	if err != nil {
		return errors.Wrapf(ErrDB, "update inventory item: %v", err)
	}

	if tag.RowsAffected() != 1 {
		return inventory.ErrNotFound
	}

	return nil
}

func (ps *InventoryPgxStore) Add(n string) (int, error) {
	var id int32
	err := ps.DB.QueryRow(context.Background(), "insert into inventory(name, enabled) values($1, true) returning id", n).Scan(&id)

	if err != nil {
		return 0, errors.Wrapf(ErrDB, "add inventory item: %v", err)
	}

	return int(id), nil
}

func (ps *InventoryPgxStore) Find(n string) (int, error) {
	var id int
	err := ps.DB.QueryRow(context.Background(), "select id from inventory where name = $1", n).Scan(&id)

	switch err {
	case nil:
		return id, nil
	case pgx.ErrNoRows:
		return 0, inventory.ErrNotFound
	default:
		return 0, errors.Wrapf(ErrDB, "find inventory item id: %v", err)
	}
}

func (ps *InventoryPgxStore) Get(id int) (inventory.DTO, error) {
	var (
		name    string
		enabled bool
	)

	sql := `select name, enabled 
		from inventory 
		where id = $1`

	err := ps.DB.QueryRow(context.Background(), sql, id).Scan(&name, &enabled)

	switch err {
	case nil:
	case pgx.ErrNoRows:
		return inventory.DTO{}, inventory.ErrNotFound
	default:
		return inventory.DTO{}, errors.Wrapf(ErrDB, "get inventory item by id: %v", err)
	}

	item := inventory.DTO{
		ID:      id,
		Name:    name,
		Enabled: enabled,
	}

	return item, nil
}

func (ps *InventoryPgxStore) List() (items []inventory.DTO, err error) {
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

		items = append(items, inventory.DTO{
			ID:      int(id),
			Name:    name,
			Enabled: enabled,
		})
	}

	return
}
