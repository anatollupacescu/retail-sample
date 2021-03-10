package inventory

import (
	"github.com/pkg/errors"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine"
	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

func (a *UseCase) Create(name string) (inventory.DTO, error) {
	var err error

	defer func() {
		if err != nil {
			a.logger.Error().Str("action", "create item").Err(err).Send()
		}
	}()

	id, err := a.inventory.Create(name)

	switch err {
	case nil:
	case
		inventory.ErrEmptyName,
		inventory.ErrDuplicateName:
		return inventory.DTO{}, errors.Wrapf(usecase.ErrBadRequest, "create item with name '%s': %v", name, err)
	default:
		return inventory.DTO{}, err
	}

	a.logger.Info().Int("id", id).Msg("successfully created inventory item")

	item, err := a.inventoryDB.Get(id)

	if err != nil {
		return inventory.DTO{}, err
	}

	return item, nil
}
