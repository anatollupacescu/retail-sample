package inventory

import (
	"strconv"

	"github.com/pkg/errors"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine"
	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

func (a *UseCase) GetByID(itemID string) (dto inventory.DTO, err error) {
	defer func() {
		if err != nil {
			a.logger.Error().Str("action", "get by id").Err(err).Send()
		}
	}()

	id, err := strconv.Atoi(itemID)

	if err != nil {
		return dto, errors.Wrapf(usecase.ErrBadRequest, "parse item ID: %v", itemID)
	}

	item, err := a.inventoryDB.Get(id)

	switch err {
	case nil:
	case inventory.ErrNotFound:
		return dto, errors.Wrapf(usecase.ErrNotFound, "find item with id: %v", id)
	default:
		return dto, err
	}

	return item, nil
}

func (a *UseCase) GetAll() ([]inventory.DTO, error) {
	items, err := a.inventoryDB.List()
	if err != nil {
		a.logger.Error().Str("action", "get all").Err(err).Send()
		return nil, err
	}

	return items, nil
}
