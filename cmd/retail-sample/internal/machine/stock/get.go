package stock

import (
	"strconv"

	"github.com/pkg/errors"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine"
	persistence "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

func (o *UseCase) GetByID(itemID string) (Position, error) {
	var err error
	defer func() {
		if err != nil {
			o.logger.Error().Str("action", "get by id").Err(err).Send()
		}
	}()

	id, err := strconv.Atoi(itemID)

	if err != nil {
		return Position{}, errors.Wrapf(usecase.ErrBadRequest, "parse item ID: %v", itemID)
	}

	item, err := o.stockDB.Get(id)

	switch err {
	case nil:
	case inventory.ErrNotFound:
		return Position{ID: id}, nil
	default:
		return Position{}, err
	}

	dto, err := o.inventoryDB.Get(item.InventoryID)

	switch err {
	case nil:
	case inventory.ErrNotFound:
		return Position{}, errors.Wrapf(usecase.ErrNotFound, "get item with id: %v", itemID)
	default:
		return Position{}, err
	}

	pos := Position{
		ID:   id,
		Name: dto.Name,
		Qty:  item.Qty,
	}

	return pos, nil
}

func (o *UseCase) GetAll() ([]Position, error) {
	var err error
	defer func() {
		if err != nil {
			o.logger.Error().Str("action", "get all").Err(err).Send()
		}
	}()

	pp, err := o.stockDB.List()

	if err != nil {
		return nil, err
	}

	var all = make([]Position, len(pp))

	for i := range pp {
		p := pp[i]

		stockItem, err := o.inventoryDB.Get(p.InventoryID)

		if err != nil {
			return nil, err
		}

		all[i] = Position{
			ID:   p.InventoryID,
			Qty:  p.Qty,
			Name: stockItem.Name,
		}
	}

	return all, nil
}

func (o *UseCase) GetProvisionLog() ([]persistence.ProvisionEntry, error) {
	items, err := o.logDB.List()

	if err != nil {
		o.logger.Error().Str("action", "get all").Err(err).Send()
		return nil, err
	}

	return items, nil
}
