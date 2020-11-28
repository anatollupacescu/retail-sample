package persistence

import (
	"context"

	pgx "github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

type PgxProvisionLog struct {
	DB pgx.Tx
}

type ProvisionEntry struct {
	ID, Qty int
}

func (pl *PgxProvisionLog) Add(itemID, qty int) (id int, err error) {
	sql := "insert into provisionlog(inventoryid, quantity) values($1, $2) returning id"

	err = pl.DB.QueryRow(context.Background(), sql, itemID, qty).Scan(&id)

	if err != nil {
		return 0, errors.Wrapf(ErrDB, "update stock quantity for item %v: %v", id, err)
	}

	return
}

func (pl *PgxProvisionLog) Get(id int) (pe ProvisionEntry, err error) {
	sql := "select inventoryid, quantity from provisionlog where id = $1"

	var itemID, qty int
	err = pl.DB.QueryRow(context.Background(), sql, id).Scan(&itemID, &qty)

	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return pe, stock.ErrItemNotFound
	default:
		return pe, errors.Wrapf(ErrDB, "get provision entry %v: %v", id, err)
	}

	return ProvisionEntry{
		ID:  itemID,
		Qty: qty,
	}, nil
}

func (pl *PgxProvisionLog) List() (ee []ProvisionEntry, err error) {
	rows, err := pl.DB.Query(context.Background(), "select inventoryid, quantity from provisionlog")

	if err != nil {
		return nil, errors.Wrapf(ErrDB, "provisionlog list: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id  int64
			qty int16
		)

		if err := rows.Scan(&id, &qty); err != nil {
			return nil, errors.Wrapf(ErrDB, "provisionlog list scan: %v", err)
		}

		ee = append(ee, ProvisionEntry{
			ID:  int(id),
			Qty: int(qty),
		})
	}

	return
}
