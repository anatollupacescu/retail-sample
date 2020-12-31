package usecase

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

type inventoryDB interface {
	Get(int) (inventory.ItemDTO, error)
}

func NewInventory(ctx context.Context, inventory inventory.Collection, db inventoryDB) Inventory {
	logger := log.Ctx(ctx).With().Str("layer", "usecase").Logger()

	return Inventory{
		ctx:         ctx,
		inventory:   inventory,
		inventoryDB: db,
		logger:      &logger,
	}
}

type Inventory struct {
	inventory   inventory.Collection
	inventoryDB inventoryDB
	ctx         context.Context
	logger      *zerolog.Logger
}

type CreateInventoryItemDTO struct {
	Name string
}

func (a *Inventory) Create(in CreateInventoryItemDTO) (item inventory.ItemDTO, err error) {
	id, err := a.inventory.Add(in.Name)

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

type UpdateInventoryItemStatusDTO struct {
	ID      int
	Enabled bool
}

func (a *Inventory) UpdateStatus(in UpdateInventoryItemStatusDTO) (dto inventory.ItemDTO, err error) {
	dto, err = a.inventoryDB.Get(in.ID)
	if err != nil {
		return
	}

	item := inventory.Item{
		ID: dto.ID, Name: dto.Name, DB: a.inventory.DB,
	}

	switch in.Enabled {
	case true:
		err = item.Enable()
	default:
		err = item.Disable()
	}

	if err != nil {
		a.logger.Error().Err(err).Msg("call domain")
		return inventory.ItemDTO{}, err
	}

	dto.Enabled = in.Enabled

	return dto, nil
}
