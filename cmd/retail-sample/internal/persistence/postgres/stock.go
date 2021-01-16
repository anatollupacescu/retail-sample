package persistence

import (
	"context"

	pgx "github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

type StockPgxStore struct {
	DB pgx.Tx
}

func (ps *StockPgxStore) Get(inventoryID int) (dto stock.PositionDTO, err error) {
	sql := `select quantity from stock where inventoryid = $1`

	err = ps.DB.QueryRow(context.Background(), sql, inventoryID).Scan(&dto.Qty)

	switch err {
	case nil:
	case pgx.ErrNoRows:
		return stock.PositionDTO{}, stock.ErrPositionNotFound
	default:
		return stock.PositionDTO{}, errors.Wrapf(ErrDB, "get stock position for item with id %v: %v", inventoryID, err)
	}

	dto.InventoryID = inventoryID

	return dto, nil
}

func (ps *StockPgxStore) Save(dto stock.PositionDTO) error {
	sql := `insert into stock (inventoryid, quantity) values ($1, $2)
		ON CONFLICT(inventoryid) DO 
		UPDATE SET quantity = $2 where stock.inventoryid = $1`

	_, err := ps.DB.Exec(context.Background(), sql, dto.InventoryID, dto.Qty)

	if err != nil {
		return errors.Wrapf(ErrDB, "save stock position: %v", err)
	}

	return nil
}

func (ps *StockPgxStore) List() ([]stock.PositionDTO, error) {
	sql := "select inventoryid, quantity from stock"

	rows, err := ps.DB.Query(context.Background(), sql)

	if err != nil {
		return nil, errors.Wrapf(ErrDB, "save stock position: %v", err)
	}

	defer rows.Close()

	var entries = make([]stock.PositionDTO, 0)

	for rows.Next() {
		var (
			id  int
			qty int
		)

		if err := rows.Scan(&id, &qty); err != nil {
			return nil, errors.Wrapf(ErrDB, "provisionlog list scan: %v", err)
		}

		entries = append(entries, stock.PositionDTO{
			InventoryID: id,
			Qty:         qty,
		})
	}

	return entries, nil
}
