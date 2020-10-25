package persistence

import (
	"context"

	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
	"github.com/anatollupacescu/retail-sample/domain/retail/stock"

	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type StockPgxDB interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type StockPgxStore struct {
	DB StockPgxDB
}

func (ps *StockPgxStore) Provision(id, qty int) (int, error) {
	sql := `insert into stock(inventoryid, quantity) 
					values ($1, $2) 
					ON CONFLICT(inventoryid) DO UPDATE 
					set quantity = stock.quantity + $2 
					where stock.inventoryid = $1
					returning quantity`

	var newQty int
	err := ps.DB.QueryRow(context.Background(), sql, id, qty).Scan(&newQty)

	if err != nil {
		return 0, errors.Wrapf(ErrDB, "provision stock: %v", err)
	}

	return newQty, nil
}

func (ps *StockPgxStore) Quantity(id int) (int, error) {
	sql := "select quantity from stock where inventoryid = $1"

	var qty int
	err := ps.DB.QueryRow(context.Background(), sql, id).Scan(&qty)

	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return 0, stock.ErrItemNotFound
	default:
		return 0, errors.Wrapf(ErrDB, "get stock quantity for item with id %v: %v", id, err)
	}

	return qty, nil
}

func (ps *StockPgxStore) Sell(ii []recipe.Ingredient, qty int) error {
	sql := "update stock set quantity = quantity - $1 where inventoryid = $2"

	for _, i := range ii {
		_, err := ps.DB.Exec(context.Background(), sql, qty*i.Qty, i.ID)

		if err != nil {
			return errors.Wrapf(ErrDB, "update stock: %v", err)
		}
	}

	return nil
}

type PgxProvisionLog struct {
	DB StockPgxDB
}

func (pl *PgxProvisionLog) Add(itemID, qty int) (id int, err error) {
	sql := "insert into provisionlog(inventoryid, quantity) values($1, $2) returning id"

	err = pl.DB.QueryRow(context.Background(), sql, itemID, qty).Scan(&id)

	if err != nil {
		return 0, errors.Wrapf(ErrDB, "update stock quantity for item %v: %v", id, err)
	}

	return
}

func (pl *PgxProvisionLog) Get(id int) (pe stock.ProvisionEntry, err error) {
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

	return stock.ProvisionEntry{
		ID:  itemID,
		Qty: qty,
	}, nil
}

func (pl *PgxProvisionLog) List() (ee []stock.ProvisionEntry, err error) {
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

		ee = append(ee, stock.ProvisionEntry{
			ID:  int(id),
			Qty: int(qty),
		})
	}

	return
}
