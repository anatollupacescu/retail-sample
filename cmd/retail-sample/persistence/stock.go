package persistence

import (
	"context"
	"log"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
)

type PgxStock struct {
	DB PgxDB
}

func (ps *PgxStock) Provision(id, qty int) int {
	sql := "update stock set qty = qty + $1 where itemid = $2 returning qty"

	var newQty int
	err := ps.DB.QueryRow(context.Background(), sql, qty, id).Scan(&newQty)

	if err != nil {
		log.Printf("got err %v", err)
	}

	return newQty
}

func (ps *PgxStock) Quantity(id int) int {
	sql := "select qty from stock where id = $1"

	var qty int
	err := ps.DB.QueryRow(context.Background(), sql, id).Scan(&qty)

	if err != nil {
		log.Printf("got err %v", err)
	}

	return qty
}

func (ps *PgxStock) Sell(ii []recipe.Ingredient, qty int) error {
	sql := "update stock set qty = qty - $1 where itemid = $2 returning qty"

	for _, i := range ii {
		_, err := ps.DB.Query(context.Background(), sql, qty*i.Qty, i.ID)

		if err != nil {
			return err
		}
	}

	return nil
}
