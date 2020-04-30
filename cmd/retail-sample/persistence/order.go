package persistence

import (
	"context"
	"time"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/order"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type PgxOrderStore struct {
	DB PgxDB
}

func (po *PgxOrderStore) Add(o order.Order) (order.ID, error) {
	sql := "insert into outbound_order(recipeid, quantity) values($1, $2) returning id"

	var id int32
	err := po.DB.QueryRow(context.Background(), sql, o.RecipeID, o.Qty).Scan(&id)

	if err != nil {
		return order.ID(0), errors.Wrapf(DBErr, "add order: %v", err)
	}

	return order.ID(id), nil
}

func (po *PgxOrderStore) List() (orders []order.Order, err error) {
	rows, err := po.DB.Query(context.Background(), "select id, recipeid, quantity, orderdate from outbound_order")

	if err != nil {
		return nil, errors.Wrapf(DBErr, "list orders: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id, recipeID int64
			qty          int16
			time         time.Time
		)

		if err := rows.Scan(&id, &recipeID, &qty, &time); err != nil {
			return nil, errors.Wrapf(DBErr, "scan orders: %v", err)
		}

		orders = append(orders, order.Order{
			ID:   order.ID(id),
			Date: time,
			OrderEntry: order.OrderEntry{
				RecipeID: int(recipeID),
				Qty:      int(qty),
			},
		})
	}

	return
}

func (po *PgxOrderStore) Get(id order.ID) (result order.Order, err error) {
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
		return result, errors.Wrapf(DBErr, "get inventory item by id: %v", err)
	}

	return order.Order{
		OrderEntry: order.OrderEntry{
			RecipeID: recipeID,
			Qty:      qty,
		},
	}, nil
}
