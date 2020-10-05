package persistence

import (
	"context"
	"time"

	pgx "github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/anatollupacescu/retail-sample/domain/retail-sample/order"
)

type OrderPgxDB interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type OrderPgxStore struct {
	DB OrderPgxDB
}

func (po *OrderPgxStore) Add(o order.Order) (order.ID, error) {
	sql := "insert into outbound_order(recipeid, quantity) values($1, $2) returning id"

	var id int32
	err := po.DB.QueryRow(context.Background(), sql, o.RecipeID, o.Qty).Scan(&id)

	if err != nil {
		return order.ID(0), errors.Wrapf(ErrDB, "add order: %v", err)
	}

	return order.ID(id), nil
}

func (po *OrderPgxStore) List() ([]order.Order, error) {
	rows, err := po.DB.Query(context.Background(), "select id, recipeid, quantity, orderdate from outbound_order")

	if err != nil {
		return nil, errors.Wrapf(ErrDB, "list orders: %v", err)
	}

	defer rows.Close()

	var orders = make([]order.Order, 0, len(rows.RawValues()))

	for rows.Next() {
		var (
			id, recipeID int64
			qty          int16
			time         time.Time
		)

		if err := rows.Scan(&id, &recipeID, &qty, &time); err != nil {
			return nil, errors.Wrapf(ErrDB, "scan orders: %v", err)
		}

		orders = append(orders, order.Order{
			ID:   order.ID(id),
			Date: time,
			Entry: order.Entry{
				RecipeID: int(recipeID),
				Qty:      int(qty),
			},
		})
	}

	return orders, nil
}

func (po *OrderPgxStore) Get(id order.ID) (result order.Order, err error) {
	sql := `
		select 
			recipeid, quantity 
		from 
			outbound_order 
		where 
			id = $1`

	var (
		recipeID int
		qty      int
	)

	err = po.DB.QueryRow(context.Background(), sql).Scan(&recipeID, &qty)

	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return result, order.ErrOrderNotFound
	default:
		return result, errors.Wrapf(ErrDB, "get inventory item by id: %v", err)
	}

	return order.Order{
		Entry: order.Entry{
			RecipeID: recipeID,
			Qty:      qty,
		},
	}, nil
}
