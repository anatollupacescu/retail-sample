package persistence

import (
	"context"
	"log"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
	retail "github.com/anatollupacescu/retail-sample/internal/retail-sample"
)

type PgxStock struct {
	DB PgxDB
}

func (ps *PgxStock) Provision(id, qty int) int {
	sql := `insert into stock(inventoryid, quantity) 
					values ($1, $2) 
					ON CONFLICT(inventoryid) DO UPDATE 
					set quantity = stock.quantity + $2 
					where stock.inventoryid = $1
					returning quantity`

	var newQty int
	err := ps.DB.QueryRow(context.Background(), sql, id, qty).Scan(&newQty)

	if err != nil {
		log.Print("stock provision", err)
	}

	return newQty
}

func (ps *PgxStock) Quantity(id int) int {
	sql := "select quantity from stock where inventoryid = $1"

	var qty int
	row := ps.DB.QueryRow(context.Background(), sql, id)

	_ = row.Scan(&qty)

	return qty
}

func (ps *PgxStock) Sell(ii []recipe.Ingredient, qty int) error {
	sql := "update stock set quantity = quantity - $1 where inventoryid = $2"

	for _, i := range ii {
		_, err := ps.DB.Exec(context.Background(), sql, qty*i.Qty, i.ID)

		if err != nil {
			return retail.ErrNotEnoughStock
		}
	}

	return nil
}
