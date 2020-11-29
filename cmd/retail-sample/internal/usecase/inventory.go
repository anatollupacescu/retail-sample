package usecase

import (
	"context"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type inventoryDB interface {
	Get(int) (inventory.Item, error)
}

func NewInventory(ctx context.Context, inventory inventory.Inventory, db inventoryDB) Inventory {
	logger := log.Ctx(ctx).With().Str("layer", "usecase").Logger()

	return Inventory{
		ctx:         ctx,
		inventory:   inventory,
		inventoryDB: db,
		logger:      &logger,
	}
}

type Inventory struct {
	inventory   inventory.Inventory
	inventoryDB inventoryDB
	ctx         context.Context
	logger      *zerolog.Logger
}

type CreateInventoryItemDTO struct {
	Name string
}

func (a *Inventory) Create(in CreateInventoryItemDTO) (item inventory.Item, err error) {
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

func (a *Inventory) UpdateStatus(in UpdateInventoryItemStatusDTO) (item inventory.Item, err error) {
	if err = a.inventory.UpdateStatus(in.ID, in.Enabled); err != nil {
		a.logger.Error().Err(err).Msg("call domain")
		return inventory.Item{}, err
	}

	item, err = a.inventory.DB.Get(in.ID)

	if err != nil {
		a.logger.Error().Err(err).Msg("fetch updated item")
		return inventory.Item{}, err
	}

	return item, nil
}
