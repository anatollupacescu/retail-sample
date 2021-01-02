package inventory

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	pg "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

func New(ctx context.Context, t pg.TX) Inventory {
	logger := log.Ctx(ctx).With().Str("layer", "usecase").Logger()
	db := &pg.InventoryPgxStore{DB: t.Tx}

	return Inventory{
		ctx:         ctx,
		inventory:   &inventory.Collection{DB: db},
		inventoryDB: db,
		logger:      &logger,
	}
}

type Inventory struct {
	ctx         context.Context
	logger      *zerolog.Logger
	inventory   *inventory.Collection
	inventoryDB *pg.InventoryPgxStore
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
