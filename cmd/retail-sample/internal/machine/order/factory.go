package order

import (
	pg "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

type extractor struct {
	stock *pg.StockPgxStore
}

func (o extractor) Extract(inventoryID int, qty int) error {
	dto, err := o.stock.Get(inventoryID)
	if err != nil {
		return err
	}

	pos := stock.Position{
		Qty:         dto.Qty,
		InventoryID: inventoryID,
		DB:          o.stock,
	}

	return pos.Extract(qty)
}
