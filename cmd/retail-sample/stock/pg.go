package stock

import (
	"context"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/stock"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

var DBErr = errors.New("postgres")

type PgxDB interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type PgxStore struct {
	DB PgxDB
}

func (ps *PgxStore) Provision(id, qty int) (int, error) {
	sql := `insert into stock(inventoryid, quantity) 
					values ($1, $2) 
					ON CONFLICT(inventoryid) DO UPDATE 
					set quantity = stock.quantity + $2 
					where stock.inventoryid = $1
					returning quantity`

	var newQty int
	err := ps.DB.QueryRow(context.Background(), sql, id, qty).Scan(&newQty)

	if err != nil {
		return 0, errors.Wrapf(DBErr, "provision stock: %v", err)
	}

	return newQty, nil
}

func (ps *PgxStore) Quantity(id int) (int, error) {
	sql := "select quantity from stock where inventoryid = $1"

	var qty int
	err := ps.DB.QueryRow(context.Background(), sql, id).Scan(&qty)

	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return 0, stock.ErrItemNotFound
	default:
		return 0, errors.Wrapf(DBErr, "get stock quantity for item with id %v: %v", id, err)
	}

	return qty, nil
}

func (ps *PgxStore) Sell(ii []recipe.Ingredient, qty int) error {
	sql := "update stock set quantity = quantity - $1 where inventoryid = $2"

	for _, i := range ii {
		_, err := ps.DB.Exec(context.Background(), sql, qty*i.Qty, i.ID)

		if err != nil {
			return errors.Wrapf(DBErr, "update stock: %v", err)
		}
	}

	return nil
}

type PgxProvisionLog struct {
	DB PgxDB
}

func (pl *PgxProvisionLog) Add(re stock.ProvisionEntry) error {
	sql := "insert into provisionlog(inventoryid, quantity) values($1, $2)"
	if _, err := pl.DB.Exec(context.Background(), sql, re.ID, re.Qty); err != nil {
		return errors.Wrapf(DBErr, "provisionlog add: %v", err)
	}

	return nil
}

func (pl *PgxProvisionLog) List() (ee []stock.ProvisionEntry, err error) {
	rows, err := pl.DB.Query(context.Background(), "select inventoryid, quantity from provisionlog")

	if err != nil {
		return nil, errors.Wrapf(DBErr, "provisionlog list: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int64
		var qty int16

		if err := rows.Scan(&id, &qty); err != nil {
			return nil, errors.Wrapf(DBErr, "provisionlog list scan: %v", err)
		}

		ee = append(ee, stock.ProvisionEntry{
			ID:  int(id),
			Qty: int(qty),
		})
	}

	return
}
