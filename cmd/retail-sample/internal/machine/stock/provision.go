package stock

import (
	pg "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

type Position struct {
	ID   int
	Name string
	Qty  int
}

func (o *UseCase) Provision(inventoryItemID, qty int) (Position, error) {
	stockPos, err := o.stockDB.Get(inventoryItemID)

	switch err {
	case nil: //continue,
	case pg.ErrStockItemNotFound:
		stockPos = stock.PositionDTO{InventoryID: inventoryItemID}
	default:
		return Position{}, err
	}

	pos := stock.Position{
		Qty:         stockPos.Qty,
		InventoryID: stockPos.InventoryID,
		DB:          o.stockDB,
	}

	err = pos.Provision(qty)
	if err != nil {
		o.logger.Error().Err(err).Msg("call domain layer")
		return Position{}, err
	}

	_, err = o.logDB.Add(inventoryItemID, qty)
	if err != nil {
		o.logger.Error().Err(err).Msg("add log entry")
		return Position{}, err
	}

	item, err := o.inventoryDB.Get(inventoryItemID)
	if err != nil {
		o.logger.Error().Err(err).Msg("get inventory item name")
		return Position{}, err
	}

	result := Position{
		ID:   item.ID,
		Name: item.Name,
		Qty:  pos.Qty,
	}

	return result, nil
}
