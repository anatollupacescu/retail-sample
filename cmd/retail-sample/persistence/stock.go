package persistence

import (
	"context"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type PgxStockStore struct {
	DB PgxDB
}

func (ps *PgxStockStore) Provision(id, qty int) (int, error) {
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

func (ps *PgxStockStore) Quantity(id int) (int, error) {
	sql := "select quantity from stock where inventoryid = $1"

	var qty int
	err := ps.DB.QueryRow(context.Background(), sql, id).Scan(&qty)

	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return 0, nil
	default:
		return 0, errors.Wrapf(DBErr, "get stock quantity for item with id %v: %v", id, err)
	}

	return qty, nil
}

func (ps *PgxStockStore) Sell(ii []recipe.Ingredient, qty int) error {
	sql := "update stock set quantity = quantity - $1 where inventoryid = $2"

	for _, i := range ii {
		_, err := ps.DB.Exec(context.Background(), sql, qty*i.Qty, i.ID)

		if err != nil {
			return errors.Wrapf(DBErr, "update stock: %v", err)
		}
	}

	return nil
}
