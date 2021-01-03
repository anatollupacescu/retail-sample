package inventory

import "github.com/anatollupacescu/retail-sample/domain/retail/inventory"

type CreateInventoryItemDTO struct {
	Name string
}

func (a *UseCase) Create(in CreateInventoryItemDTO) (item inventory.DTO, err error) {
	id, err := a.inventory.Create(in.Name)

	if err != nil {
		a.logger.Error().Err(err).Msg("call domain")
		return
	}

	if item, err = a.inventoryDB.Get(id); err != nil {
		a.logger.Error().Err(err).Msg("retrieve new item")
		return
	}

	return item, nil
}
