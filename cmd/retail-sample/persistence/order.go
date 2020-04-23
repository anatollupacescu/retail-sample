package persistence

import (
	"context"
	"log"
	"time"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/order"
)

type PgxOrderStore struct {
	DB PgxDB
}

func (po *PgxOrderStore) Add(o order.Order) order.ID {
	sql := "insert into outbound_order(recipeid, quantity) values($1, $2) returning id"

	var id int32
	err := po.DB.QueryRow(context.Background(), sql, o.RecipeID, o.Qty).Scan(&id)

	if err != nil {
		log.Print("order add", err)
	}

	return order.ID(id)
}

func (po *PgxOrderStore) List() (orders []order.Order) {
	rows, err := po.DB.Query(context.Background(), "select id, recipeid, quantity, orderdate from outbound_order")

	if err != nil {
		log.Print("order list", err)
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id, recipeID int64
			qty          int16
			time         time.Time
		)

		if err := rows.Scan(&id, &recipeID, &qty, &time); err != nil {
			log.Print("order list scan ", err)
			break
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
