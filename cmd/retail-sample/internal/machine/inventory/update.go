package inventory

import (
	"strconv"

	"github.com/pkg/errors"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine"
	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

func (a *UseCase) UpdateStatus(recipeID string, enabled bool) (dto inventory.DTO, err error) {
	defer func() {
		if err != nil {
			a.logger.Error().Str("action", "update").Err(err).Send()
		}
	}()

	id, err := strconv.Atoi(recipeID)

	if err != nil {
		return inventory.DTO{}, errors.Wrapf(usecase.ErrBadRequest, "parse item ID: %s", recipeID)
	}

	dto, err = a.inventoryDB.Get(id)

	switch err {
	case nil:
	case inventory.ErrNotFound:
		return inventory.DTO{}, errors.Wrapf(usecase.ErrNotFound, "get item with id %d: %v", id, err)
	default:
		return
	}

	item := inventory.Item{
		ID:   dto.ID,
		Name: dto.Name,
		DB:   a.inventory.DB,
	}

	switch enabled {
	case true:
		err = item.Enable()
	default:
		err = item.Disable()
	}

	if err != nil {
		return inventory.DTO{}, err
	}

	dto.Enabled = enabled

	a.logger.Info().Int("id", id).Msg("successfully update inventory item")

	return dto, nil
}
