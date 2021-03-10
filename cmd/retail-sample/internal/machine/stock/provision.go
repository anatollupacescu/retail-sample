package stock

import (
	"strconv"

	"github.com/pkg/errors"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine"
	pg "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

type Position struct {
	ID   int
	Name string
	Qty  int
}

func (o *UseCase) Provision(itemID string, qty int) (Position, error) {
	var err error

	defer func() {
		if err != nil {
			o.logger.Error().Str("action", "provision").Err(err).Send()
		}
	}()

	item, err := getInventoryItem(o.inventoryDB, itemID)

	if err != nil {
		return Position{}, err
	}

	pos, err := getPosition(o.stockDB, item.ID)

	if err != nil {
		return Position{}, err
	}

	err = pos.Provision(qty)

	switch err {
	case nil:
	case stock.ErrInvalidProvisionQuantity:
		return Position{}, errors.Wrap(usecase.ErrBadRequest, err.Error())
	default:
		return Position{}, err
	}

	o.logger.Info().Int("id", item.ID).Msg("successfully provisioned stock")

	_, err = o.logDB.Add(item.ID, qty)

	if err != nil {
		return Position{}, err
	}

	result := Position{
		ID:   item.ID,
		Name: item.Name,
		Qty:  pos.Qty,
	}

	return result, nil
}

func getPosition(stockDB *pg.StockPgxStore, itemID int) (pos stock.Position, err error) {
	pos.InventoryID = itemID
	pos.DB = stockDB

	stockPos, err := stockDB.Get(itemID)

	switch err {
	case nil:
		pos.Qty = stockPos.Qty
	case stock.ErrPositionNotFound: // first time provisioning, zero quantity
		return pos, nil
	default:
		return stock.Position{}, err
	}

	return
}

func getInventoryItem(inv *pg.InventoryPgxStore, itemID string) (inventory.DTO, error) {
	id, err := strconv.Atoi(itemID)

	if err != nil {
		return inventory.DTO{}, errors.Wrapf(usecase.ErrBadRequest, "parse item id: %s", itemID)
	}

	item, err := inv.Get(id)

	switch err {
	case nil: //continue,
	case inventory.ErrNotFound:
		return inventory.DTO{}, errors.Wrapf(usecase.ErrNotFound, "find item with id: %d", id)
	default:
		return inventory.DTO{}, err
	}

	return item, nil
}
