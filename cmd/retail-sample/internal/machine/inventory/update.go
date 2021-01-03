package inventory

import "github.com/anatollupacescu/retail-sample/domain/retail/inventory"

func (a *UseCase) UpdateStatus(id int, enabled bool) (dto inventory.DTO, err error) {
	dto, err = a.inventoryDB.Get(id)
	if err != nil {
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
		a.logger.Error().Err(err).Msg("call domain")
		return inventory.DTO{}, err
	}

	dto.Enabled = enabled

	return dto, nil
}
