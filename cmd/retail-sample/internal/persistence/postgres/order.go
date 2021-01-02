package persistence

import (
	"context"
	"time"

	pgx "github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/anatollupacescu/retail-sample/domain/retail/order"
)

type OrderPgxStore struct {
	DB pgx.Tx
}

func (po *OrderPgxStore) Add(o order.DTO) (int, error) {
	sql := "insert into outbound_order(recipeid, quantity) values($1, $2) returning id"

	var id int
	err := po.DB.QueryRow(context.Background(), sql, o.RecipeID, o.Qty).Scan(&id)

	if err != nil {
		return 0, errors.Wrapf(ErrDB, "add order: %v", err)
	}

	return id, nil
}

func (po *OrderPgxStore) List() ([]order.DTO, error) {
	rows, err := po.DB.Query(context.Background(), "select id, recipeid, quantity, orderdate from outbound_order")

	if err != nil {
		return nil, errors.Wrapf(ErrDB, "list orders: %v", err)
	}

	defer rows.Close()

	var orders = make([]order.DTO, 0, len(rows.RawValues()))

	for rows.Next() {
		var (
			id, recipeID int
			qty          int
			time         time.Time
		)

		if err := rows.Scan(&id, &recipeID, &qty, &time); err != nil {
			return nil, errors.Wrapf(ErrDB, "scan orders: %v", err)
		}

		orders = append(orders, order.DTO{
			ID:       id,
			Date:     time,
			RecipeID: recipeID,
			Qty:      qty,
		})
	}

	return orders, nil
}

func (po *OrderPgxStore) Get(id int) (order.DTO, error) {
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

	err := po.DB.QueryRow(context.Background(), sql, id).Scan(&recipeID, &qty)

	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return order.DTO{}, order.ErrOrderNotFound
	default:
		return order.DTO{}, errors.Wrapf(ErrDB, "get inventory item by id: %v", err)
	}

	result := order.DTO{
		RecipeID: recipeID,
		Qty:      qty,
	}

	return result, nil
}
